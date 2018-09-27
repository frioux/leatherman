package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/cookiejar"
	"os"
	"regexp"
	"sync"

	"github.com/frioux/mozcookiejar"
	"github.com/headzoo/surf"
	_ "github.com/mattn/go-sqlite3" // sqlite3 required
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

var jar *cookiejar.Jar
var tidyRE = regexp.MustCompile(`^\s*(.*?)\s*$`)

// ExpandURL replaces URLs from stdin with their markdown version, using a
// title from the actual page, loaded using cookies discovered via the
// MOZ_COOKIEJAR env var.
func ExpandURL(args []string, stdin io.Reader) error {
	// some cookies cause go to log warnings to stderr
	log.SetOutput(ioutil.Discard)

	var err error
	jar, err = cj()
	if err != nil {
		return errors.Wrap(err, "loading cookiejar")
	}

	scanner := bufio.NewScanner(stdin)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "reading standard input")
	}

	// tokens limits parallelism to 10
	tokens := make(chan struct{}, 10)

	// wg ensures that we block till all lines are done
	wg := sync.WaitGroup{}

	for i := range lines {
		i := i
		wg.Add(1)
		tokens <- struct{}{}

		go func() {
			lines[i] = replaceLink(lines[i])
			<-tokens
			wg.Done()
		}()
	}

	wg.Wait()

	for _, line := range lines {
		fmt.Println(line)
	}

	return nil
}

func cj() (*cookiejar.Jar, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build cookies")
	}

	path := os.Getenv("MOZ_COOKIEJAR")
	if path == "" {
		return nil, errors.New("MOZ_COOKIEJAR should be set for expand-url to work")
	}
	db, err := sql.Open("sqlite3", "file:"+path+"?cache=shared&_journal_mode=WAL")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open db")
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	err = mozcookiejar.LoadIntoJar(db, jar)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load cookies")
	}
	return jar, nil
}

func urlToLink(url string) (string, error) {
	ua := surf.NewBrowser()
	ua.SetCookieJar(jar)
	err := ua.Open(url)
	if err != nil {
		return "", fmt.Errorf("authBabmoo: %s", err)
	}
	title := tidyRE.FindStringSubmatch(ua.Title())
	if len(title) != 2 {
		return "", fmt.Errorf("Title is blank")
	}
	return fmt.Sprintf("[%s](%s)", title[1], url), nil
}

var urlFinder = regexp.MustCompile(`^(|.*\s)(https?://\S+)(\s.*|)$`)

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
