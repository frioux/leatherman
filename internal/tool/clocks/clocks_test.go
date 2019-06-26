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
	assert.Regexp(t, `here\s+:\s+tomorrow\s+20:12\s+08:12 PM\s+-8
L.A.\s+:\s+tomorrow\s+20:12\s+08:12 PM\s+-8
MS/TX:\s+tomorrow\s+22:12\s+10:12 PM\s+-6
east\s+:\s+tomorrow\s+23:12\s+11:12 PM\s+-5
TLV\s+:\s+today\s+06:12\s+06:12 AM\s+\+2
UTC\s+:\s+today\s+04:12\s+04:12 AM\s+\+0
`, buf.String())
}
