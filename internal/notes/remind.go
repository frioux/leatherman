package notes

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/frioux/amygdala/internal/dropbox"
	"github.com/frioux/amygdala/internal/personality"
	"github.com/frioux/amygdala/internal/reminders"
	"github.com/pkg/errors"
)

// remind format:
// remind me (?:to )xyz (at $time|in $duration)
func remind(cl *http.Client, tok, message string) (string, error) {
	when, what, err := reminders.Parse(time.Now(), message)
	if err != nil {
		return personality.UserErr(err), err
	}

	sha := sha1.Sum([]byte(what))
	id := hex.EncodeToString(sha[:])
	path := "/notes/.alerts/" + when.Format(time.RFC3339) + "_" + id + ".txt"

	buf := strings.NewReader(what)

	up := dropbox.UploadParams{Path: path, Autorename: true}
	if err := dropbox.Create(cl, tok, up, buf); err != nil {
		return personality.Err(), errors.Wrap(err, "dropbox.Create")
	}

	return personality.Ack(), nil
}
