package srv

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	buf := &bytes.Buffer{}
	go serve(".", buf)

	sock := regexp.MustCompile("^Serving . on (.+)\n")

	var resp *http.Response
	var err error
	for i := 1; i < 11; i++ {
		m := sock.FindStringSubmatch(buf.String())
		if len(m) == 0 {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(i)))
			continue
		}

		resp, err = http.Get("http://" + m[1] + "/srv.go")
		if err != nil {
			t.Fatalf("Couldn't fetch srv.go: %s", err)
		}
		break
	}

	if resp == nil {
		t.Fatalf("couldn't get response from server within timeout")
	}

	f, err := os.Open("./srv.go")
	if err != nil {
		t.Fatalf("Couldn't open ./srv.go: %s", err)
	}
	expected, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Couldn't read ./srv.go: %s", err)
	}

	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Couldn't read response: %s", err)
	}

	if len(expected) == 0 {
		t.Fatal("Somehow got empty test case")
	}

	assert.Equal(t, expected, got)
}
