package main

import (
	"os"
	"path"
)

var Dispatch map[string]func()

func main() {
	which := path.Base(os.Args[0])

	Dispatch = map[string]func(){
		"addrspec-to-tabs": AddrspecToTabs,
		"clocks":           Clocks,
		"help":             Help,
	}

	if which == "leatherman" {
		if len(os.Args) > 1 {
			which = os.Args[1]
		}
	}

	fn, ok := Dispatch[which]
	if !ok {
		Help()
	}
	fn()
}
