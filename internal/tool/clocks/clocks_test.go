package clocks

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	buf := &bytes.Buffer{}

	at := time.Date(2012, 12, 12, 4, 12, 12, 12, time.UTC)
	run(at, buf)
	assert.Equal(t, `
                Local  yesterday  20:12   8:12 PM  -8
  America/Los_Angeles  yesterday  20:12   8:12 PM  -8
      America/Chicago  yesterday  22:12  10:12 PM  -6
     America/New_York  yesterday  23:12  11:12 PM  -5
       Asia/Jerusalem      today  06:12   6:12 AM  +2
                  UTC      today  04:12   4:12 AM  +0
`, "\n"+buf.String())
}
