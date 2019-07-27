package slack

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/frioux/leatherman/internal/lmhttp"
)

/*
Deaddrop allows sending messages to a slack channel without looking at slack.
Typical usage is probably something like:

```bash
$ slack-deaddrop -channel general -text 'good morning!'
```

Command: slack-deaddrop
*/
func Deaddrop(args []string, _ io.Reader) error {
	// https://api.slack.com/custom-integrations/legacy-tokens
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		return errors.New("SLACK_TOKEN is required")
	}

	var channel, text, conversationType string
	var exact, dryRun bool

	flags := flag.NewFlagSet("slack-deaddrop", flag.ExitOnError)
	flags.StringVar(&channel, "channel", "", "Channel to send to")
	flags.StringVar(&conversationType, "type", "public_channel", "Type of channel to send to (public_channel, private_channel, im, msim; public_channel is default.)")
	flags.StringVar(&text, "text", "", "Text to send")
	flags.BoolVar(&exact, "exact", false, "Set to disable regexp based channel matching")
	flags.BoolVar(&dryRun, "dry-run", false, "Set to not actually send message")
	flags.Parse(args[1:])

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

	in := listConversationsInput{
		token:           token,
		limit:           200,
		excludeArchived: true,
		types:           conversationType,
	}
	req, err := listConversations(in)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	d := json.NewDecoder(resp.Body)
	var cs listConversationsOutput
	if err := d.Decode(&cs); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("list conversations failed: " + resp.Status)
	}

	channels := cs.Channels

	for cs.ResponseMetadata.NextCursor != "" {
		in.cursor = cs.ResponseMetadata.NextCursor

		// zero the struct, otherwise we get spooky action on channels
		cs = listConversationsOutput{}
		req, err := listConversations(in)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		d := json.NewDecoder(resp.Body)
		if err := d.Decode(&cs); err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return errors.New("list conversations failed: " + resp.Status)
		}
		channels = append(channels, cs.Channels...)
	}

	if !cs.OK {
		return errors.New("list conversations failed: " + cs.Error)
	}

	var channelMatches *regexp.Regexp

	if !exact {
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
	req, err = postMessage(postMessageInput{
		token:   token,
		channel: matched[0].ID,
		asUser:  true,
		text:    text,
	})
	if err != nil {
		return err
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	fmt.Println("")

	return nil
}

type postMessageInput struct {
	token, channel, text string
	asUser               bool
}

// https://api.slack.com/methods/chat.postMessage
func postMessage(i postMessageInput) (*http.Request, error) {
	v := url.Values{}
	v.Set("token", i.token)
	v.Set("channel", i.channel)
	v.Set("text", i.text)
	if i.asUser {
		v.Set("as_user", "true")
	}

	req, err := lmhttp.NewRequest("POST", "https://slack.com/api/chat.postMessage", strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

type listConversationsInput struct {
	token, cursor, types string
	excludeArchived      bool
	limit                int
}

type slackConversation struct {
	ID   string
	Name string

	IsChannel bool `json:"is_channel"`
	IsGroup   bool `json:"is_group"`
	IsIM      bool `json:"is_im"`
}

type listConversationsOutput struct {
	OK               bool
	Error            string // only set if OK is false
	Channels         []slackConversation
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

// https://api.slack.com/methods/conversations.list
func listConversations(i listConversationsInput) (*http.Request, error) {
	v := url.Values{}
	v.Set("token", i.token)
	if i.cursor != "" {
		v.Set("cursor", i.cursor)
	}
	v.Set("types", i.types)
	if i.excludeArchived {
		v.Set("exclude_archived", "true")
	}
	v.Set("limit", strconv.Itoa(i.limit))

	req, err := lmhttp.NewRequest("GET", "https://slack.com/api/conversations.list?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
