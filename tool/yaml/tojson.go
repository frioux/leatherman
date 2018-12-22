package yaml // import "github.com/frioux/leatherman/tool/yaml"

import (
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// ToJSON reads YAML on stdin and writes JSON on stdout.
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
			return errors.Wrap(err, "Couldn't decode YAML")
		}

		err = e.Encode(data)
		if err != nil {
			return errors.Wrap(err, "Couldn't encode JSON")
		}
	}

	return nil
}
