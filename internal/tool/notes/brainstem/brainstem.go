package brainstem

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/frioux/leatherman/internal/notes"
)

/*
Brainstem allows interacting with amygdala without using any of the server
components, typically for testing the personality etc, but can also be used as
a lightweight amygdala instance.
*/
func Brainstem(args []string, _ io.Reader) error {
	var tok string

	tok = os.Getenv("DROPBOX_ACCESS_TOKEN")
	if tok == "" {
		return errors.New("DROPBOX_ACCESS_TOKEN is unset")
	}

	rules, err := notes.NewRules(tok)
	if err != nil {
		return fmt.Errorf("Couldn't create rules: %s\n", err)
	}

	if len(args) < 2 {
		return fmt.Errorf("usage: %s <cmd>\n", args[0])
	}
	message, err := rules.Dispatch(args[1], nil)
	fmt.Println(message)
	return err
}
