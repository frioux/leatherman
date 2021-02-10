package srv // import "github.com/frioux/leatherman/internal/tool/srv"

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/frioux/leatherman/internal/selfupdate"
)

func Serve(args []string, _ io.Reader) error {
	var (
		port     int
		noreload bool
	)
	fs := flag.NewFlagSet("srv", flag.ContinueOnError)
	fs.IntVar(&port, "port", 0, "port to listen on; default is random")
	fs.BoolVar(&noreload, "no-autoreload", false, "disable auto-reloading")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	dir := "."
	if len(fs.Args()) > 0 {
		dir = fs.Arg(0)
	}

	dir = filepath.Clean(dir)

	ch := make(chan net.Addr)

	go func() {
		addr := <-ch
		fmt.Fprintf(os.Stderr, "Serving %s on %s\n", dir, addr)
	}()

	return serve(!noreload, dir, port, ch)
}

func logReqs(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(os.Stderr, time.Now(), r.URL)

		if r.URL.Path == "/version" {
			selfupdate.Handler.ServeHTTP(rw, r)
			return
		}

		h.ServeHTTP(rw, r)
	})
}

func serve(reload bool, dir string, port int, log chan net.Addr) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	log <- listener.Addr()

	var sinking chan error

	h := http.FileServer(http.Dir(dir))
	if reload {
		h, sinking, err = autoReload(h, dir)
		if err != nil {
			return err
		}
	}

	s := http.Server{Handler: logReqs(h)}

	if reload {
		go func() {
			err = <-sinking
			fmt.Fprintln(os.Stderr, "auto-reload:", err)
			s.Close()
		}()
	}

	return s.Serve(listener)
}
