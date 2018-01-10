package main

import (
	"fmt"
	"os"
	"path"
)

func Explode(args []string) {
	exe, err := os.Executable()
	if err != nil {
		fmt.Println("couldn't get Executable to explode", err)
		os.Exit(1)
	}
	dir := path.Dir(exe)
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
