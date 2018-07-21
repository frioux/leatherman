package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// CSVToMarkdown converts input of CSV to Markdown
func CSVToMarkdown(_ []string, stdin io.Reader) {
	reader := csv.NewReader(stdin)

	header, err := reader.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't read header, giving up")
		os.Exit(1)
	}

	fmt.Println(strings.Join(header, " | "))
	for range header[:len(header)-1] {
		fmt.Print(" --- |")
	}
	fmt.Println(" ---")

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if len(record) != len(header) {
			continue
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't parse row: %s\n", err)
			continue
		}
		fmt.Println(strings.Join(record, " | "))
	}
}
