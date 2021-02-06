// +build linux

package zine

import (
	"fmt"
	"io"
	"os"
)

/*
Run does read only operations on notes.
*/
func Run(args []string, _ io.Reader) error {
	command := "render"
	if len(args) > 1 {
		command = args[1]
	}

	cmd, ok := commands[command]
	if !ok {
		return fmt.Errorf("unknown command «%s»; valid commands are 'render' and 'q'\n", command)
	}

	if err := cmd(args[1:]); err != nil {
		return err
	}

	return nil
}

func run() error {
	command := "render"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	cmd, ok := commands[command]
	if !ok {
		return fmt.Errorf("unknown command «%s»; valid commands are 'render' and 'q'\n", command)
	}

	if err := cmd(os.Args[1:]); err != nil {
		return err
	}

	return nil
}
