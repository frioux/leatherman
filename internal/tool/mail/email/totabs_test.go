package email_test

import (
	"strings"

	"github.com/frioux/leatherman/internal/tool/mail/email"
)

func ExampleToTabs() {

	r := strings.NewReader(`"Frew Schmidt" <frew@frew.frew>`)

	email.ToTabs(nil, r)
	// Output: frew@frew.frew	Frew Schmidt
}
