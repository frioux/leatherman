package notes

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/frioux/leatherman/internal/dropbox"
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
</head>
<a href="/list">list</a> | <a href="/">now</a><br><br>
`

/*
Serve provides a web interface to parts of my notes stored in
Dropbox.  This should eventually be merged into Amygdala.

Command: notes
*/
func Serve(args []string, _ io.Reader) error {
	var (
		listen string
	)
	fs := flag.NewFlagSet("notes", flag.ContinueOnError)
	fs.StringVar(&listen, "listen", ":0", "location to listen on; default is random")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	h, err := server()
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	s := http.Server{Handler: h}
	return s.Serve(listener)
}

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

	mux.Handle("/", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		r, err := db.Download(nowPath)
		if err != nil {
			return err
		}

		b, err := parseNow(r, time.Now())
		if err != nil {
			return err
		}

		fmt.Fprintln(rw, prelude)
		return mdwn.Convert(b, rw)
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
			fmt.Fprintln(buf, " * ["+e.Name+"](/render?file="+e.Name+")")
		}

		fmt.Fprintln(rw, prelude)
		return mdwn.Convert(buf.Bytes(), rw)
	}))

	mux.Handle("/render", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		f := req.URL.Query().Get("file")
		if f == "" {
			rw.WriteHeader(302)
			rw.Header().Add("Location", "/list")
			fmt.Fprint(rw, "No file param, going to /list")
			return nil
		}

		r, err := db.Download(dir + f)
		if err != nil {
			return err
		}

		a, err := readArticle(r)
		if err != nil {
			return fmt.Errorf("readArticle: %w", err)
		}

		fmt.Fprintln(rw, prelude)
		return mdwn.Convert(a.Body, rw)
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

		rw.WriteHeader(302)
		rw.Header().Add("Location", "/now")

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

		rw.WriteHeader(302)
		rw.Header().Add("Location", "/now")

		fmt.Fprintln(rw, prelude)
		return mdwn.Convert(b, rw)
	}))

	return mux, nil
}
