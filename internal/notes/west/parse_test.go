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

func TestParseListTextOnly(t *testing.T) {
	type listTestItem struct {
		level   int
		prefix  string
		content string
	}

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
				{content: " asterisk 1", prefix: "*", level: 1},
				{content: " asterisk 2", prefix: "*", level: 1},
				{content: " asterisk 3", prefix: "*", level: 1},
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
				{content: " minus 1", prefix: "-", level: 1},
				{content: " minus 2", prefix: "-", level: 1},
				{content: " minus 3", prefix: "-", level: 1},
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
				{content: " One", prefix: "*", level: 1},
				{content: " Two", prefix: "  *", level: 2},
				{content: " Three", prefix: "    *", level: 3},
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
				{content: " One", prefix: "*", level: 1},
				{content: " Two", prefix: "  *", level: 2},
				{content: " Three", prefix: " *", level: 2},
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

			testutil.Equal(t, len(l.ListItems), len(test.items), "item count")

			for i, expectedItem := range test.items {
				gotItem := l.ListItems[i]

				gotText := gotItem.Inline.Nodes[0].(*Text).Text

				testutil.Equal(t, gotItem.Level, expectedItem.level, "level")
				testutil.Equal(t, gotItem.Prefix, expectedItem.prefix, "prefix")
				testutil.Equal(t, gotText, expectedItem.content, "content")
			}

			testutil.Equal(t, string(l.Markdown()), inputMarkdown, "roundtrip")
		})
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
