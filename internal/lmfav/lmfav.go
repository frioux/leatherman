// lmfav provides http.Handlers for generating favicons.
package lmfav

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"
)

//go:generate ./genpal
var palettes [][]color.NRGBA

func twoHoriz(pic *image.NRGBA, a, b, _, _ color.NRGBA) {
	for x := 0; x < 16; x++ {
		for y := 0; y < 8; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := 0; x < 16; x++ {
		for y := 8; y < 16; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
}

func threeHoriz(pic *image.NRGBA, a, b, c, _ color.NRGBA) {
	for x := 0; x < 16; x++ {
		for y := 0; y < 5; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := 0; x < 16; x++ {
		for y := 5; y < 10; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
	for x := 0; x < 16; x++ {
		for y := 10; y < 15; y++ {
			pic.SetNRGBA(x, y, c)
		}
	}
	for x := 0; x < 16; x++ {
		pic.SetNRGBA(x, 16, color.NRGBA{0, 0, 0, 0})
	}
}

func fourHoriz(pic *image.NRGBA, a, b, c, d color.NRGBA) {
	for x := 0; x < 16; x++ {
		for y := 0; y < 4; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := 0; x < 16; x++ {
		for y := 4; y < 8; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
	for x := 0; x < 16; x++ {
		for y := 8; y < 12; y++ {
			pic.SetNRGBA(x, y, c)
		}
	}
	for x := 0; x < 16; x++ {
		for y := 12; y < 16; y++ {
			pic.SetNRGBA(x, y, d)
		}
	}
}

func twoVert(pic *image.NRGBA, a, b, _, _ color.NRGBA) {
	for x := 0; x < 8; x++ {
		for y := 0; y < 16; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := 8; x < 16; x++ {
		for y := 0; y < 16; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
}

func threeVert(pic *image.NRGBA, a, b, c, _ color.NRGBA) {
	for x := 0; x < 5; x++ {
		for y := 0; y < 16; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := 5; x < 10; x++ {
		for y := 0; y < 16; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
	for x := 10; x < 15; x++ {
		for y := 0; y < 16; y++ {
			pic.SetNRGBA(x, y, c)
		}
	}
	for y := 0; y < 16; y++ {
		pic.SetNRGBA(16, y, color.NRGBA{0, 0, 0, 0})
	}
}

func fourVert(pic *image.NRGBA, a, b, c, d color.NRGBA) {
	for x := 0; x < 4; x++ {
		for y := 0; y < 16; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := 4; x < 8; x++ {
		for y := 0; y < 16; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
	for x := 8; x < 12; x++ {
		for y := 0; y < 16; y++ {
			pic.SetNRGBA(x, y, c)
		}
	}
	for x := 12; x < 16; x++ {
		for y := 0; y < 16; y++ {
			pic.SetNRGBA(x, y, d)
		}
	}
}

func fourSquare(pic *image.NRGBA, a, b, c, d color.NRGBA) {
	// top left
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	// bottom left
	for x := 0; x < 8; x++ {
		for y := 8; y < 16; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
	// top right
	for x := 8; x < 16; x++ {
		for y := 0; y < 8; y++ {
			pic.SetNRGBA(x, y, c)
		}
	}
	// bottom right
	for x := 8; x < 16; x++ {
		for y := 8; y < 16; y++ {
			pic.SetNRGBA(x, y, d)
		}
	}
}

var algos = [...]func(_ *image.NRGBA, _, _, _, _ color.NRGBA){
	0: twoHoriz,
	1: threeHoriz,
	2: fourHoriz,
	3: twoVert,
	4: threeVert,
	5: fourVert,
	6: fourSquare,
}

// Flag procedurally generates a 16x16 flag based on the Host of the request.
// Host can be overridden by passing a host in the query parameters.
func Flag() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		host := r.URL.Query().Get("host")
		if host == "" {
			host = r.Host
		}
		sum := sha256.Sum256([]byte(host))

		pic := image.NewNRGBA(image.Rect(0, 0, 16, 16))
		pal := palettes[int(sum[1])%len(palettes)]
		algos[int(sum[0])%len(algos)](
			pic,
			pal[int(sum[2])%len(pal)],
			pal[int(sum[3])%len(pal)],
			pal[int(sum[4])%len(pal)],
			pal[int(sum[5])%len(pal)],
		)

		rw.Header().Add("Content-Type", "image/png")
		if err := png.Encode(rw, pic); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})
}

// Emoji generates a favicon of the passed rune using SVG.
func Emoji(favicon rune) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprintf(rw, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y=".9em" font-size="90">%c</text></svg>`, favicon)
	})
}
