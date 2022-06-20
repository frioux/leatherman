package now

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	lua "github.com/yuin/gopher-lua"

	"github.com/frioux/leatherman/internal/lmfav"
	"github.com/frioux/leatherman/internal/lmfs"
	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/frioux/leatherman/internal/lmlua"
	"github.com/frioux/leatherman/internal/lmlua/luanotes"
	"github.com/frioux/leatherman/internal/notes"
	"github.com/frioux/leatherman/internal/selfupdate"
)

// 3. list and sort all files with `todo-` prefix
//   * list(db) ([]file, error)
//   * render

func handlerAddItem(z *notes.Zine, fss fs.FS, mdwn goldmark.Markdown, nowPath string) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		if err := req.ParseForm(); err != nil {
			return err
		}

		i := req.Form.Get("item")
		if i == "" {
			rw.WriteHeader(400)
			fmt.Fprint(rw, "missing item parameter")
			return nil
		}

		b, err := fs.ReadFile(fss, nowPath)
		if err != nil {
			return err
		}

		b, err = addItem(bytes.NewReader(b), time.Now(), i)
		if err != nil {
			return err
		}

		if err := lmfs.WriteFile(fss, nowPath, b, 0); err != nil {
			return err
		}

		rw.Header().Add("Location", "/")
		rw.WriteHeader(303)

		b, err = parseNow(bytes.NewReader(b), time.Now())
		if err != nil {
			return fmt.Errorf("parseNow (file: %s), %w", nowPath, err)
		}

		v := &HTMLVars{Zine: z, Title: "now"}
		if err := mdwn.Convert(b, v); err != nil {
			return err
		}

		return tpl.ExecuteTemplate(rw, "simple.html", v)
	})
}

func handlerList(z *notes.Zine, fss fs.FS, mdwn goldmark.Markdown) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		rw.Header().Set("Cache-Control", "no-cache")
		stmt, err := z.Preparex(`SELECT title, url FROM articles ORDER BY title`)
		if err != nil {
			return err
		}

		v := listVars{HTMLVars: &HTMLVars{Zine: z}}
		if err := stmt.Select(&v.Articles); err != nil {
			return err
		}

		return tpl.ExecuteTemplate(rw, "list.html", v)
	})
}

func handlerQ(z *notes.Zine, mdwn goldmark.Markdown) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		rw.Header().Set("Cache-Control", "no-cache")
		v := qVars{HTMLVars: &HTMLVars{Zine: z}}
		q := req.URL.Query().Get("q")
		if q == "" {
			q = "SELECT * FROM articles"
		}
		var err error
		v.Records, err = z.Q(q)
		if err != nil {
			return err
		}

		return tpl.ExecuteTemplate(rw, "q.html", v)
	})
}

func handlerRoot(z *notes.Zine, fss fs.FS, mdwn goldmark.Markdown, nowPath string) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		rw.Header().Set("Cache-Control", "no-cache")
		if req.URL.Path == "/" {
			b, err := fs.ReadFile(fss, nowPath)
			if err != nil {
				return err
			}

			b, err = parseNow(bytes.NewReader(b), time.Now())
			if err != nil {
				return err
			}

			v := &HTMLVars{Zine: z, Title: "now"}
			if err := mdwn.Convert(b, v); err != nil {
				return err
			}

			if err := tpl.ExecuteTemplate(rw, "simple.html", v); err != nil {
				return err
			}

			return nil
		}

		f := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/"), "/") + ".md"
		a, err := z.LoadArticle(z.DB, f)
		if err != nil {
			return fmt.Errorf("LoadArticle (%s): %w", f, err)
		}

		if l := req.URL.Query().Get("lua"); l != "" {
			L := lua.NewState()
			defer L.Close()

			lmlua.RegisterFSType(L)
			lmlua.RegisterGoqueryPackage(L)
			lmlua.RegisterHTTPPackage(L)
			lmlua.RegisterRegexpPackage(L)
			luanotes.RegisterNotesPackage(L)

			//  * article
			//    * :bytes()
			//    * :title()
			//    * :filename()
			//    * :url()
			//    * :raw()
			//    * :tags()
			//    * :reviewed_on()
			//    * :review_by()
			//    * :extra()
			//    * :body()
			//    * :markdownlua()
			//  * fs
			//    * :open()
			//    * :create()
			//    * :writefile()

			udRW := L.NewUserData()
			udRW.Value = rw
			L.SetMetatable(udRW, L.GetTypeMetatable("responsewriter"))
			L.SetGlobal("rw", udRW)

			udReq := L.NewUserData()
			udReq.Value = req
			L.SetMetatable(udReq, L.GetTypeMetatable("request"))
			L.SetGlobal("req", udReq)

			L.SetGlobal("f", lua.LString(f))

			udFS := L.NewUserData()
			udFS.Value = fss
			L.SetMetatable(udFS, L.GetTypeMetatable("fs"))
			L.SetGlobal("fss", udFS)

			a.MarkdownLua = append(a.MarkdownLua, []byte("\n"+l+"(rw, req)\n")...)
			if err := L.DoString(string(a.MarkdownLua)); err != nil {
				return fmt.Errorf("couldn't load lua: %w\n%s", err, a.MarkdownLua)
			}
		}

		v := articleVars{
			HTMLVars:     &HTMLVars{Zine: z},
			ArticleTitle: a.Title,
			Filename:     f,
		}

		b, err := z.Render(a)
		if err != nil {
			return fmt.Errorf("Render: %w", err)
		}
		buf := bytes.NewBuffer(b)
		if _, err := io.Copy(v, buf); err != nil {
			return err
		}
		if err := tpl.ExecuteTemplate(rw, "article.html", v); err != nil {
			return err
		}
		return nil
	})
}

func handlerToggle(z *notes.Zine, fss fs.FS, mdwn goldmark.Markdown, nowPath string) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		rw.Header().Set("Cache-Control", "no-cache")
		if err := req.ParseForm(); err != nil {
			return err
		}

		c := req.Form.Get("chunk")
		if c == "" {
			rw.WriteHeader(400)
			fmt.Fprint(rw, "missing chunk parameter")
			return nil
		}

		b, err := fs.ReadFile(fss, nowPath)
		if err != nil {
			return err
		}

		b, err = toggleNow(bytes.NewReader(b), time.Now(), c)
		if err != nil {
			return err
		}

		if err := lmfs.WriteFile(fss, nowPath, b, 0); err != nil {
			return err
		}

		rw.Header().Add("Location", "/")
		rw.WriteHeader(303)

		b, err = parseNow(bytes.NewReader(b), time.Now())
		if err != nil {
			return err
		}

		v := &HTMLVars{Zine: z, Title: "now"}
		if err := mdwn.Convert(b, v); err != nil {
			return err
		}

		return tpl.ExecuteTemplate(rw, "simple.html", v)

	})
}

func handlerUpdate(z *notes.Zine, fss fs.FS) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		rw.Header().Set("Cache-Control", "no-cache")
		switch {
		case req.Method == "GET":
			f := req.URL.Query().Get("file")
			if f == "" {
				return errors.New("file parameter required")
			}

			b, err := fs.ReadFile(fss, f)
			if err != nil {
				return err
			}

			v := updateVars{
				HTMLVars: &HTMLVars{Zine: z},
				File:     f,
				Content:  string(b),
			}

			return tpl.ExecuteTemplate(rw, "update.html", v)
		case req.Method == "POST" && req.FormValue("delete") != "1":
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
		case req.Method == "POST" && req.FormValue("delete") == "1":
			f := req.FormValue("file")
			if f == "" {
				return errors.New("file parameter required")
			}

			if err := lmfs.Remove(fss, f); err != nil {
				return err
			}

			rw.Header().Add("Location", "/")
			rw.WriteHeader(303)
			fmt.Fprint(rw, "Successfully deleted")
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

	mux.Handle("/favicon", lmfav.Emoji('☕'))

	mux.Handle("/version/", selfupdate.Handler)

	mux.Handle("/", handlerRoot(z, fss, mdwn, nowPath))

	mux.Handle("/list", handlerList(z, fss, mdwn))

	mux.Handle("/q", handlerQ(z, mdwn))

	mux.Handle("/sup", handlerSup(z))

	mux.Handle("/update", handlerUpdate(z, fss))

	mux.Handle("/toggle", handlerToggle(z, fss, mdwn, nowPath))

	mux.Handle("/add-item", handlerAddItem(z, fss, mdwn, nowPath))

	return autoReload(mux, generation), nil
}
