package automoji

import (
	"regexp"
	"strings"
)

var (
	nonNameRE = regexp.MustCompile(`[^a-z_]+`)
	secretRE  = regexp.MustCompile(`\|\|.+\|\|`)
)

func newEmojiSet(m string) *emojiSet {
	s := &emojiSet{
		message:  m,
		optional: make(map[string]bool),
	}

	m = strings.ToLower(m)

	if secretRE.MatchString(m) {
		s.optional["🙈"] = true
		m = secretRE.ReplaceAllString(m, " ")
	}

	m = nonNameRE.ReplaceAllString(m, " ")
	s.words = strings.Split(m, " ")

	return s
}

type emojiSet struct {
	message  string
	words    []string
	optional map[string]bool
	required []string
}

func (s *emojiSet) all(c int) []string {
	ret := make([]string, len(s.required), c+len(s.required))

	copy(ret, s.required)

	for e := range s.optional {
		if c != 0 && len(ret) == cap(ret) {
			break
		}
		ret = append(ret, e)
	}

	return ret
}
