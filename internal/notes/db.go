package notes

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const Schema = `
CREATE TABLE articles (
	title,
	url,
	filename,
	reviewed_on NULLABLE,
	review_by NULLABLE,
	body,
	markdownlua
);
CREATE TABLE article_tag ( id, tag );
CREATE VIEW _ ( id, title, url, filename, body, markdownlua, reviewed_on, review_by, tag) AS
	SELECT a.rowid, title, url, filename, body, markdownlua, reviewed_on, review_by, tag
	FROM articles a
	JOIN article_tag at ON a.rowid = at.id;
`

type DB struct{ *sqlx.DB }

func NewDB(name string) (*DB, error) {
	var (
		dbh *sqlx.DB
		err error
	)

	if name == "" {
		name = "notes"
	}

	dbh, err = sqlx.Open("sqlite", "file:"+name+"?mode=memory&_sync=OFF&_journal=OFF&_vacuum=0&cache=shared")
	if err != nil {
		return nil, err
	}

	if _, err := dbh.Exec(`
		PRAGMA journal_mode = OFF
	`); err != nil {
		return nil, err
	}

	if _, err := dbh.Exec(`
		PRAGMA synchronous = OFF
	`); err != nil {
		return nil, err
	}

	if _, err := dbh.Exec(`
		PRAGMA auto_vacuum = OFF
	`); err != nil {
		return nil, err
	}

	if _, err := dbh.Exec(Schema); err != nil {
		return nil, err
	}

	var success bool
	defer func() {
		if !success {
			dbh.Close()
		}
	}()

	d := &DB{DB: dbh}

	success = true
	return d, nil
}

func (d *DB) InsertArticle(db sqlx.Preparer, a Article) error {
	if a.Filename == "" {
		return errors.New("Filename is required")
	}
	if a.URL == "" {
		return errors.New("URL is required")
	}

	stmt, err := sqlx.Preparex(db, `INSERT INTO articles (
		title, url, filename, reviewed_on, review_by, body, markdownlua
	) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	r, err := stmt.Exec(a.Title, a.URL, a.Filename, a.ReviewedOn, a.ReviewBy, string(a.Body), string(a.MarkdownLua))
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	insertTags, err := sqlx.Preparex(db, `INSERT INTO article_tag (id, tag) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	for _, tag := range a.Tags {
		if _, err := insertTags.Exec(id, tag); err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) LoadArticle(db sqlx.Preparer, name string) (Article, error) {
	stmt, err := sqlx.Preparex(db, `
	SELECT rowid, title, url, filename, reviewed_on, review_by, body, markdownlua
	FROM articles
	WHERE filename = ?
	`)
	if err != nil {
		return Article{}, err
	}

	var ret struct {
		Article
		RowID int
	}

	if err := stmt.Get(&ret, name); err != nil {
		return Article{}, err
	}

	tagsStmt, err := sqlx.Preparex(db, `SELECT tag FROM article_tag WHERE id = ?`)
	if err != nil {
		return Article{}, err
	}

	if err := tagsStmt.Select(&ret.Tags, ret.RowID); err != nil {
		return Article{}, err
	}

	return ret.Article, nil
}

func (d *DB) DeleteArticle(db sqlx.Preparer, name string) error {
	tagStmt, err := sqlx.Preparex(db, `DELETE FROM article_tag WHERE id IN (SELECT rowid FROM articles WHERE filename = ?)`)
	if err != nil {
		return err
	}

	if _, err := tagStmt.Exec(name); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	stmt, err := sqlx.Preparex(db, `DELETE FROM articles WHERE filename = ?`)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(name); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	return nil
}

func (d *DB) ReplaceArticle(db sqlx.Preparer, a Article) (err error) {
	if err := d.DeleteArticle(db, a.Filename); err != nil {
		return err
	}
	if err := d.InsertArticle(db, a); err != nil {
		return err
	}

	return nil
}
