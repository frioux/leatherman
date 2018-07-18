package main

import (
	"os"
	"path/filepath"
)

// Dispatch is the dispatch table that maps command names to functions.
var Dispatch map[string]func([]string)

func main() {
	which := filepath.Base(os.Args[0])
	args := os.Args

	Dispatch = map[string]func([]string){
		"addrs":                Addrs,
		"addrspec-to-tabs":     AddrspecToTabs,
		"backlight":            Backlight,
		"clocks":               Clocks,
		"csv2json":             CSVToJSON,
		"debounce":             Debounce,
		"dump-mozlz4":          DumpMozLZ4,
		"ec2-resource-for-ip":  EC2ResourceForIP,
		"expand-url":           ExpandURL,
		"export-bamboohr":      ExportBambooHR,
		"export-bamboohr-tree": ExportBambooHRTree,
		"fn":                 Fn,
		"gen-pass":           GenPass,
		"group-by-date":      GroupByDate,
		"netrc-password":     NetrcPassword,
		"pomotimer":          Pomotimer,
		"prepend-emoji-hist": PrependEmojiHist,
		"render-mail":        RenderMail,
		"replace-unzip":      ReplaceUnzip,
		"rss":                RSS,
		"ssh-quote":          SSHQuote,
		"yaml2json":          YAMLToJSON,

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
		Help(os.Args)
		os.Exit(1)
	}
	fn(args)
}
