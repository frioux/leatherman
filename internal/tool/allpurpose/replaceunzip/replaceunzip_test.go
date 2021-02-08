package replaceunzip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestExtractMember(t *testing.T) {
	t.Parallel()

	zt, err := buildZip(map[string][]byte{
		"a": []byte("a"),
	})
	if err != nil {
		t.Fatalf("Couldn't make test zip: %s", err)
	}

	ms := zt.File
	testutil.Equal(t, len(ms), 1, "length not equal")
	testutil.Equal(t, ms[0].Name, "a", "name not equal")

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

	if err := extractMember(ms[0]); err != nil {
		t.Errorf("couldn't extractMember: %s", err)
		return
	}

	f, err := os.Open("a")
	if err != nil {
		t.Errorf("couldn't os.Open: %s", err)
		return
	}

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, f); err != nil {
		t.Errorf("couldn't io.Copy: %s", err)
		return
	}
	testutil.Equal(t, buf.String(), "a", "file contents not equal")
}

func TestSanitizeFilter(t *testing.T) {
	t.Parallel()

	zt, err := buildZip(map[string][]byte{
		"a":            []byte("a"),
		"__MACOSX/foo": []byte(""),
		"x/.DS_Store":  []byte(""),
	})
	if err != nil {
		t.Fatalf("Couldn't make test zip: %s", err)
	}

	ms, err := sanitize("", zt.File)
	if err != nil {
		t.Errorf("couldn't sanitize: %s", err)
		return
	}
	if testutil.Equal(t, len(ms), 1, "length not equal") {
		testutil.Equal(t, ms[0].Name, "a", "name not equal")
	}
}

func TestSanitizeSecure(t *testing.T) {
	t.Parallel()

	zt, err := buildZip(map[string][]byte{
		"a":      []byte("a"),
		"b/../c": []byte("b"),
	})
	if err != nil {
		t.Errorf("couldn't buildZip: %s", err)
		return
	}

	if _, err := sanitize("", zt.File); err == nil {
		t.Errorf("sanitize should have errored")
	}
}

func TestSanitizeSetRoot(t *testing.T) {
	t.Parallel()

	zt, err := buildZip(map[string][]byte{
		"a": []byte("a"),
	})
	if err != nil {
		t.Errorf("buildZip failed: %s", err)
		return
	}

	ms, err := sanitize("c", zt.File)
	if err != nil {
		t.Errorf("sanitize failed: %s", err)
		return
	}
	if testutil.Equal(t, len(ms), 1, "length not equal") {
		testutil.Equal(t, ms[0].Name, "c/a", "name not equal")
	}

	zt, err = buildZip(map[string][]byte{
		"a": []byte("a"),
	})
	if err != nil {
		t.Errorf("buildZip failed: %s", err)
		return
	}

	ms, err = sanitize("", zt.File)
	if err != nil {
		t.Errorf("sanitize failed: %s", err)
		return
	}
	if testutil.Equal(t, len(ms), 1, "length not equal") {
		testutil.Equal(t, ms[0].Name, "a", "name not equal")
	}
}

func TestHasRoot(t *testing.T) {
	t.Parallel()

	zt, err := buildZip(map[string][]byte{
		"a": []byte("a"),
		"b": []byte("b"),
	})
	if err != nil {
		t.Fatalf("Couldn't make test zip: %s", err)
	}

	if hasRoot(zt.File) {
		t.Errorf("hasRoot() should be false")
	}

	zt, err = buildZip(map[string][]byte{
		"a/":  []byte(""),
		"a/a": []byte("a"),
		"a/b": []byte("b"),
	})
	if err != nil {
		t.Fatalf("Couldn't make test zip: %s", err)
	}

	if !hasRoot(zt.File) {
		t.Errorf("hasRoot() should be true")
	}
}

func TestGenRoot(t *testing.T) {
	t.Parallel()

	testutil.Equal(t, genRoot("foo.zip"), "foo", "genRoot() should be foo")
	testutil.Equal(t, genRoot("bar"), "bar", "genRoot() should be bar")
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
	t.Parallel()

	zr, err := buildZip(map[string][]byte{
		"frew": []byte("frew"),
		"bar":  []byte("bar"),
	})
	if err != nil {
		t.Fatalf("Couldn't build ZR: %s", err)
	}

	testutil.Equal(t, len(zr.File), 2, "length should be 2")
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
		testutil.Equal(t, b.String(), f.Name, "wrong name")
	}
}
