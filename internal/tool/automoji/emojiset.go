package automoji

import (
	"regexp"
	"strings"
	"time"

	"github.com/hackebrot/turtle"
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
		s.add(turtle.Emojis["see_no_evil"])
		m = secretRE.ReplaceAllString(m, " ")
	}

	m = nonNameRE.ReplaceAllString(m, " ")
	s.words = strings.Split(m, " ")

	for _, word := range s.words {
		if word == "" {
			continue
		}

		if e, ok := turtle.Emojis[word]; ok {
			s.add(e)
		}

		if es := turtle.Category(word); es != nil {
			for _, e := range es {
				s.add(e)
			}
		}

		if es := turtle.Keyword(word); es != nil {
			for _, e := range es {
				s.add(e)
			}
		}
	}
	if len(s.optional) == 0 { // since this always finds too much, only use it when nothing is found
		for _, word := range s.words {
			if es := turtle.Search(word); es != nil {
				for _, e := range es {
					s.add(e)
				}
			}
		}
	}

	return s
}

type emojiSet struct {
	message  string
	words    []string
	optional map[string]bool
	required []string
}

func (s *emojiSet) len() int {
	return len(s.optional) + len(s.required)
}

func (s *emojiSet) add(e *turtle.Emoji) {
	t := time.Now().Local()
	isFlagDay := t.Month() == 6 && t.Day() == 14

	// not flag day, don't include flags
	if !isFlagDay && e.Category == "flags" {
		return
	}

	// flag day, *only* include flags
	if isFlagDay && e.Category != "flags" {
		return
	}

	s.optional[e.Char] = true
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
