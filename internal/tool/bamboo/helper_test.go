package bamboo

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

func authHandler(rw http.ResponseWriter, r *http.Request) {
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
			Value: "xyzzy",
		})
		// fmt.Println(r.Form)
		return
	}
	rw.WriteHeader(404)
}

func treeHandler(rw http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth")
	if err != nil {
		panic("Couldn't parse Cookie: " + err.Error())
	}
	if c.Value != "xyzzy" {
		rw.WriteHeader(401)
		return
	}

	rw.Write([]byte(`json = {"tree":1};` + "\n"))
}

func dirHandler(rw http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth")
	if err != nil {
		panic("Couldn't parse Cookie: " + err.Error())
	}
	if c.Value != "xyzzy" {
		rw.WriteHeader(401)
		return
	}

	rw.Write([]byte("dir"))
}

func testClientAndServer() (*client, func()) {
	mx := http.NewServeMux()
	mx.Handle("/auth", http.HandlerFunc(authHandler))
	mx.Handle("/dir", http.HandlerFunc(dirHandler))
	mx.Handle("/tree", http.HandlerFunc(treeHandler))

	ts := httptest.NewServer(mx)

	return &client{
		authURL: ts.URL + "/auth",
		dirURL:  ts.URL + "/dir",
		treeURL: ts.URL + "/tree",
	}, ts.Close
}
