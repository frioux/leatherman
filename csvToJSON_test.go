package main

import "strings"

func ExampleCSVToJSON() {
	r := strings.NewReader("foo,bar\n1,2\n2,3")

	CSVToJSON(nil, r)
	// Output:
	// {"bar":"2","foo":"1"}
	// {"bar":"3","foo":"2"}
}
