package prependemojihist

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestRun(t *testing.T) {
	t.Parallel()

	historyPath := "./testdata/hist.txt"
	history, err := os.Open(historyPath)
	if err != nil {
		t.Fatalf("Couldn't open for test: %s", err)
	}
	fi, err := os.Stat(historyPath)
	if err != nil {
		t.Fatalf("Couldn't stat for test: %s", err)
	}

	pos := int(fi.Size())

	in := strings.NewReader(`WHITE STAR
RABBIT
BEER MUG
SKULL AND CROSSBONES
`)
	out := &bytes.Buffer{}

	if err := run(history, in, pos, out); err != nil {
		t.Fatalf("Couldn't run `run`: %s", err)
	}

	testutil.Equal(t, out.String(), "SKULL AND CROSSBONES\nBEER MUG\nWHITE STAR\nRABBIT\n", "output not equal")
}
