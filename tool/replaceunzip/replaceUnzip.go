package replaceunzip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

var garbage = regexp.MustCompile(`(?:^__MACOSX/|/\.DS_Store$)`)

func hasRoot(r *zip.ReadCloser) bool {
	names := make([]string, 0, len(r.File))
	for _, f := range r.File {
		if garbage.MatchString(f.Name) {
			continue
		}
		names = append(names, f.Name)
	}
	sort.Slice(names, func(i, j int) bool { return len(names[i]) < len(names[j]) })

	if !strings.HasSuffix(names[0], "/") {
		return false
	}
	root := names[0]

	for _, member := range names[1:] {
		if !strings.HasPrefix(member, root) {
			return false
		}
	}
	return true
}

func genRoot(zipName string) string {
	file := filepath.Base(zipName)

	ext := filepath.Ext(file)
	if ext == "" {
		return file
	}
	return strings.TrimSuffix(file, ext)
}

// ReplaceUnzip acts like unzip, but leaves out .DS_Store and __MACOSX files,
// and puts all of the zip contents in a single root directory if they were not
// already.
func ReplaceUnzip(args []string, _ io.Reader) error {
	if len(args) != 2 {
		fmt.Println("Usage:", args[0], "some-zip-file.zip")
		os.Exit(1)
	}

	zipName := args[1]

	fmt.Println("Archive:", zipName)
	r, err := zip.OpenReader(zipName)
	if err != nil {
		return errors.Wrap(err, "Couldn't open zip file")
	}
	defer r.Close()
	var root string
	if !hasRoot(r) {
		root = genRoot(zipName)
	}

	for _, f := range r.File {
		err := extractMember(root, f)
		if err != nil {
			return errors.Wrap(err, "extractMember")
		}
	}
	return nil
}

func extractMember(root string, f *zip.File) error {
	if garbage.MatchString(f.Name) {
		return nil
	}
	destName := filepath.Join(root, f.Name)
	fmt.Printf("  inflating: %s\n", destName)

	rc, err := f.Open()
	if err != nil {
		return errors.Wrap(err, "Couldn't open zip file member")
	}
	defer rc.Close()

	dir := filepath.Dir(destName)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create directory to extract to: %s", err)
		return nil
	}

	file, err := os.Create(destName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create file to extract to: %s", err)
		return nil
	}

	_, err = io.Copy(file, rc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't copy zip file member (%s): %s", destName, err)
	}
	err = file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't close extracted file: %s", err)
	}

	return nil
}
