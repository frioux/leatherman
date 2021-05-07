package lmlua

import (
	"fmt"
	"io/fs"

	lua "github.com/yuin/gopher-lua"

	"github.com/frioux/leatherman/internal/lmfs"
)

func RegisterFSType(L *lua.LState) {
	fmt.Println("registering fs type")
	mt := L.NewTypeMetatable("fs")
	L.SetGlobal("fs", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), fsMethods))
}

func CheckFS(L *lua.LState, where int) fs.FS {
	ud := L.CheckUserData(where)
	if v, ok := ud.Value.(fs.FS); ok {
		return v
	}
	L.ArgError(1, fmt.Sprintf("fs expected, saw %T", ud.Value))
	return nil
}

var fsMethods = map[string]lua.LGFunction{
	// TODO
	// open
	// create
	"writefile": func(L *lua.LState) int {
		fss := CheckFS(L, 1)
		path := L.CheckString(2)
		contents := L.CheckString(3)
		lmfs.WriteFile(fss, path, []byte(contents), 0644)
		return 0
	},
}
