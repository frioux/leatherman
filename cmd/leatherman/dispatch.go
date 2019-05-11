package main

import (
	"io"

	"github.com/frioux/leatherman/internal/tool/backlight"
	"github.com/frioux/leatherman/internal/tool/bamboo"
	"github.com/frioux/leatherman/internal/tool/clocks"
	"github.com/frioux/leatherman/internal/tool/csv"
	"github.com/frioux/leatherman/internal/tool/debounce"
	"github.com/frioux/leatherman/internal/tool/deferlwn"
	"github.com/frioux/leatherman/internal/tool/dumpmozlz4"
	"github.com/frioux/leatherman/internal/tool/email"
	"github.com/frioux/leatherman/internal/tool/expandurl"
	"github.com/frioux/leatherman/internal/tool/fn"
	"github.com/frioux/leatherman/internal/tool/genpass"
	"github.com/frioux/leatherman/internal/tool/groupbydate"
	"github.com/frioux/leatherman/internal/tool/minotaur"
	"github.com/frioux/leatherman/internal/tool/netrcpassword"
	"github.com/frioux/leatherman/internal/tool/noaa"
	"github.com/frioux/leatherman/internal/tool/pomotimer"
	"github.com/frioux/leatherman/internal/tool/prependemojihist"
	"github.com/frioux/leatherman/internal/tool/replaceunzip"
	"github.com/frioux/leatherman/internal/tool/rss"
	"github.com/frioux/leatherman/internal/tool/smlist"
	"github.com/frioux/leatherman/internal/tool/srv"
	"github.com/frioux/leatherman/internal/tool/sshquote"
	"github.com/frioux/leatherman/internal/tool/toml"
	"github.com/frioux/leatherman/internal/tool/undefer"
	"github.com/frioux/leatherman/internal/tool/uni"
	"github.com/frioux/leatherman/internal/tool/yaml"
)

func init() {
	Dispatch = map[string]func([]string, io.Reader) error{
		"addrs": email.Addrs,
		"addrspec-to-tabs": email.ToTabs,
		"backlight": backlight.Run,
		"clocks": clocks.Run,
		"csv2json": csv.ToJSON,
		"csv2md": csv.ToMarkdown,
		"debounce": debounce.Run,
		"defer-lwn": deferlwn.Run,
		"dump-mozlz4": dumpmozlz4.Run,
		"email2json": email.ToJSON,
		"expand-url": expandurl.Run,
		"export-bamboohr": bamboo.ExportDirectory,
		"export-bamboohr-tree": bamboo.ExportOrgChart,
		"fn": fn.Run,
		"gen-pass": genpass.Run,
		"group-by-date": groupbydate.Run,
		"minotaur": minotaur.Run,
		"netrc-password": netrcpassword.Run,
		"noaa-proxy": noaa.Proxy,
		"pomotimer": pomotimer.Run,
		"prepend-emoji-hist": prependemojihist.Run,
		"render-mail": email.Render,
		"replace-unzip": replaceunzip.Run,
		"rss": rss.Run,
		"sm-list": smlist.Run,
		"srv": srv.Serve,
		"ssh-quote": sshquote.Run,
		"toml2json": toml.ToJSON,
		"undefer": undefer.Run,
		"uni": uni.Describe,
		"yaml2json": yaml.ToJSON,

		"help":    Help,
		"version": Version,
		"explode": Explode,
	}
}
