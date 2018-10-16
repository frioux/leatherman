// +build !linux

package groupbydate

import (
	"time"
)

func parseDate(format, input string) (time.Time, error) {
	return time.Parse(format, input)
}

func formatDate(format string, date time.Time) (string, error) {
	return date.Format(format), nil
}
