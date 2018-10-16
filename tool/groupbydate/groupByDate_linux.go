// +build linux

package groupbydate

import (
	"strings"
	"time"

	"github.com/jeffjen/datefmt"
)

func parseDate(format, input string) (time.Time, error) {
	if strings.ContainsRune(format, '%') {
		return datefmt.Strptime(format, input)
	}
	return time.Parse(format, input)
}

func formatDate(format string, date time.Time) (string, error) {
	if strings.ContainsRune(format, '%') {
		return datefmt.Strftime(format, date)
	}
	return date.Format(format), nil
}
