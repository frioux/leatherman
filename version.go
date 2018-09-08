package main

import (
	"fmt"
	"io"
	"runtime"
)

var version, when, who, where string

// Version prints current version
func Version(args []string, _ io.Reader) {
	fmt.Printf("Leatherman built from %s on %s by %s@%s with %s\n",
		version, when, who, where, runtime.Version())
}
