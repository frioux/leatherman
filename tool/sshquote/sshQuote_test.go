package sshquote

func ExampleRun() {
	Run([]string{"ssh-quote", "foo", "bar"}, nil)
	// Output: 'sh -c '\''foo bar'\'
}
