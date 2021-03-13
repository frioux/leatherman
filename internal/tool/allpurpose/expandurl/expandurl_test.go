package expandurl

import (
	"bytes"
	_ "embed"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

//go:embed testdata/test.html
var testhtml []byte

func TestRun(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if _, err := io.Copy(w, bytes.NewReader(testhtml)); err != nil {
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
