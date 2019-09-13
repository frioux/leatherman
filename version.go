package main

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
)

// Version prints current version
func Version(args []string, _ io.Reader) error {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprintln(os.Stderr, "Couldn't read version info")
		return nil
	}

	m := bi.Main
	fmt.Printf("Leatherman (%s) Version %s (%s)\n", m.Path, m.Version, m.Sum)

	for _, dep := range bi.Deps {
		fmt.Printf("%s@%s (%s)\n", dep.Path, dep.Version, dep.Sum)
		if dep.Replace != nil {
			r := dep.Replace
			fmt.Printf("   replaced by %s@%s (%s)\n", r.Path, r.Version, r.Sum)
		}
	}

	return nil
}
