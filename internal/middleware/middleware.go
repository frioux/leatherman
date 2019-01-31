package middleware

import (
	"io"
	"net/http"

	"encoding/json"
)

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

type logline struct {
	Method     string
	URL        string
	UserAgent  string
	Proto      string
	Host       string
	RemoteAddr string
	StatusCode int
}

func Log(logger io.Writer) Adapter {
	e := json.NewEncoder(logger)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			lw := &loggingResponseWriter{ResponseWriter: w}
			defer func() {
				e.Encode(logline{
					Method:     r.Method,
					URL:        r.URL.String(),
					UserAgent:  r.UserAgent(),
					Proto:      r.Proto,
					Host:       r.Host,
					RemoteAddr: r.RemoteAddr,
					StatusCode: lw.statusCode,
				})
			}()
			h.ServeHTTP(lw, r)

		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
