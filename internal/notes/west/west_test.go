package west

import (
	_ "embed"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestAST(t *testing.T) {
	raw := []byte(`# a [b](/c) d

* 1
 * 2
* 3 [4](/5) 6
`)
	eg := &Document{
		start: 0,
		end:   Pos(len(raw)),

		Nodes: []Node{
			&Header{
				start: 0,
				end:   14,
				Level: 1,
				Inline: &Inline{
					start: 2,
					end:   13,
					Nodes: []Node{
						&Text{start: 2, end: 4, Text: "a "},
						&Link{start: 4, end: 10, Body: soloInline(&Text{Text: "b"}), HRef: "/c"},
						&Text{start: 10, end: 12, Text: " d"},
					},
				},
			},
			&Text{start: 15, end: 16, Text: "\n"},
			&List{
				start: 16,
				end:   37,
				ListItems: []*ListItem{
					{
						start:  16,
						end:    20,
						Prefix: "* ",
						Inline: soloInline(&Text{
							start: 19,
							end:   20,
							Text:  "1",
						}),
					},
					{
						start:  21,
						end:    26,
						Prefix: " * ",
						Inline: soloInline(&Text{
							start: 24,
							end:   25,
							Text:  "2",
						}),
					},
					{

						start:  27,
						end:    33,
						Prefix: "* ",
						Inline: &Inline{
							start: 30,
							end:   32,
							Nodes: []Node{
								&Text{start: 30, end: 32, Text: "3 "},
								&Link{start: 32, end: 38, Body: soloInline(&Text{Text: "4"}), HRef: "/5"},
								&Text{start: 38, end: 40, Text: " 6"},
							},
						},
					},
				},
			},
		},
	}

	testutil.Equal(t, string(eg.Markdown()), string(raw), "AST Roundtrips")
}

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
