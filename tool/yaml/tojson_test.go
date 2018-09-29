package yaml

import "strings"

func ExampleToJSON() {
	r := strings.NewReader("---\n  - foo: 1\n  - bar: 2")

	ToJSON([]string{"yaml2json"}, r)
	// Output: [{"foo":1},{"bar":2}]
}
