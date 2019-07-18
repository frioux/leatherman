package reminders

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func truncateToDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func timeOffset(t time.Time) time.Duration {
	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second +
		time.Duration(t.Nanosecond())*time.Nanosecond
}

const day = 24 * time.Hour

func nextTime(start, clock time.Time) time.Time {
	soon := truncateToDay(start).Add(timeOffset(clock))
	if soon.Before(start) {
		return soon.Add(day)
	}
	return soon
}

type userErr struct {
	error
}

func (u *userErr) InvalidInput() {}
func (u *userErr) Error() string { return u.error.Error() }

var remindFormat = regexp.MustCompile(`(?i)^remind\s+me(?:\s+to)?\s+(?P<message>.+)\s+(?:at\s+(?P<when>.+)|in\s+(?P<duration>.+))$`)

const (
	MESSAGE  = 1
	WHEN     = 2
	DURATION = 3
)

var (
	errInvalidRemind   = errors.New("invalid remind format")
	errBlankEvent      = errors.New("blank event")
	errInvalidDuration = errors.New("invalid duration format")
	errImpossible      = errors.New("impossible")
)

func Parse(now time.Time, message string) (time.Time, string, error) {
	m := remindFormat.FindStringSubmatch(message)
	if len(m) < 4 {
		return time.Time{}, "", &userErr{errInvalidRemind}
	}

	if m[MESSAGE] == "" {
		return time.Time{}, "", &userErr{errBlankEvent}
	}

	if m[WHEN] != "" {
		when := strings.ToLower(m[WHEN])

		if when == "noon" {
			when = "12:00pm"
		} else if when == "midnight" {
			when = "12:00am"
		}
		for _, format := range []string{"3:04pm", "3pm"} {
			clock, err := time.Parse(format, when)
			if err != nil {
				continue
			}
			return nextTime(now, clock), m[MESSAGE], nil
		}
		return time.Time{}, "", &userErr{errInvalidRemind}
	} else if m[DURATION] != "" {
		d := parseDuration(m[DURATION])
		if d == time.Duration(0) {
			return time.Time{}, "", &userErr{errInvalidDuration}
		}
		return now.Add(d), m[MESSAGE], nil
	}

	return time.Time{}, "", errImpossible
}

// an hour, two hours, 2 hours, etc
var lazy = regexp.MustCompile(`^(\d+|a|an|one|two|three|four|five|six|seven|eight|nine)\s+(hours?|minutes?|days?)$`)

func parseDuration(s string) time.Duration {
	// beware, atypical Go; no success rail :(
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}

	var count int
	if m := lazy.FindStringSubmatch(s); len(m) > 0 {
		fmt.Println("HERE")
		switch m[1] {
		case "a", "an", "one", "1":
			count = 1
		case "two", "2":
			count = 2
		case "three", "3":
			count = 3
		case "four", "4":
			count = 4
		case "five", "5":
			count = 5
		case "six", "6":
			count = 6
		case "seven", "7":
			count = 7
		case "eight", "8":
			count = 8
		case "nine", "9":
			count = 9
		case "ten", "10":
			count = 10
		default:
			var err error
			count, err = strconv.Atoi(m[1])
			if err != nil {
				return time.Duration(0)
			}
		}

		switch m[2] {
		case "hour", "hours":
			return time.Duration(count) * time.Hour
		case "minute", "minutes":
			return time.Duration(count) * time.Minute
		case "day", "days":
			return time.Duration(count) * time.Hour * 24
		default:
			panic("impossible condition reached via " + s)
		}
	}

	return time.Duration(0)
}
