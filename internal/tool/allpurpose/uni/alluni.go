package uni

import (
	"fmt"
	"io"

	"golang.org/x/text/unicode/rangetable"
	"golang.org/x/text/unicode/runenames"
)

func All(_ []string, _ io.Reader) error {
	t := rangetable.Assigned(unicodeVersion)

	rangetable.Visit(t, func(r rune) {
		name := runenames.Name(r)
		fmt.Println(name)
	})

	return nil
}
