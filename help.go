package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"sort"
)

//go:embed README.mdwn
var readme []byte

// Help prints tool listing
func Help(args []string, _ io.Reader) error {
	flags := flag.NewFlagSet("help", flag.ExitOnError)

	var full bool
	flags.BoolVar(&full, "v", false, "show full help")

	var command string
	flags.StringVar(&command, "command", "", "show help for just command")

	err := flags.Parse(args[1:])
	if err != nil {
		return fmt.Errorf("flags.Parse: %w", err)
	}

	if full {
		comment, err := regexp.Compile(`<!--.*?-->\n?`)
		if err != nil {
			return err
		}

		fmt.Print(string(comment.ReplaceAll(readme, []byte{})))
		return nil
	}

	if command != "" {
		doc, err := fs.ReadFile(helpFS, helpPaths[command])
		if err != nil && errors.Is(err, fs.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "No such command: %s\n", command)
			os.Exit(1)
		}
		if err != nil {
			return err
		}
		fmt.Print(string(doc))
		return nil
	}

	tools := make([]string, 0, len(Dispatch))
	for k := range Dispatch {
		if k == "xyzzy" { // nothing to see here
			continue
		}
		tools = append(tools, k)
	}

	str := "Tools:\n"
	sort.Strings(tools)
	for _, k := range tools {
		str += " * " + k + "\n"
	}
	str += "\nGet more help for each tool with `leatherman help -command <tool>`, or `leatherman help -v`"
	fmt.Println(str)

	return nil
}
