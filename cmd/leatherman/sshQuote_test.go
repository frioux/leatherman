package main

func ExampleSSHQuote() {
	SSHQuote([]string{"ssh-quote", "foo", "bar"}, nil)
	// Output: 'sh -c '\''foo bar'\'
}
