package now

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jmoiron/sqlx"

	"github.com/frioux/leatherman/internal/lmfs"
	"github.com/frioux/leatherman/internal/notes"
)

func syncEventsToDB(fss fs.FS, z *notes.Zine, events []fsnotify.Event) (err error) {
	tx, err := z.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback() // how to handle error?
		} else {
			err = tx.Commit()
		}
	}()

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

			if e.Op == fsnotify.Remove {
				articles[i].deleted = true
				return
			}

			b, err := fs.ReadFile(fss, e.Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dropbox.Download: %s\n", err)
				return
			}

			articles[i].Article, err = notes.ReadArticle(bytes.NewReader(b))
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
			if err := z.DeleteArticle(tx, a.Article.Filename); err != nil {
				return err
			}
		} else {
			fmt.Fprintln(os.Stderr, "replacing", a.Filename, "...")
			if err := z.ReplaceArticle(tx, a.Article); err != nil {
				return err
			}
		}
	}

	return nil
}

func maintainDB(fss fs.FS, generation *chan bool, z *notes.Zine) {
	wfs, ok := fss.(lmfs.WatchFS)
	if !ok {
		fmt.Fprintf(os.Stderr, "cannot maintainDB on an FS that can't watch (fs is a %T)\n", fss)
		return
	}
	watcher, err := wfs.Watch(context.Background(), ".")
	if err != nil {
		panic(err)
	}
	go func() {
		for events := range watcher {
			if err := syncEventsToDB(fss, z, events); err != nil {
				fmt.Fprintf(os.Stderr, "syncEventsToDB: %s\n", err)
			}
			close(*generation)
			*generation = make(chan bool)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Hour)
			rebuildDB(fss, z)
		}
	}()
}

func rebuildDB(f fs.FS, z *notes.Zine) (err error) {
	tx, err := z.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to rebuild database, rolling back: %s\n", err)
			tx.Rollback() // how to handle error?
		} else {
			err = tx.Commit()
		}
	}()

	if err := clearDB(z, tx); err != nil {
		return err
	}

	if err := populateDB(f, z, tx); err != nil {
		return err
	}

	return nil
}

func clearDB(z *notes.Zine, tx sqlx.Preparer) error {
	stmt, err := tx.Prepare("DELETE FROM articles")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(); err != nil {
		return err
	}

	stmt, err = tx.Prepare("DELETE FROM article_tag")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}

func populateDB(f fs.FS, z *notes.Zine, tx sqlx.Preparer) error {
	t0 := time.Now()

	entries, err := fs.ReadDir(f, ".")
	if err != nil {
		return fmt.Errorf("fs.ReadDir: %w", err)
	}

	articles := make([]notes.Article, len(entries))
	wg := &sync.WaitGroup{}
	for i, e := range entries {
		wg.Add(1)

		name := e.Name()
		i := i
		go func() {
			defer wg.Done()

			// unclear what to do about errors here
			b, err := fs.ReadFile(f, name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dropbox.Download: %s\n", err)
				return
			}

			articles[i], err = notes.ReadArticle(bytes.NewReader(b))
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
		if err := z.InsertArticle(tx, a); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	fmt.Fprintf(os.Stderr, "db loaded in %s\n", time.Now().Sub(t0))

	return nil
}

func loadDB(f fs.FS, generation *chan bool) (z *notes.Zine, err error) {
	z, err = notes.NewZine("")
	if err != nil {
		return nil, err
	}

	tx, err := z.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback() // how to handle error?
		} else {
			err = tx.Commit()
		}
	}()

	if err := populateDB(f, z, tx); err != nil {
		return nil, fmt.Errorf("populateDB: %w", err)
	}

	maintainDB(f, generation, z)
	return z, nil
}
