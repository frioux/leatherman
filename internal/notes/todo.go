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
	"golang.org/x/xerrors"
)

var bodyTemplate *template.Template

type bodyArgs struct {
	Message, ID, At string
}

func init() {
	var err error
	bodyTemplate, err = template.New("xxx").Parse(`---
title: {{.Message | printf "%q"}}
date: "{{.At}}"
tags: [ private, inbox ]
guid: {{.ID}}
---

 * {{.Message}}

`)
	if err != nil {
		panic(err)
	}
}

func body(message, id string, at time.Time) io.Reader {
	buf := &bytes.Buffer{}

	bodyTemplate.Execute(buf, bodyArgs{message, id, at.Format("2006-01-02T15:04:05")})

	return buf
}

// todo creates an item tagged inbox
func todo(cl dropbox.Client) func(message string, media []twilio.Media) (string, error) {
	return func(message string, media []twilio.Media) (string, error) {
		for i, m := range media {
			if strings.HasPrefix(m.ContentType, "image/") {
				message += fmt.Sprintf(` <img alt="attachment %d" src="%s" height="128" />`, i, m.URL)
			} else {
				message += fmt.Sprintf(" [attachment %d](%s)", i, m.URL)
			}
		}

		sha := sha1.Sum([]byte(message))
		id := hex.EncodeToString(sha[:])
		path := "/notes/content/posts/todo-" + id + ".md"

		buf := body(message, id, time.Now())

		up := dropbox.UploadParams{Path: path, Autorename: true}
		if err := cl.Create(up, buf); err != nil {
			return personality.Err(), xerrors.Errorf("dropbox.Create: %w", err)
		}

		return personality.Ack(), nil
	}
}
