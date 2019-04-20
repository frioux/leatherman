package debounce

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLeading(t *testing.T) {
	buf := &bytes.Buffer{}

	l := newBouncer(false, buf, time.Millisecond)

	start := time.Date(2012, 12, 12, 0, 0, 0, 0, time.UTC)
	l.Write(start, []byte("foo\n"))
	l.Write(start.Add(time.Nanosecond), []byte("bar\n"))
	l.Write(start.Add(time.Nanosecond), []byte("baz\n"))
	l.Write(start.Add(time.Second+2*time.Nanosecond), []byte("biff\n"))

	assert.Equal(t, "foo\nbiff\n", buf.String())
}
