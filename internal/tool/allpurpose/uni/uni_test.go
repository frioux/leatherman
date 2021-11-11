package uni_test

import (
	"fmt"
	"os"

	"github.com/frioux/leatherman/internal/tool/allpurpose/uni"
)

func ExampleDescribe() {
	err := uni.Describe([]string{"uni", "⢾"}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't Describe: %s\n", err)
		os.Exit(1)
	}
	// Output: '⢾' @ 10430 aka BRAILLE PATTERN DOTS-234568 ( graphic | printable | symbol )
}
