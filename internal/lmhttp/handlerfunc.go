package lmhttp

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (f HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := f(rw, r); err != nil {
		rw.WriteHeader(500)
		fmt.Fprintln(os.Stderr, err)
	}
}

// TrimHandlerPrefix adapts handlers to a mux or possibly an alternate
// subroute.  The prefix is stripped from the url path such that the inner
// handler is unaware of the prefix.
func TrimHandlerPrefix(prefix string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		r = r.Clone(r.Context())

		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)

		h.ServeHTTP(rw, r)
	})
}
