package main

import (
	"fmt"
	"math/rand"
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
	rules, err := notes.NewRules(tok)
	if err != nil {
		fmt.Printf("Couldn't create rules: %s\n", err)
		os.Exit(1)
	}

	message, err := rules.Dispatch(os.Args[1], nil)
	fmt.Println(message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "(%s)\n", err)
		os.Exit(1)
	}
}
