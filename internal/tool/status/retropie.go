package status

import (
	"bufio"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type retropie struct {
	Console, Emu, Game, Command string
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
			v.Console = s.Text()
		case 2:
			v.Emu = s.Text()
		case 3:
			v.Game = s.Text()
		case 4:
			v.Command = s.Text()
		default:
			return errors.New("runcommand.info longer than expected")
		}

	}

	return nil
}

func (v *retropie) render(rw http.ResponseWriter) {
	e := json.NewEncoder(rw)
	e.Encode(v)
}
