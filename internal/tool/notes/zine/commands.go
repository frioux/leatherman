package zine

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/frioux/leatherman/internal/notes"
)

var commands = map[string]func([]string) error{
	"render": render,
	"q":      q,
	"debug":  debug,
}

// q runs a query against the corpus.
func q(args []string) error {
	var (
		root, sql, tpl string
		showschema     bool
	)

	flags := flag.NewFlagSet("q", flag.ContinueOnError)
	flags.StringVar(&root, "root", "./content", "root input directory")
	flags.StringVar(&sql, "sql", "SELECT * FROM _", "sql to run")
	flags.StringVar(&tpl, "tpl", `{{join . "\t"}}`, "template to run")
	flags.BoolVar(&showschema, "showschema", false, "show sql schema")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if showschema {
		fmt.Print(notes.Schema)
		return nil
	}

	t := template.New("x")
	t.Funcs(template.FuncMap{
		"join": func(is map[string]interface{}, sep string) string {
			s := make([]string, 0, len(is))
			for i := range is {
				s = append(s, fmt.Sprint(is[i]))
			}
			return strings.Join(s, sep)
		},
	})

	t, err := t.Parse(tpl)
	if err != nil {
		return err
	}

	z, err := notes.NewZine("")
	if err != nil {
		return err
	}
	z.FS = os.DirFS(root)
	if err := z.Load(nil); err != nil {
		return err
	}

	ret, err := z.Q(sql, flags.Args()...)
	if err != nil {
		return err
	}

	for _, out := range ret {
		if err := t.Execute(os.Stdout, out); err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}

// render will convert the corpus to html.
func render(args []string) error {
	var root, out, static, publicPrefix string

	flags := flag.NewFlagSet("render", flag.ContinueOnError)
	flags.StringVar(&root, "root", "./content", "root input directory")
	flags.StringVar(&out, "out", "./public", "directory to render output to")
	flags.StringVar(&static, "static", "./static", "directory to prepopulate out with")
	flags.StringVar(&publicPrefix, "public-prefix", "notes/", "public prefix to render urls with")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	z, err := notes.NewZine("")
	if err != nil {
		return err
	}
	z.PublicPrefix = publicPrefix
	z.FS = os.DirFS(root)

	metas := []notes.Article{}
	if err := z.Load(&metas); err != nil {
		return err
	}

	if err := os.RemoveAll(out); err != nil {
		return err
	}

	if err := filepath.WalkDir(static, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(static, path)
		if err != nil {
			return err
		}
		dir := filepath.Dir(rel)
		if err := os.MkdirAll(filepath.Join(out, dir), 0755); err != nil {
			return err
		}

		from, err := os.Open(path)
		if err != nil {
			return err
		}
		defer from.Close()

		to, err := os.Create(filepath.Join(out, rel))
		if err != nil {
			return err
		}
		defer to.Close()

		if _, err := io.Copy(to, from); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	for i := range metas {
		// don't be running on windows
		dir := filepath.Join(out, strings.TrimPrefix(metas[i].URL, z.PublicPrefix))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("couldn't create dir for %s: %w", metas[i].Filename, err)
		}

		b, err := z.Render(metas[i])
		if err != nil {
			return fmt.Errorf("couldn't render %s: %w", metas[i].Filename, err)
		}

		f, err := os.Create(filepath.Join(dir, "index.html"))
		if err != nil {
			return fmt.Errorf("couldn't create %s: %w", filepath.Join(dir, "index.html"), err)
		}

		if _, err := io.Copy(f, bytes.NewReader(b)); err != nil {
			return fmt.Errorf("couldn't copy: %w", err)
		}
	}

	return nil
}

func debug(args []string) error {
	var root, file string

	flags := flag.NewFlagSet("debug", flag.ContinueOnError)
	flags.StringVar(&root, "root", "./content", "root input directory")
	flags.StringVar(&file, "file", "", "file to render to markdown")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if file == "" {
		return errors.New("-file argument is required")
	}

	fileMatcher, err := regexp.Compile(file)
	if err != nil {
		return err
	}

	z, err := notes.NewZine("")
	if err != nil {
		return err
	}
	z.FS = os.DirFS(root)

	metas := []notes.Article{}
	if err := z.Load(&metas); err != nil {
		return err
	}

	for i := range metas {
		if !fileMatcher.MatchString(metas[i].Filename) {
			continue
		}

		b, err := z.RenderToMarkdown(metas[i])
		if err != nil {
			return fmt.Errorf("couldn't render %s: %w", metas[i].Filename, err)
		}

		if _, err := io.Copy(os.Stdout, bytes.NewReader(b)); err != nil {
			return fmt.Errorf("couldn't copy: %w", err)
		}
	}

	return nil
}
