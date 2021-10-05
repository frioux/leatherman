package west

import (
	"bytes"
	"regexp"
	"strconv"
)

func Parse(b []byte) *Document {
	p := NewParser(b)
	d := &Document{}
	p.Parse(d)
	return d
}

type Parser struct {
	b []byte
	o int
}

func NewParser(b []byte) *Parser { return &Parser{b: b} }

func (p *Parser) copy() *Parser { return &Parser{p.b, p.o} }

func (p *Parser) expect(prefix []byte) bool {
	if p.peak(prefix) {
		p.o += len(prefix)
		return true
	}

	return false
}

func (p *Parser) peak(prefix []byte) bool { return bytes.HasPrefix(p.rest(), prefix) }

func (p *Parser) end() bool { return p.o == len(p.b) }

func (p *Parser) loadUntil(suffix []byte, found *[]byte) bool {
	i := bytes.Index(p.rest(), suffix)
	if i == -1 {
		return false
	}

	*found = p.rest()[:i]
	p.o += i + len(suffix)
	return true
}

func (p *Parser) rest() []byte { return p.b[p.o:] }

func (p *Parser) Parse(d *Document) bool {
	*d = Document{}

	for p.o != len(p.b) {
		b := &CodeFenceBlock{}
		if p.parseCodeFenceBlock(b) {
			d.Nodes = append(d.Nodes, b)
			continue
		}

		h := &Header{}
		if p.parseHeader(h) {
			d.Nodes = append(d.Nodes, h)
			continue
		}

		l := &List{}
		if p.parseList(l) {
			d.Nodes = append(d.Nodes, l)
			continue
		}

		t := &Table{}
		if p.parseTable(t) {
			d.Nodes = append(d.Nodes, l)
			continue
		}

		para := &Inline{}
		if p.parseParagraph(para) {
			d.Nodes = append(d.Nodes, para)
			continue
		}
	}

	return true
}

func (p *Parser) parseList(l *List) bool { return false }

func (p *Parser) parseParagraph(para *Inline) bool {
	p.parseInline(para, []string{"\n\n"})
	p.inlineConsume(para, []string{"\n\n"})
	return true
}

func (p *Parser) parseLinkBody(inline *Inline) bool {
	p.parseInline(inline, []string{"]"})

	return true
}

func (p *Parser) parseLink(link *Link) bool {
	*link = Link{Body: &Inline{}}
	n := p.copy()
	if !n.expect([]byte("[")) {
		return false
	}
	if !n.parseLinkBody(link.Body) {
		return false
	}
	if !n.expect([]byte("](")) {
		return false
	}
	url := []byte{}
	if !n.loadUntil([]byte(")"), &url) {
		return false
	}
	link.HRef = string(url)
	p.o = n.o
	return true
}

func (p *Parser) parseCodeSpan(codeSpan *InlineCode) bool {
	n := p.copy()

	in := n.rest()

	if !n.expect([]byte("`")) {
		return false
	}

	if len(n.rest()) == 0 {
		return false
	}

	// find newline so we stop looking after we see it
	nlI := bytes.IndexRune(in[1:], rune('\n'))
	if nlI == -1 {
		nlI = len(in)
	}

	if nlI <= 1 {
		return false
	}

	// look within that byte slice
	in = in[1:nlI]
	*codeSpan = InlineCode{start: Pos(n.o)}
	i := bytes.IndexRune(in[1:], rune('`'))
	if i == -1 {
		return false
	}

	n.o += 1 /* ` */ + i + 1 /* ` */
	codeSpan.end = codeSpan.start + Pos(i)
	codeSpan.Text = string(in[:i+1])

	p.o = n.o
	return true
}

// matches ```ruby, 0=match, 1=fence, 2=lang
var codeFenceLang = regexp.MustCompile("^(~{3,}|`{3,})" + `(\S*)` + "\n")

func (p *Parser) parseCodeFenceBlock(cfb *CodeFenceBlock) bool {
	n := p.copy()

	in := n.rest()
	m := codeFenceLang.FindSubmatch(in)
	if m == nil {
		return false
	}
	*cfb = CodeFenceBlock{
		fence: string(m[1]),
		lang:  string(m[2]),
	}

	rein := "\n" + string(m[1][0]) + "{" + strconv.Itoa(len(m[1])) + ",}\n"
	body := in[len(m[0]):]
	re := regexp.MustCompile(rein)

	offsets := re.FindIndex(body)

	if offsets == nil {
		cfb.body = string(body)
		n.o += len(m[0]) + len(body)
	} else {
		cfb.body = string(body[:offsets[0]+1])
		cfb.endfence = string(body[offsets[0]+1 : offsets[1]])
		n.o += len(m[0]) + offsets[1]
	}

	p.o = n.o
	return true
}

func (p *Parser) parseTable(t *Table) bool {
	n := p.copy()

	*t = Table{start: Pos(p.o)}

	for {
		if t.Header == nil {
			row := &TableRow{}
			if !n.parseTableRow(row) {
				return false
			}

			t.Header = row
		} else if t.Delimiter == nil {
			row := &TableDelimiterRow{}
			if !n.parseTableDelimiterRow(row) {
				return false
			}
			t.Delimiter = row
		} else {
			row := &TableRow{}
			if !n.parseTableRow(row) {
				return false
			}
			t.Rows = append(t.Rows, row)
		}
		t.end = Pos(n.o)

		if (n.expect([]byte("\n")) || n.end()) && t.Delimiter != nil {
			break
		}
	}

	p.o = n.o
	return true
}

func (p *Parser) parseTableRow(r *TableRow) bool {
	n := p.copy()

	*r = TableRow{start: Pos(p.o)}

	for {
		cell := &Inline{}
		if !n.parseTableCell(cell) {
			return false
		}

		r.Columns = append(r.Columns, cell)
		r.end = Pos(n.o)

		if n.expect([]byte("|")) {
			continue
		}

		if n.expect([]byte("\n")) || n.end() {
			if len(r.Columns) < 2 {
				return false
			}
			if len(r.Columns) > 1 {
				break
			}
		}
	}

	p.o = n.o
	return true
}

func (p *Parser) parseTableDelimiterRow(r *TableDelimiterRow) bool {
	n := p.copy()

	*r = TableDelimiterRow{start: Pos(p.o)}

	for {
		cell := &Text{}
		if !n.parseTableDelimiterCell(cell) {
			return false
		}

		r.Delimiters = append(r.Delimiters, cell)
		r.end = Pos(n.o)

		if n.expect([]byte("|")) {
			continue
		}

		if n.expect([]byte("\n")) || n.end() {
			if len(r.Delimiters) < 2 {
				return false
			}
			if len(r.Delimiters) > 1 {
				break
			}
		}
	}

	p.o = n.o
	return true
}

var tableDelimiterCellMatcher = regexp.MustCompile(`^([\t ]*[:-]-+[:-][\t ]*)`)

func (p *Parser) parseTableDelimiterCell(c *Text) bool {
	*c = Text{start: Pos(p.o), end: Pos(p.o)}

	for {
		if p.peak([]byte("|")) || p.peak([]byte("\n")) || p.end() {
			if c.end != c.start {
				c.Text = string(p.b[int(c.start):int(c.end)])
				// I should refactor this to find the text
				// rather than verify it.
				if tableDelimiterCellMatcher.MatchString(c.Text) {
					break
				} else {
					return false
				}
			}
		}

		if p.end() {
			return false
		}

		c.end++
		p.o++
	}

	return true
}

func (p *Parser) parseTableCell(inline *Inline) bool {
	p.parseInline(inline, []string{"|", "\n"})

	return true
}

func (p *Parser) parseHeader(header *Header) bool {
	n := p.copy()

	*header = Header{start: Pos(n.o), end: Pos(n.o)}

	for {
		if !n.expect([]byte("#")) {
			if header.Level == 0 {
				return false
			} else {
				break
			}
		}

		header.Level++
		header.end++
	}

	if !n.peak([]byte(" ")) {
		return false
	}

	inline := &Inline{}
	header.Inline = inline
	n.parseInline(inline, []string{"\n\n"})
	n.inlineConsume(inline, []string{"\n\n"})

	p.o = n.o
	return true
}

func (p *Parser) parseInline(inline *Inline, terminators []string) {
	*inline = Inline{start: Pos(p.o), end: Pos(p.o)}
	text := &Text{start: Pos(p.o), end: Pos(p.o)}

parseInlineLoop:
	for {
		codeSpan := &InlineCode{}
		if p.parseCodeSpan(codeSpan) {
			if text.end != text.start {
				text.Text = string(p.b[int(text.start):int(text.end)])
				inline.Nodes = append(inline.Nodes, text)
			}
			inline.Nodes = append(inline.Nodes, codeSpan)
			text = &Text{start: Pos(p.o), end: Pos(p.o)}
			continue
		}

		link := &Link{}
		if p.parseLink(link) {
			if text.end != text.start {
				text.Text = string(p.b[int(text.start):int(text.end)])
				inline.Nodes = append(inline.Nodes, text)
			}
			inline.Nodes = append(inline.Nodes, link)
			text = &Text{start: Pos(p.o), end: Pos(p.o)}
			continue
		}

		if p.end() {
			break
		}

		for _, terminator := range terminators {
			if p.peak([]byte(terminator)) {
				break parseInlineLoop
			}
		}

		text.end++
		p.o++
	}

	if text.end != text.start {
		text.Text = string(p.b[int(text.start):int(text.end)])
		inline.Nodes = append(inline.Nodes, text)
	}
}

func (p *Parser) inlineConsume(inline *Inline, terminators []string) {
	if p.end() {
		return
	}

	for _, terminator := range terminators {
		if p.expect([]byte(terminator)) {
			var textNode *Text

			if len(inline.Nodes) > 0 {
				textNode, _ = inline.Nodes[len(inline.Nodes)-1].(*Text)
			}

			if textNode == nil {
				textNode = &Text{start: Pos(p.o - len([]byte(terminator))), end: Pos(p.o)}
				inline.Nodes = append(inline.Nodes, textNode)
			} else {
				textNode.end += Pos(len([]byte(terminator)))
			}

			textNode.Text = string(p.b[int(textNode.start):int(textNode.end)])

			return
		}
	}
}
