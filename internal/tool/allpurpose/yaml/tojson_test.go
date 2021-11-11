package yaml_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/frioux/leatherman/internal/tool/allpurpose/yaml"
)

func ExampleToJSON() {
	r := strings.NewReader("---\n  - foo: 1\n  - bar: 2\n---\nx: 1\n")

	err := yaml.ToJSON(nil, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't ToJSON: %s\n", err)
		os.Exit(1)
	}
	// Output:
	// [{"foo":1},{"bar":2}]
	// {"x":1}
}
