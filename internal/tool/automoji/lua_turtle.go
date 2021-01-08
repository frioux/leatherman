package automoji

import (
	"github.com/hackebrot/turtle"
	lua "github.com/yuin/gopher-lua"
)

func registerTurtleType(L *lua.LState) {
	mt := L.NewTypeMetatable("turtleemoji")
	L.SetGlobal("turtleemoji", mt)
	L.SetField(mt, "findbyname", L.NewFunction(findTurtleByName))
	L.SetField(mt, "findbychar", L.NewFunction(findTurtleByChar))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), turtleMethods))
}

func findTurtleByName(L *lua.LState) int {
	name := L.CheckString(1)

	v, ok := turtle.Emojis[name]
	if !ok {
		L.Push(lua.LNil)
		return 1
	}

	ud := L.NewUserData()
	ud.Value = v
	L.SetMetatable(ud, L.GetTypeMetatable("turtleemoji"))
	L.Push(ud)
	return 1
}

func findTurtleByChar(L *lua.LState) int {
	name := L.CheckString(1)

	v, ok := turtle.EmojisByChar[name]
	if !ok {
		L.Push(lua.LNil)
		return 1
	}

	ud := L.NewUserData()
	ud.Value = v
	L.SetMetatable(ud, L.GetTypeMetatable("turtleemoji"))
	L.Push(ud)
	return 1
}

func checkTurtle(L *lua.LState) *turtle.Emoji {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*turtle.Emoji); ok {
		return v
	}
	L.ArgError(1, "turtle.Emoji expected")
	return nil
}

var turtleMethods = map[string]lua.LGFunction{
	"name": func(L *lua.LState) int {
		c := checkTurtle(L)

		L.Push(lua.LString(c.Name))

		return 1
	},
	"category": func(L *lua.LState) int {
		c := checkTurtle(L)

		L.Push(lua.LString(c.Category))

		return 1
	},
	"char": func(L *lua.LState) int {
		c := checkTurtle(L)

		L.Push(lua.LString(c.Char))

		return 1
	},
	"haskeyword": func(L *lua.LState) int {
		c := checkTurtle(L)
		e := L.CheckString(2)

		for _, s := range c.Keywords {
			if s == e {
				L.Push(lua.LBool(true))
				return 1
			}
		}
		L.Push(lua.LBool(false))
		return 1
	},
}
