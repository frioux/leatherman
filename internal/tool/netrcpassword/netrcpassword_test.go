package netrcpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	pass, err := run("./testdata/basic.netrc", "foo", "bar")
	if err != nil {
		t.Fatalf("Failed to call run: %s", err)
	}

	assert.Equal(t, "baz", pass)
}
