package minotaur

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	t.Parallel()

	type row struct {
		name string

		in []string

		expectedConfig config
		expectedErr    error
	}

	include := regexp.MustCompile("")
	ignore := regexp.MustCompile("(^.git|/.git$|/.git/)")

	var table = []row{{
		name: "simple and correct",
		in: []string{
			"./foo", "./bar",
			"--",
			"foo", "bar",
		},
		expectedConfig: config{
			dirs:    []string{"./foo", "./bar"},
			script:  []string{"foo", "bar"},
			include: include,
			ignore:  ignore,
		},
	}, {
		name: "missing --",
		in: []string{
			"./foo", "./bar",
			"foo", "bar",
		},
		expectedErr: errNoScript,
	}}

	for i, test := range table {
		c, err := parseFlags(test.in)
		assert.Equal(t, test.expectedConfig, c, fmt.Sprintf("%s (%d): c", test.name, i))
		assert.True(t, errors.Is(err, test.expectedErr), fmt.Sprintf("%s (%d): err", test.name, i))
	}
}
