package toml

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	parser "github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

func ToJSON(_ []string, stdin io.Reader) error {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, stdin); err != nil {
		return errors.Wrap(err, "io.Copy")
	}
	var ret interface{}
	if err := parser.Unmarshal(buf.Bytes(), &ret); err != nil {
		return errors.Wrap(err, "toml.Unmarshal")
	}

	e := json.NewEncoder(os.Stdout)
	if err := e.Encode(ret); err != nil {
		return errors.Wrap(err, "json.Encode")
	}

	return nil
}
