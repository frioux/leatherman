package main

import (
	"fmt"
	"io"
	"sort"
)

// Help prints tool listing
func Help(args []string, _ io.Reader) {
	tools := make([]string, 0, len(Dispatch))
	for k := range Dispatch {
		tools = append(tools, k)
	}

	str := "Tools:\n"
	sort.Strings(tools)
	for _, k := range tools {
		str += " * " + k + "\n"
	}
	fmt.Println(str)
}
