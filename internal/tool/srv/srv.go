package srv // import "github.com/frioux/leatherman/internal/tool/srv"

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"golang.org/x/xerrors"
)

// Serve an html directory tree, with index.html being shown for dirs.
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

func serve(dir string, log chan net.Addr) error {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return xerrors.Errorf("net.Listen: %w", err)
	}

	log <- listener.Addr()

	return http.Serve(listener, http.FileServer(http.Dir(dir)))
}
