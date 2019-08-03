package notes

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/frioux/amygdala/internal/dropbox"
	"github.com/frioux/amygdala/internal/personality"
	"github.com/frioux/amygdala/internal/twilio"
	"golang.org/x/xerrors"
)

// jumpTo starts at the start and jumps to dest.
//
// Copied from leatherman's timeutil
func jumpTo(start time.Time, dest time.Weekday) time.Time {
	offset := (dest - start.Weekday()) % 7
	// Go's modulus is dumb?
	if offset < 0 {
		offset += 7
	}
	return start.AddDate(0, 0, int(offset))
}

var deferPattern = regexp.MustCompile(`(?i)^\s*defer\s+(?:(.*)\s+)?(?:until|till|til)\s+(\d{4}-\d\d-\d\d|mon|monday|tue|tuesday|wed|wednesday|thu|thur|thursday|fri|friday|sat|saturday|sun|sunday)\s*`)

// deferMessage creates a deferred message for future frew.  Format is
// 	defer <message> till 2006-01-02
// or
//    defer <message> till mon
func deferMessage(cl dropbox.Client) func(string, []twilio.Media) (string, error) {
	return func(input string, media []twilio.Media) (string, error) {
		m := deferPattern.FindStringSubmatch(input)
		if len(m) != 3 {
			return personality.Err(), errors.New("deferMessage: input didn't match pattern (" + input + ")")
		}

		message, when := m[1], m[2]
		weekday := time.Weekday(7) // intentionally invalid
		switch strings.ToLower(when) {
		case "mon", "monday":
			weekday = time.Monday
		case "tue", "tuesday":
			weekday = time.Tuesday
		case "wed", "wednesday":
			weekday = time.Wednesday
		case "thu", "thur", "thursday":
			weekday = time.Thursday
		case "fri", "friday":
			weekday = time.Friday
		case "sat", "saturday":
			weekday = time.Saturday
		case "sun", "sunday":
			weekday = time.Sunday
		}

		if weekday != time.Weekday(7) {
			when = jumpTo(time.Now(), weekday).Format("2006-01-02")
		}
		for i, m := range media {
			if strings.HasPrefix(m.ContentType, "image/") {
				message += fmt.Sprintf(` <img alt="attachment %d" src="%s" height="128" />`, i, m.URL)
			} else {
				message += fmt.Sprintf(" [attachment %d](%s)", i, m.URL)
			}
		}

		sha := sha1.Sum([]byte(message))
		id := hex.EncodeToString(sha[:])
		path := "/notes/.deferred/" + when + "-" + id + ".md"

		up := dropbox.UploadParams{Path: path, Autorename: true}
		if err := cl.Create(up, strings.NewReader(" * "+message+"\n")); err != nil {
			return personality.Err(), xerrors.Errorf("dropbox.Create: %w", err)
		}

		return personality.Ack(), nil
	}
}
