package yaml

import (
	"fmt"
	"os"
	"strings"
)

func ExampleToJSON() {
	r := strings.NewReader("---\n  - foo: 1\n  - bar: 2")

	err := ToJSON(nil, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't ToJSON: %s\n", err)
		os.Exit(1)
	}
	// Output:
	// [{"foo":1},{"bar":2}]
}
