package lmhttp

import (
	"fmt"
	"net/http"
	"sort"
)

type ClearMux struct {
	endpoints []string
	*http.ServeMux
}

func NewClearMux() *ClearMux {
	m := &ClearMux{ServeMux: http.NewServeMux()}
	m.ServeMux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		for _, s := range m.endpoints {
			fmt.Fprintf(rw, " * %s\n", s)
		}
	}))

	return m
}

func (m *ClearMux) Handle(pattern string, handler http.Handler) {
	m.endpoints = append(m.endpoints, pattern)
	sort.Strings(m.endpoints) // could use a heap but meh
	m.ServeMux.Handle(pattern, handler)
}
