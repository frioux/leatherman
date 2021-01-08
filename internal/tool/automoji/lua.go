package automoji

import (
	"io/ioutil"
	"regexp"

	"github.com/frioux/leatherman/internal/dropbox"
	lua "github.com/yuin/gopher-lua"
)

func loadLua(dbCl dropbox.Client, path string) (string, error) {
	r, err := dbCl.Download(path)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), err
}

func registerEmojiSetType(L *lua.LState) {
	mt := L.NewTypeMetatable("emojiset")
	L.SetGlobal("emojiset", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), emojiSetMethods))
}

func setGlobalEmojiSet(L *lua.LState, name string, es *emojiSet) int {
	ud := L.NewUserData()
	ud.Value = es
	L.SetMetatable(ud, L.GetTypeMetatable("emojiset"))
	L.SetGlobal(name, ud)
	return 1
}

func checkEmojiSet(L *lua.LState) *emojiSet {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*emojiSet); ok {
		return v
	}
	L.ArgError(1, "emojiSet expected")
	return nil
}

var emojiSetMethods = map[string]lua.LGFunction{
	"hasoptional": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		e := L.CheckString(2)
		L.Push(lua.LBool(es.optional[e]))
		return 1
	},
	"addoptional": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		e := L.CheckString(2)
		es.optional[e] = true
		return 0
	},
	"removeoptional": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		e := L.CheckString(2)
		delete(es.optional, e)
		return 0
	},

	"hasrequired": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		e := L.CheckString(2)
		for _, s := range es.required {
			if s == e {
				L.Push(lua.LBool(true))
				return 1
			}
		}
		L.Push(lua.LBool(false))
		return 1
	},
	"addrequired": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		e := L.CheckString(2)
		es.required = append(es.required, e)
		return 0
	},
	"removerequired": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		e := L.CheckString(2)
		newRequired := es.required[:0] // share backing array
		for _, v := range es.required {
			if v == e {
				continue
			}
			newRequired = append(newRequired, v)
		}
		es.required = newRequired
		return 0
	},

	"message": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		L.Push(lua.LString(es.message))
		return 1
	},
	"messagematches": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		e := L.CheckString(2)
		re := regexp.MustCompile(e)
		L.Push(lua.LBool(re.MatchString(es.message)))
		return 1
	},

	"hasword": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		e := L.CheckString(2)
		for _, s := range es.words {
			if s == e {
				L.Push(lua.LBool(true))
				return 1
			}
		}
		L.Push(lua.LBool(false))
		return 1
	},
}

func luaEval(es *emojiSet, code string) error {
	L := lua.NewState()
	defer L.Close()

	registerEmojiSetType(L)
	registerTurtleType(L)
	setGlobalEmojiSet(L, "es", es)

	return L.DoString(code)
}
