package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	startMatcher = regexp.MustCompile(`index:\s`)
	splitter     = regexp.MustCompile(`^\s*(.+?)\s*[:=]\s*(.+?)\s*$`)
)

func sound() bool {
	is, err := listSinkInputs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	for _, i := range is {
		if i["state"] == "RUNNING" {
			return true
		}
	}

	return false
}

func listSinkInputs() ([]map[string]string, error) {
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
