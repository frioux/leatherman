package srv // import "github.com/frioux/leatherman/internal/tool/srv"

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

/*
Serve will serve a directory over http, injecting javascript to have pages
reload when files change.

It takes an optional dir to serve, the default is `.`.

```bash
$ srv ~
Serving /home/frew on [::]:21873
```

You can pass -port if you care to choose the listen port.

It will set up file watchers and trigger page reloads (via SSE,) this
functionality can be disabled with -no-autoreload.

```bash
$ srv -port 8080 -no-autoreload ~
Serving /home/frew on [::]:8080
```

Command: srv
*/
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
		h.ServeHTTP(rw, r)
	})
}

func serve(reload bool, dir string, port int, log chan net.Addr) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	log <- listener.Addr()

	h := http.FileServer(http.Dir(dir))
	if reload {
		h, err = autoReload(h, dir)
		if err != nil {
			return err
		}
	}

	return http.Serve(listener, logReqs(h))
}
