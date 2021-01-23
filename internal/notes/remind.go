package notes

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/frioux/leatherman/internal/personality"
	"github.com/frioux/leatherman/internal/reminders"
	"github.com/frioux/leatherman/internal/twilio"
)

// remind format:
// remind me (?:to )xyz (at $time|in $duration)
func remind(cl dropbox.Client) func(string, []twilio.Media) (string, error) {
	return func(message string, media []twilio.Media) (string, error) {
		when, what, err := reminders.Parse(time.Now(), message)
		if err != nil {
			return personality.UserErr(err), err
		}

		for _, m := range media {
			what += " " + m.URL
		}

		sha := sha1.Sum([]byte(what))
		id := hex.EncodeToString(sha[:])
		path := "/notes/content/posts/deferred_" + id + ".md"

		const tpl = `{
"title": "deferred %s",
"tags":["deferred"],
"review_by": "%s",
}

%s
`
		buf := strings.NewReader(fmt.Sprintf(tpl, id, when.Format("2006-01-02"), what))

		up := dropbox.UploadParams{Path: path, Autorename: true}
		if err := cl.Create(up, buf); err != nil {
			return personality.Err(), fmt.Errorf("dropbox.Create: %w", err)
		}

		return personality.Ack() + "; will remind you @ " + when.Format(time.RFC3339), nil
	}
}
