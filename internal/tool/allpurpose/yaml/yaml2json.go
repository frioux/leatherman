package yaml

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func ToJSON(_ []string, stdin io.Reader) error {
	d := yaml.NewDecoder(stdin)
	e := json.NewEncoder(os.Stdout)

	var data interface{}

	for {
		err := d.Decode(&data)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Couldn't decode YAML: %w", err)
		}

		err = e.Encode(data)
		if err != nil {
			return fmt.Errorf("Couldn't encode JSON: %w", err)
		}
	}

	return nil
}
