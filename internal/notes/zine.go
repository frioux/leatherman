package notes

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Zine struct {
	Root string

	PublicPrefix string

	DB
	tpl  template.Template
	mdwn goldmark.Markdown

	Q func(string, ...string) ([]map[string]interface{}, error)
}

func NewZine() (*Zine, error) {
	d, err := NewDB()
	if err != nil {
		return nil, fmt.Errorf("newDB: %s", err)
	}

	z := &Zine{
		DB:  *d,
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
				extension.Table,
			),
		),
		Q: func(q string, more ...string) ([]map[string]interface{}, error) {
			stmt, err := d.Preparex(q)
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
<!-- header -->
{{ end }}

{{ define "footer" }}
<!-- footer -->
{{ end }}
`)
	z.tpl.Funcs(template.FuncMap{
		"q": z.Q,
		"bible": func(s string) (string, error) {
			return s, nil
		},
	})

	return z, nil
}

func (z *Zine) Load(as *[]Article) error {
	var files []string

	// parse index first so it can override header and footer
	if _, err := z.tpl.ParseFiles(z.Root + "/index.tmpl"); err != nil {
		return err
	}

	if err := filepath.Walk(z.Root, func(path string, _ os.FileInfo, err error) error {
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

	var a Article
	for _, f := range files {
		var err error

		a, err = readArticleFromFile(f)
		if err != nil {
			return fmt.Errorf("error parsing %s: %w", f, err)
		}
		a.Filename = f

		a.URL, err = filepath.Rel(z.Root, f)
		if err != nil {
			return fmt.Errorf("error getting relname for %s: %w", f, err)
		}
		a.URL = z.PublicPrefix + a.URL
		if a.Filename == "index.tmpl" || strings.HasSuffix(a.Filename, "index.tmpl") {
			a.URL = strings.TrimSuffix(a.URL, "index.tmpl")
		} else {
			a.URL = strings.TrimSuffix(a.URL, ".md")
		}

		if err := z.InsertArticle(a); err != nil {
			return fmt.Errorf("error inserting data from %s: %w", f, err)
		}

		if as != nil {
			*as = append(*as, a)
		}
	}

	return nil
}

func (z *Zine) RenderToMarkdown(a Article) ([]byte, error) {
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

func (z *Zine) Render(a Article) ([]byte, error) {
	mdwn, err := z.RenderToMarkdown(a)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	if err := z.mdwn.Convert(mdwn, out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func readArticleFromFile(f string) (Article, error) {
	r, err := os.Open(f)
	if err != nil {
		return Article{}, err
	}
	defer r.Close()

	return ReadArticle(r)
}
