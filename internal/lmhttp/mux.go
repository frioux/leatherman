package lmhttp

import (
	"net/http"
	"sort"

	"embed"
	"html/template"
)

type ClearMux struct {
	endpoints []string
	*http.ServeMux
}

//go:embed templates/*
var templateFS embed.FS

var templates = template.Must(template.New("tmpl").ParseFS(templateFS, "templates/*"))

func NewClearMux() *ClearMux {
	m := &ClearMux{ServeMux: http.NewServeMux()}
	m.ServeMux.Handle("/", HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) error {
		rw.Header().Add("Content-Type", "text/html")
		if err := templates.ExecuteTemplate(rw, "list.html", m.endpoints); err != nil {
			return err
		}

		return nil
	}))

	return m
}

func (m *ClearMux) Handle(pattern string, handler http.Handler) {
	m.endpoints = append(m.endpoints, pattern)
	sort.Strings(m.endpoints) // could use a heap but meh
	m.ServeMux.Handle(pattern, handler)
}
