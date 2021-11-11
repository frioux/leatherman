package datefmt_test

import (
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
	"github.com/frioux/leatherman/pkg/datefmt"
)

func TestTranslateFormat(t *testing.T) {
	t.Parallel()

	testutil.Equal(t, datefmt.TranslateFormat("%F"), "2006-01-02", "wrong date")
	testutil.Equal(t, datefmt.TranslateFormat("%FT%T"), "2006-01-02T15:04:05", "wrong date")
}
