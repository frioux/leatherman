package email

import (
	"encoding/json"
	"io"
	"mime"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type email struct {
	Header map[string]string
}

func toJSON(path string, w io.Writer) error {
	enc := json.NewEncoder(w)

	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	defer file.Close()

	e, err := mail.ReadMessage(file)
	if err != nil {
		return errors.Wrap(err, "mail.ReadMessage, path="+path)
	}

	dec := new(mime.WordDecoder)

	eml := email{Header: make(map[string]string)}
	for k := range e.Header {
		header, err := dec.DecodeHeader(e.Header.Get(k))
		if err != nil {
			panic(err)
		}
		eml.Header[k] = header
	}

	err = enc.Encode(eml)
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
