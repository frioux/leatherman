package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal takes t, got, expected, and a prefix, returning true if got and
// expected are expected.
func Equal(t *testing.T, got, expected interface{}, prefix string, opts ...cmp.Option) bool {
	if diff := cmp.Diff(got, expected, opts...); diff != "" {
		t.Errorf("%s (-want +got):\n%s", prefix, diff)
		return false
	}

	return true
}
