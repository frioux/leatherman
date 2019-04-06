package mozcookiejar_test

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/frioux/leatherman/pkg/mozcookiejar"
	_ "github.com/mattn/go-sqlite3" // sqlite3 required
	"golang.org/x/net/publicsuffix"
)

func Example() {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build cookies: %s\n", err)
		os.Exit(1)
	}
	db, err := sql.Open("sqlite3", os.Getenv("MOZ_COOKIEJAR"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open db: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()

	err = mozcookiejar.LoadIntoJar(db, jar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load cookies: %s\n", err)
		os.Exit(1)
	}
	ua := http.Client{Jar: jar}

	resp, err := ua.Get("https://some.authenticated.com/website")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch page: %s\n", err)
		os.Exit(1)
	}
	io.Copy(os.Stdout, resp.Body)
}
