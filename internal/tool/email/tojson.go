package email

import (
	"encoding/json"
	"io"
	"mime"
	"net/mail"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

type email struct {
	Header map[string]string
}

func toJSON(r io.Reader, w io.Writer) error {
	enc := json.NewEncoder(w)

	e, err := mail.ReadMessage(r)
	if err != nil {
		return xerrors.Errorf("mail.ReadMessage: %w", err)
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
		return xerrors.Errorf("json.Encode: %w", err)
	}

	return nil
}

func toJSONFromFile(path string, w io.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return xerrors.Errorf("os.Open: %w", err)
	}
	defer file.Close()

	return toJSON(file, w)
}

// ToJSON produces a JSON version of an email based on a list of globs.
func ToJSON(args []string, stdin io.Reader) error {
	if len(args) < 2 {
		return xerrors.New("please pass one or more globs")
	}

	for _, glob := range args[1:] {
		matches, err := filepath.Glob(glob)
		if err != nil {
			return xerrors.Errorf("filepath.Glob: %w", err)
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
