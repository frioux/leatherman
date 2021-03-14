package notes

import "io/fs"

type WatchDirFS interface {
	fs.FS

	Watch(path string) (func(), chan struct{}, error)
}
