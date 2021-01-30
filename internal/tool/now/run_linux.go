// +build linux

package now

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/frioux/leatherman/internal/dropbox"
	_ "modernc.org/sqlite"
)

/*
Serve provides a web interface to parts of my notes stored in
Dropbox.  This should eventually be merged into Amygdala.

Command: notes
*/
func Serve(args []string, _ io.Reader) error {
	var (
		listen string
	)
	fs := flag.NewFlagSet("notes", flag.ContinueOnError)
	fs.StringVar(&listen, "listen", ":0", "location to listen on; default is random")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	cl, err := dropbox.NewClient(dropbox.Client{
		Token: os.Getenv("LM_DROPBOX_TOKEN"),
	})
	if err != nil {
		return err
	}

	generation := make(chan bool)
	z, err := loadDB(cl, "/notes/content/posts/", &generation)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	fmt.Fprintf(os.Stderr, "listening on %s\n", listener.Addr())
	h, err := server(z, &generation)
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	s := http.Server{Handler: h}
	return s.Serve(listener)
}
