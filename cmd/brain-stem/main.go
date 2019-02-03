package main

import (
	"fmt"
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

	message, err := notes.Dispatch(cl, tok, os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println(message)
}
