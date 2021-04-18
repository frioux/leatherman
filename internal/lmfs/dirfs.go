package lmfs

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
)

func OpenDirFS(d string) fs.FS {
	return dirFS(d)
}

// dirFS, containsAny, and dirFS.Open were all copied verbatim from 1.16's os.

type dirFS string

func containsAny(s, chars string) bool {
	for i := 0; i < len(s); i++ {
		for j := 0; j < len(chars); j++ {
			if s[i] == chars[j] {
				return true
			}
		}
	}
	return false
}

func (dir dirFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) || runtime.GOOS == "windows" && containsAny(name, `\:`) {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
	}
	f, err := os.Open(string(dir) + "/" + name)
	if err != nil {
		return nil, err // nil fs.File
	}
	return f, nil
}

func (fss dirFS) Create(f string) (FileWriter, error) {
	return os.Create(filepath.Join(string(fss), f))
}

func (fss dirFS) Sub(d string) (fs.FS, error) {
	return OpenDirFS(filepath.Join(string(fss), d)), nil
}

func addDir(watcher *fsnotify.Watcher, path string) error {
	return filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if err := watcher.Add(path); err != nil {
			return fmt.Errorf("lmfs.dirFS: watcher.Add: %w", err)
		}
		return nil
	})
}

func (fss dirFS) Remove(name string) error { return os.Remove(string(fss) + "/" + name) }

func (fss dirFS) Watch(ctx context.Context, path string) (chan []fsnotify.Event, error) {
	watchroot := filepath.Join(string(fss), path)
	ch := make(chan []fsnotify.Event)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("fsnotify.NewWatcher: %w", err)
	}

	go func() {
		<-ctx.Done()
		watcher.Close()
	}()

	if err := addDir(watcher, watchroot); err != nil {
		return nil, err
	}

	go func() {
	LOOP:
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					fmt.Fprintln(os.Stderr, "watcher went away")
					break LOOP
				}

				r, err := filepath.Rel(watchroot, event.Name)
				if err != nil {
					panic(err)
				}
				ch <- []fsnotify.Event{{Op: event.Op, Name: r}}

				if event.Op&fsnotify.Create == fsnotify.Create {
					stat, err := os.Stat(event.Name)
					if err != nil {
						if os.IsNotExist(err) {
							continue LOOP
						}
						fmt.Fprintf(os.Stderr, "Couldn't stat created thing: %s\n", err)
					} else if stat.IsDir() {
						if err := addDir(watcher, event.Name); err != nil {
							fmt.Fprintf(os.Stderr, "failed to watch %s: %s\n", event.Name, err)
						}
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Fprintln(os.Stderr, "watcher went away")
					break LOOP
				}
				fmt.Fprintln(os.Stderr, "watcher error:", err)

			case <-ctx.Done():
				break LOOP
			}
		}
	}()

	return ch, nil
}
