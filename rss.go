package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"

	"github.com/mmcdole/gofeed"
)

func RSS(args []string) {
	if len(args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s feedURL statefile\n", args[0])
		os.Exit(1)
	}

	statePath := args[2]

	fp := gofeed.NewParser()
	feedURL, err := url.Parse(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse feed url: %s\n", err)
		os.Exit(1)
	}
	f, err := fp.ParseURL(feedURL.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't fetch feed: %s\n", err)
		os.Exit(1)
	}

	seen, err := syncRead(statePath, f.Items)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't sync read: %s\n", err)
		os.Exit(1)
	}

	items := newItems(seen, f.Items)

	fixLinks(feedURL, items)
	renderItems(os.Stdout, items)

	err = os.Rename(statePath+".tmp", statePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't rename state file: %s\n", err)
		os.Exit(1)
	}
}

// fixLinks adds hostname and schema from feed link to item links
func fixLinks(feedURL *url.URL, items []*gofeed.Item) {
	for _, i := range items {
		itemURL, _ := url.Parse(i.Link)
		if itemURL.Hostname() == "" {
			itemURL.Host = feedURL.Hostname()
		}
		if itemURL.Scheme == "" {
			itemURL.Scheme = feedURL.Scheme
		}
		i.Link = itemURL.String()
	}
}

func renderItems(out io.Writer, items []*gofeed.Item) {
	for _, i := range items {
		fmt.Fprintf(out, "[%s](%s)\n", i.Title, i.Link)
	}
}

// Return items in feed that are not in sync
func newItems(seen map[string]bool, items []*gofeed.Item) []*gofeed.Item {
	ret := make([]*gofeed.Item, 0, len(items))

	for _, i := range items {
		if !seen[i.GUID] {
			ret = append(ret, i)
		}
	}

	return ret
}

// Store JSON containing seen GUIDs for the current feed.
func syncRead(state string, items []*gofeed.Item) (map[string]bool, error) {
	ret := make(map[string]bool, len(items))

	guids, err := readState(state)
	if err != nil {
		return nil, fmt.Errorf("couldn't read state: %s", err)
	}

	for _, g := range guids {
		ret[g] = true
	}

	// Generate news state
	newState := make(map[string]bool, len(items)+len(guids))

	for _, g := range guids {
		newState[g] = true
	}
	for _, i := range items {
		newState[i.GUID] = true
	}
	toStore := make([]string, 0, len(newState))

	for k := range newState {
		toStore = append(toStore, k)
	}
	sort.Strings(toStore)

	err = writeState(state, toStore)
	if err != nil {
		return nil, fmt.Errorf("couldn't write state: %s", err)
	}
	return ret, nil
}

func readState(state string) ([]string, error) {
	file, err := os.Open(state)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("couldn't open state file: %s", err)
	}

	var guids []string

	if err == nil {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&guids)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("couldn't decode state file: %s", err)
		}
	}

	return guids, nil
}

func writeState(state string, guids []string) error {
	tmp, err := os.Create(state + ".tmp")
	if err != nil {
		return fmt.Errorf("couldn't create state file: %s", err)
	}
	encoder := json.NewEncoder(tmp)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(guids)
	if err != nil {
		return fmt.Errorf("couldn't encode state file: %s", err)
	}
	err = tmp.Close()
	if err != nil {
		return fmt.Errorf("couldn't write state file: %s", err)
	}
	return nil
}
