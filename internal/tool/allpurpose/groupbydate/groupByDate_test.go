package groupbydate_test

import (
	"strings"

	"github.com/frioux/leatherman/internal/tool/allpurpose/groupbydate"
)

func ExampleRun() {

	r := strings.NewReader(`2012-12-12T12:12:12.000Z
2012-12-12T13:12:12.000Z
2012-12-12T12:14:12.000Z
2012-12-12T12:12:22.000Z`)

	groupbydate.Run([]string{"group-by-date"}, r)
	// Output: 2012-12-12T00:00:00Z,4
}
