package bamboo

import (
	"bytes"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestTree(t *testing.T) {
	t.Parallel()

	cl, cleanup := testClientAndServer()
	defer cleanup()

	if err := cl.auth(); err != nil {
		t.Fatalf("Couldn't auth: " + err.Error())
	}

	buf := &bytes.Buffer{}
	if err := cl.tree(buf); err != nil {
		t.Fatalf("Couldn't load tree: " + err.Error())
	}

	testutil.Equal(t, buf.String(), `{"tree":1}`, "wrong tree found")
}

func TestDir(t *testing.T) {
	t.Parallel()

	cl, cleanup := testClientAndServer()
	defer cleanup()

	if err := cl.auth(); err != nil {
		t.Fatalf("Couldn't auth: " + err.Error())
	}

	buf := &bytes.Buffer{}
	if err := cl.directory(buf); err != nil {
		t.Fatalf("Couldn't load directory: " + err.Error())
	}

	testutil.Equal(t, buf.String(), `dir`, "wrong directory found")
}
