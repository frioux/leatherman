package west

import (
	"strconv"
	"strings"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestParseCodeSpan(t *testing.T) {
	p := NewParser([]byte("`chill it out`"))
	cs := &InlineCode{}
	if !p.parseCodeSpan(cs) {
		t.Error("shoulda parsed")
	}
	if cs.Text != "chill it out" {
		t.Errorf("wrong text: %q", cs.Text)
	}
}

func TestParseLinkBody(t *testing.T) {
	in := "`code` foo `morse code`"
	p := NewParser([]byte(in + "]"))
	l := &Inline{}
	if !p.parseLinkBody(l) {
		t.Error("shoulda parsed")
		return
	}
	if !testutil.Equal(t, len(l.Nodes), 3, "node length") {
		return
	}

	testutil.Equal(t, l.Nodes[0].(*InlineCode).Text, "code", "0")
	testutil.Equal(t, l.Nodes[1].(*Text).Text, " foo ", "1")
	testutil.Equal(t, l.Nodes[2].(*InlineCode).Text, "morse code", "2")
	testutil.Equal(t, string(l.Markdown()), in, "roundtrip")

	p = NewParser([]byte("foo]"))
	l = &Inline{}
	if !p.parseLinkBody(l) {
		t.Error("shoulda parsed")
		return
	}
	if !testutil.Equal(t, len(l.Nodes), 1, "node length") {
		return
	}

	testutil.Equal(t, l.Nodes[0].(*Text).Text, "foo", "0")
}

func TestParseLink(t *testing.T) {
	in := "[foo](/bar)"
	p := NewParser([]byte(in))
	l := &Link{}
	if !p.parseLink(l) {
		t.Error("shoulda parsed")
		return
	}
	testutil.Equal(t, l.HRef, "/bar", "url")
	if !testutil.Equal(t, len(l.Body.Nodes), 1, "node length") {
		return
	}
	testutil.Equal(t, string(l.Markdown()), in, "roundtrip")
}

func TestParseTable(t *testing.T) {
	in := ` foo | bar | baz
--- | :-- | --:
buzz | borp | bapzinga`
	p := NewParser([]byte(in))
	tab := &Table{}
	if !p.parseTable(tab) {
		t.Error("shoulda parsed")
		return
	}

	testutil.Equal(t, len(tab.Rows), 1, "rowcount")
	testutil.Equal(t, len(tab.Rows[0].(*TableRow).Columns), 3, "column count")
	testutil.Equal(t, string(tab.Markdown()), in, "roundtrip")

	in = ` foo
biff
buzz`
	p = NewParser([]byte(in))
	tab = &Table{}
	if p.parseTable(tab) {
		t.Error("should not parse")
		return
	}
}

func TestParseParagraph(t *testing.T) {
	in := `
this is a test
of some [link3](/url3) words
woo` + " ` code is here! ` " + `
hoo
[link](/url)  [` + "`" + `link2` + "`" + `](/url2)
	`
	p := NewParser([]byte(in))
	l := &Inline{}
	if !p.parseParagraph(l) {
		t.Error("shoulda parsed")
		return
	}
	if !testutil.Equal(t, len(l.Nodes), 9, "nodes") {
		return
	}
	testutil.Equal(t, l.Nodes[0].(*Text).Text, "\nthis is a test\nof some ", "0")
	testutil.Equal(t, l.Nodes[1].(*Link).HRef, "/url3", "1")
	testutil.Equal(t, l.Nodes[2].(*Text).Text, " words\nwoo ", "2")
	testutil.Equal(t, l.Nodes[3].(*InlineCode).Text, " code is here! ", "3")
	testutil.Equal(t, l.Nodes[4].(*Text).Text, " \nhoo\n", "4")
	testutil.Equal(t, l.Nodes[5].(*Link).HRef, "/url", "5")
	testutil.Equal(t, l.Nodes[6].(*Text).Text, "  ", "6")
	testutil.Equal(t, l.Nodes[7].(*Link).HRef, "/url2", "7")
	testutil.Equal(t, l.Nodes[8].(*Text).Text, "\n\t", "8")
	testutil.Equal(t, string(l.Markdown()), in, "roundtrip")
}

func TestParseDocument(t *testing.T) {
	in := `
this is a test
of some [link3](/url3) words
woo` + " ` code is here! ` " + `
hoo
[link](/url)  [` + "`" + `link2` + "`" + `](/url2)

xyzzy
`
	p := NewParser([]byte(in))
	d := &Document{}
	if !p.Parse(d) {
		t.Error("shoulda parsed")
		return
	}
	testutil.Equal(t, string(d.Markdown()), in, "roundtrip")
}

func TestMutateDocument(t *testing.T) {
	in := `
 * [a](/a?x=1)
 * [c](/c?x=1)
 * [b](/b?x=1)
`
	p := NewParser([]byte(in))
	d := &Document{}
	if !p.Parse(d) {
		t.Error("shoulda parsed")
		return
	}
	Walk(d, func(n Node) error {
		l, ok := n.(*Link)
		if !ok {
			return nil
		}

		if l.Body.Nodes[0].(*Text).Text == "a" {
			l.HRef = l.HRef + "&y=xyzzy"
			return WalkBreak
		}

		return nil
	})

	testutil.Equal(t, string(d.Markdown()), `
 * [a](/a?x=1&y=xyzzy)
 * [c](/c?x=1)
 * [b](/b?x=1)
`, "link mutation")

	in = `[a](/a?x=1), [c](/c?x=1), [b](/b?x=1)`
	p = NewParser([]byte(in))
	d = &Document{}
	if !p.Parse(d) {
		t.Error("shoulda parsed")
		return
	}
	par := d.Nodes[0].(*Inline)
	par.Nodes[2], par.Nodes[4] = par.Nodes[4], par.Nodes[2]
	testutil.Equal(t, string(d.Markdown()), `[a](/a?x=1), [b](/b?x=1), [c](/c?x=1)`, "link mutation")
}

func TestParseHeader(t *testing.T) {
	in := "### station `foo` [bar](/baz)\n"
	p := NewParser([]byte(in))
	h := &Header{}
	if !p.parseHeader(h) {
		t.Error("shoulda parsed")
		return
	}
	testutil.Equal(t, h.Level, 3, "level")
	testutil.Equal(t, len(h.Inline.Nodes), 5, "node length")
	testutil.Equal(t, string(h.Markdown()), in, "roundtrip")
}

type listTestItem struct {
	prefix   string
	content  string
	children []listTestItem
}

func TestParseListTextOnly(t *testing.T) {
	type listTest struct {
		name     string
		markdown string
		items    []listTestItem
	}

	listTests := []listTest{
		{
			name: "asterisks",
			markdown: `
* asterisk 1
* asterisk 2
* asterisk 3
`,

			items: []listTestItem{
				{content: " asterisk 1", prefix: "*"},
				{content: " asterisk 2", prefix: "*"},
				{content: " asterisk 3", prefix: "*"},
			},
		},
		{
			name: "minuses",
			markdown: `
- minus 1
- minus 2
- minus 3
`,

			items: []listTestItem{
				{content: " minus 1", prefix: "-"},
				{content: " minus 2", prefix: "-"},
				{content: " minus 3", prefix: "-"},
			},
		},
		{
			name: "nested",
			markdown: `
* One
  * Two
    * Three
`,

			items: []listTestItem{
				{
					content: " One",
					prefix:  "*",
					children: []listTestItem{
						{
							content: " Two",
							prefix:  "  *",
							children: []listTestItem{
								{content: " Three", prefix: "    *"},
							},
						},
					},
				},
			},
		},
		{
			name: "nested_not_strictly_increasing",
			markdown: `
* One
  * Two
 * Three
`,

			items: []listTestItem{
				{
					content: " One",
					prefix:  "*",
					children: []listTestItem{
						{content: " Two", prefix: "  *"},
						{content: " Three", prefix: " *"},
					},
				},
			},
		},
		{
			name: "mixed_introducers",
			markdown: `
* One
- Two
* Three
`,

			items: []listTestItem{
				{content: " One", prefix: "*"},
				{content: " Two", prefix: "-"},
				{content: " Three", prefix: "*"},
			},
		},
	}

	for _, test := range listTests {
		t.Run(test.name, func(t *testing.T) {
			inputMarkdown := strings.TrimPrefix(test.markdown, "\n")
			p := NewParser([]byte(inputMarkdown))
			l := &List{}
			if !p.parseList(l) {
				t.Error("shoulda parsed")
				return
			}

			listsEqual(t, l.ListItems, test.items)

			testutil.Equal(t, string(l.Markdown()), inputMarkdown, "roundtrip")
		})
	}
}

func listsEqual(t *testing.T, gotItems []*ListItem, expectedItems []listTestItem) {
	testutil.Equal(t, len(gotItems), len(expectedItems), "item count")

	for i, gotItem := range gotItems {
		expectedItem := expectedItems[i]

		gotText := gotItem.Inline.Nodes[0].(*Text).Text

		testutil.Equal(t, gotItem.Prefix, expectedItem.prefix, "prefix matches")
		testutil.Equal(t, gotText, expectedItem.content, "content matches")

		listsEqual(t, gotItem.ListItems, expectedItem.children)
	}
}

func TestParseListEmbeddedMarkup(t *testing.T) {
	markdown := strings.Trim(`
* Foo [Example](https://example.com) Bar
`, "\n")

	p := NewParser([]byte(markdown))
	l := &List{}
	if !p.parseList(l) {
		t.Error("shoulda parsed")
		return
	}

	testutil.Equal(t, len(l.ListItems), 1, "item count")
	item := l.ListItems[0]
	testutil.Equal(t, len(item.Inline.Nodes), 3, "singleton item children count")

	firstText := item.Inline.Nodes[0].(*Text)
	onlyLink := item.Inline.Nodes[1].(*Link)
	secondText := item.Inline.Nodes[2].(*Text)

	testutil.Equal(t, firstText.Text, " Foo ", "first text chunk matches")
	testutil.Equal(t, onlyLink.HRef, "https://example.com", "link target matches")
	testutil.Equal(t, onlyLink.Body.Nodes[0].(*Text).Text, "Example", "link text matches")
	testutil.Equal(t, secondText.Text, " Bar", "second text chunk matches")
}

func TestListAmongParagraphs(t *testing.T) {
	markdown := strings.Trim(`
this is a test
of some [link3](/url3) words

* One
* Two
* Three

This is some more text
`, "\n")

	p := NewParser([]byte(markdown))
	doc := &Document{}
	if !p.Parse(doc) {
		t.Error("shoulda parsed")
		return
	}

	testutil.Equal(t, len(doc.Nodes), 3, "three kids")

	{
		_, isParagraph := doc.Nodes[0].(*Inline)
		if !isParagraph {
			t.Error("expected first child to be a paragraph")
		}
	}

	onlyList := doc.Nodes[1].(*List)
	testutil.Equal(t, len(onlyList.ListItems), 3, "three list items")
	testutil.Equal(t, onlyList.ListItems[0].Nodes[0].(*Text).Text, " One", "first list item")
	testutil.Equal(t, onlyList.ListItems[1].Nodes[0].(*Text).Text, " Two", "second list item")
	testutil.Equal(t, onlyList.ListItems[2].Nodes[0].(*Text).Text, " Three", "three list item")

	{
		_, isParagraph := doc.Nodes[2].(*Inline)
		if !isParagraph {
			t.Error("expected second child to be a paragraph")
		}
	}
}

var crashers = []string{
	0: "`",
	1: "`\n",
	2: "```\n",
	// 3: "|\n---|---",
}

func TestCrashers(t *testing.T) {
	for i, in := range crashers {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p := NewParser([]byte(in))
			d := &Document{}
			if !p.Parse(d) {
				t.Error("shoulda parsed")
				return
			}
			testutil.Equal(t, string(d.Markdown()), in, "roundtrip")
		})
	}
}
