package west

import (
	_ "embed"
	"fmt"
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
						&Link{start: 4, end: 10, Text: "b", HRef: "/c"},
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
								&Link{start: 32, end: 38, Text: "4", HRef: "/5"},
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

func TestParse(t *testing.T) {
	d := Parse(_000)
	var linkCount int
	Walk(d, func(n Node) error {
		if _, ok := n.(*Link); ok {
			linkCount += 1
		}
		return nil
	})
	if string(_000) != string(d.Markdown()) {
		t.Error("didn't roundtrip")
		fmt.Printf("raw:\n~~~\n%s\n~~~\n", _000)
		fmt.Printf("from ast:\n~~~\n%s\n~~~\n", d.Markdown())
	}
	if linkCount != 2 {
		t.Error("expected link count of 2")
	}
}
