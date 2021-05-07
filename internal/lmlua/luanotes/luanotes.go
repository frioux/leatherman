package luanotes

import (
	lua "github.com/yuin/gopher-lua"

	"github.com/frioux/leatherman/internal/lmlua"
	"github.com/frioux/leatherman/internal/notes"
)

func RegisterNotesPackage(L *lua.LState) {
	ns := L.NewTable()
	L.SetGlobal("notes", ns)

	registerNotesFunctions(L, ns)
	registerZineType(L)
	registerArticleType(L)
}

func registerNotesFunctions(L *lua.LState, ns *lua.LTable) {
	ns.RawSet(lua.LString("readarticlefromfs"), L.NewFunction(func(L *lua.LState) int {
		fss := lmlua.CheckFS(L, 1)
		path := L.CheckString(2)

		a, err := notes.ReadArticleFromFS(fss, path)
		if err != nil {
			L.Error(lua.LString(err.Error()), 0)
			return 0
		}

		ud := L.NewUserData()
		ud.Value = a
		L.SetMetatable(ud, L.GetTypeMetatable("article"))
		L.Push(ud)

		return 1
	}))
}

var articleMethods = map[string]lua.LGFunction{

	"rawcontents": func(L *lua.LState) int {
		a := checkArticle(L, 1)

		L.Push(lua.LString(string(a.RawContents)))

		return 1
	},
}

func checkArticle(L *lua.LState, where int) notes.Article {
	ud := L.CheckUserData(where)
	if v, ok := ud.Value.(notes.Article); ok {
		return v
	}
	L.ArgError(1, "notes.Article expected")
	return notes.Article{}
}

func registerArticleType(L *lua.LState) {
	mt := L.NewTypeMetatable("article")
	L.SetGlobal("article", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), articleMethods))
}

var zineMethods = map[string]lua.LGFunction{}

func checkZine(L *lua.LState, where int) *notes.Zine {
	ud := L.CheckUserData(where)
	if v, ok := ud.Value.(*notes.Zine); ok {
		return v
	}
	L.ArgError(1, "*notes.Zine expected")
	return nil
}

func registerZineType(L *lua.LState) {
	mt := L.NewTypeMetatable("zine")
	L.SetGlobal("zine", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), zineMethods))
}
