package uni

import (
	"errors"
	"fmt"
	"io"

	"golang.org/x/text/unicode/rangetable"
	"golang.org/x/text/unicode/runenames"
)

func ToRune(args []string, _ io.Reader) error {
	t := rangetable.Assigned(unicodeVersion)

	if len(args) != 2 {
		return errors.New("name2rune requires a name")
	}
	search := args[1]
	var found bool
	rangetable.Visit(t, func(r rune) {
		name := runenames.Name(r)
		if name == search {
			fmt.Println(string(r))

			found = true
		}
	})

	if found {
		return nil
	}

	return errors.New("no rune found")
}
