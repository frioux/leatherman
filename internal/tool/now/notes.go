package now

import (
	"bytes"
	"errors"
	"fmt"
	corehtml "html"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/frioux/leatherman/internal/notes"
	"github.com/frioux/leatherman/internal/selfupdate"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
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

type handlerFunc func(http.ResponseWriter, *http.Request) error

func (f handlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := f(rw, r); err != nil {
		rw.WriteHeader(500)
		fmt.Fprintln(os.Stderr, err)
	}
}

func server() (http.Handler, error) {
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
	db, err := dropbox.NewClient(dropbox.Client{
		Token: os.Getenv("LM_DROPBOX_TOKEN"),
	})
	if err != nil {
		return nil, err
	}

	const (
		dir     = "/notes/content/posts/"
		nowPath = dir + "now.md"
	)

	mux.Handle("/favicon", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprintln(rw, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y=".9em" font-size="90">☕</text></svg>`)
	}))

	mux.Handle("/version", selfupdate.Handler)

	mux.Handle("/", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		if req.URL.Path == "/" {
			r, err := db.Download(nowPath)
			if err != nil {
				return err
			}

			b, err := parseNow(r, time.Now())
			if err != nil {
				return err
			}

			fmt.Fprintf(rw, prelude, "now")
			return mdwn.Convert(b, rw)
		}

		f := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/"), "/") + ".md"
		r, err := db.Download(dir + f)
		if err != nil {
			return err
		}

		a, err := notes.ReadArticle(r)
		if err != nil {
			return fmt.Errorf("ReadArticle: %w", err)
		}

		fmt.Fprintf(rw, prelude, "now: "+a.Title)
		fmt.Fprintf(rw, `<br><a href="/update?file=%s">Update %s</a><br>`, f, f)
		return mdwn.Convert(a.Body, rw)
	}))

	mux.Handle("/list", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		r, err := db.ListFolder(dropbox.ListFolderParams{Path: dir})
		if err != nil {
			return err
		}

		entries := r.Entries

		for r.HasMore {
			r, err = db.ListFolderContinue(r.Cursor)
			if err != nil {
				return err
			}

			entries = append(entries, r.Entries...)
		}

		sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })

		buf := &bytes.Buffer{}
		for _, e := range entries {
			fmt.Fprintln(buf, " * ["+e.Name+"](/"+strings.TrimSuffix(e.Name, ".md")+")")
		}

		fmt.Fprintf(rw, prelude, "now: list")
		return mdwn.Convert(buf.Bytes(), rw)
	}))

	mux.Handle("/sup", handlerFunc(sup))

	mux.Handle("/update", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		switch req.Method {
		case "GET":
			f := req.URL.Query().Get("file")
			if f == "" {
				return errors.New("file parameter required")
			}

			r, err := db.Download(dir + f)
			if err != nil {
				return err
			}

			b, err := ioutil.ReadAll(r)
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
			err = db.Create(dropbox.UploadParams{
				Path: dir + f,
				Mode: "overwrite",
			}, strings.NewReader(b))
			if err != nil {
				return err
			}
			rw.Header().Add("Location", "/"+strings.TrimSuffix(f, ".md"))
			rw.WriteHeader(303)
			fmt.Fprint(rw, "Successfully updated")
			return nil
		default:
			return errors.New("invalid method for /update")
		}
		return nil
	}))

	mux.Handle("/toggle", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		if err := req.ParseForm(); err != nil {
			return err
		}

		v := req.Form.Get("chunk")
		if v == "" {
			rw.WriteHeader(400)
			fmt.Fprint(rw, "missing chunk parameter")
			return nil
		}

		r, err := db.Download(nowPath)
		if err != nil {
			return err
		}

		b, err := toggleNow(r, time.Now(), v)
		if err != nil {
			return err
		}

		if err := db.Create(dropbox.UploadParams{
			Path: nowPath,
			Mode: "overwrite",
		}, bytes.NewReader(b)); err != nil {
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
	}))

	mux.Handle("/add-item", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		if err := req.ParseForm(); err != nil {
			return err
		}

		v := req.Form.Get("item")
		if v == "" {
			rw.WriteHeader(400)
			fmt.Fprint(rw, "missing item parameter")
			return nil
		}

		r, err := db.Download(nowPath)
		if err != nil {
			return err
		}

		b, err := addItem(r, time.Now(), v)
		if err != nil {
			return err
		}

		if err := db.Create(dropbox.UploadParams{
			Path: nowPath,
			Mode: "overwrite",
		}, bytes.NewReader(b)); err != nil {
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
	}))

	arMux, err := autoReload(db, mux, dir)
	if err != nil {
		return nil, err
	}

	return arMux, nil
}
