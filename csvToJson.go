package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
)

func CsvToJson(args []string) {
	reader := csv.NewReader(os.Stdin)
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
