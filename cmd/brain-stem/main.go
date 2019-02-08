package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/frioux/amygdala/internal/notes"
)

var tok string

func init() {
	tok = os.Getenv("DROPBOX_ACCESS_TOKEN")
	if tok == "" {
		panic("dropbox token is missing")
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	cl := &http.Client{}

	message, err := notes.Dispatch(cl, tok, os.Args[1])
	fmt.Println(message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "(%s)\n", err)
		os.Exit(1)
	}
}
