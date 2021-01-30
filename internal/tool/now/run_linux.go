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
		load   bool
	)
	fs := flag.NewFlagSet("notes", flag.ContinueOnError)
	fs.StringVar(&listen, "listen", ":0", "location to listen on; default is random")
	fs.BoolVar(&load, "load", false, "load the db, for testing reasons?")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	if load {
		cl, err := dropbox.NewClient(dropbox.Client{
			Token: os.Getenv("LM_DROPBOX_TOKEN"),
		})
		if err != nil {
			return err
		}

		_, err = loadDB(cl, "/notes/content/posts/")
		return err
	}

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	h, err := server()
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	s := http.Server{Handler: h}
	return s.Serve(listener)
}
