package mozlz4

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/pierrec/lz4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func errHasPrefix(t *testing.T, err error, prefix string) bool {
	if !strings.HasPrefix(err.Error(), prefix) {
		t.Logf("Error «%s» does not start with «%s»\n", err, prefix)
		t.Fail()
		return false
	}
	return true
}

func TestHappyPath(t *testing.T) {
	str := "abcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyz"
	r := strings.NewReader(str)
	w := &bytes.Buffer{}

	err := compress(r, w, len(str))
	if err != nil {
		t.Logf("Failed to compress data: %s\n", err)
		t.Fail()
		return
	}

	rt, err := NewReader(w)
	if err != nil {
		t.Logf("Failed to decompress data: %s\n", err)
		t.Fail()
		return
	}

	out, err := ioutil.ReadAll(rt)
	if err != nil {
		t.Logf("Failed to RedaAll data: %s\n", err)
		t.Fail()
		return
	}

	assert.Equal(t, str, string(out), "data roundtripped")
}

func TestWrongLength(t *testing.T) {
	str := "abcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyzabcdefghijklmnopqrstuvzxyz"
	r := strings.NewReader(str)
	w := &bytes.Buffer{}

	err := compress(r, w, 12+len(str))
	if err != nil {
		t.Logf("Failed to compress data: %s\n", err)
		t.Fail()
		return
	}

	_, err = NewReader(w)
	assert.Equal(t, ErrWrongSize, errors.Cause(err))
}

func TestCantReadHeader(t *testing.T) {
	r := bytes.NewReader(nil)
	_, err := NewReader(r)
	errHasPrefix(t, err, "couldn't read header")
}

func TestWrongHeader(t *testing.T) {
	r := bytes.NewReader([]byte("lol"))
	_, err := NewReader(r)
	assert.Equal(t, ErrWrongHeader, errors.Cause(err))
}

func TestCantReadSize(t *testing.T) {
	r := bytes.NewReader([]byte(magicHeader + "x"))
	_, err := NewReader(r)
	errHasPrefix(t, err, "couldn't read size")
}

func TestCantDecompress(t *testing.T) {
	w := &bytes.Buffer{}
	w.Write([]byte(magicHeader))
	var size uint32 = 12
	binary.Write(w, binary.LittleEndian, size)
	w.Write([]byte{1, 2, 3, 4})
	t.Log(w.Bytes())

	r := bytes.NewReader(w.Bytes())
	_, err := NewReader(r)

	errHasPrefix(t, err, "couldn't decompress data")
}

func TestCantReadAll(t *testing.T) {
	w := &bytes.Buffer{}
	w.Write([]byte(magicHeader))
	var size uint32 = 12
	binary.Write(w, binary.LittleEndian, size)
	w.Write([]byte{1, 2, 3, 4})

	r := &ErrReader{Reader: bytes.NewReader(w.Bytes()), errAfter: 3}
	_, err := NewReader(r)

	errHasPrefix(t, err, "couldn't read compressed data")
}

type ErrReader struct {
	io.Reader
	errAfter int
}

func (r *ErrReader) Read(p []byte) (int, error) {
	if r.errAfter == 0 {
		return 0, errors.New("faked io error")
	}
	r.errAfter--
	return r.Reader.Read(p)
}

func compress(src io.Reader, dst io.Writer, intendedSize int) error {
	_, err := dst.Write([]byte(magicHeader))
	if err != nil {
		return errors.Wrap(err, "couldn't Write header")
	}
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return errors.Wrap(err, "couldn't ReadAll to Compress")
	}

	err = binary.Write(dst, binary.LittleEndian, uint32(intendedSize))
	if err != nil {
		return errors.Wrap(err, "couldn't encode length")
	}
	dstBytes := make([]byte, 10*len(b))
	sz, err := lz4.CompressBlockHC(b, dstBytes, -1)
	if err != nil {
		return errors.Wrap(err, "couldn't CompressBlock")
	}
	if sz == 0 {
		return errors.New("data incompressible")
	}
	_, err = dst.Write(dstBytes[:sz])
	if err != nil {
		return errors.Wrap(err, "couldn't Write compressed data")
	}

	return nil
}
