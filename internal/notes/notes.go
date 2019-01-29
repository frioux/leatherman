package notes

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"html/template"
	"io"
	"net/http"

	"github.com/frioux/amygdala/internal/dropbox"
	"github.com/pkg/errors"
)

var bodyTemplate *template.Template

func init() {
	var err error
	bodyTemplate, err = template.New("xxx").Parse(`---
title: {{.}}
tags: [ private, inbox ]
---

{{.}}

`)
	if err != nil {
		panic(err)
	}
}

func body(message string) io.Reader {
	buf := &bytes.Buffer{}

	bodyTemplate.Execute(buf, message)

	return buf
}

// Todo creates an item tagged inbox
func Todo(cl *http.Client, tok, message string) error {
	buf := body(message)

	sha := sha1.Sum([]byte(message))
	path := "/notes/content/posts/todo-" + hex.EncodeToString(sha[:]) + ".md"

	up := dropbox.UploadParams{Path: path, Autorename: true}
	if err := dropbox.Create(cl, tok, up, buf); err != nil {
		return errors.Wrap(err, "dropbox.Create")
	}

	return nil
}
