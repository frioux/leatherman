package undefer

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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
	f, err := newFiles("foo", []os.FileInfo{mf("bar")}, time.Now())

	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(f))
	}

	f, err = newFiles("foo", []os.FileInfo{
		mf("silly.md"),
		mf("2018-12-12-xyzzy.md"),
		mf("2048-12-13.md"),
		mf("2018-02-13.md"),
	}, time.Date(2018, 12, 12, 0, 0, 0, 0, time.Local))

	if assert.NoError(t, err) {
		assert.Equal(t,
			[]string{"foo/2018-02-13.md", "foo/2018-12-12-xyzzy.md"}, f)
	}
}

func TestContent(t *testing.T) {
	fs = afero.NewMemMapFs()

	f, err := fs.Create("x/y")
	if assert.NoError(t, err) {
		_, err := f.WriteString("in output1\n")
		assert.NoError(t, err)
	}

	f, err = fs.Create("x/z")
	if assert.NoError(t, err) {
		_, err := f.WriteString("in output2\n")
		assert.NoError(t, err)
	}

	f, err = fs.Create("x/t")
	if assert.NoError(t, err) {
		_, err := f.WriteString("XXX\n")
		assert.NoError(t, err)
	}

	w := &bytes.Buffer{}
	err = content([]string{"x/y", "x/z"}, w)
	if assert.NoError(t, err) {
		assert.Equal(t, "in output1\nin output2\n", string(w.Bytes()))
	}
}
