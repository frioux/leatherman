package shellquote_test

import (
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
	"github.com/frioux/leatherman/pkg/shellquote"
)

func test(t *testing.T, in []string, expected string) {
	ret, err := shellquote.Quote(in)
	if err != nil {
		t.Errorf("Quote errored: %s", err)
		return
	}
	testutil.Equal(t, ret, expected, "wrong quote")
}

func TestShellQuote(t *testing.T) {
	t.Parallel()

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

	_, err := shellquote.Quote([]string{"\x00"})
	if err != shellquote.ErrNull {
		t.Errorf("err should be ErrNull; was %s", err)
	}
}
