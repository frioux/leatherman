// +build generator

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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

			Dir     string
			GoFiles []string
		}
		if err := json.Unmarshal(s.Bytes(), &c); err != nil {
			return err
		}
		dir := c.Dir

		for _, file := range c.GoFiles {
			path := dir + "/" + file
			if !toolMatch.MatchString(path) {
				continue
			}
			cmd := exec.Command("goblin", "-file", path)
			cmd.Stderr = os.Stderr
			o, err := cmd.Output()
			if err != nil {
				return err
			}
			var g struct {
				PackageName struct {
					Value string
				} `json:"package-name"`
				Declarations []struct {
					Type string
					Name struct {
						Value string
					}
				}
			}
			if err := json.Unmarshal(o, &g); err != nil {
				return err
			}

			for _, decl := range g.Declarations {
				if decl.Type != "function" {
					continue
				}
				if strings.ToUpper(decl.Name.Value[:1]) != decl.Name.Value[:1] {
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
					decl.Name.Value, c.ImportPath, g.PackageName.Value, path,
				}

				if err := e.Encode(out); err != nil {
					return err
				}
			}
		}
	}

	return s.Err()
}
