package slack

import (
	"errors"
	"flag"
	"io"
	"net/http"
	"os"
)

func Status(args []string, _ io.Reader) error {
	// https://api.slack.com/custom-integrations/legacy-tokens
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		return errors.New("SLACK_TOKEN is required")
	}

	var text, emoji string
	// var expiration easyTime
	var debug bool

	flags := flag.NewFlagSet("slack-status", flag.ExitOnError)
	flags.StringVar(&text, "text", "", "text to set status to")
	flags.StringVar(&emoji, "emoji", "", "emoji to set status to")
	flags.Parse(args[1:])

	cl := client{
		Token:  token,
		Client: &http.Client{},
		debug:  debug,
	}

	err := cl.usersProfileSet(usersProfileSetInput{
		StatusText:  text,
		StatusEmoji: emoji,
	})
	if err != nil {
		return err
	}

	return nil
}
