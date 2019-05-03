package main

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

// Explode all the tools as symlinks
func Explode(_ []string, _ io.Reader) error {
	exe, err := os.Executable()
	if err != nil {
		return xerrors.Errorf("Couldn't get Executable to explode: %w", err)
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
