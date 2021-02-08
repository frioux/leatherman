package backlight

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func runAndCheck(t *testing.T, change, newBrightness int) {
	if err := run(change); err != nil {
		t.Errorf("run failed: %s", err)
		return
	}

	f, err := os.Open("./brightness")
	if err != nil {
		t.Errorf("os.Open failed: %s", err)
		return
	}

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, f)
	if err != nil {
		t.Errorf("io.Copy failed: %s", err)
		return
	}

	raw := buf.String()
	i, err := strconv.Atoi(raw[:len(raw)-1])
	if err != nil {
		t.Errorf("strconv.Atoi failed: %s", err)
		return
	}

	testutil.Equal(t, i, newBrightness, "wrong brightness")
}

func TestRun(t *testing.T) {
	t.Parallel()

	d, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create TempDir: %s", err)
	}
	defer os.RemoveAll(d)

	err = os.Chdir(d)
	if err != nil {
		t.Fatalf("Couldn't Chdir: %s", err)
	}

	// max_brightness
	f, err := os.Create("./max_brightness")
	if err != nil {
		t.Fatalf("Couldn't create max_brightness: %s", err)
	}
	_, err = f.WriteString("1000\n")
	if err != nil {
		t.Fatalf("Couldn't write max_brightness: %s", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("Couldn't close max_brightness: %s", err)
	}

	// brightness
	f, err = os.Create("./brightness")
	if err != nil {
		t.Fatalf("Couldn't create brightness: %s", err)
	}
	_, err = f.WriteString("750\n")
	if err != nil {
		t.Fatalf("Couldn't write brightness: %s", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("Couldn't close brightness: %s", err)
	}

	runAndCheck(t, 1, 760)
	runAndCheck(t, 2, 780)
	runAndCheck(t, -5, 730)
}
