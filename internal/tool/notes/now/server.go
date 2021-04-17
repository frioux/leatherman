package now

import (
	"bytes"
	"errors"
	"fmt"
	corehtml "html"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/frioux/leatherman/internal/lmfs"
	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/frioux/leatherman/internal/notes"
	"github.com/frioux/leatherman/internal/selfupdate"
)

// 3. list and sort all files with `todo-` prefix
//   * list(db) ([]file, error)
//   * render

const prelude = `<!DOCTYPE html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1" /> 
<title>%s</title>
<link rel="icon" href="/favicon">
</head>
<a href="/list">list</a> | <a href="/sup">sup</a> | <a href="/">now</a>
<br><br>
`

func handlerAddItem(fss fs.FS, mdwn goldmark.Markdown, nowPath string) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		if err := req.ParseForm(); err != nil {
			return err
		}

		v := req.Form.Get("item")
		if v == "" {
			rw.WriteHeader(400)
			fmt.Fprint(rw, "missing item parameter")
			return nil
		}

		b, err := fs.ReadFile(fss, nowPath)
		if err != nil {
			return err
		}

		b, err = addItem(bytes.NewReader(b), time.Now(), v)
		if err != nil {
			return err
		}

		if err := lmfs.WriteFile(fss, nowPath, b, 0); err != nil {
			return err
		}

		b, err = parseNow(bytes.NewReader(b), time.Now())
		if err != nil {
			return err
		}

		rw.Header().Add("Location", "/")
		rw.WriteHeader(303)

		fmt.Fprintln(rw, prelude)
		return mdwn.Convert(b, rw)
	})
}

func handlerFavicon() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprintln(rw, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y=".9em" font-size="90">☕</text></svg>`)
	})
}

func handlerList(z *notes.Zine, fss fs.FS, mdwn goldmark.Markdown) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		stmt, err := z.Preparex(`SELECT title, url FROM articles ORDER BY title`)
		if err != nil {
			return err
		}

		articles := make([]struct{ Title, URL string }, 0, 1000)
		if err := stmt.Select(&articles); err != nil {
			return err
		}

		buf := &bytes.Buffer{}
		for _, e := range articles {
			fmt.Fprintln(buf, " * ["+e.Title+"]("+e.URL+")")
		}

		fmt.Fprintf(rw, prelude, "now: list")
		return mdwn.Convert(buf.Bytes(), rw)
	})
}

func handlerQ(z *notes.Zine, mdwn goldmark.Markdown) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		q := req.URL.Query().Get("q")
		if q == "" {
			q = "SELECT * FROM articles"
		}
		ret, err := z.Q(q)
		if err != nil {
			return err
		}

		buf := &bytes.Buffer{}
		fmt.Fprintf(buf, "```\n")
		for _, e := range ret {
			fmt.Fprintf(buf, "%v\n", e)
		}
		fmt.Fprintf(buf, "```\n")

		fmt.Fprintf(rw, prelude, "now: q")
		return mdwn.Convert(buf.Bytes(), rw)
	})
}

func handlerRoot(z *notes.Zine, fss fs.FS, mdwn goldmark.Markdown, nowPath string) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		if req.URL.Path == "/" {
			b, err := fs.ReadFile(fss, nowPath)
			if err != nil {
				return err
			}

			b, err = parseNow(bytes.NewReader(b), time.Now())
			if err != nil {
				return err
			}

			fmt.Fprintf(rw, prelude, "now")
			return mdwn.Convert(b, rw)
		}

		f := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/"), "/") + ".md"
		a, err := z.LoadArticle(z.DB, f)
		if err != nil {
			return fmt.Errorf("LoadArticle: %w", err)
		}

		fmt.Fprintf(rw, prelude, "now: "+a.Title)
		fmt.Fprintf(rw, `<br><a href="/update?file=%s">Update %s</a><br>`, f, f)

		b, err := z.Render(a)
		if err != nil {
			return fmt.Errorf("Render: %w", err)
		}
		buf := bytes.NewBuffer(b)
		_, err = io.Copy(rw, buf)
		return err
	})
}

func handlerToggle(fss fs.FS, mdwn goldmark.Markdown, nowPath string) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		if err := req.ParseForm(); err != nil {
			return err
		}

		v := req.Form.Get("chunk")
		if v == "" {
			rw.WriteHeader(400)
			fmt.Fprint(rw, "missing chunk parameter")
			return nil
		}

		b, err := fs.ReadFile(fss, nowPath)
		if err != nil {
			return err
		}

		b, err = toggleNow(bytes.NewReader(b), time.Now(), v)
		if err != nil {
			return err
		}

		if err := lmfs.WriteFile(fss, nowPath, b, 0); err != nil {
			return err
		}

		b, err = parseNow(bytes.NewReader(b), time.Now())
		if err != nil {
			return err
		}

		rw.Header().Add("Location", "/")
		rw.WriteHeader(303)

		fmt.Fprintln(rw, prelude)
		return mdwn.Convert(b, rw)
	})
}

func handlerUpdate(z *notes.Zine, fss fs.FS) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		switch req.Method {
		case "GET":
			f := req.URL.Query().Get("file")
			if f == "" {
				return errors.New("file parameter required")
			}

			b, err := fs.ReadFile(fss, f)
			if err != nil {
				return err
			}

			fmt.Fprintf(rw, prelude, "now: update "+f)
			const form = `
<form action="/update" method="post">
	<input type="hidden" name="file" value="%s" />
	<textarea rows="50" cols="80" name="value">%s</textarea>
	<button>Save</button>
</form>
			`
			fmt.Fprintf(rw, form, f, corehtml.EscapeString(string(b)))
			return nil
		case "POST":
			f := req.FormValue("file")
			if f == "" {
				return errors.New("file parameter required")
			}

			b := req.FormValue("value")
			if b == "" {
				return errors.New("value parameter required")
			}

			b = strings.ReplaceAll(b, "\r", "") // unix files only!
			if err := lmfs.WriteFile(fss, f, []byte(b), 0); err != nil {
				return err
			}
			rw.Header().Add("Location", "/"+strings.TrimSuffix(f, ".md"))
			rw.WriteHeader(303)
			fmt.Fprint(rw, "Successfully updated")
			return nil
		}
		return errors.New("invalid method for /update")
	})
}

func server(fss fs.FS, z *notes.Zine, generation *chan bool) (http.Handler, error) {
	mux := http.NewServeMux()

	mdwn := goldmark.New(
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
	)

	const nowPath = "now.md"

	mux.Handle("/favicon", handlerFavicon())

	mux.Handle("/version/", selfupdate.Handler)

	mux.Handle("/", handlerRoot(z, fss, mdwn, nowPath))

	mux.Handle("/list", handlerList(z, fss, mdwn))

	mux.Handle("/q", handlerQ(z, mdwn))

	mux.Handle("/sup", handlerSup())

	mux.Handle("/update", handlerUpdate(z, fss))

	mux.Handle("/toggle", handlerToggle(fss, mdwn, nowPath))

	mux.Handle("/add-item", handlerAddItem(fss, mdwn, nowPath))

	return autoReload(mux, generation), nil
}
