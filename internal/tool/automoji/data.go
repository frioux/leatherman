package automoji

import (
	"encoding/json"
	"regexp"

	"github.com/frioux/leatherman/internal/dropbox"
)

type JSONRE struct {
	*regexp.Regexp
}

func (r *JSONRE) UnmarshalJSON(b []byte) error {
	var (
		err error
		str string
	)

	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	r.Regexp, err = regexp.Compile(str)
	return err
}

type matcher struct {
	// required will disable the use of randomness and always
	// insert the specified response
	Required bool

	// JSONRE matches against input
	JSONRE JSONRE

	// emoji to add to the response; use just the emoji (☃)
	// to add standard responses, or a string (constanza) to
	// use one of the custom emoji
	Emoji string
}

func (m matcher) MatchString(s string) bool {
	return m.JSONRE.MatchString(s)
}

func loadMatchers(dbCl dropbox.Client, path string) ([]matcher, error) {
	r, err := dbCl.Download(path)
	if err != nil {
		return nil, err
	}

	var m []matcher

	d := json.NewDecoder(r)
	if err := d.Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}
