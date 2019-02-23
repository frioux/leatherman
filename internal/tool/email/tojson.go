package email

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type email struct {
	Header map[string]string
	Body   map[string]string
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

	eml := email{Header: make(map[string]string), Body: make(map[string]string)}
	for k := range e.Header {
		header, err := dec.DecodeHeader(e.Header.Get(k))
		if err != nil {
			panic(err)
		}
		eml.Header[k] = header
	}

	mediaType, params, err := mime.ParseMediaType(e.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(e.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}
			partMediaType, _, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
			if err != nil {
				log.Fatal(err)
			}
			if partMediaType != "text/plain" {
				continue
			}

			eml.Body[partMediaType] = string(slurp)
		}
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
