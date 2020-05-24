package status

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type value interface {
	load() error
	render(http.ResponseWriter)
}

type cacher struct {
	mu      sync.Mutex
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

type locked struct{ value bool }

func (l *locked) load() error {
	cmd := exec.Command("pgrep", "i3lock")
	_, err := cmd.Output()
	if err != nil {
		eErr := &exec.ExitError{}
		if errors.As(err, &eErr) {
			if eErr.ExitCode() == 1 {
				l.value = false
				return nil
			}
			return err
		}
		return err
	}
	l.value = true
	return nil
}

func (l *locked) render(rw http.ResponseWriter) { fmt.Fprintf(rw, "%t\n", l.value) }

type curWindow struct{ value string }

func (l *curWindow) load() error {
	// curwindow is an auxiliary leatherman tool
	cmd := exec.Command("curwindow")
	b, err := cmd.Output()
	if err != nil {
		return err
	}
	l.value = string(b)
	return nil
}

func (l *curWindow) render(rw http.ResponseWriter) { fmt.Fprint(rw, l.value) }

type tabs struct{ value string }

func (l *tabs) load() error {
	// XXX use a function since we can?
	cmd := exec.Command("dump-mozlz4", os.Getenv("MOZ_RECOVERY"))
	b, err := cmd.Output()
	if err != nil {
		return err
	}
	l.value = string(b)
	return nil
}

func (l *tabs) render(rw http.ResponseWriter) { fmt.Fprintf(rw, "%s\n", l.value) }

/*
Status runs a little web server that surfaces status information related to how
I'm using the machine.  For example, it can say which window is active, what
firefox tabs are loaded, if the screen is locked, etc.  The main benefit of the
tool is that it caches the values returned.

Command: status
*/
func Status(args []string, _ io.Reader) error {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		for _, s := range []string{"curwindow", "tabs", "locked"} {
			fmt.Fprintf(rw, " * /%s\n", s)
		}
	}))

	mux.Handle("/locked", &cacher{reloadEvery: time.Second, value: &locked{}})
	mux.Handle("/curwindow", &cacher{reloadEvery: time.Second, value: &curWindow{}})
	mux.Handle("/tabs", &cacher{reloadEvery: time.Second * 2, value: &tabs{}})

	// sound := cachedValue{}
	// mux.Handle("/sound", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	// 	sound.mu.Lock()
	// 	defer sound.mu.Unlock()

	// }))

	// camera := cachedValue{}
	// mux.Handle("/camera", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	// 	camera.mu.Lock()
	// 	defer camera.mu.Unlock()

	// }))

	port := 8081
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	s := http.Server{Handler: mux}

	return s.Serve(listener)
}
