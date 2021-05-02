package notes

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

	"github.com/tailscale/hujson"
)

type Article struct {
	Title string

	// Filename will be set after parsing.
	Filename string `json:"-"`

	// URL will be set after parsing.
	URL string `json:"-"`

	// Raw tells the parser not to include the standard header and footer.
	Raw bool

	Tags []string

	ReviewedOn *string `json:"reviewed_on" db:"reviewed_on"`

	ReviewBy *string `json:"review_by" db:"review_by"`

	Extra map[string]string

	Body []byte

	// MarkdownLua can be used both to filter the Body at render time as
	// well as allowing interactive functionality implemented in the page
	// itself
	MarkdownLua []byte
}

var mdluaMatcher = regexp.MustCompile("(?s)```mdlua\n(.*?)```\n")

func ReadArticle(r io.Reader) (Article, error) {
	var a Article
	d := hujson.NewDecoder(r)
	err := d.Decode(&a)
	if err != nil {
		return a, fmt.Errorf("hujson.Decoder.Decode: %w", err)
	}
	raw, err := ioutil.ReadAll(d.Buffered())
	if err != nil {
		return a, fmt.Errorf("hujson.Decoder.Buffered+ioutil.ReadAll: %w", err)
	}

	c, err := ioutil.ReadAll(r)
	if err != nil {
		return a, err
	}

	raw = append(raw, c...)

	found := mdluaMatcher.FindAllSubmatch(raw, -1)
	for _, f := range found {
		a.MarkdownLua = append(a.MarkdownLua, f[1]...)
	}

	a.Body = mdluaMatcher.ReplaceAll(raw, nil)

	return a, err
}
