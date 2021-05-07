package lmlua

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/frioux/leatherman/internal/lmhttp"
	lua "github.com/yuin/gopher-lua"
)

func RegisterHTTPPackage(L *lua.LState) {
	registerHTTPFunctions(L)
	registerResponseWriterType(L)
	registerRequestType(L)
	registerURLType(L)
	registerValuesType(L)
}

func registerHTTPFunctions(L *lua.LState) {
	httpNS := L.NewTable()
	L.SetGlobal("http", httpNS)

	httpNS.RawSet(lua.LString("get"), L.NewFunction(func(L *lua.LState) int {
		s := L.CheckString(1)
		r, err := lmhttp.Get(context.Background(), s)
		if err != nil {
			L.Error(lua.LString(err.Error()), 0)
			return 0
		}
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		L.Push(lua.LString(string(b)))
		return 1
	}))

	httpNS.RawSet(lua.LString("multiget"), L.NewFunction(func(L *lua.LState) int {
		t := L.CheckTable(1)

		contents := make([]string, t.Len())
		wg := &sync.WaitGroup{}
		wg.Add(t.Len())

		k, v := t.Next(lua.LNil)
		for ; k.Type() != lua.LTNil; k, v = t.Next(k) {
			k, v := k, v
			go func() {
				defer wg.Done()

				i := int(k.(lua.LNumber)) - 1
				url := string(v.(lua.LString))
				resp, err := lmhttp.Get(context.Background(), url)
				if err != nil {
					L.Error(lua.LString(err.Error()), 1)
					return
				}
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}

				contents[i] = string(b)
			}()
		}

		wg.Wait()

		resp := L.NewTable()

		for i := range contents {
			resp.RawSet(t.RawGet(lua.LNumber(float64(i)+1)), lua.LString(contents[i]))
		}

		L.Push(resp)
		return 1
	}))
}

var responsewriterMethods = map[string]lua.LGFunction{
	// TODO
	// header
	"write": func(L *lua.LState) int {
		rw := checkResponseWriter(L)
		s := L.CheckString(2)
		if _, err := rw.Write([]byte(s)); err != nil {
			L.Error(lua.LString(err.Error()), 0)
			return 0
		}
		return 0
	},
	"writeheader": func(L *lua.LState) int {
		rw := checkResponseWriter(L)
		s := L.CheckInt(2)
		rw.WriteHeader(s)
		return 0
	},
}

func registerResponseWriterType(L *lua.LState) {
	mt := L.NewTypeMetatable("responsewriter")
	L.SetGlobal("responsewriter", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), responsewriterMethods))
}

func checkResponseWriter(L *lua.LState) http.ResponseWriter {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(http.ResponseWriter); ok {
		return v
	}
	L.ArgError(1, "http.ResponseWriter expected")
	return nil
}

var requestMethods = map[string]lua.LGFunction{
	"url": func(L *lua.LState) int {
		r := checkRequest(L)
		ud := L.NewUserData()
		ud.Value = r.URL
		L.SetMetatable(ud, L.GetTypeMetatable("url"))
		L.Push(ud)
		return 1
	},
}

func registerRequestType(L *lua.LState) {
	mt := L.NewTypeMetatable("request")
	L.SetGlobal("request", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), requestMethods))
}

func checkRequest(L *lua.LState) *http.Request {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*http.Request); ok {
		return v
	}
	L.ArgError(1, "*http.Request expected")
	return nil
}

var urlMethods = map[string]lua.LGFunction{
	"query": func(L *lua.LState) int {
		u := checkURL(L)
		ud := L.NewUserData()
		ud.Value = u.Query()
		L.SetMetatable(ud, L.GetTypeMetatable("values"))
		L.Push(ud)
		return 1
	},
}

func registerURLType(L *lua.LState) {
	ns := L.NewTable()
	mt := L.NewTypeMetatable("url")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), urlMethods))
	ns.RawSet(lua.LString("url"), mt)
	ns.RawSet(lua.LString("parse"), L.NewFunction(func(L *lua.LState) int {
		s := L.CheckString(1)
		u, err := url.Parse(s)
		if err != nil {
			L.Error(lua.LString(err.Error()), 0)
			return 0
		}
		ud := L.NewUserData()
		ud.Value = u
		L.SetMetatable(ud, mt)
		L.Push(ud)
		return 1
	}))
	L.SetGlobal("url", ns)
}

func checkURL(L *lua.LState) *url.URL {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*url.URL); ok {
		return v
	}
	L.ArgError(1, "*url.URL expected")
	return nil
}

var valuesMethods = map[string]lua.LGFunction{
	"get": func(L *lua.LState) int {
		v := checkValues(L)
		k := L.CheckString(2)
		L.Push(lua.LString(v.Get(k)))
		return 1
	},
}

func registerValuesType(L *lua.LState) {
	mt := L.NewTypeMetatable("values")
	L.SetGlobal("values", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), valuesMethods))
}

func checkValues(L *lua.LState) url.Values {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(url.Values); ok {
		return v
	}
	L.ArgError(1, "url.Values expected")
	return nil
}
