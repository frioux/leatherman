package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/icza/backscanner"
	"golang.org/x/text/unicode/runenames"
)

func PrependEmojiHist(args []string) {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "you must pass a history file!")
		os.Exit(1)
	}

	file, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open history file: %s\n", err)
		os.Exit(1)
	}
	fi, err := os.Stat(args[1])
	var pos int
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Couldn't stat history file: %s\n", err)
		os.Exit(1)
	}
	if fi != nil {
		pos = int(fi.Size())
	}

	seen := map[string]bool{}
	scanner := backscanner.New(file, pos)
	for {
		line, _, err := scanner.Line()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't read line: %s\n", err)
			os.Exit(1)
		}

		for i, char := range line {
			name := strings.ToLower(runenames.Name(char))
			if seen[name] {
				continue
			}
			seen[name] = true
			fmt.Println(name)
			if i > 0 {
				fmt.Fprintln(os.Stderr, "Multiple characters on line, breaking loop")
				break
			}
		}
	}

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		line := stdin.Text()
		if seen[line] {
			continue
		}
		seen[line] = true
		fmt.Println(line)
	}
}
