// note ast is noteast is not east is west
package west

import (
	"errors"
	"strconv"
)

var Debug string

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
	ret := []byte{}

	for _, d := range d.Nodes {
		ret = append(ret, d.Markdown()...)
	}

	if Debug != "" {
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
func (d *Document) End() Pos         { return d.start }

type Header struct {
	start, end Pos

	Level int // h1 => Level: 1

	*Inline
}

func (h *Header) Start() Pos { return h.start }
func (h *Header) End() Pos   { return h.start }
func (h *Header) Markdown() []byte {
	ret := make([]byte, h.Level+1)
	for i := 0; i < h.Level; i++ {
		ret[i] = '#'
	}

	ret[h.Level] = ' '
	ret = append(ret, h.Inline.Markdown()...)
	ret = append(ret, '\n')

	if Debug != "" {
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
func (i *Inline) End() Pos         { return i.start }
func (i *Inline) Children() []Node { return i.Nodes }
func (i *Inline) Markdown() []byte {
	ret := []byte{}

	for _, i := range i.Nodes {
		ret = append(ret, i.Markdown()...)
	}

	if Debug != "" {
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

	Text, HRef string
}

func (l *Link) Start() Pos       { return l.start }
func (l *Link) End() Pos         { return l.end }
func (l *Link) Markdown() []byte { return []byte("[" + l.Text + "](" + l.HRef + ")") }

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
	ret := []byte{}

	for _, l := range l.ListItems {
		ret = append(ret, l.Markdown()...)
	}

	return ret
}

type ListItem struct {
	start, end Pos

	Prefix string
	*Inline
}

func (l *ListItem) Start() Pos { return l.start }
func (l *ListItem) End() Pos   { return l.end }
func (l *ListItem) Markdown() []byte {
	return append(append([]byte(l.Prefix), l.Inline.Markdown()...), '\n')
}

// TODO: CodeIndentBlock

type CodeFenceBlock struct {
	start, end Pos

	lang, body string
}

func (b *CodeFenceBlock) Start() Pos { return b.start }
func (b *CodeFenceBlock) End() Pos   { return b.end }
func (b *CodeFenceBlock) Markdown() []byte {
	return []byte("```" + b.lang + "\n" + b.body + "```\n")
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
