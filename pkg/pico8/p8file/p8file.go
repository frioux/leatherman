// Package p8file implements a parser and formatter for
// https://pico-8.fandom.com/wiki/P8FileFormat.
package p8file

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"strconv"
	"strings"
)

// palatte was generated from https://pico-8.fandom.com/wiki/Palette
var palette = color.Palette([]color.Color{
	0:  color.NRGBA{0x00, 0x00, 0x00, 0xFF}, // black
	1:  color.NRGBA{0x1D, 0x2B, 0x53, 0xFF}, // dark-blue
	2:  color.NRGBA{0x7E, 0x25, 0x53, 0xFF}, // dark-purple
	3:  color.NRGBA{0x00, 0x87, 0x51, 0xFF}, // dark-green
	4:  color.NRGBA{0xAB, 0x52, 0x36, 0xFF}, // brown
	5:  color.NRGBA{0x5F, 0x57, 0x4F, 0xFF}, // dark-grey
	6:  color.NRGBA{0xC2, 0xC3, 0xC7, 0xFF}, // light-grey
	7:  color.NRGBA{0xFF, 0xF1, 0xE8, 0xFF}, // white
	8:  color.NRGBA{0xFF, 0x00, 0x4D, 0xFF}, // red
	9:  color.NRGBA{0xFF, 0xA3, 0x00, 0xFF}, // orange
	10: color.NRGBA{0xFF, 0xEC, 0x27, 0xFF}, // yellow
	11: color.NRGBA{0x00, 0xE4, 0x36, 0xFF}, // green
	12: color.NRGBA{0x29, 0xAD, 0xFF, 0xFF}, // blue
	13: color.NRGBA{0x83, 0x76, 0x9C, 0xFF}, // lavender
	14: color.NRGBA{0xFF, 0x77, 0xA8, 0xFF}, // pink
	15: color.NRGBA{0xFF, 0xCC, 0xAA, 0xFF}, // light-peach
})

func (c Cart) SpriteImage() *image.Paletted {
	p := image.NewPaletted(
		image.Rect(0, 0, 128, 128),
		palette,
	)

	for i, v := range c.Spritesheet {
		p.Pix[i] = uint8(v)
	}

	return p
}

func (c Cart) LabelImage() *image.Paletted {
	p := image.NewPaletted(
		image.Rect(0, 0, 128, 128),
		palette,
	)

	for i, v := range c.Label {
		p.Pix[i] = uint8(v)
	}

	return p
}

func (c Cart) SpriteAt(i uint8) *image.Paletted {
	i_ := int(i)
	r := image.Rect(0, 0, 8, 8)
	p := image.NewPaletted(r, palette)

	if i == 0 {
		return p
	}

	raw_x := i_ % 16
	x := 8 * raw_x
	y := 8 * ((i_ - raw_x) / 16)

	draw.Draw(p, r, c.SpriteImage(), image.Point{x, y}, draw.Over)

	return p
}

func (c Cart) MapImage() *image.Paletted {
	p := image.NewPaletted(image.Rect(0, 0, 128*8, 64*8), palette)

	for i, v := range c.Map {
		raw_x := i % 128
		x := 8 * raw_x
		y := 8 * ((i - raw_x) / 128)

		draw.Draw(p, image.Rect(x, y, x+8, y+8), c.SpriteAt(uint8(v)), image.Point{}, draw.Over)
	}

	return p
}

type Cart struct {
	Version      int
	Lua          []byte
	Spritesheet  []byte
	Spriteflags  []byte
	Label        []byte
	Map          []byte
	SoundEffects []byte
	Music        []byte
}

type p8ParseState int

const (
	parseStart    p8ParseState = iota // just started
	parsedPico8                       // parsed pico 8 line "pico-8 cartridge // http://www.pico-8.com"
	parsedVersion                     // parsed version "version 8"
	parsingLua                        // parsing __lua__ section
	parsingGFX                        // parsing __gfx__ section
	parsingGFF                        // parsing __gff__ section
	parsingLabel                      // parsing __label__ section
	parsingMap                        // parsing __map__ section
	parsingSFX                        // parsing __sfx__ section
	parsingMusic                      // parsing __music__ section

	sectionLua   = "__lua__"
	sectionGFX   = "__gfx__"
	sectionGFF   = "__gff__"
	sectionLabel = "__label__"
	sectionMap   = "__map__"
	sectionSFX   = "__sfx__"
	sectionMusic = "__music__"
)

var spritebytemap = []byte{
	int('0'): 0,
	int('1'): 1,
	int('2'): 2,
	int('3'): 3,
	int('4'): 4,
	int('5'): 5,
	int('6'): 6,
	int('7'): 7,
	int('8'): 8,
	int('9'): 9,
	int('a'): 10,
	int('b'): 11,
	int('c'): 12,
	int('d'): 13,
	int('e'): 14,
	int('f'): 15,
}

func Parse(r io.Reader) (Cart, error) {
	cart := Cart{}

	s := bufio.NewScanner(r)

	var state p8ParseState

	for s.Scan() {
		l := s.Text()
		switch state {
		case parseStart:
			if l == "pico-8 cartridge // http://www.pico-8.com" {
				state = parsedPico8
				continue
			}
		case parsedPico8:
			preLen := len(l)
			if suff := strings.TrimPrefix(l, "version "); len(suff) < preLen {
				ver, err := strconv.Atoi(suff)
				if err != nil {
					return Cart{}, fmt.Errorf("parsing version: %w", err)
				}
				cart.Version = ver
				state = parsedVersion
				continue
			}
		case parsedVersion:
			switch l {
			case sectionLua:
				state = parsingLua
				continue
			case sectionGFX:
				state = parsingGFX
				continue
			case sectionGFF:
				state = parsingGFF
				continue
			case sectionLabel:
				state = parsingLabel
				continue
			case sectionMap:
				state = parsingMap
				continue
			case sectionSFX:
				state = parsingSFX
				continue
			case sectionMusic:
				state = parsingMusic
				continue
			default:
				return Cart{}, fmt.Errorf("unexpected: %s", l)
			}
		case parsingLua:
			switch l {
			case sectionLua:
				state = parsingLua
				continue
			case sectionGFX:
				state = parsingGFX
				continue
			case sectionGFF:
				state = parsingGFF
				continue
			case sectionLabel:
				state = parsingLabel
				continue
			case sectionMap:
				state = parsingMap
				continue
			case sectionSFX:
				state = parsingSFX
				continue
			case sectionMusic:
				state = parsingMusic
				continue
			default:
				cart.Lua = append(cart.Lua, []byte(l)...)
				cart.Lua = append(cart.Lua, []byte("\n")...)
			}
		case parsingGFX:
			switch l {
			case sectionLua:
				state = parsingLua
				continue
			case sectionGFX:
				state = parsingGFX
				continue
			case sectionGFF:
				state = parsingGFF
				continue
			case sectionLabel:
				state = parsingLabel
				continue
			case sectionMap:
				state = parsingMap
				continue
			case sectionSFX:
				state = parsingSFX
				continue
			case sectionMusic:
				state = parsingMusic
				continue
			default:
				b := []byte(l)
				for i, v := range b {
					b[i] = spritebytemap[v]
				}
				cart.Spritesheet = append(cart.Spritesheet, b...)
			}
		case parsingGFF:
			switch l {
			case sectionLua:
				state = parsingLua
				continue
			case sectionGFX:
				state = parsingGFX
				continue
			case sectionGFF:
				state = parsingGFF
				continue
			case sectionLabel:
				state = parsingLabel
				continue
			case sectionMap:
				state = parsingMap
				continue
			case sectionSFX:
				state = parsingSFX
				continue
			case sectionMusic:
				state = parsingMusic
				continue
			default:
				b, err := hex.DecodeString(l)
				if err != nil {
					return Cart{}, fmt.Errorf("parsing __gff__ section: %w", err)
				}
				cart.Spriteflags = append(cart.Spriteflags, b...)
			}
		case parsingLabel:
			switch l {
			case sectionLua:
				state = parsingLua
				continue
			case sectionGFX:
				state = parsingGFX
				continue
			case sectionGFF:
				state = parsingGFF
				continue
			case sectionLabel:
				state = parsingLabel
				continue
			case sectionMap:
				state = parsingMap
				continue
			case sectionSFX:
				state = parsingSFX
				continue
			case sectionMusic:
				state = parsingMusic
				continue
			default:
				b := []byte(l)
				for i, v := range b {
					b[i] = spritebytemap[v]
				}
				cart.Label = append(cart.Label, b...)
			}
		case parsingMap:
			switch l {
			case sectionLua:
				state = parsingLua
				continue
			case sectionGFX:
				state = parsingGFX
				continue
			case sectionGFF:
				state = parsingGFF
				continue
			case sectionLabel:
				state = parsingLabel
				continue
			case sectionMap:
				state = parsingMap
				continue
			case sectionSFX:
				state = parsingSFX
				continue
			case sectionMusic:
				state = parsingMusic
				continue
			default:
				b, err := hex.DecodeString(l)
				if err != nil {
					return Cart{}, fmt.Errorf("parsing __map__ section: %w", err)
				}
				cart.Map = append(cart.Map, b...)
			}
		case parsingSFX:
			switch l {
			case sectionLua:
				state = parsingLua
				continue
			case sectionGFX:
				state = parsingGFX
				continue
			case sectionGFF:
				state = parsingGFF
				continue
			case sectionLabel:
				state = parsingLabel
				continue
			case sectionMap:
				state = parsingMap
				continue
			case sectionSFX:
				state = parsingSFX
				continue
			case sectionMusic:
				state = parsingMusic
				continue
			default:
				b, err := hex.DecodeString(l)
				if err != nil {
					return Cart{}, fmt.Errorf("parsing __sfx__ section: %w", err)
				}
				cart.SoundEffects = append(cart.SoundEffects, b...)
			}
		case parsingMusic:
			switch l {
			case sectionLua:
				state = parsingLua
				continue
			case sectionGFX:
				state = parsingGFX
				continue
			case sectionGFF:
				state = parsingGFF
				continue
			case sectionLabel:
				state = parsingLabel
				continue
			case sectionMap:
				state = parsingMap
				continue
			case sectionSFX:
				state = parsingSFX
				continue
			case sectionMusic:
				state = parsingMusic
				continue
			default:
			b, err := hex.DecodeString(strings.ReplaceAll(l, " ", ""))
			if err != nil {
				return Cart{}, fmt.Errorf("parsing __music__ section: %w", err)
			}
			cart.Music = append(cart.Music, b...)
			}
		}
	}

	if err := s.Err(); err != nil {
		return Cart{}, err
	}

	return cart, nil
}
