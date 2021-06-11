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

		l := &List{}
		if p.parseList(l) {
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
	*para = Inline{start: Pos(p.o)}

	text := &Text{start: Pos(p.o), end: Pos(p.o)}

	for {
		codeSpan := &InlineCode{}
		if p.parseCodeSpan(codeSpan) {
			if text.end != text.start {
				text.Text = string(p.b[int(text.start):int(text.end)])
				para.Nodes = append(para.Nodes, text)
			}
			para.Nodes = append(para.Nodes, codeSpan)
			text = &Text{start: Pos(p.o), end: Pos(p.o)}
			continue
		}

		link := &Link{}
		if p.parseLink(link) {
			if text.end != text.start {
				text.Text = string(p.b[int(text.start):int(text.end)])
				para.Nodes = append(para.Nodes, text)
			}
			para.Nodes = append(para.Nodes, link)
			text = &Text{start: Pos(p.o), end: Pos(p.o)}
			continue
		}

		if p.end() {
			break
		}

		if p.expect([]byte("\n\n")) {
			text.end += 2
			break
		}

		text.end++
		p.o++
	}

	if text.end != text.start {
		text.Text = string(p.b[int(text.start):int(text.end)])
		para.Nodes = append(para.Nodes, text)
	}

	return true
}

func (p *Parser) parseLinkBody(inline *Inline) bool {
	*inline = Inline{start: Pos(p.o)}

	text := &Text{start: Pos(p.o), end: Pos(p.o)}

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

		if p.peak([]byte("]")) {
			if text.end != text.start {
				text.Text = string(p.b[int(text.start):int(text.end)])
				inline.Nodes = append(inline.Nodes, text)
			}
			break
		}

		if p.end() {
			return false
		}

		text.end++
		p.o++
	}

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
	in := p.rest()

	if !p.expect([]byte("`")) {
		return false
	}

	// find newline so we stop looking after we see it
	nlI := bytes.IndexRune(in[1:], rune('\n'))
	if nlI == -1 {
		nlI = len(in)
	}
	// look within that byte slice
	in = in[1:nlI]
	*codeSpan = InlineCode{start: Pos(p.o)}
	i := bytes.IndexRune(in[1:], rune('`'))
	if i == -1 {
		return false
	}

	p.o += 1 /* ` */ + i + 1 /* ` */
	codeSpan.end = codeSpan.start + Pos(i)
	codeSpan.Text = string(in[:i+1])

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
		n.o += len(m[0]) + offsets[1]
	}

	p.o = n.o
	return true
}
