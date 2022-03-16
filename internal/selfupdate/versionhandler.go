package selfupdate

import (
	"fmt"
	"net/http"

	"github.com/frioux/leatherman/internal/version"
)

var Handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("content-type", "text/plain")

	if mostRecentFailure != nil {
		fmt.Fprintf(rw, "update failure: %s\n\n", mostRecentFailure)
	}

	if invalidToken {
		fmt.Fprintf(rw, "token is invalid, only updating hourly\n\n")
	}

	version.Render(rw)
})

func init() {
	http.DefaultServeMux.Handle("/version", Handler)
}
