package bamboo

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

func authHandler() (string, http.HandlerFunc) {
	pass := "123432"
	return pass, func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			f, err := os.Open("./testdata/auth.html")
			if err != nil {
				panic("Couldn't load testdata: " + err.Error())
			}

			if _, err := io.Copy(rw, f); err != nil {
				panic("Couldn't copy testdata: " + err.Error())
			}
			return
		}
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				panic("Couldn't ParseForm: " + err.Error())
			}
			http.SetCookie(rw, &http.Cookie{
				Name:  "auth",
				Value: r.Form.Get("password"),
			})
			return
		}
		rw.WriteHeader(404)
	}
}

func treeHandler(pass string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("auth")
		if err != nil {
			panic("Couldn't parse Cookie: " + err.Error())
		}
		if c.Value != pass {
			rw.WriteHeader(401)
			return
		}

		f, err := os.Open("./testdata/tree.html")
		if err != nil {
			panic("Couldn't load testdata: " + err.Error())
		}

		if _, err := io.Copy(rw, f); err != nil {
			panic("Couldn't copy testdata: " + err.Error())
		}
	}
}

func dirHandler(pass string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("auth")
		if err != nil {
			panic("Couldn't parse Cookie: " + err.Error())
		}
		if c.Value != pass {
			rw.WriteHeader(401)
			return
		}

		rw.Write([]byte("dir"))
	}
}

func testClientAndServer() (*client, func()) {
	mx := http.NewServeMux()
	pass, auth := authHandler()
	mx.Handle("/auth", http.HandlerFunc(auth))
	mx.Handle("/dir", http.HandlerFunc(dirHandler(pass)))
	mx.Handle("/tree", http.HandlerFunc(treeHandler(pass)))

	ts := httptest.NewServer(mx)

	return &client{
		authURL: ts.URL + "/auth",
		dirURL:  ts.URL + "/dir",
		treeURL: ts.URL + "/tree",

		user:     "ignored",
		password: pass,
	}, ts.Close
}
