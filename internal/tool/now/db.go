package now

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/frioux/leatherman/internal/notes"
)

func loadDB(db dropbox.Client, dir string) (*notes.Zine, error) {
	t0 := time.Now()

	r, err := db.ListFolder(dropbox.ListFolderParams{Path: dir})
	if err != nil {
		return nil, err
	}

	entries := r.Entries

	for r.HasMore {
		r, err = db.ListFolderContinue(r.Cursor)
		if err != nil {
			return nil, err
		}

		entries = append(entries, r.Entries...)
	}

	z, err := notes.NewZine()
	if err != nil {
		return nil, err
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
			r, err := db.Download(dir + name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			articles[i], err = notes.ReadArticle(r)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}()
	}
	wg.Wait()

	for _, a := range articles {
		if err := z.InsertArticle(a); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	fmt.Fprintf(os.Stderr, "db loaded in %s\n", time.Now().Sub(t0))

	return z, nil
}
