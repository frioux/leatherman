package srv

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	t.Parallel()

	ch := make(chan net.Addr)
	go serve(".", ch)

	var addr net.Addr
	timer := time.NewTimer(time.Second)
	select {
	case <-timer.C:
		t.Fatalf("couldn't get response from server within timeout")
	case addr = <-ch:
	}

	resp, err := http.Get("http://" + string(addr.String()) + "/srv.go")
	if err != nil {
		t.Fatalf("Couldn't fetch srv.go: %s", err)
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
