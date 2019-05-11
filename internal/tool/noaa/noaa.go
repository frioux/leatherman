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

/*
Proxy creates a proxy for https://tgftp.nws.noaa.gov, but http, and listening on
9090.  This is because Ubuntu 18.04 ships with [taffybar]() 0.4.6, which only
supports http for the weather widgets and has hardcoded URLs.

To install, add this line to your hosts file:

```
127.0.0.1       tgftp.nws.noaa.gov
```

And run this iptables command:

```
iptables -t nat -A OUTPUT -o lo -p tcp --dport 80 -j REDIRECT --to-port 9090
```

Command: noaa-proxy
*/
func Proxy(_ []string, _ io.Reader) error {
	http.Handle("/", httputil.NewSingleHostReverseProxy(upstream))

	return http.ListenAndServe(":9090", nil)
}
