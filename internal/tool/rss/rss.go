package rss

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/mmcdole/gofeed"
	"golang.org/x/sync/errgroup"
)

/*
Run is a minimalist rss client.  Outputs JSON on STDOUT.  Takes urls
to feeds and path to state file. Example usage:

```bash
$ rss -state feed.json https://blog.afoolishmanifesto.com/index.xml | jq -r '" * [" + .title + "](" +.link+")"'
 * [Announcing shellquote](https://blog.afoolishmanifesto.com/posts/announcing-shellquote/)
 * [Detecting who used the EC2 metadata server with BCC](https://blog.afoolishmanifesto.com/posts/detecting-who-used-ec2-metadata-server-bcc/)
 * [Centralized known_hosts for ssh](https://blog.afoolishmanifesto.com/posts/centralized-known-hosts-for-ssh/)
 * [Buffered Channels in Golang](https://blog.afoolishmanifesto.com/posts/buffered-channels-in-golang/)
 * [C, Golang, Perl, and Unix](https://blog.afoolishmanifesto.com/posts/c-golang-perl-and-unix/)
```

Command: rss
*/
func Run(args []string, _ io.Reader) error {
	flags := flag.NewFlagSet("rss", flag.ExitOnError)

	var statePath string

	flags.StringVar(&statePath, "state", "", "location to store state")
	if err := flags.Parse(args[1:]); err != nil {
		return fmt.Errorf("flags.Parse: %w", err)
	}

	if len(flags.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s -state rss.json <url> [<url>...]\n", args[0])
		os.Exit(1)
	}

	if statePath == "" {
		fmt.Fprintln(os.Stderr, "-state is required")
		os.Exit(1)
	}

	return run(statePath, flags.Args(), os.Stdout)
}

func loadFeed(fp *gofeed.Parser, urlString string) ([]*gofeed.Item, error) {
	feedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse feed url (%s): %w", urlString, err)
	}

	resp, err := lmhttp.Get(urlString)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get feed: %w", err)
	}

	f, err := fp.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Couldn't fetch feed (%s): %w", feedURL, err)
	}
	fixItems(feedURL, f.Items)

	return f.Items, nil
}

func syncFeed(state indexedStates, items []*gofeed.Item, urlString string, w io.Writer) error {
	if state[urlString] == nil {
		state[urlString] = make(map[string]bool, len(items))
	}

	items = newItems(state[urlString], items)

	for _, i := range items {
		state[urlString][i.GUID] = true
	}

	return renderItems(w, items)
}

func run(statePath string, urls []string, w io.Writer) error {
	state, err := readState(statePath)
	if err != nil {
		return fmt.Errorf("couldn't read state: %w", err)
	}
	fp := gofeed.NewParser()

	results := make([][]*gofeed.Item, len(urls))
	g, _ := errgroup.WithContext(context.Background())

	for i, urlString := range urls {
		i, urlString := i, urlString
		g.Go(func() error { // O(n) goroutines
			items, err := loadFeed(fp, urlString)
			if err != nil {
				return err
			}
			results[i] = items
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	for i, items := range results {
		if err := syncFeed(state, items, urls[i], w); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	if err := writeState(statePath, state); err != nil {
		return fmt.Errorf("Couldn't save state file: %w", err)
	}
	if err := os.Rename(statePath+".tmp", statePath); err != nil {
		return fmt.Errorf("Couldn't rename state file: %w", err)
	}

	return nil
}

// fixItems ensures GUID is set and adds hostname and schema from feed link to
// item links
func fixItems(feedURL *url.URL, items []*gofeed.Item) {
	for _, i := range items {
		if i.GUID == "" {
			i.GUID = i.Link
		}

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

func renderItems(out io.Writer, items []*gofeed.Item) error {
	e := json.NewEncoder(out)

	for _, i := range items {
		if err := e.Encode(i); err != nil {
			return err
		}
	}

	return nil
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
