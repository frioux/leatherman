package minotaur

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
	"github.com/google/go-cmp/cmp"
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

	opt := cmp.Comparer(func(x, y config) bool {
		var sameInclude bool
		if x.include == nil && y.include == nil {
			sameInclude = true
		} else if x.include.String() == y.include.String() {
			sameInclude = true
		}

		var sameIgnore bool
		if x.ignore == nil && y.ignore == nil {
			sameIgnore = true
		} else if x.ignore.String() == y.ignore.String() {
			sameIgnore = true
		}

		return sameInclude && sameIgnore &&
			x.verbose == y.verbose &&
			cmp.Equal(x.dirs, y.dirs) &&
			cmp.Equal(x.script, y.script)
	})

	for i, test := range table {
		c, err := parseFlags(test.in)
		testutil.Equal(t, c, test.expectedConfig, fmt.Sprintf("%s (%d): c", test.name, i), opt)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("expected err to be %s, instead it's: %s", test.expectedErr, err)
		}
	}
}
