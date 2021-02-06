package selfupdate

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/frioux/leatherman/internal/version"
)

var Handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("content-type", "text/plain")

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		rw.WriteHeader(500)
	}

	if mostRecentFailure != nil {
		fmt.Fprintf(rw, "update failure: %s\n\n", mostRecentFailure)
	}

	if invalidToken {
		fmt.Fprintf(rw, "token is invalid, only updating hourly\n\n")
	}

	fmt.Fprintf(rw, "Leatherman built from %s on %s by with %s\n",
		version.Version, version.When, runtime.Version())

	for _, dep := range bi.Deps {
		fmt.Fprintf(rw, "%s@%s (%s)\n", dep.Path, dep.Version, dep.Sum)
		if dep.Replace != nil {
			r := dep.Replace
			fmt.Fprintf(rw, "   replaced by %s@%s (%s)\n", r.Path, r.Version, r.Sum)
		}
	}
})

func init() {
	http.DefaultServeMux.Handle("/version", Handler)
}
