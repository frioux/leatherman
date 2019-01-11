package minotaur

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	type row struct {
		name string

		in []string

		expectedDirs, expectedCmd []string
		expectedErr               error
	}

	var table = []row{{
		name: "simple and correct",
		in: []string{
			"./foo", "./bar",
			"--",
			"foo", "bar",
		},
		expectedDirs: []string{"./foo", "./bar"},
		expectedCmd:  []string{"foo", "bar"},
	}, {
		name: "missing --",
		in: []string{
			"./foo", "./bar",
			"foo", "bar",
		},
		expectedErr: errNoScript,
	}, {
		name: "early --",
		in: []string{
			"--",
			"./foo", "./bar",
			"foo", "bar",
		},
		expectedErr: errNoDirs,
	}}

	for i, test := range table {
		dirs, cmd, err := parseFlags(test.in)
		assert.Equal(t, test.expectedDirs, dirs, fmt.Sprintf("%s (%d): dirs", test.name, i))
		assert.Equal(t, test.expectedCmd, cmd, fmt.Sprintf("%s (%d): cmd", test.name, i))
		assert.Equal(t, test.expectedErr, errors.Cause(err), fmt.Sprintf("%s (%d): err", test.name, i))
	}
}
