package bamboo

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	cl, cleanup := testClientAndServer()
	defer cleanup()

	if err := cl.auth(); err != nil {
		t.Fatalf("Couldn't auth: " + err.Error())
	}

	buf := &bytes.Buffer{}
	if err := cl.tree(buf); err != nil {
		t.Fatalf("Couldn't load tree: " + err.Error())
	}

	assert.Equal(t, `{"tree":1}`, buf.String())
}

func TestDir(t *testing.T) {
	cl, cleanup := testClientAndServer()
	defer cleanup()

	if err := cl.auth(); err != nil {
		t.Fatalf("Couldn't auth: " + err.Error())
	}

	buf := &bytes.Buffer{}
	if err := cl.directory(buf); err != nil {
		t.Fatalf("Couldn't load directory: " + err.Error())
	}

	assert.Equal(t, `dir`, buf.String())
}
