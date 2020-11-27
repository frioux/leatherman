package status

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/frioux/leatherman/internal/lmhttp"
)

type steambox struct{ game string }

var event = regexp.MustCompile(`^.*AppID (\d+) state changed : (.*),$`)

func (l *steambox) gameName(id string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := lmhttp.Get(ctx, "https://store.steampowered.com/api/appdetails/?filters=basic&appids="+id)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var details map[string]struct {
		Data struct {
			Name string
		}
	}
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&details); err != nil {
		return "", err
	}

	return details[id].Data.Name, nil
}

func (l *steambox) runningGame(r io.Reader) (string, error) {
	var ret string

	s := bufio.NewScanner(r)
LINE:
	for s.Scan() {
		m := event.FindSubmatch(s.Bytes())
		if m == nil {
			continue
		}

		appID := string(m[1])

		events := strings.Split(string(m[2]), ",")
		for _, e := range events {
			if e == "App Running" {
				ret = appID
				continue LINE
			}
		}
		if ret == appID { // didn't see App Running and since ret was already appID it stopped running
			ret = ""
		}
	}

	return ret, s.Err()

}

func (l *steambox) load() error {
	const (
		path = "/home/steam/.local/share/Steam/logs/content_log.txt"
	)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	appID, err := l.runningGame(f)
	if err != nil {
		return err
	}

	l.game, err = l.gameName(appID)
	if err != nil {
		return err
	}

	return nil
}

func (l *steambox) render(rw http.ResponseWriter) {
	fmt.Fprintf(rw, "%s\n", l.game)
}
