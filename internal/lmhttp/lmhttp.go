package lmhttp

import (
	"io"
	"net/http"
	"runtime/debug"
)

// UserAgent is the canonical UserAgent string for the leatherman.
var UserAgent = "leatherman/"

func init() {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		UserAgent += bi.Main.Version
	} else {
		UserAgent += "devel"
	}
}

// NewRequest returns an *http.Request with the UserAgent header properly set.
func NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	return req, err
}

// Get requests the url with http.DefaultClient, using NewRequest
func Get(url string) (*http.Response, error) {
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}