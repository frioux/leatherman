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
	"optional": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		optional := L.CreateTable(0, len(es.optional))
		for w, v := range es.optional {
			optional.RawSetString(w, lua.LBool(v))
		}
		L.Push(optional)
		return 1
	},
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
	"optionallen": func(L *lua.LState) int {
		es := checkEmojiSet(L)

		L.Push(lua.LNumber(float64(len(es.optional))))
		return 1
	},

	"required": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		required := L.CreateTable(len(es.required), 0)
		for i, s := range es.required {
			required.Insert(i+1, lua.LString(s))
		}
		L.Push(required)
		return 1
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

	"len": func(L *lua.LState) int {
		es := checkEmojiSet(L)

		L.Push(lua.LNumber(len(es.optional) + len(es.required)))
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
	"words": func(L *lua.LState) int {
		es := checkEmojiSet(L)
		words := L.CreateTable(len(es.words), 0)
		for i, s := range es.words {
			words.Insert(i+1, lua.LString(s))
		}
		L.Push(words)
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
