package deferlwn

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeferLink(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create dir: %s", err)
	}
	defer os.RemoveAll(d)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		f, err := os.Open("./page.html")
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(w, f)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	buf := &bytes.Buffer{}
	bufErr := &bytes.Buffer{}
	err = run(d, strings.NewReader("[title]("+ts.URL+")"), buf, bufErr)
	assert.NoError(t, err)
	assert.Empty(t, buf.Bytes())
	assert.Empty(t, bufErr.Bytes())

	dh, err := os.Open(d)
	if !assert.NoError(t, err) {
		return
	}

	n, err := dh.Readdirnames(10)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []string{"2019-04-05-lwn.md"}, n)

	f, err := os.Open(filepath.Join(d, "2019-04-05-lwn.md"))
	if !assert.NoError(t, err) {
		return
	}

	buf = &bytes.Buffer{}
	_, err = io.Copy(buf, f)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "[title]("+ts.URL+")\n", buf.String())
}
