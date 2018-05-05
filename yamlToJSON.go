package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/frioux/yaml"
)

func YAMLToJSON(args []string) {
	d := yaml.NewDecoder(os.Stdin)
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
