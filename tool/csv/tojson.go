package csv // import "github.com/frioux/leatherman/tool/csv"

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

// ToJSON converts input of CSV to JSON.
func ToJSON(_ []string, stdin io.Reader) error {
	reader := csv.NewReader(stdin)
	writer := json.NewEncoder(os.Stdout)

	header, err := reader.Read()
	if err != nil {
		return errors.New("can't read header, giving up")
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
			return errors.Wrap(err, "json.Encode")
		}
	}

	return nil
}
