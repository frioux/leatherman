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
		"addrs":               Addrs,
		"addrspec-to-tabs":    AddrspecToTabs,
		"clocks":              Clocks,
		"csv2json":            CsvToJson,
		"debounce":            Debounce,
		"ec2-resource-for-ip": Ec2ResourceForIp,
		"netrc-password":      NetrcPassword,
		"pomotimer":           Pomotimer,
		"render-mail":         RenderMail,

		"help":    Help,
		"explode": Explode,
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
