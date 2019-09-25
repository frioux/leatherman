package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"
)

//go:generate maint/generate-README
//go:generate maint/generate-dispatch

// Dispatch is the dispatch table that maps command names to functions.
var Dispatch map[string]func([]string, io.Reader) error

func main() {
	startDebug()

	which := filepath.Base(os.Args[0])
	args := os.Args

	if _, ok := Dispatch[which]; !ok && len(args) > 1 {
		args = args[1:]
		which = args[0]
	}

	fn, ok := Dispatch[which]
	if !ok {
		_ = Help(os.Args, os.Stdin)
		stopDebug()
		os.Exit(1)
	}
	var err error

	trace.WithRegion(context.Background(), which, func() {
		err = fn(args, os.Stdin)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", which, err)
		stopDebug()
		os.Exit(1)
	}
	stopDebug()
}