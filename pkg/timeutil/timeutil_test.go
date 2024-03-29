package timeutil_test

import (
	"testing"
	"time"

	"github.com/frioux/leatherman/internal/testutil"
	"github.com/frioux/leatherman/pkg/timeutil"
)

func TestJumpTo(t *testing.T) {
	t.Parallel()

	testutil.Equal(t,
		timeutil.JumpTo(time.Date(2018, 9, 23, 0, 0, 0, 0, time.UTC), time.Friday),
		time.Date(2018, 9, 28, 0, 0, 0, 0, time.UTC),
		"Sun -> Fri",
	)

	testutil.Equal(t,
		timeutil.JumpTo(time.Date(2018, 9, 22, 0, 0, 0, 0, time.UTC), time.Friday),
		time.Date(2018, 9, 28, 0, 0, 0, 0, time.UTC),
		"Sat -> Fri",
	)

}
