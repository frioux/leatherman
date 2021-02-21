package main

import "embed"

//go:embed internal/tool/allpurpose/uni/alluni.md
//go:embed internal/tool/allpurpose/clocks/clocks.md
//go:embed internal/tool/allpurpose/csv/csv2json.md
//go:embed internal/tool/allpurpose/csv/csv2md.md
//go:embed internal/tool/allpurpose/debounce/debounce.md
//go:embed internal/tool/allpurpose/dumpmozlz4/dump-mozlz4.md
//go:embed internal/tool/allpurpose/expandurl/expand-url.md
//go:embed internal/tool/allpurpose/fn/fn.md
//go:embed internal/tool/allpurpose/groupbydate/group-by-date.md
//go:embed internal/tool/allpurpose/minotaur/minotaur.md
//go:embed internal/tool/allpurpose/uni/name2rune.md
//go:embed internal/tool/allpurpose/netrcpassword/netrc-password.md
//go:embed internal/tool/allpurpose/pomotimer/pomotimer.md
//go:embed internal/tool/allpurpose/replaceunzip/replace-unzip.md
//go:embed internal/tool/allpurpose/srv/srv.md
//go:embed internal/tool/allpurpose/toml/toml2json.md
//go:embed internal/tool/allpurpose/uni/uni.md
//go:embed internal/tool/allpurpose/yaml/yaml2json.md
//go:embed internal/tool/chat/automoji/auto-emote.md
//go:embed internal/tool/chat/slack/slack-deaddrop.md
//go:embed internal/tool/chat/slack/slack-open.md
//go:embed internal/tool/chat/slack/slack-status.md
//go:embed internal/tool/leatherman/update/update.md
//go:embed internal/tool/mail/email/addrs.md
//go:embed internal/tool/mail/email/addrspec-to-tabs.md
//go:embed internal/tool/mail/email/email2json.md
//go:embed internal/tool/mail/email/render-mail.md
//go:embed internal/tool/misc/backlight/backlight.md
//go:embed internal/tool/misc/img/draw.md
//go:embed internal/tool/misc/bamboo/export-bamboohr.md
//go:embed internal/tool/misc/bamboo/export-bamboohr-tree.md
//go:embed internal/tool/misc/prependhist/prepend-hist.md
//go:embed internal/tool/misc/smlist/sm-list.md
//go:embed internal/tool/misc/status/status.md
//go:embed internal/tool/misc/twilio/twilio.md
//go:embed internal/tool/misc/wuphf/wuphf.md
//go:embed internal/tool/notes/amygdala/amygdala.md
//go:embed internal/tool/notes/brainstem/brainstem.md
//go:embed internal/tool/notes/now/notes.md
//go:embed internal/tool/notes/proj/proj.md
//go:embed internal/tool/notes/zine/zine.md
var helpFS embed.FS

var helpPaths = map[string]string{
	"alluni": "internal/tool/allpurpose/uni/alluni.md",

	"clocks": "internal/tool/allpurpose/clocks/clocks.md",

	"csv2json": "internal/tool/allpurpose/csv/csv2json.md",

	"csv2md": "internal/tool/allpurpose/csv/csv2md.md",

	"debounce": "internal/tool/allpurpose/debounce/debounce.md",

	"dump-mozlz4": "internal/tool/allpurpose/dumpmozlz4/dump-mozlz4.md",

	"expand-url": "internal/tool/allpurpose/expandurl/expand-url.md",

	"fn": "internal/tool/allpurpose/fn/fn.md",

	"group-by-date": "internal/tool/allpurpose/groupbydate/group-by-date.md",

	"minotaur": "internal/tool/allpurpose/minotaur/minotaur.md",

	"name2rune": "internal/tool/allpurpose/uni/name2rune.md",

	"netrc-password": "internal/tool/allpurpose/netrcpassword/netrc-password.md",

	"pomotimer": "internal/tool/allpurpose/pomotimer/pomotimer.md",

	"replace-unzip": "internal/tool/allpurpose/replaceunzip/replace-unzip.md",

	"srv": "internal/tool/allpurpose/srv/srv.md",

	"toml2json": "internal/tool/allpurpose/toml/toml2json.md",

	"uni": "internal/tool/allpurpose/uni/uni.md",

	"yaml2json": "internal/tool/allpurpose/yaml/yaml2json.md",

	"auto-emote": "internal/tool/chat/automoji/auto-emote.md",

	"slack-deaddrop": "internal/tool/chat/slack/slack-deaddrop.md",

	"slack-open": "internal/tool/chat/slack/slack-open.md",

	"slack-status": "internal/tool/chat/slack/slack-status.md",

	"update": "internal/tool/leatherman/update/update.md",

	"addrs": "internal/tool/mail/email/addrs.md",

	"addrspec-to-tabs": "internal/tool/mail/email/addrspec-to-tabs.md",

	"email2json": "internal/tool/mail/email/email2json.md",

	"render-mail": "internal/tool/mail/email/render-mail.md",

	"backlight": "internal/tool/misc/backlight/backlight.md",

	"draw": "internal/tool/misc/img/draw.md",

	"export-bamboohr": "internal/tool/misc/bamboo/export-bamboohr.md",

	"export-bamboohr-tree": "internal/tool/misc/bamboo/export-bamboohr-tree.md",

	"prepend-hist": "internal/tool/misc/prependhist/prepend-hist.md",

	"sm-list": "internal/tool/misc/smlist/sm-list.md",

	"status": "internal/tool/misc/status/status.md",

	"twilio": "internal/tool/misc/twilio/twilio.md",

	"wuphf": "internal/tool/misc/wuphf/wuphf.md",

	"amygdala": "internal/tool/notes/amygdala/amygdala.md",

	"brainstem": "internal/tool/notes/brainstem/brainstem.md",

	"notes": "internal/tool/notes/now/notes.md",

	"proj": "internal/tool/notes/proj/proj.md",

	"zine": "internal/tool/notes/zine/zine.md",
}
