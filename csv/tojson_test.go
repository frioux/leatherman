package csv // import "github.com/frioux/leatherman/csv"

import "strings"

func ExampleToJSON() {
	r := strings.NewReader("foo,bar\n1,2\n2,3")

	ToJSON(nil, r)
	// Output:
	// {"bar":"2","foo":"1"}
	// {"bar":"3","foo":"2"}
}
