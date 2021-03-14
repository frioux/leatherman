// +build linux

package now

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/frioux/leatherman/internal/lmfs"
)

func Serve(args []string, _ io.Reader) error {
	var (
		listen string
	)
	flags := flag.NewFlagSet("notes", flag.ContinueOnError)
	flags.StringVar(&listen, "listen", ":0", "location to listen on; default is random")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	path := os.Getenv("LM_NOTES_PATH")
	if path == "" {
		return errors.New("must set LM_NOTES_PATH")
	}

	var f fs.FS
	if t := os.Getenv("LM_DROPBOX_TOKEN"); t != "" {
		cl, err := dropbox.NewClient(dropbox.Client{
			Token: os.Getenv("LM_DROPBOX_TOKEN"),
		})
		if err != nil {
			return err
		}
		f = cl.AsFS(context.TODO())
		f, err = fs.Sub(f, path)
		if err != nil {
			return err
		}
	} else {
		f = lmfs.OpenDirFS(path)
	}

	generation := make(chan bool)
	z, err := loadDB(f, &generation)
	if err != nil {
		return fmt.Errorf("loadDB: %w", err)
	}

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	fmt.Fprintf(os.Stderr, "listening on %s\n", listener.Addr())
	h, err := server(f, z, &generation)
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	s := http.Server{Handler: h}
	return s.Serve(listener)
}
