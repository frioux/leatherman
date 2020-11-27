package status

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type retropie struct {
	console, emu, game, command string
}

func (v *retropie) load() error {
	f, err := os.Open("/run/shm/runcommand.info")
	if err != nil {
		return err
	}

	s := bufio.NewScanner(f)

	i := 0
	for s.Scan() {
		i++

		switch i {
		case 1:
			v.console = s.Text()
		case 2:
			v.emu = s.Text()
		case 3:
			v.game = s.Text()
		case 4:
			v.command = s.Text()
		default:
			return errors.New("runcommand.info longer than expected")
		}

	}

	return nil
}

func (v *retropie) render(rw http.ResponseWriter) {
	fmt.Fprintf(rw, "console=%s, emu=%s, game=%s, command=%s\n", v.console, v.emu, v.game, v.command)
}
