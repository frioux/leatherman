package csv_test

import (
	"strings"

	"github.com/frioux/leatherman/internal/tool/allpurpose/csv"
)

func ExampleToMarkdown() {
	r := strings.NewReader("foo,bar,baz\n1,2,3\n3,2,1")
	csv.ToMarkdown(nil, r)
	// Output:
	// foo | bar | baz
	//  --- | --- | ---
	// 1 | 2 | 3
	// 3 | 2 | 1
}
