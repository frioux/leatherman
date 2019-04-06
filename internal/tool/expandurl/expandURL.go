package expandurl

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/frioux/leatherman/pkg/mozcookiejar"
	_ "github.com/mattn/go-sqlite3" // sqlite3 required
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

var tidyRE = regexp.MustCompile(`^\s*(.*?)\s*$`)

// Run replaces URLs from stdin with their markdown version, using a
// title from the actual page, loaded using cookies discovered via the
// MOZ_COOKIEJAR env var.
func Run(args []string, stdin io.Reader) error {
	return run(stdin, os.Stdout)
}

func run(r io.Reader, w io.Writer) error {
	// some cookies cause go to log warnings to stderr
	log.SetOutput(ioutil.Discard)

	jar, err := cj()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "loading cookiejar"))
		jar, _ = cookiejar.New(nil)
	}
	ua := &http.Client{Jar: jar}

	scanner := bufio.NewScanner(r)
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
			lines[i] = replaceLink(ua, lines[i])
			<-tokens
			wg.Done()
		}()
	}

	wg.Wait()

	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	return nil
}

func cj() (*cookiejar.Jar, error) {
	j, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build cookies")
	}

	path := os.Getenv("MOZ_COOKIEJAR")
	if path == "" {
		return nil, errors.New("MOZ_COOKIEJAR should be set for expand-url to work")
	}

	orig, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open for copying")
	}

	dest, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.TempFile for copying")
	}

	_, err = io.Copy(dest, orig)
	if err != nil {
		return nil, errors.Wrap(err, "io.Copy for copying")
	}
	err = dest.Close()
	if err != nil {
		return nil, errors.Wrap(err, "dest.Close for copying")
	}
	err = orig.Close()
	if err != nil {
		return nil, errors.Wrap(err, "orig.Close for copying")
	}

	db, err := sql.Open("sqlite3", "file:"+dest.Name())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open db")
	}
	defer db.Close()

	err = mozcookiejar.LoadIntoJar(db, j)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load cookies")
	}
	err = os.Remove(dest.Name())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to clean up db copy")
	}

	return j, nil
}

func urlToLink(ua *http.Client, url string) (string, error) {
	resp, err := ua.Get(url)
	if err != nil {
		return "", fmt.Errorf("ua.Get: %s", err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("goquery.NewDocumentFromReader: %s", err)
	}
	title := tidyRE.FindStringSubmatch(doc.Find("title").Text())
	if len(title) != 2 {
		return "", fmt.Errorf("title is blank")
	}
	return fmt.Sprintf("[%s](%s)", title[1], url), nil
}

var urlFinder = regexp.MustCompile(`^(|.*\s)(https?://\S+)(\s.*|)$`)

func replaceLink(ua *http.Client, line string) string {
	for {
		if match := urlFinder.FindStringSubmatch(line); len(match) > 0 {
			md, err := urlToLink(ua, match[2])
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
