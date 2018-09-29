package timeutil // import "github.com/frioux/leatherman/pkg/timeutil"

import "time"

// JumpTo starts at the start and jumps to dest.
func JumpTo(start time.Time, dest time.Weekday) time.Time {
	offset := time.Duration((dest - start.Weekday()) % 7)
	// Go's modulus is dumb?
	if offset < 0 {
		offset += 7
	}
	return start.Add(time.Hour * 24 * offset)
}
