package fn

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/frioux/leatherman/pkg/shellquote"
	"github.com/pkg/errors"
)

var dir = os.Getenv("HOME") + "/code/dotfiles/bin"

// Run generates shell scripts based on the args
func Run(args []string, _ io.Reader) error {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s $scriptname [-f] $command $tokens\n", args[0])
		os.Exit(1)
	}

	script := dir + "/" + args[1]

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

	if err := ioutil.WriteFile(script, []byte("#!/bin/sh\n\n"+body+"\n"), 0755); err != nil {
		return errors.Wrap(err, "Couldn't create new script")
	}

	return nil
}
