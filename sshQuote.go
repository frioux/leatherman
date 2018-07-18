package main

import (
	"fmt"
	"os"

	"github.com/frioux/shellquote"
)

// SSHQuote takes a command and prints how you would need to quote it for ssh to
// execute it for you.
func SSHQuote(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [some tokens to quote]\n", args[0])
		os.Exit(1)
	}
	quoted, err := shellquote.Quote(args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't quote input: %s\n", err)
		os.Exit(1)
	}
	double, _ := shellquote.Quote([]string{"sh", "-c", quoted})
	triple, _ := shellquote.Quote([]string{double})
	fmt.Println(triple)
}
