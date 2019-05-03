package lwn // import "github.com/frioux/leatherman/pkg/lwn"

import (
	"bufio"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"golang.org/x/xerrors"
)

var re = regexp.MustCompile(`^\s+available on (\w+ \d+, \d{4})\)$`)

// ErrNotFound means date couldn't be found in page.
var ErrNotFound = xerrors.New("couldn't find date, already released or never private")

// AvailableOn returns date the passed page will be free to read.
func AvailableOn(page *url.URL) (time.Time, error) {
	res, err := http.Get(page.String())
	if err != nil {
		return time.Time{}, xerrors.Errorf("http.Get: %w", err)
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
			return time.Time{}, xerrors.Errorf("time.Parse: %w", err)
		}
		return date, nil
	}
	if s.Err() != nil {
		return time.Time{}, xerrors.Errorf("Scanner.Scan: %w", s.Err())
	}

	return time.Time{}, ErrNotFound
}
