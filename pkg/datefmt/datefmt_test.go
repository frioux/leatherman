package datefmt

import (
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestTranslateFormat(t *testing.T) {
	t.Parallel()

	testutil.Equal(t, TranslateFormat("%F"), "2006-01-02", "wrong date")
	testutil.Equal(t, TranslateFormat("%FT%T"), "2006-01-02T15:04:05", "wrong date")
}
