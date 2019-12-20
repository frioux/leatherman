package undefer

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/frioux/leatherman/internal/testutil"
)

// for mock file
type mf string

func (f mf) Name() string       { return string(f) }
func (f mf) Size() int64        { return 0 }
func (f mf) Mode() os.FileMode  { return 0 }
func (f mf) ModTime() time.Time { return time.Unix(0, 0) }
func (f mf) IsDir() bool        { return false }
func (f mf) Sys() interface{}   { return nil }

func TestNewFiles(t *testing.T) {
	t.Parallel()

	f, err := newFiles("foo", []os.FileInfo{mf("bar")}, time.Now())
	if err != nil {
		t.Errorf("newFiles failed: %s", err)
		return
	}

	testutil.Equal(t, len(f), 0, "wrong length")

	f, err = newFiles("foo", []os.FileInfo{
		mf("silly.md"),
		mf("2018-12-12-xyzzy.md"),
		mf("2048-12-13.md"),
		mf("2018-02-13.md"),
	}, time.Date(2018, 12, 12, 0, 0, 0, 0, time.Local))
	if err != nil {
		t.Errorf("newFiles failed: %s", err)
		return
	}

	testutil.Equal(t, f, []string{"foo/2018-02-13.md", "foo/2018-12-12-xyzzy.md"}, "wrong paths")
}

func TestContent(t *testing.T) {
	t.Parallel()

	d, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("Couldn't set up TempDir: %s", err)
		return
	}
	defer os.RemoveAll(d)

	f, err := os.Create(filepath.Join(d, "y"))
	if err != nil {
		t.Errorf("couldn't os.Create: %s", err)
		return
	}

	if _, err := f.WriteString("in output1\n"); err != nil {
		t.Errorf("couldn't os.File.WriteString: %s", err)
		return
	}

	f, err = os.Create(filepath.Join(d, "z"))
	if err != nil {
		t.Errorf("couldn't os.Create: %s", err)
		return
	}

	if _, err := f.WriteString("in output2\n"); err != nil {
		t.Errorf("couldn't os.File.WriteString: %s", err)
		return
	}

	f, err = os.Create(filepath.Join(d, "t"))
	if err != nil {
		t.Errorf("couldn't os.Create: %s", err)
		return
	}

	if _, err := f.WriteString("XXX\n"); err != nil {
		t.Errorf("couldn't os.File.WriteString: %s", err)
		return
	}

	w := &bytes.Buffer{}
	if err := content([]string{filepath.Join(d, "y"), filepath.Join(d, "z")}, w); err != nil {
		t.Errorf("couldn't content(): %s", err)
		return
	}

	testutil.Equal(t, string(w.Bytes()), "in output1\nin output2\n", "wrong output")
}
