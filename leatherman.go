package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/frioux/leatherman/bamboo"
	"github.com/frioux/leatherman/csv"
	"github.com/frioux/leatherman/email"
	"github.com/frioux/leatherman/yaml"
)

// Dispatch is the dispatch table that maps command names to functions.
var Dispatch map[string]func([]string, io.Reader)

func main() {
	which := filepath.Base(os.Args[0])
	args := os.Args

	Dispatch = map[string]func([]string, io.Reader){
		"addrs":                email.Addrs,
		"addrspec-to-tabs":     email.ToTabs,
		"backlight":            Backlight,
		"clocks":               Clocks,
		"csv2json":             csv.ToJSON,
		"csv2md":               csv.ToMarkdown,
		"debounce":             Debounce,
		"dump-mozlz4":          DumpMozLZ4,
		"ec2-resource-for-ip":  EC2ResourceForIP,
		"expand-url":           ExpandURL,
		"export-bamboohr":      bamboo.ExportDirectory,
		"export-bamboohr-tree": bamboo.ExportOrgChart,
		"fn":                   Fn,
		"gen-pass":             GenPass,
		"group-by-date":        GroupByDate,
		"netrc-password":       NetrcPassword,
		"pomotimer":            Pomotimer,
		"prepend-emoji-hist":   PrependEmojiHist,
		"render-mail":          email.Render,
		"replace-unzip":        ReplaceUnzip,
		"rss":                  RSS,
		"ssh-quote":            SSHQuote,
		"undefer":              Undefer,
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
		Help(os.Args, os.Stdin)
		os.Exit(1)
	}
	fn(args, os.Stdin)
}
