package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
)

// CSVToJSON converts input of CSV to JSON.
func CSVToJSON(args []string, stdin io.Reader) {
	reader := csv.NewReader(stdin)
	writer := json.NewEncoder(os.Stdout)

	header, err := reader.Read()
	if err != nil {
		log.Fatal("Can't read header, giving up")
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

		writer.Encode(toEncode)
	}
}
