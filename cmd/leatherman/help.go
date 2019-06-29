package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"golang.org/x/xerrors"
)

// Help prints tool listing
func Help(args []string, _ io.Reader) error {
	flags := flag.NewFlagSet("help", flag.ExitOnError)

	var full bool
	flags.BoolVar(&full, "v", false, "show full help")

	var command string
	flags.StringVar(&command, "command", "", "show help for just command")

	err := flags.Parse(args[1:])
	if err != nil {
		return xerrors.Errorf("flags.Parse: %w", err)
	}

	if full {
		fmt.Print(string(readme))
		return nil
	}

	if command != "" {
		readme, ok := commandReadme[command]
		if !ok {
			fmt.Fprintf(os.Stderr, "No such command: %s\n", command)
			os.Exit(1)
		}
		fmt.Print(string(readme))
		return nil
	}

	tools := make([]string, 0, len(Dispatch))
	for k := range Dispatch {
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
