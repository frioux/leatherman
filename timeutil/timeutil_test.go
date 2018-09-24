package timeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJumpTo(t *testing.T) {
	assert.Equal(t,
		time.Date(2018, 9, 28, 0, 0, 0, 0, time.UTC),
		JumpTo(time.Date(2018, 9, 23, 0, 0, 0, 0, time.UTC), time.Friday),
		"Sun -> Fri",
	)

	assert.Equal(t,
		time.Date(2018, 9, 28, 0, 0, 0, 0, time.UTC),
		JumpTo(time.Date(2018, 9, 22, 0, 0, 0, 0, time.UTC), time.Friday),
		"Sat -> Fri",
	)

}
