package rss

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/mmcdole/gofeed"
)

type feedState struct {
	URL   string
	GUIDs []string
}

type allStates []feedState

func (a allStates) toIndexedStates() indexedStates {
	var i indexedStates = make(map[string]map[string]bool, len(a))

	for _, f := range a {
		i[f.URL] = make(map[string]bool, len(f.GUIDs))

		for _, g := range f.GUIDs {
			i[f.URL][g] = true
		}
	}

	return i
}

// indexedStates is a map from url to guid to seen
type indexedStates map[string]map[string]bool

func (i indexedStates) toAllStates() allStates {
	var a allStates = make([]feedState, 0, len(i))

	f := make([]string, 0, len(a))
	for k := range i {
		f = append(f, k)
	}

	sort.Strings(f)

	for _, u := range f {
		guids := make([]string, 0, len(i[u]))

		for g := range i[u] {
			guids = append(guids, g)
		}
		sort.Strings(guids)

		a = append(a, feedState{
			URL:   u,
			GUIDs: guids,
		})
	}

	return a
}

func readState(state string) (indexedStates, error) {
	file, err := os.Open(state)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("couldn't open state file: %w", err)
	}

	var a allStates

	if err != nil {
		return a.toIndexedStates(), nil
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&a); err != nil {
		return nil, fmt.Errorf("couldn't decode state file: %w", err)
	}

	return a.toIndexedStates(), nil
}

func writeState(state string, i indexedStates) error {
	tmp, err := os.Create(state + ".tmp")
	if err != nil {
		return fmt.Errorf("couldn't create state file: %w", err)
	}
	encoder := json.NewEncoder(tmp)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(i.toAllStates()); err != nil {
		return fmt.Errorf("couldn't encode state file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("couldn't write state file: %w", err)
	}
	return nil
}

// Store JSON containing seen GUIDs for the current feed.
func syncRead(f map[string]bool, items []*gofeed.Item) {
	for _, i := range items {
		f[i.GUID] = true
	}
}
