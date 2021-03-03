package desktop

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/godbus/dbus"
)

func MediaRemote(args []string, _ io.Reader) error {
	var (
		play, pause, next, prev, playPause bool

		raw string
	)

	fs := flag.NewFlagSet("media-remote", flag.ContinueOnError)
	fs.BoolVar(&play, "play", false, "calls org.mpris.MediaPlayer2.Player.Play method")
	fs.BoolVar(&pause, "pause", false, "calls org.mpris.MediaPlayer2.Player.Pause method")
	fs.BoolVar(&playPause, "play-pause", false, "calls org.mpris.MediaPlayer2.Player.PlayPause method")
	fs.BoolVar(&next, "next", false, "calls org.mpris.MediaPlayer2.Player.Next method")
	fs.BoolVar(&prev, "prev", false, "calls org.mpris.MediaPlayer2.Player.Prev method")
	fs.StringVar(&raw, "raw", "", "calls whatever method you pass with no arguments")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("dbus.SessionBus: %w", err)
	}

	var s []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		return fmt.Errorf("ListNames: %w", err)
	}

	var found string
	for _, v := range s {
		if !strings.HasPrefix(v, "org.mpris.MediaPlayer2.") {
			continue
		}
		found = v
	}

	obj := conn.Object(found, "/org/mpris/MediaPlayer2")

	method := "org.mpris.MediaPlayer2."
	switch {
	case play:
		method += "Player.Play"
	case pause:
		method += "Player.Pause"
	case playPause:
		method += "Player.PlayPause"
	case next:
		method += "Player.Next"
	case prev:
		method += "Player.Prev"
	case raw != "":
		method = raw
	default:
		fmt.Println(found)
		return nil
	}

	call := obj.Call(method, 0)
	return call.Err
}
