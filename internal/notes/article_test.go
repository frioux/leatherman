package notes

import (
	"strings"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestReadArticle(t *testing.T) {
	str := `
{
	"title": "frew",
	"tags": ["foo", "bar"],
	"id": "xyzzy",
	"extra": { "foo": "bar" }
}
# markdown

goes here`

	a, err := ReadArticle(strings.NewReader(str))

	if err != nil {
		t.Fatalf("couldn't readMetadata: %s", err)
	}
	testutil.Equal(t, a, Article{
		Title:       "frew",
		Tags:        []string{"foo", "bar"},
		Extra:       map[string]string{"foo": "bar"},
		RawContents: []byte(str),
		Body: []byte(`
# markdown

goes here`),
	}, "basic")
}

func TestReadArticleAndLua(t *testing.T) {
	str := `
{
	"title": "frew",
	"tags": ["foo", "bar"],
	"id": "xyzzy",
	"extra": { "foo": "bar" }
}
# markdown

goes here

` + "```mdlua\n" + `
function foo()

end

function bar()

end
` + "```\n"
	a, err := ReadArticle(strings.NewReader(str))

	if err != nil {
		t.Fatalf("couldn't readMetadata: %s", err)
	}
	testutil.Equal(t, a, Article{
		Title:       "frew",
		Tags:        []string{"foo", "bar"},
		Extra:       map[string]string{"foo": "bar"},
		RawContents: []byte(str),
		Body: []byte(`
# markdown

goes here

`),
		MarkdownLua: []byte(`
function foo()

end

function bar()

end
`),
	}, "basic")
}

var A Article

func BenchmarkReadArticle(b *testing.B) {
	var a Article
	for i := 0; i < b.N; i++ {
		a, _ = ReadArticle(strings.NewReader(`
		{
			"title": "frew",
			"tags": ["foo", "bar"],
			"id": "xyzzy",
			"extra": { "foo": "bar" }
		}
		# markdown

		goes here
		`))
	}
	A = a
}
