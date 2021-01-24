package main

import (
	"context"
	"fmt"
	"hash/maphash"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime/trace"

	"github.com/frioux/leatherman/internal/selfupdate"
)

//go:generate maint/generate-README
//go:generate maint/generate-dispatch

// Dispatch is the dispatch table that maps command names to functions.
var Dispatch map[string]func([]string, io.Reader) error

// run returns false when an error occurred
func run() bool {
	startDebug()
	defer stopDebug()
	h := &maphash.Hash{}
	h.WriteByte(byte(os.Getpid()))
	h.WriteByte(byte(os.Getppid()))
	if n, err := os.Hostname(); err == nil {
		h.WriteString(n)
	}

	rand.Seed(int64(h.Sum64()))

	selfupdate.AutoUpdate()

	which := filepath.Base(os.Args[0])
	args := os.Args

	Dispatch["xyzzy"] = func([]string, io.Reader) error { fmt.Println("nothing happens"); return nil }
	if _, ok := Dispatch[which]; !ok && len(args) > 1 {
		args = args[1:]
		which = args[0]
	}

	fn, ok := Dispatch[which]
	if !ok {
		_ = Help(os.Args, os.Stdin)
		return false
	}
	var err error

	trace.WithRegion(context.Background(), which, func() {
		err = fn(args, os.Stdin)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", which, err)
		return false
	}
	return true
}

func main() {
	if !run() {
		os.Exit(1)
	}
}
