package netrcpassword

import (
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestRun(t *testing.T) {
	t.Parallel()

	pass, err := run("./testdata/basic.netrc", "foo", "bar")
	if err != nil {
		t.Fatalf("Failed to call run: %s", err)
	}

	testutil.Equal(t, pass, "baz", "passwords not equal")
}
