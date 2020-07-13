package status

import (
	"fmt"
	"net/http"
)

type vpn struct{ value bool }

func (l *vpn) load() error {
	v, err := exec1Fail("pgrep", []string{"openvpn"})
	if err != nil {
		return err
	}
	l.value = v
	return nil
}

func (l *vpn) render(rw http.ResponseWriter) { fmt.Fprintf(rw, "%t\n", l.value) }
