package minotaur

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
)

var (
	errNoScript = errors.New("no script passed, forgot -- ?")
	errNoDirs   = errors.New("no dirs passed")
	errUsage    = errors.New("usage: minotaur <dir1> [dir2 dir3] -- <cmd> [args to cmd]")
)

type config struct {
	dirs   []string
	script []string

	include, ignore *regexp.Regexp

	verbose, suppressArgs, runAtStart bool
}

func parseFlags(args []string) (config, error) {
	flags := flag.NewFlagSet("minotaur", flag.ExitOnError)

	var c config

	var ignoreStr, includeStr string

	flags.StringVar(&includeStr, "include", "", "regexp matching directories to include")
	flags.StringVar(&ignoreStr, "ignore", "(^.git|/.git$|/.git/)", "regexp matching directories to include")
	flags.BoolVar(&c.verbose, "verbose", false, "enable verbose output")
	flags.BoolVar(&c.suppressArgs, "suppress-args", false, "suppress event args args to script")
	flags.BoolVar(&c.runAtStart, "run-at-start", false, "run the script when you start")

	err := flags.Parse(args)
	if err != nil {
		return config{}, fmt.Errorf("flags.Parse: %w", err)
	}

	include := regexp.MustCompile(includeStr)
	ignore := regexp.MustCompile(ignoreStr)

	args = flags.Args()

	if len(args) < 3 {
		return config{}, errUsage
	}

	var token string

	token, args = args[0], args[1:]
	for len(args) > 0 && token != "--" {
		c.dirs = append(c.dirs, token)

		token, args = args[0], args[1:]
	}

	c.script = args

	if len(c.script) == 0 {
		return config{}, errNoScript
	}

	if len(c.dirs) == 0 {
		return config{}, errNoDirs
	}

	c.include = include
	c.ignore = ignore
	return c, nil
}
