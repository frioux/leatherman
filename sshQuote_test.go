package main

func ExampleSSHQuote() {
	SSHQuote([]string{"ssh-quote", "foo", "bar"})
	// Output: 'sh -c '\''foo bar'\'
}
