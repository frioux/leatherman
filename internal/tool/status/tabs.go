package status

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type tabs struct{ value string }

func (l *tabs) load() error {
	// XXX use a function since we can?
	cmd := exec.Command("dump-mozlz4", os.Getenv("MOZ_RECOVERY"))
	b, err := cmd.Output()
	if err != nil {
		return err
	}
	l.value = string(b)
	return nil
}

func (l *tabs) render(rw http.ResponseWriter) { fmt.Fprintf(rw, "%s\n", l.value) }
