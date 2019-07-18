package notes

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/frioux/amygdala/internal/twilio"
)

type rule struct {
	*regexp.Regexp
	action func(*http.Client, string, string, []twilio.Media) (string, error)
}

var rules []rule

func init() {
	rules = []rule{
		{Regexp: regexp.MustCompile(`(?i)^\s*inspire\s+me\s*$`), action: inspireMe},
		{Regexp: regexp.MustCompile(`(?i)^\s*remind\s+me\s*`), action: remind},
		{Regexp: deferPattern, action: deferMessage},
		{Regexp: regexp.MustCompile(``), action: todo},
	}
}

var errNoRule = errors.New("no rules matched")

func Dispatch(cl *http.Client, tok, input string, media []twilio.Media) (string, error) {
	for _, r := range rules {
		if !r.MatchString(input) {
			continue
		}
		return r.action(cl, tok, input, media)
	}

	return "", errNoRule
}
