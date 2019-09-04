package email

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/mail"
	"os"
	"path/filepath"
)

type email struct {
	Header map[string]string
}

func toJSON(r io.Reader, w io.Writer) error {
	enc := json.NewEncoder(w)

	e, err := mail.ReadMessage(r)
	if err != nil {
		return fmt.Errorf("mail.ReadMessage: %w", err)
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
		return fmt.Errorf("json.Encode: %w", err)
	}

	return nil
}

func toJSONFromFile(path string, w io.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer file.Close()

	return toJSON(file, w)
}

/*
ToJSON produces a JSON representation of an email from a list of globs.  Only
headers are currently supported, patches welcome to support bodies.

```bash
$ email2json '/home/frew/var/mail/mitsi/cur/*' | head -1 | jq .
{
  "Header": {
    "Content-Type": "multipart/alternative; boundary=00163642688b8ef3070464661533",
    "Date": "Thu, 5 Mar 2009 15:45:17 -0600",
    "Delivered-To": "xxx",
    "From": "fREW Schmidt <xxx>",
    "Message-Id": "<fb3648c60903051345o728960f5l6cfb9e1f324bbf50@mail.gmail.com>",
    "Mime-Version": "1.0",
    "Received": "by 10.103.115.8 with HTTP; Thu, 5 Mar 2009 13:45:17 -0800 (PST)",
    "Subject": "STATION",
    "To": "xyzzy@googlegroups.com"
  }
}
```

Command: email2json
*/
func ToJSON(args []string, stdin io.Reader) error {
	if len(args) < 2 {
		return errors.New("please pass one or more globs")
	}

	for _, glob := range args[1:] {
		matches, err := filepath.Glob(glob)
		if err != nil {
			return fmt.Errorf("filepath.Glob: %w", err)
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
