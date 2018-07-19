package main

import "strings"

func ExampleAddrspecToTabs() {

	r := strings.NewReader(`"Frew Schmidt" <frew@frew.frew>`)

	AddrspecToTabs(nil, r)
	// Output: frew@frew.frew	Frew Schmidt
}
