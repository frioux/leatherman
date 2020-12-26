package zine

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tailscale/hujson"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func run() error {
	command := "render"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	cmd, ok := commands[command]
	if !ok {
		return fmt.Errorf("unknown command «%s»; valid commands are 'render' and 'q'\n", command)
	}

	if err := cmd(os.Args[1:]); err != nil {
		return err
	}

	return nil
}

type article struct {
	Title string

	// Filename will be set after parsing.
	Filename string `json:"-"`

	// URL will be set after parsing.
	URL string `json:"-"`

	// Raw tells the parser not to include the standard header and footer.
	Raw bool

	Tags []string

	ReviewedOn *string `json:"reviewed_on"`

	ReviewBy *string `json:"review_by"`

	Extra map[string]string

	Body []byte
}

type zine struct {
	root string

	db
	tpl  template.Template
	mdwn goldmark.Markdown

	q func(string, ...string) ([]map[string]interface{}, error)
}

func newZine() (*zine, error) {
	d, err := newDB()
	if err != nil {
		return nil, fmt.Errorf("newDB: %s", err)
	}

	z := &zine{
		db:  *d,
		tpl: *template.New(""),
		mdwn: goldmark.New(
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
			goldmark.WithExtensions(
				extension.Strikethrough,
			),
		),
		q: func(q string, more ...string) ([]map[string]interface{}, error) {
			stmt, err := d.prepareCached(q)
			if err != nil {
				return nil, err
			}

			is := make([]interface{}, len(more))
			for i := range is {
				is[i] = more[i]
			}
			rows, err := stmt.Queryx(is...)
			if err != nil {
				return nil, err
			}
			cols, err := rows.Columns()
			if err != nil {
				return nil, err
			}

			ret := []map[string]interface{}{}
			for rows.Next() {
				m := make(map[string]interface{}, len(cols))
				if err := rows.MapScan(m); err != nil {
					return nil, err
				}
				ret = append(ret, m)
			}

			return ret, nil
		},
	}
	z.tpl.Parse(`
{{ define "header" }}
start
{{ end }}

{{ define "footer" }}
end
{{ end }}
`)
	z.tpl.Funcs(template.FuncMap{
		"q": z.q,
		"bible": func(s string) (string, error) {
			return s, nil
		},
	})

	return z, nil
}

func (z *zine) load(as *[]article) error {
	var files []string

	// parse index first so it can override header and footer
	if _, err := z.tpl.ParseFiles(z.root + "/index.tmpl"); err != nil {
		return err
	}

	if err := filepath.Walk(z.root, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, "/index.tmpl") {
			files = append(files, path)
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return err
	}

	var a article
	for _, f := range files {
		var err error

		a, err = readArticleFromFile(f)
		if err != nil {
			return fmt.Errorf("error parsing %s: %w", f, err)
		}
		a.Filename = f

		a.URL, err = filepath.Rel(z.root, f)
		if err != nil {
			return fmt.Errorf("error getting relname for %s: %w", f, err)
		}
		if a.Filename == "index.tmpl" || strings.HasSuffix(a.Filename, "index.tmpl") {
			a.URL = strings.TrimSuffix(a.URL, "index.tmpl")
		} else {
			a.URL = strings.TrimSuffix(a.URL, ".md")
		}

		if err := z.insertArticle(a); err != nil {
			return fmt.Errorf("error inserting data from %s: %w", f, err)
		}

		if as != nil {
			*as = append(*as, a)
		}
	}

	return nil
}

func (z *zine) renderToMarkdown(a article) ([]byte, error) {
	// XXX this may be expensive, but fixes the new error introduced here:
	// https://github.com/golang/go/commit/604146ce8961d32f410949015fc8ee31f9052209
	t, err := z.tpl.Clone()
	if err != nil {
		return nil, err
	}
	t = t.New("x")

	str := string(a.Body)
	if !a.Raw {
		str = "{{ template \"header\" . }}\n" + str + "\n{{ template \"footer\" . }}"
	}
	t, err = t.Parse(str)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, a); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (z *zine) render(a article) ([]byte, error) {
	mdwn, err := z.renderToMarkdown(a)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	if err := z.mdwn.Convert(mdwn, out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func readArticleFromFile(f string) (article, error) {
	r, err := os.Open(f)
	if err != nil {
		return article{}, err
	}
	defer r.Close()

	return readArticle(r)
}

func readArticle(r io.Reader) (article, error) {
	var a article
	d := hujson.NewDecoder(r)
	err := d.Decode(&a)
	if err != nil {
		return a, fmt.Errorf("hujson.Decoder.Decode: %w", err)
	}
	a.Body, err = ioutil.ReadAll(d.Buffered())
	if err != nil {
		return a, fmt.Errorf("hujson.Decoder.Buffered+ioutil.ReadAll: %w", err)
	}

	c, err := ioutil.ReadAll(r)
	if err != nil {
		return a, err
	}

	a.Body = append(a.Body, c...)

	return a, err
}