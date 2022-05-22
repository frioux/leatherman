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

const size = 128

func twoHoriz(pic *image.NRGBA, a, b, _, _ color.NRGBA) {
	for x := 0; x < size; x++ {
		for y := 0; y < size/2; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := 0; x < size; x++ {
		for y := size/2; y < size; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
}

func threeHoriz(pic *image.NRGBA, a, b, c, _ color.NRGBA) {
	for x := 0; x < size; x++ {
		for y := 0; y < size/3; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := 0; x < size; x++ {
		for y := size/3; y < 2*(size/3); y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
	for x := 0; x < size; x++ {
		for y := 2*(size/3); y < 3*(size/3); y++ {
			pic.SetNRGBA(x, y, c)
		}
	}
}

func fourHoriz(pic *image.NRGBA, a, b, c, d color.NRGBA) {
	for x := 0; x < size; x++ {
		for y := 0; y < size/4; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := 0; x < size; x++ {
		for y := size/4; y < size/2; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
	for x := 0; x < size; x++ {
		for y := size/2; y < 3*size/4; y++ {
			pic.SetNRGBA(x, y, c)
		}
	}
	for x := 0; x < size; x++ {
		for y := 3*size/4; y < size; y++ {
			pic.SetNRGBA(x, y, d)
		}
	}
}

func twoVert(pic *image.NRGBA, a, b, _, _ color.NRGBA) {
	for x := 0; x < size/2; x++ {
		for y := 0; y < size; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := size/2; x < size; x++ {
		for y := 0; y < size; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
}

func threeVert(pic *image.NRGBA, a, b, c, _ color.NRGBA) {
	for x := 0; x < size/3; x++ {
		for y := 0; y < size; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := size/3; x < 2*(size/3); x++ {
		for y := 0; y < size; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
	for x := 2*(size/3); x < 3*(size/3); x++ {
		for y := 0; y < size; y++ {
			pic.SetNRGBA(x, y, c)
		}
	}
}

func fourVert(pic *image.NRGBA, a, b, c, d color.NRGBA) {
	for x := 0; x < size/4; x++ {
		for y := 0; y < size; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	for x := size/4; x < size/2; x++ {
		for y := 0; y < size; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
	for x := size/2; x < 3*size/4; x++ {
		for y := 0; y < size; y++ {
			pic.SetNRGBA(x, y, c)
		}
	}
	for x := 3*size/4; x < size; x++ {
		for y := 0; y < size; y++ {
			pic.SetNRGBA(x, y, d)
		}
	}
}

func fourSquare(pic *image.NRGBA, a, b, c, d color.NRGBA) {
	// top left
	for x := 0; x < size/2; x++ {
		for y := 0; y < size/2; y++ {
			pic.SetNRGBA(x, y, a)
		}
	}
	// bottom left
	for x := 0; x < size/2; x++ {
		for y := size/2; y < size; y++ {
			pic.SetNRGBA(x, y, b)
		}
	}
	// top right
	for x := size/2; x < size; x++ {
		for y := 0; y < size/2; y++ {
			pic.SetNRGBA(x, y, c)
		}
	}
	// bottom right
	for x := size/2; x < size; x++ {
		for y := size/2; y < size; y++ {
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

// Flag procedurally generates a size x size flag based on the Host of the request.
// Host can be overridden by passing a host in the query parameters.
func Flag() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		host := r.URL.Query().Get("host")
		if host == "" {
			host = r.Host
		}
		sum := sha256.Sum256([]byte(host))

		pic := image.NewNRGBA(image.Rect(0, 0, size, size))
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
