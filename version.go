package main

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"

	"github.com/frioux/leatherman/internal/version"
)

// Version prints current version
func Version(args []string, _ io.Reader) error {
	fmt.Printf("Leatherman built from %s on %s by with %s\n",
		version.Version, version.When, runtime.Version())

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}

	for _, dep := range bi.Deps {
		fmt.Printf("%s@%s (%s)\n", dep.Path, dep.Version, dep.Sum)
		if dep.Replace != nil {
			r := dep.Replace
			fmt.Printf("   replaced by %s@%s (%s)\n", r.Path, r.Version, r.Sum)
		}
	}

	return nil
}
