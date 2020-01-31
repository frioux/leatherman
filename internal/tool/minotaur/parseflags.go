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

type regexpFlag struct {
	regexp.Regexp
}

func (f *regexpFlag) Set(s string) error {
	re, err := regexp.Compile(s)
	if err != nil {
		return err
	}

	f.Regexp = *re

	return nil
}

type config struct {
	dirs   []string
	script []string

	include, ignore regexpFlag

	verbose, report, includeArgs, noRunAtStart bool
}

func parseFlags(args []string) (config, error) {
	flags := flag.NewFlagSet("minotaur", flag.ExitOnError)

	var c config

	if err := c.ignore.Set("(^.git|/.git$|/.git/)"); err != nil {
		return config{}, fmt.Errorf("couldn't create default ignore value: %w", err)
	}

	if err := c.include.Set(""); err != nil {
		return config{}, fmt.Errorf("couldn't create default include value: %w", err)
	}

	flags.Var(&c.include, "include", "regexp matching directories to include")
	flags.Var(&c.ignore, "ignore", "regexp matching directories to include")
	flags.BoolVar(&c.verbose, "verbose", false, "enable verbose output")
	flags.BoolVar(&c.includeArgs, "include-args", false, "include event args args to script")
	flags.BoolVar(&c.noRunAtStart, "no-run-at-start", false, "do not run the script when you start")
	flags.BoolVar(&c.report, "report", false, "wrap script runs with an ascii report")

	err := flags.Parse(args)
	if err != nil {
		return config{}, fmt.Errorf("flags.Parse: %w", err)
	}

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

	return c, nil
}
