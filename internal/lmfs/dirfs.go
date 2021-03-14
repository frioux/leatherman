package lmfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
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
