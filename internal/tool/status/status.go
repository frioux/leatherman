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

/*
Status runs a little web server that surfaces status information related to how
I'm using the machine.  For example, it can say which window is active, what
firefox tabs are loaded, if the screen is locked, etc.  The main benefit of the
tool is that it caches the values returned.

In the background, it interact swith the [blink(1)](http://blink1.thingm.com/).
It turns the light green when I'm in a meeting and red when audio is playing.

Command: status
*/
func Status(args []string, _ io.Reader) error {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		for _, s := range []string{"cam", "curwindow", "sound", "tabs", "vpn", "locked"} {
			fmt.Fprintf(rw, " * /%s\n", s)
		}
	}))

	mux.Handle("/locked", &cacher{reloadEvery: time.Second, value: &locked{}, mu: &sync.Mutex{}})
	mux.Handle("/curwindow", &cacher{reloadEvery: time.Second, value: &curWindow{}, mu: &sync.Mutex{}})
	mux.Handle("/tabs", &cacher{reloadEvery: time.Second * 2, value: &tabs{}, mu: &sync.Mutex{}})
	mux.Handle("/vpn", &cacher{reloadEvery: time.Second, value: &vpn{}, mu: &sync.Mutex{}})
	mux.Handle("/retropie", &cacher{reloadEvery: time.Second, value: &retropie{}, mu: &sync.Mutex{}})

	s := &sound{}
	soundCacher := &cacher{reloadEvery: time.Second, value: s, mu: &sync.Mutex{}}
	mux.Handle("/sound", soundCacher)

	c := &cam{}
	camCacher := &cacher{reloadEvery: time.Minute, value: c, mu: &sync.Mutex{}}
	mux.Handle("/cam", camCacher)

	go func() {
		for {
			if err := manageLight(soundCacher.mu, camCacher.mu, c, s); err != nil {
				fmt.Fprintf(os.Stderr, "couldn't manage light: %s\n", err)
			}
			time.Sleep(time.Second)
		}
	}()

	port := 8081
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	srv := http.Server{Handler: logReqs(mux)}

	return srv.Serve(listener)
}

func logReqs(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(os.Stderr, time.Now(), r.URL)
		h.ServeHTTP(rw, r)
	})
}

func exec1Fail(cmd string, rest ...string) (bool, error) {
	c := exec.Command(cmd, rest...)
	_, err := c.Output()
	if err != nil {
		eErr := &exec.ExitError{}
		if errors.As(err, &eErr) {
			if eErr.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
