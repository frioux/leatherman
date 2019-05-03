package deferlwn

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/frioux/leatherman/pkg/lwn"
	"github.com/frioux/leatherman/pkg/timeutil"
	"golang.org/x/xerrors"
)

// Run writes the input that contained links to be made available in the
// future on the relevant day, and otherwise prints the lines on standard
// output.
func Run(args []string, stdin io.Reader) error {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <dir>\n", args[0])
		os.Exit(1)
	}
	dir := args[1]

	return run(dir, stdin, os.Stdout, os.Stderr)
}

func run(dir string, r io.Reader, w, wErr io.Writer) error {
	// tokens limits parallelism to 10
	tokens := make(chan struct{}, 10)

	// wg ensures that we block till all lines are done
	wg := sync.WaitGroup{}

	s := bufio.NewScanner(r)

	for s.Scan() {
		line := s.Text()

		wg.Add(1)
		tokens <- struct{}{}

		go func() {
			err := deferLink(line, dir)
			if err != nil {
				fmt.Fprintln(wErr, err)
				fmt.Fprintln(w, line)
			}
			<-tokens
			wg.Done()
		}()
	}

	wg.Wait()
	if s.Err() != nil {
		return xerrors.Errorf("bufio.Scanner.Scan: %w", s.Err())
	}

	return nil

}

var (
	errNoMatch = xerrors.Errorf("no matching url")
	linkFinder = regexp.MustCompile(`^(.*\()(https?://\S+)(\).*)$`)
)

func deferLink(line, dir string) error {
	match := linkFinder.FindStringSubmatch(line)
	if len(match) == 0 {
		return errNoMatch
	}

	page, err := url.Parse(match[2])
	if err != nil {
		return xerrors.Errorf("url.Parse: %w", err)
	}

	date, err := lwn.AvailableOn(page)
	if err != nil {
		return xerrors.Errorf("lwn.AvailableOn: %w", err)
	}
	date = timeutil.JumpTo(date, time.Friday)

	filename := dir + "/" + date.Format("2006-01-02") + "-lwn.md"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return xerrors.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(line + "\n")
	if err != nil {
		return xerrors.Errorf("os.File.WriteString: %w", err)
	}

	return nil
}
