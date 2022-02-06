package status

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type value interface {
	load() error
	render(http.ResponseWriter)
}

type cacher struct {
	mu      *sync.Mutex
	timeout time.Time

	reloadEvery time.Duration
	value       value
}

func (v *cacher) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.reloadEvery == 0 {
		if err := v.value.load(); err != nil {
			rw.WriteHeader(500)
			fmt.Fprint(rw, err.Error())
			return
		}

		rw.Header().Set("Cache-Control", "no-cache")
		v.value.render(rw)
		return
	}

	if v.timeout.Before(time.Now()) {
		if err := v.value.load(); err != nil {
			rw.WriteHeader(500)
			fmt.Fprint(rw, err.Error())
			return
		}

		v.timeout = time.Now().Add(v.reloadEvery)
	}

	rw.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d, immutable", int(v.timeout.Sub(time.Now()).Seconds())))
	v.value.render(rw)
}
