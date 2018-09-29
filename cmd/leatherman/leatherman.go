package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/frioux/leatherman/tool/backlight"
	"github.com/frioux/leatherman/tool/bamboo"
	"github.com/frioux/leatherman/tool/clocks"
	"github.com/frioux/leatherman/tool/csv"
	"github.com/frioux/leatherman/tool/debounce"
	"github.com/frioux/leatherman/tool/deferlwn"
	"github.com/frioux/leatherman/tool/dumpmozlz4"
	"github.com/frioux/leatherman/tool/ec2resourceforip"
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
	which := filepath.Base(os.Args[0])
	args := os.Args

	Dispatch = map[string]func([]string, io.Reader) error{
		"addrs":                email.Addrs,
		"addrspec-to-tabs":     email.ToTabs,
		"backlight":            backlight.Backlight,
		"clocks":               clocks.Clocks,
		"csv2json":             csv.ToJSON,
		"csv2md":               csv.ToMarkdown,
		"debounce":             debounce.Debounce,
		"defer-lwn":            deferlwn.DeferLWN,
		"dump-mozlz4":          dumpmozlz4.DumpMozLZ4,
		"ec2-resource-for-ip":  ec2resourceforip.EC2ResourceForIP,
		"expand-url":           expandurl.ExpandURL,
		"export-bamboohr":      bamboo.ExportDirectory,
		"export-bamboohr-tree": bamboo.ExportOrgChart,
		"fn":                 fn.Fn,
		"gen-pass":           genpass.GenPass,
		"group-by-date":      groupbydate.GroupByDate,
		"netrc-password":     netrcpassword.NetrcPassword,
		"pomotimer":          pomotimer.Pomotimer,
		"prepend-emoji-hist": prependemojihist.PrependEmojiHist,
		"render-mail":        email.Render,
		"replace-unzip":      replaceunzip.ReplaceUnzip,
		"rss":                rss.RSS,
		"ssh-quote":          sshquote.SSHQuote,
		"sm-list":            smlist.SMList,
		"undefer":            undefer.Undefer,
		"yaml2json":          yaml.ToJSON,

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
	err := errors.Wrap(fn(args, os.Stdin), which)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
