package notes

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/frioux/leatherman/internal/lmhttp"
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

	changed := make(chan struct{})
	go longpoll(db, dir, changed)

	mux.Handle("/favicon", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprintln(rw, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y=".9em" font-size="90">☕</text></svg>`)
	}))

	mux.Handle("/version", selfupdate.Handler)

	mux.Handle("/", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
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

		fmt.Fprintf(rw, prelude, "now: list")
		return mdwn.Convert(buf.Bytes(), rw)
	}))

	mux.Handle("/sup", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
		defer cancel()

		wg := &sync.WaitGroup{}
		wg.Add(2)

		var rpi struct{ Game string }
		go func() {
			defer wg.Done()
			resp, err := lmhttp.Get(ctx, "http://retropie:8081/retropie")
			if err != nil {
				// whyyyyy
				rpi.Game = err.Error()
				return
			}
			defer resp.Body.Close()

			dec := json.NewDecoder(resp.Body)
			if err := dec.Decode(&rpi); err != nil {
				// ugh wtf
				rpi.Game = err.Error()
			}
		}()

		var steamos []byte
		go func() {
			defer wg.Done()
			resp, err := lmhttp.Get(ctx, "http://steamos:8081/steambox")
			if err != nil {
				// I don't like it
				steamos = []byte(err.Error())
				return
			}
			defer resp.Body.Close()

			steamos, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				// I should have thought this through more carefully
				steamos = []byte(err.Error())
			}
		}()

		wg.Wait()

		fmt.Fprintf(rw, prelude, "now: sup")
		fmt.Fprintf(rw, "retropie: %s<br>steamos: %s", rpi.Game, steamos)

		return nil
	}))

	mux.Handle("/render", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		f := req.URL.Query().Get("file")
		if f == "" {
			rw.Header().Add("Location", "/list")
			rw.WriteHeader(303)
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

		fmt.Fprintf(rw, prelude, "now: "+a.Title)
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
