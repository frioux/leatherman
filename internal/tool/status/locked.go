package status

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
)

type locked struct{ value bool }

func (l *locked) load() error {
	cmd := exec.Command("pgrep", "i3lock")
	_, err := cmd.Output()
	if err != nil {
		eErr := &exec.ExitError{}
		if errors.As(err, &eErr) {
			if eErr.ExitCode() == 1 {
				l.value = false
				return nil
			}
			return err
		}
		return err
	}
	l.value = true
	return nil
}

func (l *locked) render(rw http.ResponseWriter) { fmt.Fprintf(rw, "%t\n", l.value) }
