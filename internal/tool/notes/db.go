package notes

import (
	"io/ioutil"
	"os"

	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/jmoiron/sqlx"
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

	dbh, err := sqlx.Open("sqlite", "file:"+d+"/.posts.db")
	if err != nil {
		return nil, cleanup, err
	}

	if _, err := dbh.Exec(`
		PRAGMA journal_mode = OFF
	`); err != nil {
		return nil, cleanup, err
	}

	if _, err := dbh.Exec(`
		PRAGMA synchronous = OFF
	`); err != nil {
		return nil, cleanup, err
	}

	if _, err := dbh.Exec(`
		PRAGMA auto_vacuum = OFF
	`); err != nil {
		return nil, cleanup, err
	}

	return dbh, cleanup, nil
}
