package notes

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/frioux/amygdala/internal/dropbox"
	"github.com/frioux/amygdala/internal/personality"
	"github.com/pkg/errors"
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
func todo(cl *http.Client, tok, message string) (string, error) {
	sha := sha1.Sum([]byte(message))
	id := hex.EncodeToString(sha[:])
	path := "/notes/content/posts/todo-" + id + ".md"

	buf := body(message, id, time.Now())

	up := dropbox.UploadParams{Path: path, Autorename: true}
	if err := dropbox.Create(cl, tok, up, buf); err != nil {
		return personality.Err(), errors.Wrap(err, "dropbox.Create")
	}

	return personality.Ack(), nil
}
