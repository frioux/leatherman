package fn

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/frioux/shellquote"
	"github.com/pkg/errors"
)

// Run generates shell scripts based on the args
func Run(args []string, _ io.Reader) error {
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
			return errors.Wrap(err, "Couldn't quote args to script script")
		}
	}

	// If script exists or we can't stat it
	stat, err := os.Stat(script)
	if stat != nil {
		return errors.Wrap(err, "Script ("+script+") already exists")
	} else if !os.IsNotExist(err) {
		return errors.Wrap(err, "Couldn't stat new script")
	}

	file, err := os.Create(script)
	if err != nil {
		return errors.Wrap(err, "Couldn't create new script")
	}

	w := bufio.NewWriter(file)
	_, err = w.WriteString("#!/bin/sh\n\n" + body + "\n")
	if err != nil {
		return errors.Wrap(err, "Couldn't write to new script")
	}
	err = w.Flush()
	if err != nil {
		return errors.Wrap(err, "Couldn't flush new script")
	}

	err = file.Close()
	if err != nil {
		return errors.Wrap(err, "Couldn't save new script")
	}
	err = os.Chmod(script, os.FileMode(0755))
	if err != nil {
		return errors.Wrap(err, "Couldn't chown new script")
	}

	return nil
}
