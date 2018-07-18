package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Explode all the tools as symlinks
func Explode(args []string, _ io.Reader) {
	exe, err := os.Executable()
	if err != nil {
		fmt.Println("couldn't get Executable to explode", err)
		os.Exit(1)
	}
	dir := filepath.Dir(exe)
	for k := range Dispatch {
		if k == "help" {
			continue
		}
		if k == "explode" {
			continue
		}
		err := os.Symlink(exe, dir+"/"+k)
		if err != nil {
			// ignore for now
		}
	}
}
