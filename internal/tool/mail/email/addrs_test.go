package email

import (
	"net/mail"
	"strings"
	"testing"
	"time"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestAllAddrs(t *testing.T) {
	t.Parallel()

	message, err := mail.ReadMessage(strings.NewReader(`To: "A B" <a.b@c.com>
Date: Sat, 21 Dec 2012 00:59:51 +0000
Subject: Station

Foo bar`))
	if err != nil {
		t.Errorf("Couldn't set up test: %s", err)
		return
	}

	a := allAddrs(message)
	if testutil.Equal(t, len(a), 1, "wrong address count") {
		testutil.Equal(t, a[0].Name, "A B", "wrong name")
		testutil.Equal(t, a[0].Address, "a.b@c.com", "wrong address")
	}
}
func TestScoreEmail(t *testing.T) {
	t.Parallel()

	message, err := mail.ReadMessage(strings.NewReader(`To: "A B" <a.b@c.com>
Date: Sat, 21 Dec 2012 00:59:51 +0000
Subject: Station

Foo bar`))
	if err != nil {
		t.Errorf("Couldn't set up test: %s", err)
		return
	}

	m := newFrecencyMap()
	m["a.b@c.com"] = &addr{
		addr:  "a.b@c.com",
		name:  "A B",
		score: 0,
	}

	m.scoreEmail(message, time.Date(2012, 12, 12, 12, 12, 12, 12, time.UTC))
	if s := m["a.b@c.com"].score; s > 1.217+0.001 || s < 1.217-0.001 {
		t.Errorf("Score should have been about 1.217, was %f", s)
	}
}

func TestBuildAddrMap(t *testing.T) {
	t.Parallel()

	m := buildAddrMap(strings.NewReader("\n" +
		"frew@me.com\tFrew\n" +
		"frew@me.com\tFrew2\n",
	))

	if testutil.Equal(t, len(m), 1, "incorrectly sized map") {
		a, ok := m["frew@me.com"]
		if !ok {
			t.Errorf("didn't find address frew@me.com")
			return
		}
		testutil.Equal(t, a.name, "Frew", "incorrect name")
		testutil.Equal(t, a.addr, "frew@me.com", "incorrect address")
		testutil.Equal(t, a.score, float64(0), "incorrect score")
	}
}

func TestSortAddrMap(t *testing.T) {
	t.Parallel()

	m := newFrecencyMap()
	m["a@b.com"] = &addr{name: "a", addr: "a@b.com", score: 1}
	m["b@b.com"] = &addr{name: "b", addr: "b@b.com", score: -1}
	m["c@b.com"] = &addr{name: "c", addr: "c@b.com", score: 0}
	m["d@b.com"] = &addr{name: "d", addr: "d@b.com", score: 3}

	a := sortAddrMap(m)
	if testutil.Equal(t, len(a), 4, "incorrectly sized slice") {
		testutil.Equal(t, a[0].name, "d", "wrong name")
		testutil.Equal(t, a[1].name, "a", "wrong name")
		testutil.Equal(t, a[2].name, "c", "wrong name")
		testutil.Equal(t, a[3].name, "b", "wrong name")
	}
}

func TestAddrString(t *testing.T) {
	t.Parallel()

	a := &addr{name: "frew", addr: "a@b.com"}
	testutil.Equal(t, a.String(), "a@b.com\tfrew", "wrong mutt address format")
}
