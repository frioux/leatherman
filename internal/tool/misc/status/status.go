package status

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/frioux/leatherman/internal/selfupdate"
)

func Status(args []string, _ io.Reader) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var (
		listen                string
		enableLightManagement bool
	)
	fs.StringVar(&listen, "listen", ":8081", "addres:port to listen on")
	fs.BoolVar(&enableLightManagement, "manage-light", false, "manage the blink(1) light")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	mux := lmhttp.NewClearMux()

	mux.Handle("/version", selfupdate.Handler)
	mux.Handle("/locked", &cacher{reloadEvery: time.Second, value: &locked{}, mu: &sync.Mutex{}})
	mux.Handle("/curwindow", &cacher{reloadEvery: time.Second, value: &curWindow{}, mu: &sync.Mutex{}})
	mux.Handle("/tabs", &cacher{reloadEvery: time.Second * 2, value: &tabs{}, mu: &sync.Mutex{}})
	mux.Handle("/vpn", &cacher{reloadEvery: time.Second, value: &vpn{}, mu: &sync.Mutex{}})
	mux.Handle("/retropie", &cacher{reloadEvery: time.Second, value: &retropie{}, mu: &sync.Mutex{}})
	mux.Handle("/steambox", &cacher{reloadEvery: time.Second, value: &steambox{}, mu: &sync.Mutex{}})
	mux.Handle("/x11title", &cacher{reloadEvery: 0, value: &x11title{}, mu: &sync.Mutex{}})
	mux.Handle("/x11shot", &cacher{value: &x11shot{}, mu: &sync.Mutex{}})
	mux.Handle("/favicon.ico", faviconHandler())

	s := &sound{}
	soundCacher := &cacher{reloadEvery: time.Second, value: s, mu: &sync.Mutex{}}
	mux.Handle("/sound", soundCacher)

	c := &cam{}
	camCacher := &cacher{reloadEvery: time.Minute, value: c, mu: &sync.Mutex{}}
	mux.Handle("/cam", camCacher)

	if enableLightManagement {
		go func() {
			for {
				if err := manageLight(soundCacher.mu, camCacher.mu, c, s); err != nil {
					fmt.Fprintf(os.Stderr, "couldn't manage light: %s\n", err)
				}
				time.Sleep(time.Second)
			}
		}()
	}

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}
	fmt.Fprintln(os.Stderr, "status listening on", listener.Addr())

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
