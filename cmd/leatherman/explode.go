package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Explode all the tools as symlinks
func Explode(_ []string, _ io.Reader) error {
	exe, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "couldn't get Executable to explode")
	}
	dir := filepath.Dir(exe)
	for k := range Dispatch {
		if k == "help" {
			continue
		}
		if k == "explode" {
			continue
		}
		_ = os.Symlink(exe, dir+"/"+k)
	}

	return nil
}
