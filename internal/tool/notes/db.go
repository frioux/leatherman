package notes

import (
	"io/ioutil"
	"os"

	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func loadDB(cl dropbox.Client) (*sqlx.DB, func(), error) {
	d, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, func() {}, err
	}
	cleanup := func() { os.RemoveAll(d) }

	r, err := cl.Download("/notes/.posts.db")
	if err != nil {
		return nil, cleanup, err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, cleanup, err
	}

	if err := ioutil.WriteFile(d+"/.posts.db", b, 0600); err != nil {
		return nil, cleanup, err
	}

	dbh, err := sqlx.Open("sqlite3", "file:"+d+"/.posts.db?_sync=OFF&_journal=OFF&_vacuum=0")
	if err != nil {
		return nil, cleanup, err
	}

	return dbh, cleanup, nil
}
