package clocks

import (
	"bytes"
	_ "embed"
	"testing"
	"time"

	"github.com/frioux/leatherman/internal/testutil"
)

//go:embed expect.txt
var expect string

func TestRun(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	at := time.Date(2012, 12, 12, 4, 12, 12, 12, time.UTC)
	run(at, []string{"America/Los_Angeles", "UTC"}, buf)
	testutil.Equal(t, buf.String(), expect, "wrong report")
}
