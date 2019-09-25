package clocks

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	at := time.Date(2012, 12, 12, 4, 12, 12, 12, time.UTC)
	run(at, []string{"America/Los_Angeles", "UTC"}, buf)
	assert.Equal(t, `
  America/Los_Angeles  yesterday  20:12  8:12 PM  -8
                  UTC      today  04:12  4:12 AM  +0
`, "\n"+buf.String())
}
