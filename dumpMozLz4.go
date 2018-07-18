package main

import (
	"fmt"
	"io"
	"os"

	"github.com/frioux/mozlz4"
)

const magicHeader = "mozLz40\x00"

// DumpMozLZ4 writes the uncompressed mozlz4 file from the first argument to stdout
func DumpMozLZ4(args []string, _ io.Reader) {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s session.jsonlz4\n", args[0])
		os.Exit(1)
	}
	file, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open: %s\n", err)
		os.Exit(1)
	}

	r, err := mozlz4.NewReader(file)
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't copy: %s\n", err)
		os.Exit(1)
	}
}
