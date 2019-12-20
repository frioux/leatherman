package rss

import (
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestToAllStates(t *testing.T) {
	var i indexedStates = map[string]map[string]bool{
		"aurl": {
			"a": true,
			"b": true,
			"c": true,
		},
		"curl": {
			"1": true,
			"2": true,
			"3": true,
		},
	}

	testutil.Equal(t, i.toAllStates(), allStates([]feedState{
		{URL: "aurl", GUIDs: []string{"a", "b", "c"}},
		{URL: "curl", GUIDs: []string{"1", "2", "3"}},
	}), "wrong states")
}
