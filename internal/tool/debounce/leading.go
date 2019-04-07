package debounce

import (
	"io"
	"time"
)

type leadingBouncer struct {
	w        io.Writer
	next     time.Time
	duration time.Duration
}

func (l *leadingBouncer) Write(t time.Time, s []byte) error {
	oldNext := l.next
	l.next = t.Add(l.duration)
	if t.After(oldNext) {
		_, err := l.w.Write([]byte(s))
		return err
	}
	return nil
}
