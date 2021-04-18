package dropbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	errNotADir  = errors.New("not a directory")
	errNotAFile = errors.New("not a file")
)

func (c Client) AsFS(ctx context.Context) FS { return FS{c, ctx, "/"} }

type FS struct {
	cl  Client
	ctx context.Context
	dir string
}

func (f FS) getmeta(path string) (*DropboxDirEntry, error) {
	s, err := f.cl.GetMetadata(GetMetadataParams{Path: path})
	if err != nil {
		return nil, err
	}
	return &DropboxDirEntry{
		fs:       f,
		Metadata: s,
		dir:      strings.TrimSuffix(path, s.Name),
	}, nil
}

func (f FS) Open(name string) (fs.File, error) { return f.getmeta(name) }

func (f FS) Stat(name string) (fs.FileInfo, error) { return f.getmeta(name) }

func (f FS) ReadDir(name string) ([]fs.DirEntry, error) { return f.readDir(name, -1) }

func (f FS) readDir(name string, n int) ([]fs.DirEntry, error) {
	name = filepath.Join(f.dir, name)
	var (
		ret             []fs.DirEntry
		listFolderLimit uint32
	)
	if n > 0 {
		ret = make([]fs.DirEntry, 0, n)
		listFolderLimit = uint32(n)
	}
	r, err := f.cl.ListFolder(ListFolderParams{Path: name, Limit: listFolderLimit})
	if err != nil {
		return nil, fmt.Errorf("path=%s %w", name, err)
	}

	entries := r.Entries

	for _, en := range entries {
		ret = append(ret, &DropboxDirEntry{
			fs:       f,
			Metadata: en,
			dir:      name,
		})
		if n > 0 && len(ret) == n {
			return ret, nil
		}
	}

	for r.HasMore {
		r, err = f.cl.ListFolderContinue(r.Cursor)
		if err != nil {
			return nil, err
		}

		for _, en := range r.Entries {
			ret = append(ret, &DropboxDirEntry{
				fs:       f,
				Metadata: en,
				dir:      name,
			})
			if n > 0 && len(ret) == n {
				return ret, nil
			}
		}

	}

	return ret, nil
}

func (f FS) ReadFile(name string) ([]byte, error) {
	b, err := f.cl.Download(filepath.Join(f.dir, name))
	if err != nil {
		return nil, fmt.Errorf("path=%s %w", name, err)
	}

	return b, nil
}

func (f FS) Remove(name string) error { return f.cl.Delete(name) }

func (f FS) Watch(ctx context.Context, dir string) (chan []fsnotify.Event, error) {

	dropboxCh := make(chan []Metadata)
	go f.cl.Longpoll(ctx, filepath.Join(f.dir, dir), dropboxCh)

	ch := make(chan []fsnotify.Event)
	go func() {
		for metadata := range dropboxCh {
			events := make([]fsnotify.Event, len(metadata))
			for i, m := range metadata {
				var err error
				events[i].Name, err = filepath.Rel(filepath.Join(f.dir, dir), m.PathLower)
				if err != nil {
					fmt.Fprintln(os.Stderr, "uhhh", err)
					continue
				}
				switch m.Tag {
				case "deleted":
					events[i].Op = fsnotify.Remove
				case "file":
					events[i].Op = fsnotify.Create
				case "folder":
					events[i].Op = fsnotify.Create
				default:
					panic("unknown longpoll tag: " + m.Tag)
				}
			}
			ch <- events
		}
	}()

	return ch, nil
}

func (f FS) Sub(dir string) (fs.FS, error) {
	return FS{
		cl:  f.cl,
		ctx: f.ctx,
		dir: filepath.Join(f.dir, dir),
	}, nil
}

func (f FS) WriteFile(name string, contents []byte, _ fs.FileMode) error {
	return f.cl.Create(UploadParams{
		Path: filepath.Join(f.dir, name),
		Mode: "overwrite",
	}, bytes.NewReader(contents))
}

type DropboxDirEntry struct {
	fs FS

	dir string

	Metadata
	io.Reader
}

func (e *DropboxDirEntry) Read(b []byte) (int, error) {
	if e.IsDir() {
		return 0, errNotAFile
	}

	if e.Reader == nil {
		bs, err := e.fs.ReadFile(e.dir + "/" + e.Name())
		if err != nil {
			return 0, err
		}

		e.Reader = bytes.NewReader(bs)
		return e.Reader.Read(b)
	}

	return e.Reader.Read(b)
}

func (e *DropboxDirEntry) Close() error { return nil }

func (e *DropboxDirEntry) Stat() (fs.FileInfo, error) { return e, nil }

func (e *DropboxDirEntry) ReadDir(n int) ([]fs.DirEntry, error) {
	if !e.IsDir() {
		return nil, errNotADir
	}

	return e.fs.readDir(e.dir+e.Name(), n)
}

func (e *DropboxDirEntry) Name() string { return e.Metadata.Name }

func (e *DropboxDirEntry) IsDir() bool { return e.Metadata.Tag == "folder" }

func (e *DropboxDirEntry) Size() int64 { return int64(e.Metadata.Size) }

func (e *DropboxDirEntry) Type() fs.FileMode {
	// copied from embed/embed.go
	if e.IsDir() {
		return fs.ModeDir | 0555
	}
	return 0444
}

func (e *DropboxDirEntry) Mode() fs.FileMode { return e.Type() }

func (e *DropboxDirEntry) ModTime() time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05Z", e.ServerModified)
	return t
}

func (e *DropboxDirEntry) Info() (fs.FileInfo, error) { return e, nil }

func (e *DropboxDirEntry) Sys() interface{} { return e }
