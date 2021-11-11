package toml_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/frioux/leatherman/internal/tool/allpurpose/toml"
)

func ExampleToJSON() {
	r := strings.NewReader("foo = \"bar\"\n")

	err := toml.ToJSON(nil, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't ToJSON: %s\n", err)
		os.Exit(1)
	}
	// Output:
	// {"foo":"bar"}
}
