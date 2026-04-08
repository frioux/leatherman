//go:build generator
// +build generator

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var toolMatch = regexp.MustCompile(os.Getenv("LM_TOOL"))

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cmd := exec.Command("/bin/sh", "-c", "go list -json ./internal/tool/... | jq -c .")
	cmd.Stderr = os.Stderr
	o, err := cmd.Output()
	if err != nil {
		return err
	}

	s := bufio.NewScanner(bytes.NewReader(o))
	e := json.NewEncoder(os.Stdout)

	for s.Scan() {
		var c struct {
			ImportPath string

			Dir            string
			GoFiles        []string
			IgnoredGoFiles []string
		}
		if err := json.Unmarshal(s.Bytes(), &c); err != nil {
			return err
		}
		dir := c.Dir

		allFiles := append(c.GoFiles, c.IgnoredGoFiles...)
		for _, file := range allFiles {
			path := dir + "/" + file
			if !toolMatch.MatchString(path) {
				continue
			}
			mdPath := strings.TrimSuffix(path, ".go") + ".md"
			if _, err := os.Stat(mdPath); err != nil {
				continue
			}
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return err
			}

			pkgName := f.Name.Name

			for _, decl := range f.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				if fn.Recv != nil {
					continue
				}
				if !fn.Name.IsExported() {
					continue
				}

				if os.Getenv("LM_TOOL") != "" {
					fmt.Fprintf(os.Stderr, "# %s matched LM_TOOL\n", path)
				}

				out := struct {
					Func    string `json:"func"`
					Import  string `json:"import"`
					Package string `json:"package"`
					Path    string `json:"path"`
				}{
					fn.Name.Name, c.ImportPath, pkgName, path,
				}

				if err := e.Encode(out); err != nil {
					return err
				}
			}
		}
	}

	return s.Err()
}
