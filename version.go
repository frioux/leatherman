package main

import (
	"io"
	"os"

	"github.com/frioux/leatherman/internal/version"
)

// Version prints current version
func Version(args []string, _ io.Reader) error {
	version.Render(os.Stdout)

	return nil
}
