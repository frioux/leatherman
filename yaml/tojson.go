package yaml // import "github.com/frioux/leatherman/yaml"

import (
	"encoding/json"
	"io"
	"os"

	"github.com/frioux/yaml"
	"github.com/pkg/errors"
)

// ToJSON reads YAML on stdin and writes JSON on stdout.
func ToJSON(_ []string, stdin io.Reader) error {
	d := yaml.NewDecoder(stdin)
	e := json.NewEncoder(os.Stdout)

	var data interface{}
	err := d.Decode(&data)
	if err != nil {
		return errors.Wrap(err, "Couldn't decode YAML")
	}

	err = e.Encode(data)
	if err != nil {
		return errors.Wrap(err, "Couldn't encode JSON")
	}

	return nil
}
