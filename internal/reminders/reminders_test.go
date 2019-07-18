package reminders

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAssertRegexLocations(t *testing.T) {
	sn := remindFormat.SubexpNames()
	assert.Equal(t, sn[MESSAGE], "message")
	assert.Equal(t, sn[WHEN], "when")
	assert.Equal(t, sn[DURATION], "duration")
}

var LA *time.Location

func init() {
	var err error
	LA, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
}
func TestNextTime(t *testing.T) {
	now := time.Date(2012, 12, 12, 06, 00, 00, 00, LA)

	type assertion struct {
		start, clock, result time.Time
	}
	assertions := []assertion{
		{now, time.Date(0, 0, 0, 7, 0, 0, 0, LA), time.Date(2012, 12, 12, 7, 0, 0, 0, LA)},
		{now, time.Date(0, 0, 0, 5, 0, 0, 0, LA), time.Date(2012, 12, 13, 5, 0, 0, 0, LA)},
	}
	for _, a := range assertions {
		assert.Equal(t, a.result, nextTime(a.start, a.clock))
	}
}

func TestParse(t *testing.T) {
	now := time.Date(2012, 12, 12, 00, 00, 00, 00, LA)

	type assertion struct {
		in, message string

		t   time.Time
		err bool
	}

	assertions := []assertion{
		{"", "", time.Time{}, true},
		{"remind me to frew in an hour", "frew", now.Add(time.Hour), false},
		{"remind me to frew in 10m", "frew", now.Add(10 * time.Minute), false},
		{"remind me to frioux at 10am", "frioux", now.Add(10 * time.Hour), false},
		{"remind me to frioux at 10AM", "frioux", now.Add(10 * time.Hour), false},
		{"remind me to frioux at 10:01am", "frioux", now.Add(10*time.Hour + time.Minute), false},
		{"remind me to frioux at noon", "frioux", time.Date(2012, 12, 12, 12, 0, 0, 0, LA), false},
		{"remind me to frioux at midnight", "frioux", time.Date(2012, 12, 12, 0, 0, 0, 0, LA), false},
	}

	for _, a := range assertions {
		t.Run(a.in, func(t *testing.T) {
			when, mess, err := Parse(now, a.in)
			if a.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, a.t, when)
			assert.Equal(t, a.message, mess)
		})
	}

	// deferred
	// when, mess, err = parseReminder(time.Time{}, "remind me to frew on Wed")
	// when, mess, err = parseReminder(time.Time{}, "remind me to frew on Wednesday")
}
