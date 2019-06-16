package srv // import "github.com/frioux/leatherman/internal/tool/srv"

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/xerrors"
)

/*
Serve will serve a directory over http; takes a single optional parameter which
is the dir to serve, default is `.`:

```bash
$ srv ~
Serving /home/frew on [::]:21873
```

Command: srv
*/
func Serve(args []string, _ io.Reader) error {
	dir := "."
	if len(args) > 1 {
		dir = args[1]
	}

	ch := make(chan net.Addr)

	go func() {
		addr := <-ch
		fmt.Fprintf(os.Stderr, "Serving %s on %s\n", dir, addr)
	}()

	return serve(dir, ch)
}

func logReqs(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(os.Stderr, time.Now(), r.URL)
		h.ServeHTTP(rw, r)
	})
}

func serve(dir string, log chan net.Addr) error {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return xerrors.Errorf("net.Listen: %w", err)
	}

	log <- listener.Addr()

	return http.Serve(listener, logReqs(http.FileServer(http.Dir(dir))))
}
