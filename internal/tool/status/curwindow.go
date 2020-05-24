package status

import (
	"fmt"
	"net/http"
	"os/exec"
)

type curWindow struct{ value string }

func (l *curWindow) load() error {
	// curwindow is an auxiliary leatherman tool
	cmd := exec.Command("curwindow")
	b, err := cmd.Output()
	if err != nil {
		return err
	}
	l.value = string(b)
	return nil
}

func (l *curWindow) render(rw http.ResponseWriter) { fmt.Fprint(rw, l.value) }
