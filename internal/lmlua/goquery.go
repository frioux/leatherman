package lmlua

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	lua "github.com/yuin/gopher-lua"
)

func RegisterGoqueryPackage(L *lua.LState) {
	registerSelectionType(L)
	registerGoqueryFunctions(L)
}

func registerGoqueryFunctions(L *lua.LState) {
	ns := L.NewTable()
	L.SetGlobal("goquery", ns)

	ns.RawSet(lua.LString("newdocumentfromstring"), L.NewFunction(func(L *lua.LState) int {
		s := L.CheckString(1)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
		if err != nil {
			L.Error(lua.LString(err.Error()), 0)
			return 0
		}
		ud := L.NewUserData()
		ud.Value = doc.Selection
		L.SetMetatable(ud, L.GetTypeMetatable("goqueryselection"))
		L.Push(ud)
		return 1
	}))
}

func registerSelectionType(L *lua.LState) {
	mt := L.NewTypeMetatable("goqueryselection")
	L.SetGlobal("goqueryselection", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), goqueryselectionMethods))
}

func checkGoQuerySelection(L *lua.LState) *goquery.Selection {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*goquery.Selection); ok {
		return v
	}
	L.ArgError(1, "*goquery.Selection expected")
	return nil
}

var goqueryselectionMethods = map[string]lua.LGFunction{

	"attr": func(L *lua.LState) int {
		sel := checkGoQuerySelection(L)
		s := L.CheckString(2)
		found, _ := sel.Attr(s)
		L.Push(lua.LString(found))

		return 1
	},

	"each": func(L *lua.LState) int {
		sel := checkGoQuerySelection(L)
		f := L.CheckFunction(2)

		ud := L.NewUserData()
		ud.Value = sel.Each(func(i int, sel *goquery.Selection) {
			L.Push(f)
			ud := L.NewUserData()
			ud.Value = sel
			L.SetMetatable(ud, L.GetTypeMetatable("goqueryselection"))
			L.Push(ud)
			L.Push(lua.LNumber(float64(i)))
			L.Call(2, 0)
		})
		L.SetMetatable(ud, L.GetTypeMetatable("goqueryselection"))
		L.Push(ud)
		return 1
	},

	"find": func(L *lua.LState) int {
		sel := checkGoQuerySelection(L)
		s := L.CheckString(2)

		ud := L.NewUserData()
		ud.Value = sel.Find(s)
		L.SetMetatable(ud, L.GetTypeMetatable("goqueryselection"))
		L.Push(ud)
		return 1
	},

	"text": func(L *lua.LState) int {
		sel := checkGoQuerySelection(L)
		L.Push(lua.LString(sel.Text()))
		return 1
	},
}
