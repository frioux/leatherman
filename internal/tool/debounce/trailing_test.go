package debounce

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTrailing(t *testing.T) {
	t.Skip("flaky test")

	buf := &bytes.Buffer{}

	l := newBouncer(true, buf, 5*time.Millisecond)

	l.Write(time.Now(), []byte("1\n"))

	time.Sleep(time.Millisecond)
	l.Write(time.Now(), []byte("2\n"))

	time.Sleep(time.Millisecond)
	l.Write(time.Now(), []byte("3\n"))

	time.Sleep(100 * time.Millisecond) // print 3
	l.Write(time.Now(), []byte("4\n"))

	time.Sleep(100 * time.Millisecond) // print 4
	l.Write(time.Now(), []byte("5\n"))

	time.Sleep(100 * time.Millisecond) // print 5

	assert.Equal(t, "3\n4\n5\n", buf.String())
}
