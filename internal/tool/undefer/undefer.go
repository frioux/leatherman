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
			return nil, fmt.Errorf("couldn't parse %s: %w", file.Name(), err)
		}
		month, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, fmt.Errorf("couldn't parse %s: %w", file.Name(), err)
		}
		day, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("couldn't parse %s: %w", file.Name(), err)
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
			return fmt.Errorf("couldn't Open: %w", err)
		}
		_, err = io.Copy(stdout, file)
		if err != nil {
			return fmt.Errorf("couldn't Copy: %w", err)
		}
		err = file.Close()
		if err != nil {
			return fmt.Errorf("couldn't Close: %w", err)
		}
		err = os.Remove(path)
		if err != nil {
			return fmt.Errorf("couldn't Remove: %w", err)
		}
	}

	return nil
}

/*
Run takes a directory argument, prints contents of each file named before the
current date, and then deletes the file.

If the current date were `2018-06-07` the starred files would be printed and
then deleted:

```
 * 2018-01-01.txt
 * 2018-06-06-awesome-file.md
   2018-07-06.txt
```

Command: undefer
*/
func Run(args []string, _ io.Reader) error {
	if len(args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s $dir\n", args[0])
		os.Exit(1)
	}

	dir, err := os.Open(args[1])
	if err != nil {
		return fmt.Errorf("Couldn't Open: %w", err)
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("Couldn't ReadDir: %w", err)
	}

	paths, err := newFiles(args[1], files, time.Now())
	if err != nil {
		return fmt.Errorf("Couldn't get newFiles: %w", err)
	}
	err = content(paths, os.Stdout)
	if err != nil {
		return fmt.Errorf("Couldn't write content: %w", err)
	}

	return nil
}
