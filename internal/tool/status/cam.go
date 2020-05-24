package status

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
)

type cam struct{ value bool }

func (v *cam) load() error {
	c := exec.Command("lsof", "/dev/video0")
	_, err := c.Output()
	if err != nil {
		eErr := &exec.ExitError{}
		if errors.As(err, &eErr) {
			if eErr.ExitCode() == 1 {
				v.value = false
				return nil
			}
			return err
		}
		return err
	}
	v.value = true
	return nil
}

func (v *cam) render(rw http.ResponseWriter) { fmt.Fprintf(rw, "%t\n", v.value) }
