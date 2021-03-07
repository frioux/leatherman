package lmhttp

import (
	"fmt"
	"net/http"
	"os"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (f HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := f(rw, r); err != nil {
		rw.WriteHeader(500)
		fmt.Fprintln(os.Stderr, err)
	}
}
