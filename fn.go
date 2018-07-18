package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/frioux/shellquote"
)

// Fn generates shell scripts based on the args
func Fn(args []string, _ io.Reader) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s $scriptname [-f] $command $tokens\n", args[0])
		os.Exit(1)
	}

	script := os.Getenv("HOME") + "/code/dotfiles/bin/" + args[1]

	if args[2] == "-f" {
		os.Remove(script)
		args = append(args[:2], args[3:]...)
	}

	var body string
	if len(args[2:]) == 1 {
		body = args[2]
	} else {
		var err error
		body, err = shellquote.Quote(args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't quote args to script script: %s\n", err)
			os.Exit(1)
		}
	}

	// If script exists or we can't stat it
	stat, err := os.Stat(script)
	if stat != nil {
		fmt.Fprintf(os.Stderr, "Script (%s) already exists\n", script)
		os.Exit(1)
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Couldn't stat new script: %s\n", err)
		os.Exit(1)
	}

	file, err := os.Create(script)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create new script: %s\n", err)
		os.Exit(1)
	}

	w := bufio.NewWriter(file)
	_, err = w.WriteString("#!/bin/sh\n\n" + body + "\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't write to new script: %s\n", err)
		os.Exit(1)
	}
	err = w.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't flush new script: %s\n", err)
		os.Exit(1)
	}

	err = file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't save new script: %s\n", err)
		os.Exit(1)
	}
	err = os.Chmod(script, os.FileMode(0755))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't chown new script: %s\n", err)
		os.Exit(1)
	}
}
