// +build gofuzzbeta

package west

import (
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func FuzzParse(f *testing.F) {
	f.Add([]byte(""))
	f.Fuzz(func(t *testing.T, mdwn []byte) {
		d := Parse(mdwn)
		testutil.Equal(t, string(d.Markdown()), string(mdwn), "roundtrips")
	})
}
