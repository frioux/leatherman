package expandurl

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestRun(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		f, err := os.Open("./testdata/test.html")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err := io.Copy(w, f); err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	buf := &bytes.Buffer{}
	err := run(strings.NewReader(ts.URL), buf)
	if err != nil {
		t.Errorf("run errored: %s", err)
		return
	}
	testutil.Equal(t, buf.String(), "[fREW Schmidt's Foolish Manifesto]("+ts.URL+")\n", "wrong output")
}
