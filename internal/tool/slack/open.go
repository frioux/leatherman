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
Open opens a channel, group message, or direct message by name:

```bash
$ slack-open -channel general
```

Command: slack-open
*/
func Open(args []string, _ io.Reader) error {

	// https://api.slack.com/custom-integrations/legacy-tokens
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		return errors.New("SLACK_TOKEN is required")
	}

	cl := client{
		Token:  token,
		Client: &http.Client{},
	}

	var channel, conversationType string
	var exact, dryRun bool

	flags := flag.NewFlagSet("slack-open", flag.ExitOnError)
	flags.StringVar(&channel, "channel", "", "Channel to send open")
	flags.StringVar(&conversationType, "type", "public_channel", "Type of channel to send to (public_channel, private_channel, im, msim; public_channel is default.)")
	flags.BoolVar(&exact, "exact", false, "Set to disable regexp based channel matching")
	flags.BoolVar(&dryRun, "dry-run", false, "Set to not actually send message")
	flags.Parse(args[1:])

	if channel == "" {
		fmt.Fprint(os.Stderr, "-channel is required\n\n")
		flags.Usage()
		os.Exit(2)
	}

	team, err := cl.teamInfo("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't get team info: %s\n", err)
		os.Exit(1)
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

	fmt.Println("https://slack.com/app_redirect?team=" + team + "&channel=" + matched[0].ID)

	return nil
}
