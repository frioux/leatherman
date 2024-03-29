package toml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	parser "github.com/BurntSushi/toml"
)

func ToJSON(_ []string, stdin io.Reader) error {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, stdin); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	var ret interface{}
	if err := parser.Unmarshal(buf.Bytes(), &ret); err != nil {
		if terr, ok := err.(parser.ParseError); ok {
			return fmt.Errorf("toml.Unmarshal: %s", terr.ErrorWithPosition())
		} else {
			return fmt.Errorf("toml.Unmarshal: %w", err)
		}
	}

	e := json.NewEncoder(os.Stdout)
	if err := e.Encode(ret); err != nil {
		return fmt.Errorf("json.Encode: %w", err)
	}

	return nil
}
