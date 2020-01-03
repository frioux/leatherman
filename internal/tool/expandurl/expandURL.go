package expandurl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/frioux/leatherman/internal/lmhttp"
)

var tidyRE = regexp.MustCompile(`^\s*(.*?)\s*$`)

/*
Run reads text on STDIN and writes the same text back, converting any links to
Markdown links, with the title of the page as the title of the link.

Command: expand-url
*/
func Run(args []string, stdin io.Reader) error {
	return run(stdin, os.Stdout)
}

func run(r io.Reader, w io.Writer) error {
	// some cookies cause go to log warnings to stderr
	log.SetOutput(ioutil.Discard)

	ua := &http.Client{}

	scanner := bufio.NewScanner(r)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading standard input: %w", err)
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

func urlToLink(ua *http.Client, url string) (string, error) {
	resp, err := lmhttp.Get(context.TODO(), url)
	if err != nil {
		return "", fmt.Errorf("lmhttp.Get: %s", err)
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
