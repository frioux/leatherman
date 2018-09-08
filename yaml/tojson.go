package yaml // import "github.com/frioux/leatherman/yaml"

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/frioux/yaml"
)

// ToJSON reads YAML on stdin and writes JSON on stdout.
func ToJSON(_ []string, stdin io.Reader) {
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
