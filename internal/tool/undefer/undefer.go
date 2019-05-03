package undefer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"golang.org/x/xerrors"
)

var dateRE = regexp.MustCompile(`^(\d{4})-(\d\d)-(\d\d)`)

func newFiles(dir string, files []os.FileInfo, now time.Time) ([]string, error) {
	ret := make([]string, 0, len(files))

	for _, file := range files {
		matches := dateRE.FindStringSubmatch(file.Name())
		if len(matches) != 4 {
			continue
		}
		year, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, xerrors.Errorf("couldn't parse %s: %w", file.Name(), err)
		}
		month, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, xerrors.Errorf("couldn't parse %s: %w", file.Name(), err)
		}
		day, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, xerrors.Errorf("couldn't parse %s: %w", file.Name(), err)
		}

		then := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
		if then.After(now) {
			continue
		}
		ret = append(ret, filepath.Join(dir, file.Name()))
	}
	sort.Strings(ret)
	return ret, nil
}

func content(paths []string, stdout io.Writer) error {
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return xerrors.Errorf("couldn't Open: %w", err)
		}
		_, err = io.Copy(stdout, file)
		if err != nil {
			return xerrors.Errorf("couldn't Copy: %w", err)
		}
		err = file.Close()
		if err != nil {
			return xerrors.Errorf("couldn't Close: %w", err)
		}
		err = os.Remove(path)
		if err != nil {
			return xerrors.Errorf("couldn't Remove: %w", err)
		}
	}

	return nil
}

// Run prints the contents of files in the passed directory that have a
// prefix of a date in the past, and then deletes the files.
func Run(args []string, _ io.Reader) error {
	if len(args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s $dir\n", args[0])
		os.Exit(1)
	}

	dir, err := os.Open(args[1])
	if err != nil {
		return xerrors.Errorf("Couldn't Open: %w", err)
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return xerrors.Errorf("Couldn't ReadDir: %w", err)
	}

	paths, err := newFiles(args[1], files, time.Now())
	if err != nil {
		return xerrors.Errorf("Couldn't get newFiles: %w", err)
	}
	err = content(paths, os.Stdout)
	if err != nil {
		return xerrors.Errorf("Couldn't write content: %w", err)
	}

	return nil
}
