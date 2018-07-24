package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()

func newFiles(dir string, files []os.FileInfo, now time.Time) ([]string, error) {
	dateRE := regexp.MustCompile(`^(\d{4})-(\d\d)-(\d\d)`)

	ret := make([]string, 0, len(files))

	for _, file := range files {
		matches := dateRE.FindStringSubmatch(file.Name())
		if len(matches) != 4 {
			continue
		}
		year, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("couldn't parse %s: %s", file.Name(), err)
		}
		month, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, fmt.Errorf("couldn't parse %s: %s", file.Name(), err)
		}
		day, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("couldn't parse %s: %s", file.Name(), err)
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
		file, err := fs.Open(path)
		if err != nil {
			return fmt.Errorf("couldn't Open: %s", err)
		}
		_, err = io.Copy(stdout, file)
		if err != nil {
			return fmt.Errorf("couldn't Copy: %s", err)
		}
		err = file.Close()
		if err != nil {
			return fmt.Errorf("couldn't Close: %s", err)
		}
		err = fs.Remove(path)
		if err != nil {
			return fmt.Errorf("couldn't Remove: %s", err)
		}
	}

	return nil
}

// Undefer prints the contents of files in the passed directory that have a
// prefix of a date in the past, and then deletes the files.
func Undefer(args []string, _ io.Reader) {
	if len(args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s $dir\n", args[0])
		os.Exit(1)
	}

	files, err := afero.Afero{Fs: fs}.ReadDir(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't ReadDir: %s", err)
		os.Exit(1)
	}

	paths, err := newFiles(args[1], files, time.Now())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't get newFiles: %s", err)
		os.Exit(1)
	}
	err = content(paths, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't write content: %s", err)
		os.Exit(1)
	}
}
