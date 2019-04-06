package shellquote_test

import (
	"fmt"
	"os"

	"github.com/frioux/leatherman/pkg/shellquote"
)

func Example() {
	fmt.Println("#!/bin/sh")
	fmt.Println("")
	quoted, err := shellquote.Quote(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't quote input: %s\n", err)
		os.Exit(1)
	}
	// error won't happen if the first input didn't error
	doublequoted, _ := shellquote.Quote([]string{quoted})
	fmt.Println("ssh superserver", doublequoted)
}
