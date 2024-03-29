//go:build gofuzzbeta
// +build gofuzzbeta

package west_test

import (
	"testing"

	"github.com/frioux/leatherman/internal/notes/west"
	"github.com/frioux/leatherman/internal/testutil"
)

func FuzzParse(f *testing.F) {
	f.Add([]byte(""))
	for _, d := range crashers {
		f.Add([]byte(d))
	}
	f.Fuzz(func(t *testing.T, mdwn []byte) {
		d := west.Parse(mdwn)
		testutil.Equal(t, string(d.Markdown()), string(mdwn), "roundtrips")
	})
}
