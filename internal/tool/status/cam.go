package status

import (
	"fmt"
	"net/http"
)

type cam struct{ value bool }

func (v *cam) load() error {
	val, err := exec1Fail("lsof", "/dev/video0")
	if err != nil {
		return err
	}

	v.value = val
	return nil
}

func (v *cam) render(rw http.ResponseWriter) { fmt.Fprintf(rw, "%t\n", v.value) }
