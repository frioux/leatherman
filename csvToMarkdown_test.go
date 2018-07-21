package main

import "strings"

func ExampleCSVToMarkdown() {
	r := strings.NewReader("foo,bar,baz\n1,2,3\n3,2,1")
	CSVToMarkdown(nil, r)
	// Output:
	// foo | bar | baz
	//  --- | --- | ---
	// 1 | 2 | 3
	// 3 | 2 | 1
}
