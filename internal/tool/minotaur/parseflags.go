package minotaur

import "errors"

var (
	errNoScript = errors.New("no script passed, forgot -- ?")
	errNoDirs   = errors.New("no dirs passed")
)

func parseFlags(args []string) ([]string, []string, error) {
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
