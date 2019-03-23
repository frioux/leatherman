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

func toJSON(r io.Reader, w io.Writer) error {
	enc := json.NewEncoder(w)

	e, err := mail.ReadMessage(r)
	if err != nil {
		return errors.Wrap(err, "mail.ReadMessage")
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

func toJSONFromFile(path string, w io.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	defer file.Close()

	return toJSON(file, w)
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
			err = toJSONFromFile(path, os.Stdout)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
