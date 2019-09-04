package mozlz4

// Package mozlz4 implements the undocumented format used by Mozilla Firefox.

// The mozlz4 format (also known as jsonlz4 and json.lz4) is used by Firefox for
// various storage backends.  The format is a magic header, a length, and an lz4
// compressed body.

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pierrec/lz4"
)

const magicHeader = "mozLz40\x00"

// Errors
var (
	ErrWrongHeader = errors.New("no mozLz4 header")
	ErrWrongSize   = errors.New("header size incorrect")
)

// NewReader returns an io.Reader that decompresses the data from r.
func NewReader(r io.Reader) (io.Reader, error) {
	header := make([]byte, len(magicHeader))
	_, err := r.Read(header)
	if err != nil {
		return nil, fmt.Errorf("couldn't read header: %w", err)
	}
	if string(header) != magicHeader {
		return nil, ErrWrongHeader
	}

	var size uint32
	err = binary.Read(r, binary.LittleEndian, &size)
	if err != nil {
		return nil, fmt.Errorf("couldn't read size: %w", err)
	}

	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("couldn't read compressed data: %w", err)
	}

	out := make([]byte, size)
	sz, err := lz4.UncompressBlock(src, out)

	if err != nil {
		return nil, fmt.Errorf("couldn't decompress data: %w", err)
	}
	// This could maybe be a warning or ignored entirely
	if sz != int(size) {
		return nil, fmt.Errorf("Header size %d, got %d: %w", size, sz, ErrWrongSize)
	}

	return bytes.NewReader(out), nil
}
