package csv_test

import (
	"strings"

	"github.com/frioux/leatherman/internal/tool/allpurpose/csv"
)

func ExampleToJSON() {
	r := strings.NewReader("foo,bar\n1,2\n2,3")

	csv.ToJSON(nil, r)
	// Output:
	// {"bar":"2","foo":"1"}
	// {"bar":"3","foo":"2"}
}
