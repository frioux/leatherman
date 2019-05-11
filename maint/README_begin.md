# Leatherman - fREW's favorite multitool

This is a little project simply to make trivial tools in Go effortless for my
personal usage.  These tools are almost surely of low utility to most people,
but may be instructive nonetheless.

[I have CI/CD to build this into a single
binary](https://github.com/frioux/leatherman/blob/master/.travis.yml) and [an
`explode` tool that builds
symlinks](https://github.com/frioux/leatherman/blob/master/cmd/leatherman/explode.go)
for each tool in the busybox style.

[I have automation in my
dotfiles](https://github.com/frioux/dotfiles/blob/bef8303c19e2cefac7dfbec420ad8d45b95415b8/install.sh#L133-L141)
to pull the latest binary at install time and run the `explode` tool.

## Installation

Here's a copy pasteable script to install the leatherman on OSX or Linux:

``` bash
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
```

This asssumes that `~/bin` is in your path.  The `explode` command will create a
symlink for each of the tools listed below.

## Usage

Each tool takes different args, but to run a tool you can either use a symlink
(presumably created by `explode`):

``` bash
$ echo "---\nfoo: 1" | yaml2json
{"foo":1}
```

or use it as a subcommand:

``` bash
echo "---\nfoo: 1" | leatherman yaml2json
{"foo":1}
```

## Current tools

