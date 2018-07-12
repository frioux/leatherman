package main

import "fmt"

var version string

// Version prints current version
func Version(args []string) {
	fmt.Printf("Leatherman built from %s\n", version)
}
