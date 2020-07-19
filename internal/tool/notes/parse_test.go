package notes

import (
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
	b, err := parseNow(strings.NewReader(eg), time.Date(2020, 7, 19, 2, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expect := `

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


`
	testutil.Equal(t, string(b), expect, "parseNow renders properly")
}
