package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// ToMarkdown converts input of CSV to Markdown
func ToMarkdown(_ []string, stdin io.Reader) error {
	reader := csv.NewReader(stdin)

	header, err := reader.Read()
	if err != nil {
		return errors.New("can't read header, giving up")
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

	return nil
}
