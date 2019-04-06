package shellquote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func test(t *testing.T, in []string, expected string) {

	ret, err := Quote(in)
	if assert.Nil(t, err, "No error") {
		assert.Equal(t, expected, ret)
	}
}

func TestShellQuote(t *testing.T) {
	test(t, []string{""}, `''`)
	test(t, []string{"foo"}, `foo`)
	test(t, []string{"foo", "bar"}, `foo bar`)
	test(t, []string{"foo*"}, `'foo*'`)
	test(t, []string{"foo bar"}, `'foo bar'`)
	test(t, []string{"foo'bar"}, `'foo'\''bar'`)
	test(t, []string{"'foo"}, `\''foo'`)
	test(t, []string{"foo", "bar*"}, `foo 'bar*'`)
	test(t, []string{"foo'foo", "bar", "baz'"}, `'foo'\''foo' bar 'baz'\'`)
	test(t, []string{`\`}, `'\'`)
	test(t, []string{"'"}, `\'`)
	test(t, []string{`\'`}, `'\'\'`)
	test(t, []string{"a''b"}, `'a'"''"'b'`)
	test(t, []string{"azAZ09_!%+,-./:@^"}, `azAZ09_!%+,-./:@^`)
	test(t, []string{"foo=bar", "command"}, `'foo=bar' command`)
	test(t, []string{"foo=bar", "baz=quux", "command"}, `'foo=bar' 'baz=quux' command`)

	_, err := Quote([]string{"\x00"})
	assert.Equal(t, err, ErrNull)
}
