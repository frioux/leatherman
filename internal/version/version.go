package version

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Version is the git version that produced this binary.
var Version string

// When is the datestamp that produced this binary.
var When string

var Handler = http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("content-type", "text/plain")

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		rw.WriteHeader(500)
	}

	fmt.Fprintln(rw, "version:", Version)

	for _, dep := range bi.Deps {
		fmt.Fprintf(rw, "%s@%s (%s)\n", dep.Path, dep.Version, dep.Sum)
		if dep.Replace != nil {
			r := dep.Replace
			fmt.Fprintf(rw, "   replaced by %s@%s (%s)\n", r.Path, r.Version, r.Sum)
		}
	}
})
