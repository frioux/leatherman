package csv

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"

	"golang.org/x/xerrors"
)

/*
ToJSON reads CSV on stdin and writes JSON on stdout; first line of input is the
header, and thus the keys of the JSON.

Command: csv2json
*/
func ToJSON(_ []string, stdin io.Reader) error {
	reader := csv.NewReader(stdin)
	writer := json.NewEncoder(os.Stdout)

	header, err := reader.Read()
	if err != nil {
		return xerrors.Errorf("can't read header, giving up")
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if len(record) != len(header) {
			continue
		}
		if err != nil {
			log.Println(err)
		}
		toEncode := map[string]string{}
		for v, x := range header {
			toEncode[x] = record[v]
		}

		err = writer.Encode(toEncode)
		if err != nil {
			return xerrors.Errorf("json.Encode: %w", err)
		}
	}

	return nil
}
