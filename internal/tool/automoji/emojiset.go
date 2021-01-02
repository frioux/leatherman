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
	defer matchersMu.Unlock()
	matchersMu.Lock()

	s := &emojiSet{optional: make(map[string]bool)}

	for _, r := range matchers {
		if r.MatchString(m) {
			if r.Required {
				s.required = append(s.required, r.Emoji)
			} else {
				s.optional[r.Emoji] = true
			}
		}
	}

	m = strings.ToLower(m)

	if secretRE.MatchString(m) {
		s.add(turtle.Emojis["see_no_evil"])
		m = secretRE.ReplaceAllString(m, " ")
	}

	m = nonNameRE.ReplaceAllString(m, " ")
	words := strings.Split(m, " ")

	for _, word := range words {
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
		for _, word := range words {
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

func (s *emojiSet) alwaysAdd(e string) {

}

func (s *emojiSet) all(c int) []string {
	ret := make([]string, 0, c+len(s.required))

	for _, e := range s.required {
		ret = append(ret, e)
	}

	for e := range s.optional {
		if c != 0 && len(ret) == cap(ret) {
			break
		}
		ret = append(ret, e)
	}

	return ret
}
