package mozcookiejar

import (
	"database/sql"
	"net/http/cookiejar"
	"net/url"
	"testing"

	_ "github.com/mattn/go-sqlite3" // sqlite3 required
)

const year3018 = 33059825142

func freshDBAndJar(t *testing.T) (*sql.DB, *cookiejar.Jar) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Couldn't create db: %s", err)
	}
	_, err = db.Exec(`
   CREATE TABLE moz_cookies (
		id INTEGER PRIMARY KEY,
      name TEXT,
      value TEXT,
      host TEXT,
      path TEXT,
      expiry INTEGER,
      isSecure INTEGER
	)`)
	if err != nil {
		t.Fatalf("Couldn't create table: %s", err)
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("Couldn't create cookiejar: %s", err)
	}
	return db, jar
}

func insert(t *testing.T, db *sql.DB, name, value, host, path string, expiry int, isSecure int) {
	_, err := db.Exec(
		`INSERT INTO moz_cookies
			(name, value, host, path, expiry, isSecure)
  VALUES (?, ?, ?, ?, ?, ?)`,
		name, value, host, path, expiry, isSecure)
	if err != nil {
		t.Fatalf("Couldn't create record: %s", err)
	}
}

func urlFromString(t *testing.T, urlStr string) *url.URL {
	ret, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("Couldn't parse URL: %s", err)
	}
	return ret
}

func TestBasic(t *testing.T) {
	db, jar := freshDBAndJar(t)
	insert(t, db, "a", "b", "www.foo.com", "/frew", year3018, 1)
	insert(t, db, "c", "d", "www.bar.com", "/frioux", year3018, 0)
	err := LoadIntoJar(db, jar)
	if err != nil {
		t.Fatalf("Couldn't populate cookiejar: %s", err)
	}

	c := jar.Cookies(urlFromString(t, "https://x.y.foo.com/frew"))
	if len(c) > 0 {
		t.Error("Domain worked")
	}

	c = jar.Cookies(urlFromString(t, "https://www.foo.com/"))
	if len(c) > 0 {
		t.Error("path worked")
	}

	c = jar.Cookies(urlFromString(t, "http://www.foo.com/frew"))
	if len(c) > 0 {
		t.Error("isSecure worked")
	}

	c = jar.Cookies(urlFromString(t, "https://www.foo.com/frew"))
	if len(c) == 0 {
		t.Error("Normal access works")
	}
}

func TestMutli(t *testing.T) {
	db, jar := freshDBAndJar(t)
	insert(t, db, "a", "b", "www.foo.com", "/frew", year3018, 1)
	insert(t, db, "c", "d", "www.foo.com", "/frew", year3018, 1)
	err := LoadIntoJar(db, jar)
	if err != nil {
		t.Fatalf("Couldn't populate cookiejar: %s", err)
	}

	c := jar.Cookies(urlFromString(t, "https://x.y.foo.com/frew"))
	if len(c) > 0 {
		t.Error("Domain worked")
		t.Log(c)
	}

	c = jar.Cookies(urlFromString(t, "https://www.foo.com/"))
	if len(c) > 0 {
		t.Error("path worked")
	}

	c = jar.Cookies(urlFromString(t, "http://www.foo.com/frew"))
	if len(c) > 0 {
		t.Error("isSecure worked")
	}

	c = jar.Cookies(urlFromString(t, "https://www.foo.com/frew"))
	if len(c) != 2 {
		t.Error("Normal access works")
	}
}

func TestInvalidDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Couldn't create db: %s", err)
	}
	_, err = db.Exec(`
   CREATE TABLE invalid_db (
		id INTEGER PRIMARY KEY,
      name TEXT,
      value TEXT
	)`)
	if err != nil {
		t.Fatalf("Couldn't create table: %s", err)
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("Couldn't create cookiejar: %s", err)
	}
	err = LoadIntoJar(db, jar)
	if err == nil {
		t.Error("Invalid db should surface an error")
	}
}

func TestInvalidData(t *testing.T) {
	db, jar := freshDBAndJar(t)
	_, err := db.Exec(
		`INSERT INTO moz_cookies
			(name, value, host, path, expiry, isSecure)
  VALUES (?, ?, ?, ?, ?, ?)`,
		"a", "b", ".frew.com", "/", "explode", 0)
	if err != nil {
		t.Fatalf("Couldn't create record: %s", err)
	}
	err = LoadIntoJar(db, jar)
	if err == nil {
		t.Error("Invalid record should surface an error")
	}
}
