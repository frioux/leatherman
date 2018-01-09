package main

import (
	"fmt"
	"os"
)

func Help() {
	str := "Tools:\n"
	for k := range Dispatch {
		str += " * " + k + "\n"
	}
	fmt.Println(str)
	os.Exit(1)
}
