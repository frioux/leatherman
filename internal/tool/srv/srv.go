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
Serve will serve a directory over http; takes an optional parameter which is
the dir to serve, and -port if you care to choose the serving port,
default is `.`:

```bash
$ srv ~
Serving /home/frew on [::]:21873
```

Command: srv
*/
func Serve(args []string, _ io.Reader) error {
	var port int
	fs := flag.NewFlagSet("srv", flag.ContinueOnError)
	fs.IntVar(&port, "port", 0, "port to listen on; default is random")
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

	return serve(dir, port, ch)
}

func logReqs(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(os.Stderr, time.Now(), r.URL)
		h.ServeHTTP(rw, r)
	})
}

func serve(dir string, port int, log chan net.Addr) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	log <- listener.Addr()

	return http.Serve(listener, logReqs(autoReload(http.FileServer(http.Dir(dir)), dir)))
}
