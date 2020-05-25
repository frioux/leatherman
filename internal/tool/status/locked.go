package status

import (
	"fmt"
	"net/http"
)

type locked struct{ value bool }

func (l *locked) load() error {
	v, err := exec1Fail("pgrep", "i3lock")
	if err != nil {
		return err
	}
	l.value = v
	return nil
}

func (l *locked) render(rw http.ResponseWriter) { fmt.Fprintf(rw, "%t\n", l.value) }
