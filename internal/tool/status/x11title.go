package status

import (
	"fmt"
	"net/http"
	"os/exec"
)

type x11title struct {
	value string
}

func (v *x11title) load() error {
	val, err := exec.Command("xdotool", "getactivewindow", "getwindowname").CombinedOutput()
	if err != nil {
		return err
	}

	v.value = string(val)
	return nil
}

func (v *x11title) render(rw http.ResponseWriter) {
	fmt.Fprintln(rw, v.value)
}
