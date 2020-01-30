package main

var readme = []byte(`# Leatherman - fREW's favorite multitool

This is a little project simply to make trivial tools in Go effortless for my
personal usage.  These tools are almost surely of low utility to most people,
but may be instructive nonetheless.

[I have CI/CD to build this into a single
binary](https://github.com/frioux/leatherman/blob/master/.travis.yml) and [an
` + "`" + `explode` + "`" + ` tool that builds
symlinks](https://github.com/frioux/leatherman/blob/master/explode.go) for each
tool in the busybox style.

[I have automation in my
dotfiles](https://github.com/frioux/dotfiles/blob/bef8303c19e2cefac7dfbec420ad8d45b95415b8/install.sh#L133-L141)
to pull the latest binary at install time and run the ` + "`" + `explode` + "`" + ` tool.

## Installation

Here's a copy pasteable script to install the leatherman on OSX or Linux:

` + "`" + `` + "`" + `` + "`" + ` bash
OS=$([ $(uname -s) = "Darwin" ] && echo "-osx")
LMURL="$(curl -s https://api.github.com/repos/frioux/leatherman/releases/latest |
   grep browser_download_url |
   cut -d '"' -f 4 |
   grep -F leatherman${OS}.xz )"
mkdir -p ~/bin
curl -sL "$LMURL" > ~/bin/leatherman.xz
xz -d -f ~/bin/leatherman.xz
chmod +x ~/bin/leatherman
~/bin/leatherman explode
` + "`" + `` + "`" + `` + "`" + `

This asssumes that ` + "`" + `~/bin` + "`" + ` is in your path.  The ` + "`" + `explode` + "`" + ` command will create a
symlink for each of the tools listed below.

## Usage

Each tool takes different args, but to run a tool you can either use a symlink
(presumably created by ` + "`" + `explode` + "`" + `):

` + "`" + `` + "`" + `` + "`" + ` bash
$ echo "---\nfoo: 1" | yaml2json
{"foo":1}
` + "`" + `` + "`" + `` + "`" + `

or use it as a subcommand:

` + "`" + `` + "`" + `` + "`" + ` bash
echo "---\nfoo: 1" | leatherman yaml2json
{"foo":1}
` + "`" + `` + "`" + `` + "`" + `

## Tools

### ` + "`" + `addrs` + "`" + `

` + "`" + `addrs` + "`" + ` sorts the addresses passed on stdin (in the mutt addrbook format, see
` + "`" + `addrspec-to-tabs` + "`" + `) and sorts them based on how recently they were used, from
the glob passed on the arguments.  The tool exists so that you can create an
address list either with an export tool (like ` + "`" + `goobook` + "`" + `), a subset of your sent
addresses, or whatever else, and then you can sort it based on your sent mail
folder.

` + "`" + `` + "`" + `` + "`" + ` bash
$ <someaddrs.txt addrs "$HOME/mail/gmail.sent/cur/*" >sortedaddrs.txt
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `addrspec-to-tabs` + "`" + `

` + "`" + `addrspec-to-tabs` + "`" + ` converts email addresses from the standard format (` + "`" + `"Hello Friend"
<foo@bar>` + "`" + `) from stdin to the mutt address book format, ie tab separated fields,
on stdout.

` + "`" + `` + "`" + `` + "`" + ` bash
$ <someaddrs.txt addrs "$HOME/mail/gmail.sent/cur/*" | addrspec-to-tabs >addrbook.txt
` + "`" + `` + "`" + `` + "`" + `

This tool ignores the comment because, after actually auditing my addressbook,
most comments are incorrectly recognized by all tools. (for example:
` + "`" + `<5555555555@vzw.com> (555) 555-5555` + "`" + ` should not have a comment of ` + "`" + `(555)` + "`" + `.)

It exists to be combined with ` + "`" + `addrs` + "`" + ` and mutt.

### ` + "`" + `alluni` + "`" + `

` + "`" + `alluni` + "`" + ` prints all unicode character names.

` + "`" + `` + "`" + `` + "`" + `bash
$ alluni | grep DENTISTRY
DENTISTRY SYMBOL LIGHT VERTICAL AND TOP RIGHT
DENTISTRY SYMBOL LIGHT VERTICAL AND BOTTOM RIGHT
DENTISTRY SYMBOL LIGHT VERTICAL WITH CIRCLE
DENTISTRY SYMBOL LIGHT DOWN AND HORIZONTAL WITH CIRCLE
DENTISTRY SYMBOL LIGHT UP AND HORIZONTAL WITH CIRCLE
DENTISTRY SYMBOL LIGHT VERTICAL WITH TRIANGLE
DENTISTRY SYMBOL LIGHT DOWN AND HORIZONTAL WITH TRIANGLE
DENTISTRY SYMBOL LIGHT UP AND HORIZONTAL WITH TRIANGLE
DENTISTRY SYMBOL LIGHT VERTICAL AND WAVE
DENTISTRY SYMBOL LIGHT DOWN AND HORIZONTAL WITH WAVE
DENTISTRY SYMBOL LIGHT UP AND HORIZONTAL WITH WAVE
DENTISTRY SYMBOL LIGHT DOWN AND HORIZONTAL
DENTISTRY SYMBOL LIGHT UP AND HORIZONTAL
DENTISTRY SYMBOL LIGHT VERTICAL AND TOP LEFT
DENTISTRY SYMBOL LIGHT VERTICAL AND BOTTOM LEFT
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `backlight` + "`" + `

` + "`" + `backlight` + "`" + ` is a faster version of ` + "`" + `xbacklight` + "`" + ` by directly writing to ` + "`" + `/sys` + "`" + `.  Example:

Increase by 10%:

` + "`" + `` + "`" + `` + "`" + ` bash
$ backlight 10
` + "`" + `` + "`" + `` + "`" + `

Decrease by 10%:

` + "`" + `` + "`" + `` + "`" + ` bash
$ backlight -10
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `clocks` + "`" + `

` + "`" + `clocks` + "`" + ` shows my personal, digital, wall of clocks.  Pass one or more timezone names
to choose which timezones are shown.

` + "`" + `` + "`" + `` + "`" + `bash
clocks Africa/Johannesburg America/Los_Angeles Europe/Copenhagen
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `csv2json` + "`" + `

` + "`" + `csv2json` + "`" + ` reads CSV on stdin and writes JSON on stdout; first line of input is the
header, and thus the keys of the JSON.

### ` + "`" + `csv2md` + "`" + `

` + "`" + `csv2md` + "`" + ` reads CSV on stdin and writes Markdown on stdout.

### ` + "`" + `debounce` + "`" + `

` + "`" + `debounce` + "`" + ` debounces input from stdin to stdout

The default lockout time is one second, you can override that with the
` + "`" + `--lockoutTime` + "`" + ` argument.  By default the trailing edge triggers output, so
output is emitted after there is no input for the lockout time.  You can change
this behavior by passing the ` + "`" + `--leadingEdge` + "`" + ` flag.

### ` + "`" + `dump-mozlz4` + "`" + `

` + "`" + `dump-mozlz4` + "`" + ` dumps the contents of a ` + "`" + `mozlz4` + "`" + ` (aka ` + "`" + `jsonlz4` + "`" + `) file commonly used by
Firefox.  Just takes the name of the file to dump and writes to standard out.

### ` + "`" + `email2json` + "`" + `

` + "`" + `email2json` + "`" + ` produces a JSON representation of an email from a list of globs.  Only
headers are currently supported, patches welcome to support bodies.

` + "`" + `` + "`" + `` + "`" + `bash
$ email2json '/home/frew/var/mail/mitsi/cur/*' | head -1 | jq .
{
  "Header": {
    "Content-Type": "multipart/alternative; boundary=00163642688b8ef3070464661533",
    "Date": "Thu, 5 Mar 2009 15:45:17 -0600",
    "Delivered-To": "xxx",
    "From": "fREW Schmidt <xxx>",
    "Message-Id": "<fb3648c60903051345o728960f5l6cfb9e1f324bbf50@mail.gmail.com>",
    "Mime-Version": "1.0",
    "Received": "by 10.103.115.8 with HTTP; Thu, 5 Mar 2009 13:45:17 -0800 (PST)",
    "Subject": "STATION",
    "To": "xyzzy@googlegroups.com"
  }
}
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `expand-url` + "`" + `

` + "`" + `expand-url` + "`" + ` reads text on STDIN and writes the same text back, converting any links to
Markdown links, with the title of the page as the title of the link.

### ` + "`" + `export-bamboohr` + "`" + `

` + "`" + `export-bamboohr` + "`" + ` exports entire company directory as JSON.

### ` + "`" + `export-bamboohr-tree` + "`" + `

` + "`" + `export-bamboohr-tree` + "`" + ` exports company org chart as JSON.

### ` + "`" + `fn` + "`" + `

` + "`" + `fn` + "`" + ` creates persistent functions by actually writing scripts.  Example usage:

` + "`" + `` + "`" + `` + "`" + `
fn count-users 'wc -l < /etc/passwd'
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `group-by-date` + "`" + `

` + "`" + `group-by-date` + "`" + ` creates time series data by counting lines and grouping them by a given date
format.  takes dates on stdin in format -i, will group them by format -g, and
write them in format -o.

### ` + "`" + `minotaur` + "`" + `

` + "`" + `minotaur` + "`" + ` watches one or more directories (before the ` + "`" + `--` + "`" + `) and runs a script when
events in those directories occur.

` + "`" + `` + "`" + `` + "`" + `bash
minotaur -include-args -include internal -ignore yaml \
   ~/code/leatherman -- \
   go test ~/code/leatherman/...
` + "`" + `` + "`" + `` + "`" + `

If the ` + "`" + `-include-args` + "`" + ` flag is set, the script receives the events as
arguments, so you can exit early if only irrelevant files changed.

The arguments are of the form ` + "`" + `$event\t$filename` + "`" + `; for example ` + "`" + `CREATE	x.pl` + "`" + `.
As far as I know the valid events are;

 * ` + "`" + `CHMOD` + "`" + `
 * ` + "`" + `CREATE` + "`" + `
 * ` + "`" + `REMOVE` + "`" + `
 * ` + "`" + `RENAME` + "`" + `
 * ` + "`" + `WRITE` + "`" + `

The events are deduplicated and also debounced, so your script will never fire
more often than once a second.  If events are happening every half second the
debouncing will cause the script to never run.

The underlying library supports emitting multiple events in a single line (ie
` + "`" + `CREATE|CHMOD` + "`" + `) though I've not seen that in Linux.

` + "`" + `minotaur` + "`" + ` reëmits all output (both stderr and stdout) of the passed script to
standard out, so you could make a script like this to experiment with the
events with timestamps:

` + "`" + `` + "`" + `` + "`" + `bash
#!/bin/sh

for x in "$@"; do
   echo "$x"
done | ts
` + "`" + `` + "`" + `` + "`" + `

You can do all kinds of interesting things in the script, for example you could
verify that the events deserve a restart, then restart a service, then block till
the service can serve traffic, then restart some other related service.

The ` + "`" + `-include` + "`" + ` and ` + "`" + `-ignore` + "`" + ` arguments are optional; by default ` + "`" + `-include` + "`" + ` is
empty, so matches everything, and ` + "`" + `-ignore` + "`" + ` matches ` + "`" + `.git` + "`" + `.  You can also pass
` + "`" + `-verbose` + "`" + ` to include output about minotaur itself, like which directories it's
watching.

The flag ` + "`" + `-no-run-at-start` + "`" + ` will not the the script until there are any events.

The flag ` + "`" + `-report` + "`" + ` will decorate output with a text wrapper to clarify when the
script is run.

### ` + "`" + `name2rune` + "`" + `

` + "`" + `name2rune` + "`" + ` takes the name of a unicode character and prints out any found
characters.

### ` + "`" + `netrc-password` + "`" + `

` + "`" + `netrc-password` + "`" + ` prints password for the passed hostname and login.

` + "`" + `` + "`" + `` + "`" + `bash
$ netrc-password google.com me@gmail.com
supersecretpassword
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `pomotimer` + "`" + `

` + "`" + `pomotimer` + "`" + ` starts a timer for 25m or the duration expressed in the first
argument.

` + "`" + `` + "`" + `` + "`" + `
pomotimer 2.5m
` + "`" + `` + "`" + `` + "`" + `

or

` + "`" + `` + "`" + `` + "`" + `
pomotimer 3m12s
` + "`" + `` + "`" + `` + "`" + `

Originally a timer for use with [the pomodoro][1] [technique][2].  Handy timer in any case
since you can pass it arbitrary durations, pause it, reset it, and see it's
progress.

[1]: https://blog.afoolishmanifesto.com/posts/the-pomodoro-technique/
[2]: https://blog.afoolishmanifesto.com/posts/the-pomodoro-technique-three-years-later/

### ` + "`" + `prepend-hist` + "`" + `

` + "`" + `prepend-hist` + "`" + ` prints out deduplicated lines from the history file in reverse order and
then prints out the lines from STDIN, filtering out what's already been printed.

` + "`" + `` + "`" + `` + "`" + `bash
$ alluni | prefix-hist ~/.uni_history
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `proj` + "`" + `

` + "`" + `proj` + "`" + ` integrates my various project management tools.

Usage is:

` + "`" + `` + "`" + `` + "`" + `bash
$ cd my/cool-project
$ proj init cool-project
` + "`" + `` + "`" + `` + "`" + `

The following flags are supported before the project name:

 * ` + "`" + `-skip-vim` + "`" + `: does not create vim session
 * ` + "`" + `-force-vim` + "`" + `: overwrites any existing vim session
 * ` + "`" + `-skip-note` + "`" + `: does not create note
 * ` + "`" + `-force-note` + "`" + `: overwrites any existing note
 * ` + "`" + `-skip-smartcd` + "`" + `: does not create smartcd enter script
 * ` + "`" + `-force-smartcd` + "`" + `: overwrites any existing smartcd enter script

Once a project has been initialized, you can run:

` + "`" + `` + "`" + `` + "`" + `bash
$ proj vim
` + "`" + `` + "`" + `` + "`" + `

To open the vim session for that project.

I use [vim sessions][vim], [a notes system][notes], and of course checkouts of
code all over the place.  Proj is meant to make creation of a vim session and a
note easy and eventually allow jumping back and forth between the three.  As of
2019-12-02 it is almost painfully specific to my personal setup, but as I
discover the actual patterns I'll probably generalize.

Proj uses uses [smartcd][smartcd] both as a mechanism and as the means to
add functionality to projects within shell sessions.

[vim]: https://blog.afoolishmanifesto.com/posts/vim-session-workflow/
[notes]: https://blog.afoolishmanifesto.com/posts/a-love-letter-to-plain-text/#notes
[smartcd]: https://github.com/cxreg/smartcd

### ` + "`" + `render-mail` + "`" + `

` + "`" + `render-mail` + "`" + ` is a scrappy tool to render email with a Local-Date included, if Date is
not already in local time.

### ` + "`" + `replace-unzip` + "`" + `

` + "`" + `replace-unzip` + "`" + ` extracts zipfiles, but does not extract ` + "`" + `.DS_Store` + "`" + ` or ` + "`" + `__MACOSX/` + "`" + `.
Automatically extracts into a directory named after the zipfile if there is not
a single root for all files in the zipfile.

### ` + "`" + `rss` + "`" + `

` + "`" + `rss` + "`" + ` is a minimalist rss client.  Outputs JSON on STDOUT.  Takes urls
to feeds and path to state file. Example usage:

` + "`" + `` + "`" + `` + "`" + `bash
$ rss -state feed.json https://blog.afoolishmanifesto.com/index.xml | jq -r '" * [" + .title + "](" +.link+")"'
 * [Announcing shellquote](https://blog.afoolishmanifesto.com/posts/announcing-shellquote/)
 * [Detecting who used the EC2 metadata server with BCC](https://blog.afoolishmanifesto.com/posts/detecting-who-used-ec2-metadata-server-bcc/)
 * [Centralized known_hosts for ssh](https://blog.afoolishmanifesto.com/posts/centralized-known-hosts-for-ssh/)
 * [Buffered Channels in Golang](https://blog.afoolishmanifesto.com/posts/buffered-channels-in-golang/)
 * [C, Golang, Perl, and Unix](https://blog.afoolishmanifesto.com/posts/c-golang-perl-and-unix/)
` + "`" + `` + "`" + `` + "`" + `

Optionally takes -timeout to limit how long to wait for feeds to sync.  Passing
0 will disable timeout.  Default is 15s.

### ` + "`" + `slack-deaddrop` + "`" + `

` + "`" + `slack-deaddrop` + "`" + ` allows sending messages to a slack channel without looking at slack.
Typical usage is probably something like:

` + "`" + `` + "`" + `` + "`" + `bash
$ slack-deaddrop -channel general -text 'good morning!'
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `slack-open` + "`" + `

` + "`" + `slack-open` + "`" + ` opens a channel, group message, or direct message by name:

` + "`" + `` + "`" + `` + "`" + `bash
$ slack-open -channel general
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `sm-list` + "`" + `

` + "`" + `sm-list` + "`" + ` lists all of the available [Sweet Maria's](https://www.sweetmarias.com/) coffees
as json documents per line.  Here's how you might see the top ten coffees by
rating:

` + "`" + `` + "`" + `` + "`" + `bash
$ sm-list | jq -r '[.Score, .Title, .URL ] | @tsv' | sort -n | tail -10
` + "`" + `` + "`" + `` + "`" + `
### ` + "`" + `srv` + "`" + `

` + "`" + `srv` + "`" + ` will serve a directory over http; takes a single optional parameter which
is the dir to serve, default is ` + "`" + `.` + "`" + `:

` + "`" + `` + "`" + `` + "`" + `bash
$ srv ~
Serving /home/frew on [::]:21873
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `toml2json` + "`" + `

` + "`" + `toml2json` + "`" + ` reads [TOML](https://github.com/toml-lang/toml) on stdin and writes JSON
on stdout.


` + "`" + `` + "`" + `` + "`" + `bash
$ echo 'foo = "bar"` + "`" + ` | toml2json
{"foo":"bar"}
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `undefer` + "`" + `

` + "`" + `undefer` + "`" + ` takes a directory argument, prints contents of each file named before the
current date, and then deletes the file.

If the current date were ` + "`" + `2018-06-07` + "`" + ` the starred files would be printed and
then deleted:

` + "`" + `` + "`" + `` + "`" + `
 * 2018-01-01.txt
 * 2018-06-06-awesome-file.md
   2018-07-06.txt
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `uni` + "`" + `

` + "`" + `uni` + "`" + ` will describe the characters in the args.

` + "`" + `` + "`" + `` + "`" + `bash
$ uni ⢾
'⢾' @ 10430 aka BRAILLE PATTERN DOTS-234568 ( graphic | printable | symbol )
` + "`" + `` + "`" + `` + "`" + `

### ` + "`" + `yaml2json` + "`" + `

` + "`" + `yaml2json` + "`" + ` reads YAML on stdin and writes JSON on stdout.

## Debugging

In an effort to make debugging simpler, I've created three ways to see what
` + "`" + `leatherman` + "`" + ` is doing:

### Tracing

` + "`" + `LMTRACE=$somefile` + "`" + ` will write an execution trace to ` + "`" + `$somefile` + "`" + `; look at that with ` + "`" + `go tool trace $somefile` + "`" + `

Since so many of the tools are short lived my assumption is that the execution
trace will be the most useful.

### Profiling

` + "`" + `LMPROF=$somefile` + "`" + ` will write a cpu profile to ` + "`" + `$somefile` + "`" + `; look at that with ` + "`" + `go tool pprof -http localhost:10123 $somefile` + "`" + `

If you have a long running tool, the pprof http endpoint is exposed on
` + "`" + `localhost:6060/debug/pprof` + "`" + ` but picks a random port if that port is in use; the
port can be overridden by setting ` + "`" + `LMHTTPPROF=$someport` + "`" + `.
`)

var commandReadme map[string][]byte

func init() {
	commandReadme = map[string][]byte{
		"addrs": readme[1567:2064],

		"addrspec-to-tabs": readme[2064:2642],

		"alluni": readme[2642:3468],

		"backlight": readme[3468:3670],

		"clocks": readme[3670:3886],

		"csv2json": readme[3886:4026],

		"csv2md": readme[4026:4100],

		"debounce": readme[4100:4444],

		"dump-mozlz4": readme[4444:4627],

		"email2json": readme[4627:5342],

		"expand-url": readme[5342:5518],

		"export-bamboohr": readme[5518:5602],

		"export-bamboohr-tree": readme[5602:5689],

		"fn": readme[5689:5825],

		"group-by-date": readme[5825:6043],

		"minotaur": readme[6043:7862],

		"name2rune": readme[7862:7967],

		"netrc-password": readme[7967:8132],

		"pomotimer": readme[8132:8624],

		"prepend-hist": readme[8624:8865],

		"proj": readme[8865:10184],

		"render-mail": readme[10184:10318],

		"replace-unzip": readme[10318:10548],

		"rss": readme[10548:11474],

		"slack-deaddrop": readme[11474:11694],

		"slack-open": readme[11694:11828],

		"sm-list": readme[11828:12104],

		"srv": readme[12104:12287],

		"toml2json": readme[12287:12462],

		"undefer": readme[12462:12770],

		"uni": readme[12770:12932],

		"yaml2json": readme[12932:13009],
	}
}
