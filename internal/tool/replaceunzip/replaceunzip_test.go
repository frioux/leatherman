package replaceunzip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractMember(t *testing.T) {
	zt, err := buildZip(map[string][]byte{
		"a": []byte("a"),
	})
	if err != nil {
		t.Fatalf("Couldn't make test zip: %s", err)
	}

	ms := zt.File
	assert.Equal(t, 1, len(ms))
	assert.Equal(t, "a", ms[0].Name)

	d, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't make tempdir: %s", err)
	}
	defer os.RemoveAll(d)

	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("Couldn't get working dir: %s", err)
	}

	if err := os.Chdir(d); err != nil {
		t.Fatalf("Couldn't chdir: %s", err)
	}
	defer os.Chdir(orig)

	err = extractMember(ms[0])
	assert.NoError(t, err)

	f, err := os.Open("a")
	assert.NoError(t, err)

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, f)
	assert.NoError(t, err)
	assert.Equal(t, "a", buf.String())
}

func TestSanitizeFilter(t *testing.T) {
	zt, err := buildZip(map[string][]byte{
		"a":            []byte("a"),
		"__MACOSX/foo": []byte(""),
		"x/.DS_Store":  []byte(""),
	})
	if err != nil {
		t.Fatalf("Couldn't make test zip: %s", err)
	}

	ms, err := sanitize("", zt.File)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ms))
	assert.Equal(t, "a", ms[0].Name)
}

func TestSanitizeSecure(t *testing.T) {
	zt, err := buildZip(map[string][]byte{
		"a":      []byte("a"),
		"b/../c": []byte("b"),
	})
	assert.NoError(t, err)

	_, err = sanitize("", zt.File)
	assert.Error(t, err) // No .. segments
}

func TestSanitizeSetRoot(t *testing.T) {
	zt, err := buildZip(map[string][]byte{
		"a": []byte("a"),
	})
	assert.NoError(t, err)

	ms, err := sanitize("c", zt.File)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ms))
	assert.Equal(t, "c/a", ms[0].Name)

	zt, err = buildZip(map[string][]byte{
		"a": []byte("a"),
	})
	assert.NoError(t, err)

	ms, err = sanitize("", zt.File)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ms))
	assert.Equal(t, "a", ms[0].Name)
}

func TestHasRoot(t *testing.T) {
	zt, err := buildZip(map[string][]byte{
		"a": []byte("a"),
		"b": []byte("b"),
	})
	if err != nil {
		t.Fatalf("Couldn't make test zip: %s", err)
	}

	assert.Equal(t, false, hasRoot(zt.File))

	zt, err = buildZip(map[string][]byte{
		"a/":  []byte(""),
		"a/a": []byte("a"),
		"a/b": []byte("b"),
	})
	if err != nil {
		t.Fatalf("Couldn't make test zip: %s", err)
	}

	assert.Equal(t, true, hasRoot(zt.File))
}

func TestGenRoot(t *testing.T) {
	assert.Equal(t, "foo", genRoot("foo.zip"))
	assert.Equal(t, "bar", genRoot("bar"))
}

func buildZip(files map[string][]byte) (*zip.Reader, error) {
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)

	for name, contents := range files {
		w, err := zw.Create(name)
		if err != nil {
			return nil, fmt.Errorf("zip.Create: %w", err)
		}
		if _, err := w.Write(contents); err != nil {
			return nil, fmt.Errorf("zipmember.Write: %w", err)
		}
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("zr.Close: %w", err)
	}

	return zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
}

func TestBuildZip(t *testing.T) {
	zr, err := buildZip(map[string][]byte{
		"frew": []byte("frew"),
		"bar":  []byte("bar"),
	})
	if err != nil {
		t.Fatalf("Couldn't build ZR: %s", err)
	}

	assert.Equal(t, 2, len(zr.File))
	for _, f := range zr.File {
		r, err := f.Open()
		if err != nil {
			t.Fatalf("Couldn't open member: %s", err)
		}
		b := &bytes.Buffer{}
		_, err = io.Copy(b, r)
		if err != nil {
			t.Fatalf("Couldn't copy member: %s", err)
		}
		err = r.Close()
		if err != nil {
			t.Fatalf("Couldn't close member: %s", err)
		}
		assert.Equal(t, f.Name, b.String())
	}
}
