package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var drivers = map[string]driver{
	"wall":     wall,
	"pushover": pushover,
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "usage: %s <message>\n", os.Args[0])
		os.Exit(1)
	}

	message := strings.Join(os.Args[1:], " ")
	var failures int
	for n, d := range drivers {
		err := d(message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s failed: %s\n", n, message)
			failures++
		}
	}
	os.Exit(failures)
}

type driver func(string) error

func pushover(message string) error {
	return errors.New("not yet implemented")
}

func wall(message string) error {
	cmd := exec.Command("wall", "-t", "2", message)
	return cmd.Run()
}
