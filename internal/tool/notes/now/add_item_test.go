package now

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestAddItem(t *testing.T) {
	type testCase struct {
		item string
		when time.Time
	}

	cases := []testCase{
		{item: "xyzzy"},
		{
			item: "create-section",
			when: time.Date(2020, 7, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, c := range cases {
		t.Run(c.item, func(t *testing.T) {
			expectFile, err := ioutil.ReadFile("testdata/add_item/" + c.item + ".txt")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if c.when.IsZero() {
				c.when = time.Date(2020, 7, 19, 0, 0, 0, 0, time.UTC)
			}

			b, err := addItem(strings.NewReader(eg), c.when, c.item)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			testutil.Equal(t, string(b), string(expectFile), "adding worked")
		})
	}

}
