package notes

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/frioux/amygdala/internal/dropbox"
	"github.com/frioux/amygdala/internal/personality"
	"github.com/pkg/errors"
)

var isItem = regexp.MustCompile(`^\s?\*\s+(.*?)\s*$`)
var mdLink = regexp.MustCompile(`^\[(.*)\]\((.*)\)$`)

func beerMe(r io.Reader) (string, error) {
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

	fmt.Println(mdLink.FindStringSubmatch(o[0]))
	if l := mdLink.FindStringSubmatch(o[0]); len(l) == 3 {
		return fmt.Sprintf("[%s]( %s )", l[1], l[2]), nil
	}

	return o[0], nil
}

func inspireMe(cl *http.Client, tok, _ string) (string, error) {
	r, err := dropbox.Download(cl, tok, "/notes/content/posts/inspiration.md")
	if err != nil {
		return personality.Err(), errors.Wrap(err, "dropbox.Download")
	}
	n, err := beerMe(r)
	if err != nil {
		return personality.Err(), err
	}
	return n, nil
}
