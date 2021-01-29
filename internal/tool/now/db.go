package now

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/jmoiron/sqlx"
)

func loadDB(db dropbox.Client, dir string) (*sqlx.DB, func(), error) {
	d, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, func() {}, err
	}
	cleanup := func() { os.RemoveAll(d) }

	r, err := db.ListFolder(dropbox.ListFolderParams{Path: dir})
	if err != nil {
		return nil, cleanup, err
	}

	entries := r.Entries

	for r.HasMore {
		r, err = db.ListFolderContinue(r.Cursor)
		if err != nil {
			return nil, cleanup, err
		}

		entries = append(entries, r.Entries...)
	}

	wg := &sync.WaitGroup{}
	for _, e := range entries {
		wg.Add(1)

		name := e.Name
		fmt.Println(name)
		go func() {
			defer wg.Done()

			// unclear what to do about errors here
			_, err := db.Download(dir + name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}()
	}

	wg.Wait()

	return nil, cleanup, nil
}
