package minotaur

import "errors"

var (
	errNoScript = errors.New("no script passed, forgot -- ?")
	errNoDirs   = errors.New("no dirs passed")
	errUsage    = errors.New("usage: minotaur <dir1> [dir2 dir3] -- <cmd> [args to cmd]")
)

func parseFlags(args []string) ([]string, []string, error) {
	if len(args) < 3 {
		return nil, nil, errUsage
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
		return nil, nil, errNoScript
	}

	if len(dirs) == 0 {
		return nil, nil, errNoDirs
	}

	return dirs, script, nil
}
