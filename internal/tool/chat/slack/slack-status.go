package slack

import (
	"errors"
	"flag"
	"io"
	"net/http"
	"os"
	"time"
)

func Status(args []string, _ io.Reader) error {
	// https://api.slack.com/custom-integrations/legacy-tokens
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		return errors.New("SLACK_TOKEN is required")
	}

	var (
		text, emoji string
		expiration  time.Duration
		debug       bool
	)

	flags := flag.NewFlagSet("slack-status", flag.ExitOnError)
	flags.StringVar(&text, "text", "", "text to set status to")
	flags.StringVar(&emoji, "emoji", "", "emoji to set status to")
	flags.DurationVar(&expiration, "expiration", time.Duration(0), "when to expire status")
	flags.Parse(args[1:])

	cl := client{
		Token:  token,
		Client: &http.Client{},
		debug:  debug,
	}
	i := usersProfileSetInput{
		StatusText:  text,
		StatusEmoji: emoji,
	}

	if expiration != time.Duration(0) {
		i.StatusExpiration = time.Now().Add(expiration).Unix()
	}

	err := cl.usersProfileSet(i)
	return err
}
