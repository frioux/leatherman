package srv // import "github.com/frioux/leatherman/internal/tool/srv"

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// Serve an html directory tree, with index.html being shown for dirs.
func Serve(args []string, _ io.Reader) error {
	dir := "."
	if len(args) > 1 {
		dir = args[1]
	}

	return serve(dir, os.Stderr)
}

func serve(dir string, log io.Writer) error {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}

	fmt.Fprintf(log, "Serving %s on %s\n", dir, listener.Addr())

	return http.Serve(listener, http.FileServer(http.Dir(dir)))
}
