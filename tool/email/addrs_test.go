package email

import (
	"net/mail"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAllAddrs(t *testing.T) {
	message, err := mail.ReadMessage(strings.NewReader(`To: "A B" <a.b@c.com>
Date: Sat, 21 Dec 2012 00:59:51 +0000
Subject: Station

Foo bar`))
	if err != nil {
		assert.NoError(t, err, "Couldn't set up test")
		return
	}

	a := allAddrs(message)
	if assert.Equal(t, 1, len(a)) {
		assert.Equal(t, "A B", a[0].Name)
		assert.Equal(t, "a.b@c.com", a[0].Address)
	}
}
func TestScoreEmail(t *testing.T) {
	message, err := mail.ReadMessage(strings.NewReader(`To: "A B" <a.b@c.com>
Date: Sat, 21 Dec 2012 00:59:51 +0000
Subject: Station

Foo bar`))
	if err != nil {
		assert.NoError(t, err, "Coudln't set up test")
		return
	}

	m := newFrecencyMap()
	m["a.b@c.com"] = &addr{
		addr:  "a.b@c.com",
		name:  "A B",
		score: 0,
	}

	m.scoreEmail(message, time.Date(2012, 12, 12, 12, 12, 12, 12, time.UTC))
	assert.InDelta(t, 1.217, m["a.b@c.com"].score, 0.001, "Scored address")
}

func TestBuildAddrMap(t *testing.T) {
	m := buildAddrMap(strings.NewReader("\n" +
		"frew@me.com\tFrew\n" +
		"frew@me.com\tFrew2\n",
	))

	if assert.Equal(t, 1, len(m)) {
		a := m["frew@me.com"]
		if !assert.NotNil(t, a) {
			return
		}
		assert.Equal(t, "Frew", a.name)
		assert.Equal(t, "frew@me.com", a.addr)
		assert.Equal(t, float64(0), a.score)
	}
}

func TestSortAddrMap(t *testing.T) {
	m := newFrecencyMap()
	m["a@b.com"] = &addr{name: "a", addr: "a@b.com", score: 1}
	m["b@b.com"] = &addr{name: "b", addr: "b@b.com", score: -1}
	m["c@b.com"] = &addr{name: "c", addr: "c@b.com", score: 0}
	m["d@b.com"] = &addr{name: "d", addr: "d@b.com", score: 3}

	a := sortAddrMap(m)
	if assert.Equal(t, 4, len(a)) {
		assert.Equal(t, "d", a[0].name)
		assert.Equal(t, "a", a[1].name)
		assert.Equal(t, "c", a[2].name)
		assert.Equal(t, "b", a[3].name)
	}
}

func TestAddrString(t *testing.T) {
	a := &addr{name: "frew", addr: "a@b.com"}
	assert.Equal(t, "a@b.com\tfrew", a.String())
}
