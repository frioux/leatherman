package notes

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/frioux/amygdala/internal/dropbox"
	"github.com/frioux/amygdala/internal/personality"
	"github.com/frioux/amygdala/internal/twilio"
)

var bodyTemplate *template.Template

type bodyArgs struct {
	Message, At string
}

func init() {
	var err error
	bodyTemplate, err = template.New("xxx").Parse(`{
"title": {{.Message | printf "%q"}},
"date": "{{.At}}",
"tags": [ "private", "inbox" ]
}
 * {{.Message}}
`)
	if err != nil {
		panic(err)
	}
}

func body(message string, at time.Time) io.Reader {
	buf := &bytes.Buffer{}

	bodyTemplate.Execute(buf, bodyArgs{message, at.Format("2006-01-02T15:04:05")})

	return buf
}

// todo creates an item tagged inbox
func todo(cl dropbox.Client) func(message string, media []twilio.Media) (string, error) {
	return func(message string, media []twilio.Media) (string, error) {
		sum := sha1.New()
		// it's impossible for sha1 to emit an error
		sum.Write([]byte(message))
		for _, m := range media {
			sum.Write([]byte(m.URL))
		}
		sha := sum.Sum([]byte(""))
		id := hex.EncodeToString(sha[:])

		for i, m := range media {
			if strings.HasPrefix(m.ContentType, "image/") {
				message += fmt.Sprintf(` <img src="%s" height="128" /> attachment %d on %s`, m.URL, i, id)
			} else {
				message += fmt.Sprintf(" [attachment %d on %s](%s)", i, id, m.URL)
			}
		}

		buf := body(message, time.Now())

		path := "/notes/content/posts/todo-" + id + ".md"
		up := dropbox.UploadParams{Path: path, Autorename: true}
		if err := cl.Create(up, buf); err != nil {
			return personality.Err(), fmt.Errorf("dropbox.Create: %w", err)
		}

		return personality.Ack(), nil
	}
}
