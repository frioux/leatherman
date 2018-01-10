package main

import (
	"fmt"
)

// Help prints tool listing
func Help(args []string) {
	str := "Tools:\n"
	for k := range Dispatch {
		str += " * " + k + "\n"
	}
	fmt.Println(str)
}
