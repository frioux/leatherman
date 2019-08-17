package slack

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/frioux/leatherman/internal/lmhttp"
)

type slackConversation struct {
	ID   string
	Name string
}

type client struct {
	Token string
	*http.Client
}

type usersListInput struct {
	cursor string
	limit  int
}

type usersListOutput struct {
	OK               bool
	Error            string // only set if OK is false
	Members          []slackConversation
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

// https://api.slack.com/methods/users.list
func (c client) usersList(i usersListInput) (usersListOutput, error) {
	v := url.Values{}
	v.Set("token", c.Token)
	if i.cursor != "" {
		v.Set("cursor", i.cursor)
	}
	v.Set("limit", strconv.Itoa(i.limit))

	req, err := lmhttp.NewRequest("GET", "https://slack.com/api/users.list?"+v.Encode(), nil)
	if err != nil {
		return usersListOutput{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return usersListOutput{}, err
	}

	if resp.StatusCode != 200 {
		return usersListOutput{}, errors.New("list conversations failed: " + resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	var cs usersListOutput
	if err := d.Decode(&cs); err != nil {
		return usersListOutput{}, err
	}

	if !cs.OK {
		return usersListOutput{}, errors.New("list conversations failed: " + cs.Error)
	}

	return cs, nil
}

type conversationsListInput struct {
	cursor, types   string
	excludeArchived bool
	limit           int
}

type conversationsListOutput struct {
	OK               bool
	Error            string // only set if OK is false
	Channels         []slackConversation
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

// https://api.slack.com/methods/conversations.list
func (c client) conversationsList(i conversationsListInput) (conversationsListOutput, error) {
	v := url.Values{}
	v.Set("token", c.Token)
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
		return conversationsListOutput{}, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return conversationsListOutput{}, err
	}

	if resp.StatusCode != 200 {
		return conversationsListOutput{}, errors.New("list conversations failed: " + resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	var cs conversationsListOutput
	if err := d.Decode(&cs); err != nil {
		return conversationsListOutput{}, err
	}

	if !cs.OK {
		return conversationsListOutput{}, errors.New("list conversations failed: " + cs.Error)
	}

	return cs, nil
}

type chatPostMessageInput struct {
	channel, text string
	asUser        bool
}

// https://api.slack.com/methods/chat.postMessage
func (c client) chatPostMessage(i chatPostMessageInput) (*http.Response, error) {
	v := url.Values{}
	v.Set("token", c.Token)
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

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
