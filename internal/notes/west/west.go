// note ast is noteast is not east is west
package west

import (
	"errors"
	"strconv"
)

func makeBytes(n Node) []byte {
	if n.End() < n.Start() {
		if debug {
			panic("impossible start/end")
		} else {
			return make([]byte, 0, 0)
		}
	}
	return make([]byte, 0, int(n.End()-n.Start()))
}

type Pos int

type Block interface {
	Node
	Children() []Node
}

type Node interface {
	Start() Pos
	End() Pos
	Markdown() []byte
}

type Document struct {
	start, end Pos

	Nodes []Node
}

func (d *Document) Markdown() []byte {
	ret := makeBytes(d)

	for _, d := range d.Nodes {
		ret = append(ret, d.Markdown()...)
	}

	if debug {
		got := len(ret)
		expect := int(d.end - d.start)
		if got != expect {
			panic("final len (" + strconv.Itoa(got) + ") not the expected value (" + strconv.Itoa(expect) + ")")
		}
	}

	return ret
}

func (d *Document) Children() []Node { return d.Nodes }
func (d *Document) Start() Pos       { return d.start }
func (d *Document) End() Pos         { return d.end }

type Header struct {
	start, end Pos

	Level int // h1 => Level: 1

	*Inline
}

func (h *Header) Start() Pos { return h.start }
func (h *Header) End() Pos   { return h.end }
func (h *Header) Markdown() []byte {
	ret := makeBytes(h)
	for i := 0; i < h.Level; i++ {
		ret = append(ret, '#')
	}

	ret = append(ret, h.Inline.Markdown()...)

	if debug {
		got := len(ret)
		expect := int(h.end - h.start)
		if got != expect {
			panic("header len (" + strconv.Itoa(got) + ") not the expected value (" + strconv.Itoa(expect) + "): " + string(ret))
		}
	}

	return ret
}

func soloInline(n Node) *Inline {
	return &Inline{
		start: n.Start(),
		end:   n.End(),
		Nodes: []Node{n},
	}
}

// Inline contains multiple inline elements
type Inline struct {
	start, end Pos

	Nodes []Node
}

func (i *Inline) Start() Pos       { return i.start }
func (i *Inline) End() Pos         { return i.end }
func (i *Inline) Children() []Node { return i.Nodes }
func (i *Inline) Markdown() []byte {
	ret := makeBytes(i)

	for _, i := range i.Nodes {
		ret = append(ret, i.Markdown()...)
	}

	if debug {
		got := len(ret)
		expect := int(i.end - i.start)
		if got != expect {
			panic("inline len (" + strconv.Itoa(got) + ") not the expected value (" + strconv.Itoa(expect) + "): " + string(ret))
		}
	}

	return ret
}

type InlineCode struct {
	start, end Pos

	Text string
}

func (c *InlineCode) Start() Pos       { return c.start }
func (c *InlineCode) End() Pos         { return c.end }
func (c *InlineCode) Markdown() []byte { return []byte("`" + c.Text + "`") }

type Text struct {
	start, end Pos

	Text string
}

func (t *Text) Start() Pos       { return t.start }
func (t *Text) End() Pos         { return t.end }
func (t *Text) Markdown() []byte { return []byte(t.Text) }

type Link struct {
	start, end Pos

	Body *Inline
	HRef string
}

func (l *Link) Start() Pos { return l.start }
func (l *Link) End() Pos   { return l.end }
func (l *Link) Markdown() []byte {
	return []byte("[" + string(l.Body.Markdown()) + "](" + l.HRef + ")")
}

type List struct {
	start, end Pos

	ListItems []*ListItem
}

func (l *List) Start() Pos { return l.start }
func (l *List) End() Pos   { return l.end }
func (l *List) Children() []Node {
	ret := make([]Node, len(l.ListItems))
	for i, v := range l.ListItems {
		ret[i] = v
	}
	return ret

}
func (l *List) Markdown() []byte {
	ret := makeBytes(l)

	for _, l := range l.ListItems {
		ret = append(ret, l.Markdown()...)
	}

	return ret
}

type ListItem struct {
	start, end Pos

	Prefix    string
	Level     int
	ListItems []*ListItem
	*Inline
}

func (l *ListItem) Start() Pos { return l.start }
func (l *ListItem) End() Pos   { return l.end }
func (l *ListItem) Markdown() []byte {
	ret := append(append([]byte(l.Prefix), l.Inline.Markdown()...), '\n')

	for _, l := range l.ListItems {
		ret = append(ret, l.Markdown()...)
	}

	return ret
}

type CodeFenceBlock struct {
	start, end Pos

	fence, lang, body, endfence string
}

func (b *CodeFenceBlock) Start() Pos { return b.start }
func (b *CodeFenceBlock) End() Pos   { return b.end }
func (b *CodeFenceBlock) Markdown() []byte {
	return []byte(b.fence + b.lang + "\n" + b.body + b.endfence)
}

// Table is defined at
// https://github.github.com/gfm/#tables-extension-
type Table struct {
	start, end Pos

	Header    *TableRow
	Delimiter *TableDelimiterRow
	Rows      []Node
}

func (t *Table) Start() Pos       { return t.start }
func (t *Table) End() Pos         { return t.end }
func (t *Table) Children() []Node { return t.Rows }
func (t *Table) Markdown() []byte {
	b := makeBytes(t)
	b = append(b, t.Header.Markdown()...)
	b = append(b, '\n')
	b = append(b, t.Delimiter.Markdown()...)
	b = append(b, '\n')
	for i, row := range t.Rows {
		b = append(b, row.Markdown()...)
		if i != len(t.Rows)-1 {
			b = append(b, '\n')
		}
	}
	return b
}

type TableDelimiterRow struct {
	start, end Pos

	Delimiters []Node
}

func (r *TableDelimiterRow) Start() Pos       { return r.start }
func (r *TableDelimiterRow) End() Pos         { return r.end }
func (r *TableDelimiterRow) Children() []Node { return r.Delimiters }
func (r *TableDelimiterRow) Markdown() []byte {
	b := makeBytes(r)
	for i, cell := range r.Delimiters {
		b = append(b, cell.Markdown()...)
		if i != len(r.Delimiters)-1 {
			b = append(b, '|')
		}
	}
	return b
}

type TableRow struct {
	start, end Pos

	Columns []Node
}

func (r *TableRow) Start() Pos       { return r.start }
func (r *TableRow) End() Pos         { return r.end }
func (r *TableRow) Children() []Node { return r.Columns }
func (r *TableRow) Markdown() []byte {
	b := makeBytes(r)
	for i, cell := range r.Columns {
		b = append(b, cell.Markdown()...)
		if i != len(r.Columns)-1 {
			b = append(b, '|')
		}
	}
	return b
}

var (
	WalkBreak     = errors.New("walk break")
	WalkNoRecurse = errors.New("walk no recurse")
)

type WalkFn func(Node) error

func Walk(n Node, w WalkFn) error {
	err := walk(n, w)
	switch err {
	case WalkBreak, WalkNoRecurse:
		return nil
	}
	return err
}

func walk(n Node, w WalkFn) error {
	err := w(n)
	switch err {
	case WalkBreak:
		return WalkBreak
	case WalkNoRecurse:
		return nil
	}
	if err != nil {
		return err
	}

	if b, ok := n.(Block); ok {
		for _, m := range b.Children() {
			err := walk(m, w)
			switch err {
			case WalkBreak:
				return WalkBreak
			case WalkNoRecurse:
				return nil
			default:
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}
