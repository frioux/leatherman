package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net/http/cookiejar"
	"os"
	"regexp"

	"github.com/frioux/mozcookiejar"
	"github.com/headzoo/surf"
	_ "github.com/mattn/go-sqlite3" // sqlite3 required
	"golang.org/x/net/publicsuffix"
)

func ExpandURL(args []string) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Println(replaceLink(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func cj() *cookiejar.Jar {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build cookies: %s\n", err)
		os.Exit(1)
	}

	path := os.Getenv("MOZ_COOKIEJAR")
	if path == "" {
		fmt.Fprintln(os.Stderr, "MOZ_COOKIEJAR should be set for expand-url to work", err)
		return jar
	}
	db, err := sql.Open("sqlite3", path)
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
	return jar
}

func urlToLink(url string) (string, error) {
	ua := surf.NewBrowser()
	ua.SetCookieJar(jar)
	err := ua.Open(url)
	if err != nil {
		return "", fmt.Errorf("authBabmoo: %s", err)
	}
	return fmt.Sprintf("[%s](%s)", ua.Title(), url), nil
}

var urlFinder *regexp.Regexp
var jar *cookiejar.Jar

func init() {
	urlFinder = regexp.MustCompile(`^(|.*\s)(https?://\S+)(\s.*|)$`)
	jar = cj()
}

func replaceLink(line string) string {
	for {
		if match := urlFinder.FindStringSubmatch(line); len(match) > 0 {
			md, err := urlToLink(match[2])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				break
			}
			line = match[1] + md + match[3]
			continue
		}
		break
	}
	return line
}
