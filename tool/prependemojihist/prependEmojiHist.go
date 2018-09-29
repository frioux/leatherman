package prependemojihist

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/icza/backscanner"
	"github.com/pkg/errors"
	"golang.org/x/text/unicode/runenames"
)

// PrependEmojiHist reads a history file from the first argument and reproduces
// it on standard out, but with the names of the characters per line instead of
// the characters themselves.  Reproduces stdin on stdout, leaving out anything
// already printed.
func PrependEmojiHist(args []string, stdin io.Reader) error {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "you must pass a history file!")
		os.Exit(1)
	}

	file, err := os.Open(args[1])
	if err != nil {
		return errors.Wrap(err, "Couldn't open history file")
	}
	fi, err := os.Stat(args[1])
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "Couldn't stat history file")
	}

	var pos int
	if fi != nil {
		pos = int(fi.Size())
	}

	seen := map[string]bool{}
	scanner := backscanner.New(file, pos)
	for {
		line, _, err := scanner.Line()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "Couldn't read line")
		}

		for i, char := range line {
			name := strings.ToLower(runenames.Name(char))
			if seen[name] {
				continue
			}
			seen[name] = true
			fmt.Println(name)
			if i > 0 {
				fmt.Fprintln(os.Stderr, "Multiple characters on line, breaking loop")
				break
			}
		}
	}

	r := bufio.NewScanner(stdin)
	for r.Scan() {
		line := r.Text()
		if seen[line] {
			continue
		}
		seen[line] = true
		fmt.Println(line)
	}

	return nil
}
