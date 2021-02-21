package img

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"

	lua "github.com/yuin/gopher-lua"
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

	if err := luaEval(img, args[1:]); err != nil {
		return err
	}

	return png.Encode(os.Stdout, img)
}

func luaEval(img *image.NRGBA, code []string) error {
	L := lua.NewState()
	defer L.Close()

	registerImageFunctions(L, img)

	for _, c := range code {
		if err := L.DoString(c); err != nil {
			return err
		}
	}

	return nil
}

func checkColor(L *lua.LState, w int) color.Color {
	ud := L.CheckUserData(w)
	if v, ok := ud.Value.(color.Color); ok {
		return v
	}
	L.ArgError(w, "image.Color expected")
	return nil
}

func registerImageFunctions(L *lua.LState, img *image.NRGBA) {
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
		ud.Value = color.RGBA{uint8(r), uint8(g), uint8(b), 255}

		L.Push(ud)
		return 1
	}))

	{
		black := L.NewUserData()
		black.Value = color.Black
		L.SetGlobal("black", black)

		white := L.NewUserData()
		white.Value = color.RGBA{0, 0, 0, 255}
		L.SetGlobal("white", white)

		red := L.NewUserData()
		red.Value = color.RGBA{255, 0, 0, 255}
		L.SetGlobal("red", red)

		blue := L.NewUserData()
		blue.Value = color.RGBA{0, 0, 255, 255}
		L.SetGlobal("blue", blue)

		yellow := L.NewUserData()
		yellow.Value = color.RGBA{255, 255, 0, 255}
		L.SetGlobal("yellow", yellow)

		green := L.NewUserData()
		green.Value = color.RGBA{0, 255, 0, 255}
		L.SetGlobal("green", green)

		orange := L.NewUserData()
		orange.Value = color.RGBA{255, 165, 0, 255}
		L.SetGlobal("orange", orange)

		purple := L.NewUserData()
		purple.Value = color.RGBA{128, 0, 128, 255}
		L.SetGlobal("purple", purple)

		cyan := L.NewUserData()
		cyan.Value = color.RGBA{0, 255, 255, 255}
		L.SetGlobal("cyan", cyan)

		magenta := L.NewUserData()
		magenta.Value = color.RGBA{255, 0, 255, 255}
		L.SetGlobal("magenta", magenta)
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
		m := (y2 - y1) / (x2 - x1)
		// y = m*x + b
		// y - m*x = b
		// b = y - m*x
		b := y1 - m*x1
		l := math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))

		if m == math.Inf(1) || m == math.Inf(-1) {
			start, end := y1, y2
			if start > end {
				start, end = end, start
			}
			for y := start; y <= end; y += l / 1000 {
				img.Set(int(math.Round(x1)), int(math.Round(y)), c)
			}
		} else {
			start, end := x1, x2
			if start > end {
				start, end = end, start
			}

			for x := start; x <= end; x += l / 1000 {
				y := m*x + b
				img.Set(int(math.Round(x)), int(math.Round(y)), c)
			}
		}
	}

	L.SetGlobal("circ", L.NewFunction(func(L *lua.LState) int {
		x := int(L.CheckNumber(1))
		y := int(L.CheckNumber(2))
		r := float64(L.CheckNumber(3))
		border := checkColor(L, 4)
		fill := checkColor(L, 5)

		// draw borders
		for t := 0.0; t < 2*math.Pi; t += 0.001 /* uhh */ {
			xt := r*math.Cos(t) + float64(x)
			yt := r*math.Sin(t) + float64(y)

			line(float64(x), float64(y), xt, yt, fill)
		}
		for t := 0.0; t < 2*math.Pi; t += 0.001 /* uhh */ {
			xt := r*math.Cos(t) + float64(x)
			yt := r*math.Sin(t) + float64(y)

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
}
