package prependemojihist

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/icza/backscanner"
	"golang.org/x/text/unicode/runenames"
)

/*
Run prints out the deduplicated lines from the passed file, converting characters to
unicode names, and then printing out the lines from STDIN, filtering out what's
already been printed.  Note that `alluni.pl` is in [my
dotfiles](https://github.com/frioux/dotfiles.)

```bash
$ alluni.pl | prefix-emoji-hist ~/.uni_history
```

Command: prepend-emoji-hist
*/
func Run(args []string, stdin io.Reader) error {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "you must pass a history file!")
		os.Exit(1)
	}

	file, err := os.Open(args[1])
	if err != nil {
		return fmt.Errorf("Couldn't open history file: %w", err)
	}
	fi, err := os.Stat(args[1])
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Couldn't stat history file: %w", err)
	}

	var pos int
	if fi != nil {
		pos = int(fi.Size())
	}

	return run(file, os.Stdin, pos, os.Stdout)
}

func run(history io.ReaderAt, in io.Reader, historyLength int, out io.Writer) error {
	seen := map[string]bool{}
	scanner := backscanner.New(history, historyLength)
	for {
		line, _, err := scanner.Line()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Couldn't read line: %w", err)
		}

		for i, char := range line {
			name := strings.ToLower(runenames.Name(char))
			if seen[name] {
				continue
			}
			seen[name] = true
			fmt.Fprintln(out, name)
			if i > 0 {
				fmt.Fprintln(os.Stderr, "Multiple characters on line, breaking loop")
				break
			}
		}
	}

	r := bufio.NewScanner(in)
	for r.Scan() {
		line := r.Text()
		if seen[line] {
			continue
		}
		seen[line] = true
		fmt.Fprintln(out, line)
	}

	return nil
}
