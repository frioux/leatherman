package testutil

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal takes t, got, expected, and a prefix, returning true if got and
// expected are expected.
func Equal(t *testing.T, got, expected interface{}, prefix string, opts ...cmp.Option) bool {
	t.Helper()
	if diff := cmp.Diff(expected, got, opts...); diff != "" {
		t.Errorf("%s (-want +got):\n%s", prefix, diff)
		return false
	}

	return true
}

// JSONEqual takes a got and expected string of json and compares the parsed values with Equal.
func JSONEqual(t *testing.T, got, expected string, prefix string, opts ...cmp.Option) bool {
	t.Helper()
	var gotValue, expectedValue interface{}
	if err := json.NewDecoder(strings.NewReader(got)).Decode(&gotValue); err != nil {
		t.Errorf("Couldn't decode got: %s", err)
		return false
	}

	if err := json.NewDecoder(strings.NewReader(expected)).Decode(&expectedValue); err != nil {
		t.Errorf("Couldn't decode expected: %s", err)
		return false
	}

	return Equal(t, gotValue, expectedValue, prefix, opts...)
}
