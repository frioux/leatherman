package wuphf

import (
	"errors"
	"fmt"
	"io"
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

/*
Wuphf sends alerts via both `wall` and [pushover](https://pushover.net).  All
arguments are concatenated to produce the sent message.

The following environment variables should be set:

 * LM_PUSHOVER_TOKEN
 * LM_PUSHOVER_USER
 * LM_PUSHOVER_DEVICE

```bash
$ wuphf 'the shoes are on sale'
```
*/
func Wuphf(args []string, _ io.Reader) error {
	if len(args) == 1 {
		return fmt.Errorf("usage: %s <message>\n", args[0])
	}

	message := strings.Join(args[1:], " ")
	for n, d := range drivers {
		err := d(message)
		if err != nil {
			return fmt.Errorf("%s failed: %s\n", n, err)
		}
	}

	return nil
}

type driver func(string) error

var (
	errNoPushoverToken  = errors.New("LM_PUSHOVER_TOKEN not set")
	errNoPushoverUser   = errors.New("LM_PUSHOVER_USER not set")
	errNoPushoverDevice = errors.New("LM_PUSHOVER_DEVICE not set")
)

func pushover(message string) error {
	token := os.Getenv("LM_PUSHOVER_TOKEN")
	if token == "" {
		return errNoPushoverToken
	}
	user := os.Getenv("LM_PUSHOVER_USER")
	if user == "" {
		return errNoPushoverUser
	}
	device := os.Getenv("LM_PUSHOVER_DEVICE")
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
