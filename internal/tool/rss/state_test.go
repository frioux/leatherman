package rss

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, allStates([]feedState{
		{URL: "aurl", GUIDs: []string{"a", "b", "c"}},
		{URL: "curl", GUIDs: []string{"1", "2", "3"}},
	}), i.toAllStates())
}
