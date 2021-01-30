package now

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/frioux/leatherman/internal/notes"
)

func syncEventsToDB(cl dropbox.Client, z *notes.Zine, events []dropbox.Metadata) (err error) {
	wg := &sync.WaitGroup{}
	wg.Add(len(events))

	articles := make([]struct {
		notes.Article
		deleted bool
	}, len(events))

	for i, e := range events {
		i := i
		e := e

		go func() {
			defer wg.Done()

			defer func() {
				articles[i].Filename = e.Name
				articles[i].URL = strings.TrimSuffix(e.Name, ".md")
			}()

			if e.Tag == "deleted" {
				articles[i].deleted = true
				return
			}

			r, err := cl.Download(e.PathLower)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			articles[i].Article, err = notes.ReadArticle(r)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}()
	}

	wg.Wait()

	for _, a := range articles {
		if a.deleted {
			fmt.Fprintln(os.Stderr, "deleting", a.Filename, "...")
			if err := z.DeleteArticle(a.Article.Filename); err != nil {
				return err
			}
		} else {
			fmt.Fprintln(os.Stderr, "replacing", a.Filename, "...")
			if err := z.ReplaceArticle(a.Article); err != nil {
				return err
			}
		}
	}

	return nil
}

func maintainDB(cl dropbox.Client, dir string, generation *chan bool, z *notes.Zine) {
	watcher := make(chan []dropbox.Metadata)
	go func() { cl.Longpoll(context.Background(), dir, watcher) }()
	go func() {
		for events := range watcher {
			if err := syncEventsToDB(cl, z, events); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			close(*generation)
			*generation = make(chan bool)
		}
	}()
}

func loadDB(cl dropbox.Client, dir string, generation *chan bool) (z *notes.Zine, err error) {
	z, err = notes.NewZine()
	if err != nil {
		return nil, err
	}

	t0 := time.Now()

	var r dropbox.ListFolderResult
	r, err = cl.ListFolder(dropbox.ListFolderParams{Path: dir})
	if err != nil {
		return nil, err
	}

	entries := r.Entries

	for r.HasMore {
		r, err = cl.ListFolderContinue(r.Cursor)
		if err != nil {
			return nil, err
		}

		entries = append(entries, r.Entries...)
	}

	articles := make([]notes.Article, len(entries))
	wg := &sync.WaitGroup{}
	for i, e := range entries {
		wg.Add(1)

		name := e.Name
		i := i
		go func() {
			defer wg.Done()

			// unclear what to do about errors here
			r, err := cl.Download(dir + name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			articles[i], err = notes.ReadArticle(r)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			articles[i].Filename = name
			articles[i].URL = "/" + strings.TrimSuffix(name, ".md")
		}()
	}
	wg.Wait()

	for _, a := range articles {
		if err := z.InsertArticle(a); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	fmt.Fprintf(os.Stderr, "db loaded in %s\n", time.Now().Sub(t0))

	maintainDB(cl, dir, generation, z)
	return z, nil
}
