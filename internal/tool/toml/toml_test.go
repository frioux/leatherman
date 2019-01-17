package toml

import (
	"fmt"
	"os"
	"strings"
)

func ExampleToJSON() {
	r := strings.NewReader("foo = \"bar\"\n")

	err := ToJSON(nil, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't ToJSON: %s\n", err)
		os.Exit(1)
	}
	// Output:
	// {"foo":"bar"}
}
