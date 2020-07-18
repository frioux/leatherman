package notes

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

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

	const nowPath = "/notes/content/posts/now.md"

	mux.Handle("/", handlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		r, err := db.Download(nowPath)
		if err != nil {
			return err
		}

		b, err := parseNow(r, time.Now())
		if err != nil {
			return err
		}

		return mdwn.Convert(b, rw)
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

		return mdwn.Convert(b, rw)
	}))

	return mux, nil
}

// parseNow reads markdown and returns html.  The main difference from normal
// markdown is that a section titled ## 2020-02-02 ## will get special
// rendering treatment.
//
// A json header is discarded for now.
func parseNow(r io.Reader, when time.Time) ([]byte, error) {
	desiredHeader := "## " + when.Format("2006-01-02") + " ##"
	ret := &strings.Builder{}

	var inToday bool
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		switch {
		case !inToday && line == desiredHeader:
			inToday = true
		case inToday && strings.HasPrefix(line, "## "):
			inToday = false
		case inToday && strings.HasPrefix(line, " * "):
			md := md5.Sum([]byte(line))
			linkable := hex.EncodeToString(md[:])
			line += ` <form action="/toggle" method="POST"><input type="hidden" name="chunk" value="` + linkable + `"><button>Toggle</button></form>`
		}

		ret.WriteString(line)
		ret.WriteRune('\n')
	}

	return []byte(ret.String()), nil
}

// toggleNow will mark a list item done (surround with ~~'s) if it's in the
// section for when and it's md5sum matches sum.  If the item has already been
// done, this function will mark it undone.
func toggleNow(r io.Reader, when time.Time, sum string) ([]byte, error) {
	desiredHeader := "## " + when.Format("2006-01-02") + " ##"
	ret := &strings.Builder{}

	var inToday bool
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		switch {
		case !inToday && line == desiredHeader:
			inToday = true
		case inToday && strings.HasPrefix(line, "## "):
			inToday = false
		case inToday && strings.HasPrefix(line, " * "):
			md := md5.Sum([]byte(line))
			linkable := hex.EncodeToString(md[:])
			if sum == linkable {
				if strings.HasPrefix(line, " * ~~") && strings.HasSuffix(line, "~~") { // already done, undo
					line = " * " + line[5:len(line)-2]
				} else { // not done, mark done
					line = " * ~~" + line[3:] + "~~"
				}
			}
		}

		ret.WriteString(line)
		ret.WriteRune('\n')
	}

	return []byte(ret.String()), nil
}
