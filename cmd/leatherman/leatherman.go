package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"

	"github.com/pkg/errors"

	"github.com/frioux/leatherman/tool/backlight"
	"github.com/frioux/leatherman/tool/bamboo"
	"github.com/frioux/leatherman/tool/clocks"
	"github.com/frioux/leatherman/tool/csv"
	"github.com/frioux/leatherman/tool/debounce"
	"github.com/frioux/leatherman/tool/deferlwn"
	"github.com/frioux/leatherman/tool/dumpmozlz4"
	"github.com/frioux/leatherman/tool/email"
	"github.com/frioux/leatherman/tool/expandurl"
	"github.com/frioux/leatherman/tool/fn"
	"github.com/frioux/leatherman/tool/genpass"
	"github.com/frioux/leatherman/tool/groupbydate"
	"github.com/frioux/leatherman/tool/netrcpassword"
	"github.com/frioux/leatherman/tool/pomotimer"
	"github.com/frioux/leatherman/tool/prependemojihist"
	"github.com/frioux/leatherman/tool/replaceunzip"
	"github.com/frioux/leatherman/tool/rss"
	"github.com/frioux/leatherman/tool/smlist"
	"github.com/frioux/leatherman/tool/sshquote"
	"github.com/frioux/leatherman/tool/undefer"
	"github.com/frioux/leatherman/tool/yaml"
)

// Dispatch is the dispatch table that maps command names to functions.
var Dispatch map[string]func([]string, io.Reader) error

func main() {
	startDebug()

	which := filepath.Base(os.Args[0])
	args := os.Args

	Dispatch = map[string]func([]string, io.Reader) error{
		"addrs":                email.Addrs,
		"addrspec-to-tabs":     email.ToTabs,
		"backlight":            backlight.Run,
		"clocks":               clocks.Run,
		"csv2json":             csv.ToJSON,
		"csv2md":               csv.ToMarkdown,
		"debounce":             debounce.Run,
		"defer-lwn":            deferlwn.Run,
		"dump-mozlz4":          dumpmozlz4.Run,
		"expand-url":           expandurl.Run,
		"export-bamboohr":      bamboo.ExportDirectory,
		"export-bamboohr-tree": bamboo.ExportOrgChart,
		"fn":                   fn.Run,
		"gen-pass":             genpass.Run,
		"group-by-date":        groupbydate.Run,
		"netrc-password":       netrcpassword.Run,
		"pomotimer":            pomotimer.Run,
		"prepend-emoji-hist":   prependemojihist.Run,
		"render-mail":          email.Render,
		"replace-unzip":        replaceunzip.Run,
		"rss":                  rss.Run,
		"ssh-quote":            sshquote.Run,
		"sm-list":              smlist.Run,
		"undefer":              undefer.Run,
		"yaml2json":            yaml.ToJSON,

		"help":    Help,
		"version": Version,
		"explode": Explode,
	}

	if which == "leatherman" && len(args) > 1 {
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
		err = errors.Wrap(fn(args, os.Stdin), which)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		stopDebug()
		os.Exit(1)
	}
	stopDebug()
}
