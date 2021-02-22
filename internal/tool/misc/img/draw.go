package img

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"os"

	"github.com/frioux/leatherman/internal/drawlua"
)

func Draw(args []string, _ io.Reader) error {
	if len(args) == 1 {
		args = append(args, "")
	}

	img := image.NewNRGBA(image.Rect(0, 0, 128, 128))
	for x := 0; x < 128; x++ {
		for y := 0; y < 128; y++ {
			img.Set(x, y, color.Black)
		}
	}

	if err := drawlua.Eval(img, args[1:]); err != nil {
		return err
	}

	return png.Encode(os.Stdout, img)
}
