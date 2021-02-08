package debounce

import (
	"bytes"
	"testing"
	"time"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestLeading(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	l := newBouncer(false, buf, time.Millisecond)

	start := time.Date(2012, 12, 12, 0, 0, 0, 0, time.UTC)
	l.Write(start, []byte("foo\n"))
	l.Write(start.Add(time.Nanosecond), []byte("bar\n"))
	l.Write(start.Add(time.Nanosecond), []byte("baz\n"))
	l.Write(start.Add(time.Second+2*time.Nanosecond), []byte("biff\n"))

	testutil.Equal(t, buf.String(), "foo\nbiff\n", "wrong output")
}
