package email

import "strings"

func ExampleToTabs() {

	r := strings.NewReader(`"Frew Schmidt" <frew@frew.frew>`)

	ToTabs(nil, r)
	// Output: frew@frew.frew	Frew Schmidt
}
