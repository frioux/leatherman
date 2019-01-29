package main

import (
	"net/http"
	"os"

	"github.com/frioux/amygdala/internal/notes"
)

var tok string

func init() {
	tok = os.Getenv("DROPBOX_ACCESS_TOKEN")
	if tok == "" {
		panic("dropbox token is missing")
	}
}

func main() {
	cl := &http.Client{}

	if err := notes.Todo(cl, tok, os.Args[1]); err != nil {
		panic(err)
	}
}
