package main

import "strings"

func ExampleYAMLToJSON() {
	r := strings.NewReader("---\n  - foo: 1\n  - bar: 2")

	YAMLToJSON([]string{"yaml2json"}, r)
	// Output: [{"foo":1},{"bar":2}]
}
