package now

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/frioux/leatherman/internal/testutil"
)

var eg = `{
"title": "Now",
"tags": [ "private", "reference", "project" ],
"reviewed_on": "2020-07-15",
}

## Stash

 * foo
 * bar
 * baz

## 2020-07-19 ##

 * bong
 * biff
 * barp

## 2020-07-18 ##

 * ~~herp~~
 * ~~dong~~


`

func TestParseNow(t *testing.T) {
	type Test struct {
		in, expect string
		when       time.Time
	}

	tests := []Test{
		0: {
			when: time.Date(2020, 7, 19, 2, 0, 0, 0, time.UTC),
			in:   eg,
			expect: `

## Stash

 * foo
 * bar
 * baz

## 2020-07-19 ##

 * bong <form action="/toggle" method="POST"><input type="hidden" name="chunk" value="acf51c06e604dc806b4ec4d9f68371f5"><button>Toggle</button></form>
 * biff <form action="/toggle" method="POST"><input type="hidden" name="chunk" value="42d0788e089bcabc1b7fe94397f5de34"><button>Toggle</button></form>
 * barp <form action="/toggle" method="POST"><input type="hidden" name="chunk" value="d83809da49df6c78073909ee01a0a3dc"><button>Toggle</button></form>

<form action="/add-item" method="POST"><input type="input" name="item"><button>Add Item</button></form>

## 2020-07-18 ##

 * ~~herp~~
 * ~~dong~~


`,
		},
		1: {
			when: time.Date(2020, 7, 18, 2, 0, 0, 0, time.UTC),
			in:   eg,
			expect: `

## Stash

 * foo
 * bar
 * baz

## 2020-07-19 ##

 * bong
 * biff
 * barp

## 2020-07-18 ##

 * ~~herp~~ <form action="/toggle" method="POST"><input type="hidden" name="chunk" value="3ac3845115fb4ee703f3c170eb9ba368"><button>Toggle</button></form>
 * ~~dong~~ <form action="/toggle" method="POST"><input type="hidden" name="chunk" value="e23b23e871f27237c8d5a28960121cb7"><button>Toggle</button></form>


<form action="/add-item" method="POST"><input type="input" name="item"><button>Add Item</button></form>

`,
		},
		2: {
			when: time.Date(2020, 7, 20, 2, 0, 0, 0, time.UTC),
			in:   eg,
			expect: `

## Stash

 * foo
 * bar
 * baz

## 2020-07-20 ##

<form action="/add-item" method="POST"><input type="input" name="item"><button>Add Item</button></form>

## 2020-07-19 ##

 * bong
 * biff
 * barp

## 2020-07-18 ##

 * ~~herp~~
 * ~~dong~~


`,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, err := parseNow(strings.NewReader(tt.in), tt.when)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			testutil.Equal(t, string(b), tt.expect, "")
		})
	}
}
