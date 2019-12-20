package rss

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestRun(t *testing.T) {
	t.Parallel()

	f, err := ioutil.TempFile("", "*.js")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		f, err := os.Open("./testdata/rss.xml")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, _ = io.Copy(w, f)
	}))
	defer ts.Close()

	_, err = f.WriteString("[]")
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	err = run(f.Name(), []string{ts.URL}, buf)
	if err != nil {
		t.Errorf("Failed to run: %s", err)
		return
	}
	expected, err := ioutil.ReadFile("./testdata/output.json")
	if err != nil {
		t.Errorf("couldn't open test data file: %s", err)
		return
	}
	testutil.Equal(t, string(expected), buf.String(), "wrong json")
}
