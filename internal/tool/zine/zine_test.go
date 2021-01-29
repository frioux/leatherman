package zine

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestFullRender(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("couldn't create test: %s", err)
	}
	defer os.RemoveAll(d)

	if err := render([]string{"render", "-static", "./testdata", "-root", "./testdata", "-out", d}); err != nil {
		t.Errorf("Rendered produced unexpected error: %s", err)
	}

	b, err := ioutil.ReadFile(filepath.Join(d, "cats", "index.html"))
	if err != nil {
		t.Errorf("Couldn't read output: %s", err)
	}

	testutil.Equal(t, `<p>This is the header!</p>
<h1 id="cats">cats</h1>
<p>cats are the best.</p>
<p>this is the footer!</p>
`, string(b), "cats generated correctly")

	// XXX add test for index
}
