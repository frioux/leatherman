package west

import (
	_ "embed"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

//go:embed testdata/000.md
var _000 []byte

func TestParseBasic(t *testing.T) {
	d := Parse(_000)
	var linkCount int
	Walk(d, func(n Node) error {
		if _, ok := n.(*Link); ok {
			linkCount += 1
		}
		return nil
	})
	testutil.Equal(t, string(d.Markdown()), string(_000), "roundtrips")
	if linkCount != 3 {
		t.Error("expected link count of 3")
	}
}

func TestParseEarlyExit(t *testing.T) {
	d := Parse(_000)
	var runCount int
	Walk(d, func(n Node) error {
		runCount += 1
		if l, ok := n.(*Link); ok {
			if l.HRef == "http://frew.co" {
				return WalkBreak
			}

			if l.HRef == "http://afoolishmanifesto.com" {
				t.Error("walk didn't break")
			}
		}
		return nil
	})
}

func TestParseCodeFenceBlock(t *testing.T) {
	p := NewParser([]byte(`~~~
this a test
`))

	c := &CodeFenceBlock{}
	p.parseCodeFenceBlock(c)
	if c.body != "this a test\n" {
		t.Errorf("uhh: %q", c.body)
	}

	p = NewParser([]byte(`~~~
this a test 2
~~~
rest
`))

	p.parseCodeFenceBlock(c)
	if c.body != "this a test 2\n" {
		t.Errorf("uhh: %q", c.body)
	}

	if string(p.rest()) != "rest\n" {
		t.Errorf("why: %q", p.rest())
	}
}
