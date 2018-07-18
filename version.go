package main

import (
	"fmt"
	"io"
)

var version string

// Version prints current version
func Version(args []string, _ io.Reader) {
	fmt.Printf("Leatherman built from %s\n", version)
}
