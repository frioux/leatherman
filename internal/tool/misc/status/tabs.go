package status

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/frioux/leatherman/pkg/mozlz4"
)

type tabs struct{ value string }

func (l *tabs) load() error {
	f, err := os.Open(os.Getenv("MOZ_RECOVERY"))
	if err != nil {
		return err
	}

	r, err := mozlz4.NewReader(f)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	l.value = string(b)
	return nil
}

func (l *tabs) render(rw http.ResponseWriter) { fmt.Fprintf(rw, "%s\n", l.value) }
