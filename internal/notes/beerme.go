package notes

import (
	"bufio"
	"errors"
	"io"
	"math/rand"
	"regexp"
)

var isItem = regexp.MustCompile(`^\s+\*\s+(.*)$`)

func BeerMe(r io.Reader) (string, error) {
	s := bufio.NewScanner(r)

	o := []string{}
	for s.Scan() {
		m := isItem.FindStringSubmatch(s.Text())
		if len(m) != 2 {
			continue
		}
		o = append(o, m[1])
	}

	if len(o) == 0 {
		return "", errors.New("never found anything")
	}

	rand.Shuffle(len(o), func(i, j int) { o[i], o[j] = o[j], o[i] })

	return o[0], nil
}
