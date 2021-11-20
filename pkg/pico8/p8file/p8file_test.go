package p8file_test

import (
	"embed"
	"os"
	"testing"

	"github.com/frioux/leatherman/pkg/pico8/p8file"
)

//go:embed testdata
var carts embed.FS

func TestParse(t *testing.T) {
	f, err := os.Open("/home/frew/.lexaloffle/pico-8/carts/breakout.p8")
	if err != nil {
		t.Errorf("couldn't open breakout.p8: %s", err)
		return
	}
	_, err = p8file.Parse(f)
	if err != nil {
		t.Errorf("couldn't parse breakout.p8: %s", err)
		return
	}
	// i := c.SpriteImage()
	// o, err := os.Create("x.png")
	// if err != nil {
	// 	panic(err)
	// }
	// defer o.Close()
	// if err := png.Encode(o, i); err != nil {
	// 	t.Errorf("couldn't write png: %s", err)
	// }

	// for i := 0; i < 8; i++ {
	// 	o2, err := os.Create(fmt.Sprintf("y%d.png", i))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer o2.Close()
	// 	if err := png.Encode(o2, c.SpriteAt(uint8(i))); err != nil {
	// 		t.Errorf("couldn't write png: %s", err)
	// 	}
	// }

	// o3, err := os.Create("z.png")
	// if err != nil {
	// 	panic(err)
	// }
	// defer o3.Close()
	// if err := png.Encode(o3, c.MapImage()); err != nil {
	// 	t.Errorf("couldn't write png: %s", err)
	// }
}
