package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pierrec/lz4"
)

const magicHeader = "mozLz40\x00"

func DumpMozLZ4(args []string) {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s session.jsonlz4\n", args[0])
		os.Exit(1)
	}
	file, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open: %s\n", err)
		os.Exit(1)
	}

	header := make([]byte, len(magicHeader))
	_, err = file.Read(header)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read header: %s\n", err)
		os.Exit(1)
	}
	if string(header) != magicHeader {
		fmt.Fprintf(os.Stderr, "Incorrect header: %s\n", err)
		os.Exit(1)
	}
	b := make([]byte, 4)
	file.Read(b)

	var size uint32
	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read size: %s\n", err)
		os.Exit(1)
	}

	src, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read compressed data: %s\n", err)
		os.Exit(1)
	}

	out := make([]byte, size)
	_, err = lz4.UncompressBlock(src, out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't decompress data: %s\n", err)
		os.Exit(1)
	}
	fmt.Print(string(out))
}
