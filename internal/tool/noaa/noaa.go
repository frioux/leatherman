package noaa

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var upstream *url.URL

func init() {
	// Manually resolved CNAME of tgftp.nws.noaa.gov
	u, err := url.Parse("https://tgftp.cp.ncep.noaa.gov")
	if err != nil {
		panic("Couldn't parse url: " + err.Error())
	}
	upstream = u
}

// Proxy starts a proxy that can pretend to be the old noaa on http while
// actually proxying to noaa on https.
func Proxy(_ []string, _ io.Reader) error {
	http.Handle("/", httputil.NewSingleHostReverseProxy(upstream))

	return http.ListenAndServe(":9090", nil)
}
