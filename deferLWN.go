package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"sync"

	"github.com/frioux/leatherman/lwn"
	"github.com/pkg/errors"
)

// DeferLWN writes the input that contained links to be made available in the
// future on the relevant day, and otherwise prints the lines on standard
// output.
func DeferLWN(args []string, stdin io.Reader) error {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <dir>\n", args[0])
		os.Exit(1)
	}
	dir := args[1]

	// tokens limits parallelism to 10
	tokens := make(chan struct{}, 10)

	// wg ensures that we block till all lines are done
	wg := sync.WaitGroup{}

	s := bufio.NewScanner(stdin)

	for s.Scan() {
		line := s.Text()

		wg.Add(1)
		tokens <- struct{}{}

		go func() {
			err := deferLink(line, dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				fmt.Println(line)
			}
			<-tokens
			wg.Done()
		}()
	}

	wg.Wait()
	if s.Err() != nil {
		return errors.Wrap(s.Err(), "bufio.Scanner.Scan()")
	}

	return nil
}

var (
	errNoMatch = errors.New("No matching url")
	linkFinder = regexp.MustCompile(`^(.*\()(https?://\S+)(\).*)$`)
)

func deferLink(line, dir string) error {
	match := linkFinder.FindStringSubmatch(line)
	if len(match) == 0 {
		return errNoMatch
	}

	page, err := url.Parse(match[2])
	if err != nil {
		return errors.Wrap(err, "url.Parse")
	}

	date, err := lwn.AvailableOn(page)
	if err != nil {
		return errors.Wrap(err, "lwn.AvailableOn")
	}

	filename := dir + "/" + date.Format("2006-01-02") + "-lwn.md"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return errors.Wrap(err, "os.OpenFile")
	}
	defer file.Close()

	_, err = file.WriteString(line + "\n")
	if err != nil {
		return errors.Wrap(err, "os.File.WriteString")
	}

	return nil
}
