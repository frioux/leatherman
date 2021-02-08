package groupbydate

import "strings"

func ExampleRun() {

	r := strings.NewReader(`2012-12-12T12:12:12.000Z
2012-12-12T13:12:12.000Z
2012-12-12T12:14:12.000Z
2012-12-12T12:12:22.000Z`)

	Run([]string{"group-by-date"}, r)
	// Output: 2012-12-12T00:00:00Z,4
}
