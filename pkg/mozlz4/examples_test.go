package mozlz4_test

import (
	"fmt"
	"io"
	"os"

	"github.com/frioux/leatherman/pkg/mozlz4"
)

func Example() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open: %s\n", err)
		os.Exit(1)
	}

	r, err := mozlz4.NewReader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create reader: %s\n", err)
		os.Exit(1)
	}
	if _, err = io.Copy(os.Stdout, r); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't copy: %s\n", err)
		os.Exit(1)
	}
}
