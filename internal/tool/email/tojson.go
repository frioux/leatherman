package email

import (
	"encoding/json"
	"io"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func toJSON(path string, w io.Writer) error {
	e := json.NewEncoder(w)

	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	defer file.Close()

	email, err := mail.ReadMessage(file)
	if err != nil {
		return errors.Wrap(err, "mail.ReadMessage, path="+path)
	}
	err = e.Encode(email)
	if err != nil {
		return errors.Wrap(err, "json.Encode")
	}

	return nil
}

// ToJSON produces a JSON version of an email based on a list of globs.
func ToJSON(args []string, stdin io.Reader) error {
	if len(args) < 2 {
		return errors.New("please pass one or more globs")
	}

	for _, glob := range args[1:] {
		matches, err := filepath.Glob(glob)
		if err != nil {
			return errors.Wrap(err, "filepath.Glob")
		}

		for _, path := range matches {
			err = toJSON(path, os.Stdout)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
