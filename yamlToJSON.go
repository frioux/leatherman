package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/frioux/yaml"
)

// YAMLToJSON reads YAML on stdin and writes JSON on stdout.
func YAMLToJSON(args []string, stdin io.Reader) {
	d := yaml.NewDecoder(stdin)
	e := json.NewEncoder(os.Stdout)

	var data interface{}
	err := d.Decode(&data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't decode YAML: %s\n", err)
		os.Exit(1)
	}

	err = e.Encode(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't encode JSON: %s\n", err)
		os.Exit(1)
	}
}
