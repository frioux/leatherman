package west

import (
	"bytes"
	"fmt"
	"regexp"
)

var (
	// matches ```ruby, 0=start, 1=end, 2=langstart, 3=langend
	codeFenceLang = regexp.MustCompile("^```" + `(\S+)?$`)

	// matches ```, 0=start, 1=end
	codeFence = regexp.MustCompile("^```$")

	// matches  * foo, 0=start, 1=end, 2=bulstart, 3=bulend, 4=instart, 5=inend
	bullet = regexp.MustCompile(`^(\s*[*-]\s+)(.*)`)

	// matches [link](/to_content), 0=start, 1=end, 2=instart, 3=inend, 4=linstart, 5=linend
	link = regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)

	// matches foo | bar | baz, 0=start, 1=end, TODO document the rest
	tableRow = regexp.MustCompile(`^([^|]+\|)+`)

	// matches `inline text`, 0=start, 1=end, 2=textstart, 3=textend
	inlineCode = regexp.MustCompile("`([^`]+)`")
)

type State int

const (
	raw State = iota
	list
	table
	code
)

func parseInnerLine(in []byte, startPos Pos) []Node {
	ret := []Node{}
	if offsets := link.FindAllSubmatchIndex(in, -1); offsets != nil {
		seen := 0
		for _, o := range offsets {
			if o[0] > seen {
				ret = append(ret, parseInnerLine(in[seen:o[0]], startPos+Pos(seen))...) // non-Link
			}
			ret = append(ret, &Link{start: startPos + Pos(o[0]), end: startPos + Pos(o[1]), Text: string(in[o[2]:o[3]]), HRef: string(in[o[4]:o[5]])})
			seen = o[1]
		}

		if offsets[len(offsets)-1][1] < len(in) {
			ret = append(ret, parseInnerLine(in[offsets[len(offsets)-1][1]:], startPos+Pos(offsets[len(offsets)-1][1]))...)
		}
	} else {
		ret = append(ret, &Text{start: startPos, end: startPos + Pos(len(in)), Text: string(in)})
	}

	return ret
}

func parseInline(in []byte, startPos Pos) *Inline {
	ret := &Inline{
		start: startPos,
		end:   startPos + Pos(len(in)),
		Nodes: []Node{},
	}

	if offsets := inlineCode.FindAllSubmatchIndex(in, -1); offsets != nil {
		seen := 0
		for _, o := range offsets {
			if o[0] > seen {
				ret.Nodes = append(ret.Nodes, parseInnerLine(in[seen:o[0]], startPos+Pos(seen))...) // non-Code
			}
			ret.Nodes = append(ret.Nodes, &InlineCode{start: startPos + Pos(o[0]), end: startPos + Pos(o[1]), Text: string(in[o[2]:o[3]])})
			seen = o[1]
		}

		if offsets[len(offsets)-1][1] < len(in) {
			ret.Nodes = append(ret.Nodes, parseInnerLine(in[offsets[len(offsets)-1][1]:], startPos+Pos(offsets[len(offsets)-1][1]))...)
		}
	} else {
		ret.Nodes = append(ret.Nodes, parseInnerLine(in, startPos)...)
	}

	// 1. find `inline preformatted chunks` and call parseInline for the remainder
	// 2. find [links](/to_stuff) and call the remainder text
	return ret
}

func Parse(in []byte) *Document {
	var (
		curPos Pos
		state  State
		ret    = &Document{}
		lines  = bytes.Split(in, []byte("\n"))
	)

	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	for _, line := range lines {
		if state == code {
			cfb := ret.Nodes[len(ret.Nodes)-1].(*CodeFenceBlock)
			if offsets := codeFence.FindIndex(line); offsets != nil {
				state = 0
				// closing ```
				cfb.end = curPos + Pos(offsets[1]+1)
				ret.end = cfb.end
			} else {
				cfb.body += string(line) + "\n"
			}
		} else {
			if offsets := codeFenceLang.FindSubmatchIndex(line); offsets != nil {
				// opening ```
				state = code
				ret.Nodes = append(ret.Nodes, &CodeFenceBlock{
					start: curPos + Pos(offsets[0]),
					lang:  string(line[offsets[2]:offsets[3]]),
				})
			} else {
				ret.Nodes = append(ret.Nodes, parseInline(append(line, byte('\n')), curPos))
				ret.end = curPos + Pos(len(line)+1)
			}
		}
		curPos += Pos(len(line) + 1 /* \n */)
	}

	if debug {
		if len(in) != int(ret.end) {
			panic(fmt.Sprintf("length of bytes (%d) doesn't match final position (%d)", len(in), int(ret.end)))
		}
	}

	return ret
}
