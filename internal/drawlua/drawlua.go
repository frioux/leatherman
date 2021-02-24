package drawlua

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math"
	"os"
	"regexp"

	lua "github.com/yuin/gopher-lua"
)

type ImageSetter interface {
	image.Image
	Set(int, int, color.Color)
}

// Eval runs code against the image.
func Eval(img ImageSetter, code []string) error {
	L := lua.NewState()
	defer L.Close()

	cleanup := RegisterDrawFunctions(L, img)

	for _, c := range code {
		if err := L.DoString(c); err != nil {
			return err
		}
	}

	return cleanup()
}

func checkColor(L *lua.LState, w int) color.Color {
	ud := L.CheckUserData(w)
	if v, ok := ud.Value.(color.Color); ok {
		return v
	}
	L.ArgError(w, "image.Color expected")
	return nil
}

func RegisterDrawFunctions(L *lua.LState, img ImageSetter) (cleanup func() error) {
	palette := color.Palette([]color.Color{
		color.Black,
		color.RGBA{0, 0, 0, 255},     // white
		color.RGBA{255, 0, 0, 255},   // red
		color.RGBA{0, 0, 255, 255},   // blue
		color.RGBA{255, 255, 0, 255}, // yellow
		color.RGBA{0, 255, 0, 255},   // green
		color.RGBA{255, 165, 0, 255}, // orange
		color.RGBA{128, 0, 128, 255}, // purple
		color.RGBA{0, 255, 255, 255}, // cyan
		color.RGBA{255, 0, 255, 255}, // magenta
		color.RGBA{0, 0, 0, 0},       // clear
	})
	debugDraw := func(string, image.Image) error { return nil }
	cleanup = func() error { return nil }

	if d := os.Getenv("LM_DEBUG_DRAW"); d != "" {
		dgif := &gif.GIF{}
		shouldDebug := regexp.MustCompile(d)
		e, err := os.Create("debug.log")
		if err != nil {
			panic(err)
		}

		debugDraw = func(name string, img image.Image) error {
			if !shouldDebug.MatchString(name) {
				return nil
			}

			fmt.Fprintln(e, name)
			frame := image.NewPaletted(img.Bounds(), palette)
			draw.Over.Draw(frame, img.Bounds(), img, image.Point{})
			dgif.Image = append(dgif.Image, frame)
			dgif.Delay = append(dgif.Delay, 1) // 10ms, minimum delay
			return nil
		}
		cleanup = func() error {
			defer e.Close()
			f, err := os.Create("debug.gif")
			if err != nil {
				return err
			}
			defer f.Close()

			if err := gif.EncodeAll(f, dgif); err != nil {
				return err
			}

			return nil
		}
	}

	L.SetGlobal("set", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckNumber(1)
		y := L.CheckNumber(2)
		c := checkColor(L, 3)

		img.Set(int(x), int(y), c)

		return 0
	}))

	L.SetGlobal("rgb", L.NewFunction(func(L *lua.LState) int {
		r := L.CheckNumber(1)
		g := L.CheckNumber(2)
		b := L.CheckNumber(3)

		ud := L.NewUserData()

		if r >= 0 && r <= 1 && g >= 0 && g <= 1 && b >= 0 && b <= 1 {
			ud.Value = color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
		} else {
			ud.Value = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
		}

		L.Push(ud)
		return 1
	}))

	{
		black := L.NewUserData()
		black.Value = palette[0]
		L.SetGlobal("black", black)

		white := L.NewUserData()
		white.Value = palette[1]
		L.SetGlobal("white", white)

		red := L.NewUserData()
		red.Value = palette[2]
		L.SetGlobal("red", red)

		blue := L.NewUserData()
		blue.Value = palette[3]
		L.SetGlobal("blue", blue)

		yellow := L.NewUserData()
		yellow.Value = palette[4]
		L.SetGlobal("yellow", yellow)

		green := L.NewUserData()
		green.Value = palette[5]
		L.SetGlobal("green", green)

		orange := L.NewUserData()
		orange.Value = palette[6]
		L.SetGlobal("orange", orange)

		purple := L.NewUserData()
		purple.Value = palette[7]
		L.SetGlobal("purple", purple)

		cyan := L.NewUserData()
		cyan.Value = palette[8]
		L.SetGlobal("cyan", cyan)

		magenta := L.NewUserData()
		magenta.Value = palette[9]
		L.SetGlobal("magenta", magenta)

		clear := L.NewUserData()
		clear.Value = palette[10]
		L.SetGlobal("clear", clear)
	}

	L.SetGlobal("sin", L.NewFunction(func(L *lua.LState) int {
		t := L.CheckNumber(1)

		L.Push(lua.LNumber(math.Sin(float64(t))))
		return 1
	}))

	L.SetGlobal("cos", L.NewFunction(func(L *lua.LState) int {
		t := L.CheckNumber(1)

		L.Push(lua.LNumber(math.Cos(float64(t))))
		return 1
	}))

	L.SetGlobal("tan", L.NewFunction(func(L *lua.LState) int {
		t := L.CheckNumber(1)

		L.Push(lua.LNumber(math.Tan(float64(t))))
		return 1
	}))

	L.SetGlobal("PI", lua.LNumber(math.Pi))

	L.SetGlobal("rect", L.NewFunction(func(L *lua.LState) int {
		x1 := int(L.CheckNumber(1))
		y1 := int(L.CheckNumber(2))
		x2 := int(L.CheckNumber(3))
		y2 := int(L.CheckNumber(4))
		border := checkColor(L, 5)
		fill := checkColor(L, 6)

		// draw borders
		for x := x1; x <= x2; x++ {
			img.Set(x, y1, border)
			img.Set(x, y2, border)
		}
		for y := y1 + 1; y < y2; y++ {
			img.Set(x1, y, border)
			img.Set(x2, y, border)
		}

		// draw fill
		for x := x1 + 1; x < x2; x++ {
			for y := y1 + 1; y < y2; y++ {
				img.Set(x, y, fill)
			}
		}

		return 0
	}))

	line := func(x1, y1, x2, y2 float64, c color.Color) {
		debugDraw(fmt.Sprintf("line(%f, %f, %f, %f, <c>)", x1, y1, x2, y2), img)

		if math.Round(x1) == math.Round(x2) {
			for y := y1; y < y2; y++ {
				img.Set(int(math.Round(x1)), int(math.Round(y)), c)
			}
			return
		} else if math.Round(y1) == math.Round(y2) {
			for x := x1; x < x2; x++ {
				img.Set(int(math.Round(x)), int(math.Round(y1)), c)
			}
			return
		}

		m := (y2 - y1) / (x2 - x1)

		if m >= -1 && m <= 1 {
			y := y1
			start, end := x1, x2
			if start > end {
				start, end = end, start
				y = y2
			}

			for x := start; x <= end; x++ {
				img.Set(int(math.Round(x)), int(math.Round(y)), c)
				y += m
			}
		} else {
			m1 := (x2 - x1) / (y2 - y1)
			x := x1
			start, end := y1, y2
			if start > end {
				start, end = end, start
				x = x2
			}

			for y := start; y <= end; y++ {
				img.Set(int(math.Round(x)), int(math.Round(y)), c)
				x += m1
			}

		}
	}

	L.SetGlobal("arc", L.NewFunction(func(L *lua.LState) int {
		x := int(L.CheckNumber(1))
		y := int(L.CheckNumber(2))
		r := float64(L.CheckNumber(3))
		s := float64(L.CheckNumber(4))
		e := float64(L.CheckNumber(5))

		border := checkColor(L, 6)
		fill := checkColor(L, 7)

		for t := 2 * math.Pi * s; t < 2*math.Pi*e; t += 1 / r {
			xt := r*math.Cos(t) + float64(x)
			yt := r*math.Sin(t) + float64(y)

			line(float64(x), float64(y), xt, yt, fill)
			img.Set(int(math.Round(xt)), int(math.Round(yt)), border)
		}

		return 0
	}))

	L.SetGlobal("circ", L.NewFunction(func(L *lua.LState) int {
		xc := int(L.CheckNumber(1))
		yc := int(L.CheckNumber(2))
		r := float64(L.CheckNumber(3))

		border := checkColor(L, 4)
		fill := checkColor(L, 5)

		for x := int(-r); x <= int(r); x++ {
			for y := int(-r); y <= int(r); y++ {
				if x*x+y*y <= int(r*r) {
					img.Set(x+xc, y+yc, fill)
				}
			}
		}

		// improve with http://weber.itn.liu.se/~stegu/circle/circlealgorithm.pdf
		for t := 0.0; t < 2*math.Pi; t += 1 / r {
			xt := r*math.Cos(t) + float64(xc)
			yt := r*math.Sin(t) + float64(yc)

			img.Set(int(math.Round(xt)), int(math.Round(yt)), border)
		}

		return 0
	}))

	L.SetGlobal("line", L.NewFunction(func(L *lua.LState) int {
		x1 := float64(L.CheckNumber(1))
		y1 := float64(L.CheckNumber(2))
		x2 := float64(L.CheckNumber(3))
		y2 := float64(L.CheckNumber(4))
		c := checkColor(L, 5)

		line(x1, y1, x2, y2, c)

		return 0
	}))

	return
}
