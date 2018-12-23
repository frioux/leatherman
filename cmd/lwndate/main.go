package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/frioux/leatherman/pkg/lwn"
)

func main() {
	for _, s := range os.Args[1:] {
		u, err := url.Parse(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't parse %s: %s\n", s, err)
		}
		t, err := lwn.AvailableOn(u)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't get date for %s: %s\n", u, err)
		}
		fmt.Println(u, t)
	}
}
