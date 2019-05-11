package yaml

import (
	"encoding/json"
	"io"
	"os"

	"golang.org/x/xerrors"
	yaml "gopkg.in/yaml.v2"
)

/*
ToJSON reads YAML on stdin and writes JSON on stdout.

Command: yaml2json
*/
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
			return xerrors.Errorf("Couldn't decode YAML: %w", err)
		}

		err = e.Encode(data)
		if err != nil {
			return xerrors.Errorf("Couldn't encode JSON: %w", err)
		}
	}

	return nil
}
