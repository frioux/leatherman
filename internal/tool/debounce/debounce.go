package debounce

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

/*
Run debounces input from stdin to stdout

The default lockout time is one second, you can override that with the
`--lockoutTime` argument.  By default the trailing edge triggers output, so
output is emitted after there is no input for the lockout time.  You can change
this behavior by passing the `--leadingEdge` flag.

Command: debounce
*/
func Run(args []string, stdin io.Reader) error {
	var timeout time.Duration
	var leading, h, help bool

	flags := flag.NewFlagSet("debounce", flag.ExitOnError)

	flags.DurationVar(&timeout, "lockoutTime", time.Second, "amount of time between output")
	flags.BoolVar(&leading, "leadingEdge", false, "trigger at leading edge of cycle")
	flags.BoolVar(&h, "h", false, "help for debounce")
	flags.BoolVar(&help, "help", false, "help for debounce")

	err := flags.Parse(args[1:])
	if err != nil {
		return fmt.Errorf("flags.Parse: %w", err)
	}

	if h || help {
		fmt.Println("\n" +
			" debounce          [--leadingEdge] [--lockoutTime 2s]\n" +
			"                   [-h|--help]\n" +
			"\n" +
			"    --leadingEdge   pass this flag to output at the leading edge of a cycle\n" +
			"                    (off by default)\n" +
			"    --lockoutTime   set the lockout time in seconds, default is 1 second\n" +
			"\n" +
			"    -h --help       print usage message and exit\n" +
			"\n" +
			"\n" +
			"debounce creates cycles based on the lockout time.  The cycle\n" +
			"starts on the first line sent and stops after no lines are sent\n" +
			"within a period of the lockout time\n" +
			"\n" +
			"\n" +
			"The following would run tests after a second of 'silence' after a\n" +
			"save\n" +
			"\n" +
			" inotifywait -mr -e modify,move . | debounce | xargs -i{} make test\n" +
			"",
		)
		return nil
	}

	b := newBouncer(!leading, os.Stdout, timeout)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()

		b.Write(time.Now(), []byte(line+"\n"))
	}

	return s.Err()
}
