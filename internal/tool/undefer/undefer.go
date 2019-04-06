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

	"github.com/pkg/errors"
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
			return nil, errors.Wrap(err, "couldn't parse "+file.Name())
		}
		month, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, errors.Wrap(err, "couldn't parse "+file.Name())
		}
		day, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, errors.Wrap(err, "couldn't parse "+file.Name())
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
			return errors.Wrap(err, "couldn't Open")
		}
		_, err = io.Copy(stdout, file)
		if err != nil {
			return errors.Wrap(err, "couldn't Copy")
		}
		err = file.Close()
		if err != nil {
			return errors.Wrap(err, "couldn't Close")
		}
		err = os.Remove(path)
		if err != nil {
			return errors.Wrap(err, "couldn't Remove")
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
		return errors.Wrap(err, "Couldn't Open")
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return errors.Wrap(err, "Couldn't ReadDir")
	}

	paths, err := newFiles(args[1], files, time.Now())
	if err != nil {
		return errors.Wrap(err, "Couldn't get newFiles")
	}
	err = content(paths, os.Stdout)
	if err != nil {
		return errors.Wrap(err, "Couldn't write content")
	}

	return nil
}
