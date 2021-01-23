package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var drivers = map[string]driver{
	"wall":     wall,
	"pushover": pushover,
	"notify":   notify,
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
			fmt.Fprintf(os.Stderr, "%s failed: %s\n", n, err)
			failures++
		}
	}
	os.Exit(failures)
}

type driver func(string) error

var (
	errNoPushoverToken  = errors.New("PUSHOVER_TOKEN not set")
	errNoPushoverUser   = errors.New("PUSHOVER_USER not set")
	errNoPushoverDevice = errors.New("PUSHOVER_DEVICE not set")
)

func pushover(message string) error {
	token := os.Getenv("PUSHOVER_TOKEN")
	if token == "" {
		return errNoPushoverToken
	}
	user := os.Getenv("PUSHOVER_USER")
	if user == "" {
		return errNoPushoverUser
	}
	device := os.Getenv("PUSHOVER_DEVICE")
	if device == "" {
		return errNoPushoverDevice
	}

	resp, err := http.PostForm("https://api.pushover.net/1/messages.json", url.Values{
		"token":   {token},
		"user":    {user},
		"message": {message},
		"device":  {device},
	})
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to pushover: %s (%s)", err, resp.Status)
	}

	return nil
}

func wall(message string) error {
	cmd := exec.Command("wall", "-t", "2", message)
	return cmd.Run()
}

func notify(message string) error {
	cmd := exec.Command("notify-send", "-u", "critical", "wuphf", message)
	return cmd.Run()
}
