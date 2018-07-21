package main

import "strings"

func ExampleRenderMail() {
	r := strings.NewReader(`foo
bar
baz
Date: Wed, 18 Jul 2019 16:00:00 +0000`)

	RenderMail(nil, r)
	// Output:
	// foo
	// bar
	// baz
	// Local-Date: Thu, 18 Jul 2019 09:00:00
	// Date: Wed, 18 Jul 2019 16:00:00 +0000
}
