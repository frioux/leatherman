package slack

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

/*
Deaddrop allows sending messages to a slack channel without looking at slack.
Typical usage is probably something like:

```bash
$ slack-deaddrop -channel general -text 'good morning!'
```
*/
func Deaddrop(args []string, _ io.Reader) error {
	// https://api.slack.com/custom-integrations/legacy-tokens
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		return errors.New("SLACK_TOKEN is required")
	}

	var channel, text, conversationType string
	var debug, exact, dryRun bool

	flags := flag.NewFlagSet("slack-deaddrop", flag.ExitOnError)
	flags.StringVar(&channel, "channel", "", "Channel to send to")
	flags.StringVar(&conversationType, "type", "public_channel", "Type of channel to send to (public_channel, private_channel, im, msim; public_channel is default.)")
	flags.StringVar(&text, "text", "", "Text to send")
	flags.BoolVar(&exact, "exact", false, "Set to disable regexp based channel matching")
	flags.BoolVar(&dryRun, "dry-run", false, "Set to not actually send message")
	flags.BoolVar(&debug, "debug", false, "Print full HTTP conversation")
	flags.Parse(args[1:])

	cl := client{
		Token:  token,
		Client: &http.Client{},
		debug:  debug,
	}

	if channel == "" {
		fmt.Fprint(os.Stderr, "-channel is required\n\n")
		flags.Usage()
		os.Exit(2)
	}

	if text == "" {
		fmt.Fprint(os.Stderr, "-text is required\n\n")
		flags.Usage()
		os.Exit(2)
	}

	var channels []slackConversation
	if conversationType == "im" {
		var err error
		channels, err = cl.autopageUsersList(usersListInput{limit: 200})
		if err != nil {
			return err
		}
	} else {
		var err error
		channels, err = cl.autopageConversationsList(conversationsListInput{
			limit:           200,
			excludeArchived: true,
			types:           conversationType,
		})
		if err != nil {
			return err
		}
	}

	var channelMatches *regexp.Regexp

	if !exact {
		var err error
		channelMatches, err = regexp.Compile(channel)
		if err != nil {
			return err
		}
	}
	matched := make([]slackConversation, 0, 1)
	for _, c := range channels {
		if !exact && channelMatches.MatchString(c.Name) {
			matched = append(matched, c)
		}
		if exact && c.Name == channel {
			matched = append(matched, c)
		}
	}
	if len(matched) == 0 {
		return errors.New("no channels matched " + channel)
	}
	if len(matched) != 1 {
		names := make([]string, 0, len(matched))
		for _, m := range matched {
			names = append(names, " * "+m.Name+"\n")
		}
		sort.Strings(names)
		return errors.New("too many channels matched: \n" + strings.Join(names, ""))
	}

	if dryRun {
		fmt.Fprintf(os.Stderr, "Would send «%s» to #%s...\n", text, matched[0].Name)
		return nil
	}
	fmt.Fprintf(os.Stderr, "Sending «%s» to #%s...\n", text, matched[0].Name)
	resp, err := cl.chatPostMessage(chatPostMessageInput{
		channel: matched[0].ID,
		asUser:  true,
		text:    text,
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	fmt.Println("")

	return nil
}
