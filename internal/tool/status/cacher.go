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

	if v.timeout.Before(time.Now()) {
		if err := v.value.load(); err != nil {
			rw.WriteHeader(500)
			fmt.Fprint(rw, err.Error())
			return
		}

		v.timeout = time.Now().Add(v.reloadEvery)
	}

	v.value.render(rw)
}
