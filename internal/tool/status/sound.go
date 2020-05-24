package status

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	startMatcher = regexp.MustCompile(`index:\s`)
	splitter     = regexp.MustCompile(`^\s*(.+?)\s*[:=]\s*(.+?)\s*$`)
)

type sound struct{ value bool }

func (v *sound) load() error {
	is, err := v.listSinkInputs()
	if err != nil {
		return err
	}

	for _, i := range is {
		if i["state"] == "RUNNING" {
			v.value = true
			return nil
		}
	}

	v.value = false
	return nil
}

func (v *sound) listSinkInputs() ([]map[string]string, error) {
	c := exec.Command("pacmd", "list-sink-inputs")
	c.Stderr = os.Stderr
	out, err := c.Output()
	if err != nil {
		return nil, fmt.Errorf("pacmd list-sink-inputs: %w", err)
	}

	r := bytes.NewBuffer(out)
	s := bufio.NewScanner(r)

	var ret []map[string]string
	var current map[string]string
	for s.Scan() {
		l := s.Text()
		if startMatcher.MatchString(l) {
			if current != nil {
				ret = append(ret, current)
			}
			current = map[string]string{}
		}

		matches := splitter.FindStringSubmatch(l)
		if len(matches) == 3 {
			matches[2] = strings.TrimSuffix(matches[2], `"`)
			matches[2] = strings.TrimPrefix(matches[2], `"`)
			current[matches[1]] = matches[2]
		}
	}
	ret = append(ret, current)

	return ret, nil
}

func (v *sound) render(rw http.ResponseWriter) { fmt.Fprintf(rw, "%t\n", v.value) }
