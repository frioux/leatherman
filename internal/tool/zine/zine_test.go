package zine

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestReadArticle(t *testing.T) {
	a, err := readArticle(strings.NewReader(`
{
	"title": "frew",
	"tags": ["foo", "bar"],
	"id": "xyzzy",
	"extra": { "foo": "bar" }
}
# markdown

goes here`))

	if err != nil {
		t.Fatalf("couldn't readMetadata: %s", err)
	}
	testutil.Equal(t, a, article{
		Title: "frew",
		Tags:  []string{"foo", "bar"},
		Extra: map[string]string{"foo": "bar"},
		Body: []byte(`
# markdown

goes here`),
	}, "basic")
}

var A article

func BenchmarkReadArticle(b *testing.B) {
	var a article
	for i := 0; i < b.N; i++ {
		a, _ = readArticle(strings.NewReader(`
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

func BenchmarkInsertMetadata(b *testing.B) {
	b.StopTimer()

	db, err := newDB()
	if err != nil {
		b.Fatalf("couldn't create db: %s", err)
	}
	defer db.Close()

	a := article{
		Title: "frew",
		Tags:  []string{"foo", "bar"},
		Extra: map[string]string{"foo": "bar"},
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if err := db.insertArticle(a); err != nil {
			b.Fatalf("couldn't insert article: %s", err)
		}
	}
}

func TestQuery(t *testing.T) {
	z, err := newZine()
	if err != nil {
		t.Fatalf("couldn't create zine: %s", err)
	}
	defer z.DB.Close()

	a := article{
		Title: "frew",
		Tags:  []string{"foo", "bar"},
		Extra: map[string]string{"foo": "bar"},
	}
	for i := 0; i < 1000; i++ {
		if err := z.insertArticle(a); err != nil {
			t.Fatalf("couldn't insert article: %s", err)
		}
	}

	_, err = z.q("SELECT * FROM _ LIMIT 5")
	if err != nil {
		t.Fatalf("Invalid query: %s", err)
	}
}

var C int

func BenchmarkQuery(b *testing.B) {
	b.StopTimer()
	z, err := newZine()
	if err != nil {
		b.Fatalf("couldn't create zine: %s", err)
	}

	a := article{
		Title: "frew",
		Tags:  []string{"foo", "bar"},
		Extra: map[string]string{"foo": "bar"},
	}
	for i := 0; i < 1000; i++ {
		if err := z.insertArticle(a); err != nil {
			b.Fatalf("couldn't insert article: %s", err)
		}
	}

	b.StartTimer()
	var c int
	for i := 0; i < b.N; i++ {
		r, err := z.q("SELECT * FROM _ LIMIT 5")
		if err != nil {
			b.Fatalf("Invalid query: %s", err)
		}
		c += len(r)
	}

	C = c
}

func TestRender(t *testing.T) {
	z, err := newZine()
	if err != nil {
		t.Fatalf("couldn't create db: %s", err)
	}

	a := article{
		Title: "frew",
		Tags:  []string{"foo", "bar"},
		Extra: map[string]string{"foo": "bar"},
	}
	for i := 0; i < 1000; i++ {
		if err := z.insertArticle(a); err != nil {
			t.Fatalf("couldn't insert article: %s", err)
		}
	}
	got, err := z.render(article{Title: "x", Body: []byte(`hello! *{{ with $r := (q "SELECT COUNT(*) AS c FROM _")}}{{ index $r 0 "c" }}{{end}}*`)})
	if err != nil {
		t.Errorf("should not have gotten an error: %s", err)
		return
	}

	testutil.Equal(t, string(got), "<p>start</p>\n<p>hello! <em>2000</em></p>\n<p>end</p>\n", "simple")
}

var S string

func BenchmarkRender(b *testing.B) {
	b.StopTimer()
	z, err := newZine()
	if err != nil {
		b.Fatalf("couldn't create db: %s", err)
	}

	a := article{
		Title: "frew",
		Tags:  []string{"foo", "bar"},
		Extra: map[string]string{"foo": "bar"},
	}
	for i := 0; i < 1000; i++ {
		if err := z.insertArticle(a); err != nil {
			b.Fatalf("couldn't insert article: %s", err)
		}
	}

	var out []byte
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var err error
		out, err = z.render(article{Title: "X", Body: []byte(`hello! *{{ with $r := (q "SELECT COUNT(*) AS c FROM _")}}{{ index $r 0 "c" }}{{end}}*`)})
		if err != nil {
			b.Errorf("should not have gotten an error: %s", err)
			return
		}
	}

	S = string(out)
}

func BenchmarkLoadNilNil(b *testing.B) {
	var (
		z   *zine
		err error
	)
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		z, err = newZine()
		if err != nil {
			b.Fatalf("couldn't create zine: %s", err)
		}
		z.root = "testdata"

		b.StartTimer()
		if err := z.load(nil); err != nil {
			b.Fatalf("couldn't load: %s", err)
		}
	}
	b.StopTimer()

	S = z.root
}

func BenchmarkLoadXY(b *testing.B) {
	var (
		z   *zine
		err error
		c   int
	)
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		z, err = newZine()
		if err != nil {
			b.Fatalf("couldn't create zine: %s", err)
		}
		z.root = "testdata"

		b.StartTimer()
		var as []article
		if err := z.load(&as); err != nil {
			b.Fatalf("couldn't load: %s", err)
		}
		c += len(as)
	}
	b.StopTimer()

	S = z.root
	C = c
}

func TestFullRender(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("couldn't create test: %s", err)
	}
	defer os.RemoveAll(d)

	if err := render([]string{"render", "-static", "./testdata", "-root", "./testdata", "-out", d}); err != nil {
		t.Errorf("Rendered produced unexpected error: %s", err)
	}

	b, err := ioutil.ReadFile(filepath.Join(d, "cats", "index.html"))
	if err != nil {
		t.Errorf("Couldn't read output: %s", err)
	}

	testutil.Equal(t, `<p>This is the header!</p>
<h1 id="cats">cats</h1>
<p>cats are the best.</p>
<p>this is the footer!</p>
`, string(b), "cats generated correctly")

	// XXX add test for index
}