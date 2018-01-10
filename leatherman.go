package main

import (
	"os"
	"path"
)

var Dispatch map[string]func([]string)

func main() {
	which := path.Base(os.Args[0])
	args := os.Args

	Dispatch = map[string]func([]string){
		"addrspec-to-tabs": AddrspecToTabs,
		"clocks":           Clocks,
		"help":             Help,
	}

	if which == "leatherman" && len(args) > 1 {
		args = args[1:]
		which = args[0]
	}

	fn, ok := Dispatch[which]
	if !ok {
		Help(os.Args)
		os.Exit(1)
	}
	fn(args)
}
