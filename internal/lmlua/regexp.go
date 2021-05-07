package lmlua

import (
	"fmt"
	"regexp"

	lua "github.com/yuin/gopher-lua"
)

func RegisterRegexpPackage(L *lua.LState) {
	ns := RegisterRegexpNS(L)
	RegisterRegexpFunctions(L, ns)
	RegisterRegexpType(L, ns)
}

func RegisterRegexpFunctions(L *lua.LState, ns *lua.LTable) {
	ns.RawSet(lua.LString("compile"), L.NewFunction(func(L *lua.LState) int {
		s := L.CheckString(1)
		re, err := regexp.Compile(s)
		if err != nil {
			L.Error(lua.LString(err.Error()), 0)
			return 0
		}
		ud := L.NewUserData()
		ud.Value = re
		L.SetMetatable(ud, ns.RawGet(lua.LString("regexp")))
		L.Push(ud)
		return 1
	}))
}

func RegisterRegexpNS(L *lua.LState) *lua.LTable {
	ns := L.NewTable()
	L.SetGlobal("regexp", ns)
	return ns
}

func RegisterRegexpType(L *lua.LState, ns *lua.LTable) {
	mt := L.NewTypeMetatable("")
	ns.RawSet(lua.LString("regexp"), mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), regexpMethods))
}

func checkRegexp(L *lua.LState, which int) *regexp.Regexp {
	ud := L.CheckUserData(which)
	if v, ok := ud.Value.(*regexp.Regexp); ok {
		return v
	}
	L.ArgError(1, fmt.Sprintf("*regexp.Regexp expected, saw %T", ud.Value))
	return nil
}

var regexpMethods = map[string]lua.LGFunction{

	"findallstringsubmatch": func(L *lua.LState) int {
		re := checkRegexp(L, 1)
		s := L.CheckString(2)
		i := L.CheckNumber(3)

		found := re.FindAllStringSubmatch(s, int(i))
		ret := L.NewTable()

		for i := range found {
			t := L.NewTable()
			for j := range found[i] {
				t.RawSet(lua.LNumber(float64(j+1)), lua.LString(found[i][j]))
			}
			ret.RawSet(lua.LNumber(float64(i+1)), t)
		}
		L.Push(ret)

		return 1
	},

	"replaceallstringfunc": func(L *lua.LState) int {
		re := checkRegexp(L, 1)
		s := L.CheckString(2)
		f := L.CheckFunction(3)

		str := re.ReplaceAllStringFunc(s, func(in string) string {
			L.Push(f)
			L.Push(lua.LString(in))
			L.Call(1, 1)

			return L.CheckString(4)
		})

		L.Push(lua.LString(str))

		return 1
	},
}
