package lwn // import "github.com/frioux/leatherman/lwn"

import (
	"bufio"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

var re = regexp.MustCompile(`^\s+available on (\w+ \d+, \d{4})\)$`)

// ErrNotFound means date couldn't be found in page.
var ErrNotFound = errors.New("Couldn't find date, already released or never private")

// AvailableOn returns date the passed page will be free to read.
func AvailableOn(page *url.URL) (time.Time, error) {
	res, err := http.Get(page.String())
	if err != nil {
		return time.Time{}, errors.Wrap(err, "http.Get")
	}
	defer res.Body.Close()

	s := bufio.NewScanner(res.Body)
	for s.Scan() {
		match := re.FindStringSubmatch(s.Text())
		if len(match) < 2 {
			continue
		}
		date, err := time.Parse(`January 2, 2006`, match[1])
		if err != nil {
			return time.Time{}, errors.Wrap(err, "time.Parse")
		}
		return date, nil
	}
	if s.Err() != nil {
		return time.Time{}, errors.Wrap(s.Err(), "Scanner.Scan")
	}

	return time.Time{}, ErrNotFound
}
