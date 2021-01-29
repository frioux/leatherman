package zine

import (
	"fmt"
	"os"
)

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
