package minotaur

import (
	"flag"
	"regexp"

	"github.com/pkg/errors"
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

	verbose bool
}

func parseFlags(args []string) (config, error) {
	flags := flag.NewFlagSet("minotaur", flag.ExitOnError)

	var ignoreStr, includeStr string
	var verbose bool

	flags.StringVar(&includeStr, "include", "", "regexp matching directories to include")
	flags.StringVar(&ignoreStr, "ignore", "(^.git|/.git$|/.git/)", "regexp matching directories to include")
	flags.BoolVar(&verbose, "verbose", false, "enable verbose output")

	err := flags.Parse(args)
	if err != nil {
		return config{}, errors.Wrap(err, "flags.Parse")
	}

	include := regexp.MustCompile(includeStr)
	ignore := regexp.MustCompile(ignoreStr)

	args = flags.Args()

	if len(args) < 3 {
		return config{}, errUsage
	}

	var dirs, script []string

	var token string

	token, args = args[0], args[1:]
	for len(args) > 0 && token != "--" {
		dirs = append(dirs, token)

		token, args = args[0], args[1:]
	}

	script = args

	if len(script) == 0 {
		return config{}, errNoScript
	}

	if len(dirs) == 0 {
		return config{}, errNoDirs
	}

	return config{
		dirs:    dirs,
		script:  script,
		include: include,
		ignore:  ignore,
		verbose: verbose,
	}, nil
}
