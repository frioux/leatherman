package debounce

import (
	"io"
	"time"
)

type trailingBouncer chan struct {
	t  time.Time
	in []byte
}

func (t trailingBouncer) Write(at time.Time, in []byte) error {
	t <- struct {
		t  time.Time
		in []byte
	}{at, in}

	return nil
}

func newTrailingBouncer(w io.Writer, duration time.Duration) trailingBouncer {
	ch := make(chan struct {
		t  time.Time
		in []byte
	})

	go func() {
		v := <-ch
		timeout := time.NewTimer(duration)

		for {
			select {
			case <-timeout.C:
				w.Write(v.in)
			case v = <-ch:
				timeout.Reset(duration)
			}
		}
	}()

	return trailingBouncer(ch)
}
