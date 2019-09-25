package prependemojihist

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

	in := strings.NewReader(`white star
rabbit
beer mug
skull and crossbones
`)
	out := &bytes.Buffer{}

	if err := run(history, in, pos, out); err != nil {
		t.Fatalf("Couldn't run `run`: %s", err)
	}

	assert.Equal(t, "skull and crossbones\nbeer mug\nwhite star\nrabbit\n", out.String())
}
