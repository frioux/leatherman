// +build linux

package zine

import (
	"fmt"
	"io"

	_ "modernc.org/sqlite"
)

/*
Run does read only operations on notes.

Command: zine
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
