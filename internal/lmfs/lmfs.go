package lmfs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"

	"github.com/fsnotify/fsnotify"
)

type WriteFileFS interface {
	fs.FS

	WriteFile(string, []byte, fs.FileMode) error
}

type CreateFS interface {
	fs.FS

	Create(string) (FileWriter, error)
}

// don't love the name here, but want a WriteFile function
type FileWriter interface {
	fs.File

	Write([]byte) (int, error)
}

var errUnsupported = errors.New("filesystem does not support requested operation")

func WriteFile(fss fs.FS, name string, contents []byte, mode fs.FileMode) (err error) {
	if wffs, ok := fss.(WriteFileFS); ok {
		return wffs.WriteFile(name, contents, mode)
	}

	// XXX this leg doesn't handle the mode (yet)
	if cfs, ok := fss.(CreateFS); ok {
		f, err := cfs.Create(name)
		if err != nil {
			return err
		}
		defer func() {
			if err == nil {
				err = f.Close()
			} else {
				f.Close()
			}
		}()

		if _, err := io.Copy(f, bytes.NewReader(contents)); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("couldn't write file via %T: %w", fss, errUnsupported)
}

// XXX: add lmfs.DirFS for WriteFile (and watch) support

type WatchFS interface {
	fs.FS

	Watch(context.Context, string) (chan []fsnotify.Event, error)
}

func Watch(ctx context.Context, fss fs.FS, dir string) (chan []fsnotify.Event, error) {
	wfs, ok := fss.(WatchFS)
	if !ok {
		return nil, fmt.Errorf("couldn't watch dir via %T: %w", fss, errUnsupported)
	}

	return wfs.Watch(ctx, dir)
}

type RemoveFS interface {
	fs.FS
	Remove(string) error
}

func Remove(fss fs.FS, path string) error {
	rfs, ok := fss.(RemoveFS)
	if !ok {
		return errUnsupported
	}

	return rfs.Remove(path)
}
