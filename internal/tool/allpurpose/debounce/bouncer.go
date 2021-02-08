package debounce

import (
	"io"
	"time"
)

type bouncer interface {
	Write(time.Time, []byte) error
}

func newBouncer(trailing bool, w io.Writer, duration time.Duration) bouncer {
	if trailing {
		return newTrailingBouncer(w, duration)
	}
	return &leadingBouncer{w: w, duration: duration}
}
