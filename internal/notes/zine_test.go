package notes

import (
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func BenchmarkInsertMetadata(b *testing.B) {
	b.StopTimer()

	db, err := NewDB()
	if err != nil {
		b.Fatalf("couldn't create db: %s", err)
	}
	defer db.Close()

	a := Article{
		Title: "frew",
		Tags:  []string{"foo", "bar"},
		Extra: map[string]string{"foo": "bar"},
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if err := db.InsertArticle(a); err != nil {
			b.Fatalf("couldn't insert article: %s", err)
		}
	}
}

func TestQuery(t *testing.T) {
	z, err := NewZine()
	if err != nil {
		t.Fatalf("couldn't create zine: %s", err)
	}
	defer z.Close()

	a := Article{
		Title:    "frew",
		Filename: "frew.md",
		URL:      "/frew/",
		Tags:     []string{"foo", "bar"},
		Extra:    map[string]string{"foo": "bar"},
	}
	for i := 0; i < 1000; i++ {
		if err := z.InsertArticle(a); err != nil {
			t.Fatalf("couldn't insert article: %s", err)
		}
	}

	_, err = z.Q("SELECT * FROM _ LIMIT 5")
	if err != nil {
		t.Fatalf("Invalid query: %s", err)
	}
}

var C int

func BenchmarkQuery(b *testing.B) {
	b.StopTimer()
	z, err := NewZine()
	if err != nil {
		b.Fatalf("couldn't create zine: %s", err)
	}

	a := Article{
		Title: "frew",
		Tags:  []string{"foo", "bar"},
		Extra: map[string]string{"foo": "bar"},
	}
	for i := 0; i < 1000; i++ {
		if err := z.InsertArticle(a); err != nil {
			b.Fatalf("couldn't insert article: %s", err)
		}
	}

	b.StartTimer()
	var c int
	for i := 0; i < b.N; i++ {
		r, err := z.Q("SELECT * FROM _ LIMIT 5")
		if err != nil {
			b.Fatalf("Invalid query: %s", err)
		}
		c += len(r)
	}

	C = c
}

func TestRender(t *testing.T) {
	z, err := NewZine()
	if err != nil {
		t.Fatalf("couldn't create db: %s", err)
	}

	a := Article{
		Title:    "frew",
		Filename: "frew.md",
		URL:      "/frew/",
		Tags:     []string{"foo", "bar"},
		Extra:    map[string]string{"foo": "bar"},
	}
	for i := 0; i < 1000; i++ {
		if err := z.InsertArticle(a); err != nil {
			t.Fatalf("couldn't insert article: %s", err)
		}
	}
	got, err := z.Render(Article{Title: "x", Body: []byte(`hello! *{{ with $r := (q "SELECT COUNT(*) AS c FROM _")}}{{ index $r 0 "c" }}{{end}}*`)})
	if err != nil {
		t.Errorf("should not have gotten an error: %s", err)
		return
	}

	testutil.Equal(t, string(got), "<!-- header -->\n<p>hello! <em>2000</em></p>\n<!-- footer -->\n", "simple")
}

var S string

func BenchmarkRender(b *testing.B) {
	b.StopTimer()
	z, err := NewZine()
	if err != nil {
		b.Fatalf("couldn't create db: %s", err)
	}

	a := Article{
		Title: "frew",
		Tags:  []string{"foo", "bar"},
		Extra: map[string]string{"foo": "bar"},
	}
	for i := 0; i < 1000; i++ {
		if err := z.InsertArticle(a); err != nil {
			b.Fatalf("couldn't insert article: %s", err)
		}
	}

	var out []byte
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var err error
		out, err = z.Render(Article{Title: "X", Body: []byte(`hello! *{{ with $r := (q "SELECT COUNT(*) AS c FROM _")}}{{ index $r 0 "c" }}{{end}}*`)})
		if err != nil {
			b.Errorf("should not have gotten an error: %s", err)
			return
		}
	}

	S = string(out)
}

func BenchmarkLoadNilNil(b *testing.B) {
	var (
		z   *Zine
		err error
	)
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		z, err = NewZine()
		if err != nil {
			b.Fatalf("couldn't create zine: %s", err)
		}
		z.Root = "testdata"

		b.StartTimer()
		if err := z.Load(nil); err != nil {
			b.Fatalf("couldn't load: %s", err)
		}
	}
	b.StopTimer()

	S = z.Root
}

func BenchmarkLoadXY(b *testing.B) {
	var (
		z   *Zine
		err error
		c   int
	)
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		z, err = NewZine()
		if err != nil {
			b.Fatalf("couldn't create zine: %s", err)
		}
		z.Root = "testdata"

		b.StartTimer()
		var as []Article
		if err := z.Load(&as); err != nil {
			b.Fatalf("couldn't load: %s", err)
		}
		c += len(as)
	}
	b.StopTimer()

	S = z.Root
	C = c
}
