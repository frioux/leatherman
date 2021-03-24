package desktop

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/godbus/dbus"

	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/frioux/leatherman/internal/middleware"
)

func MediaRemote(args []string, _ io.Reader) error {
	var (
		play, pause, next, prev, playPause, selectPlayerFlag bool

		raw string

		timeout time.Duration
	)

	fs := flag.NewFlagSet("media-remote", flag.ContinueOnError)
	fs.BoolVar(&play, "play", false, "calls org.mpris.MediaPlayer2.Player.Play method")
	fs.BoolVar(&pause, "pause", false, "calls org.mpris.MediaPlayer2.Player.Pause method")
	fs.BoolVar(&playPause, "play-pause", false, "calls org.mpris.MediaPlayer2.Player.PlayPause method")
	fs.BoolVar(&next, "next", false, "calls org.mpris.MediaPlayer2.Player.Next method")
	fs.BoolVar(&prev, "prev", false, "calls org.mpris.MediaPlayer2.Player.Previous method")
	fs.StringVar(&raw, "raw", "", "calls whatever method you pass with no arguments")
	fs.BoolVar(&selectPlayerFlag, "select-player", false, "selects  player")
	fs.DurationVar(&timeout, "timeout", 3*time.Minute, "after this amount of time the program will shut down")

	go func() {
		time.Sleep(timeout)
		fmt.Fprintln(os.Stderr, "never chose a player, giving up and shutting down")
		os.Exit(1)
	}()

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("dbus.SessionBus: %w", err)
	}

	s := selectedPlayer()
	players, err := allPlayers(conn)
	if err != nil {
		return err
	}

	var found string
	for _, p := range players {
		found = p
		if s == found {
			break
		}
	}
	if s != "" && s != found {
		if err := os.Remove(selectedPlayerPath); err != nil {
			fmt.Fprintln(os.Stderr, "couldn't remove selected player:", err)
		}
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
		method += "Player.Previous"
	case raw != "":
		method = raw
	case selectPlayerFlag:
		return selectPlayer(conn)
	default:
		fmt.Println(found)
		return nil
	}

	call := obj.Call(method, 0)
	return call.Err
}

var selectedPlayerPath string

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't load current user")
		return
	}

	selectedPlayerPath = filepath.Join(u.HomeDir, ".selected-player.txt")
}

func selectedPlayer() string {
	b, err := os.ReadFile(selectedPlayerPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		fmt.Fprintln(os.Stderr, "couldn't read selected player", err)
		return ""
	}

	return string(b)
}

func storeSelectedPlayer(p string) error {
	fmt.Println(selectedPlayerPath)
	return os.WriteFile(selectedPlayerPath, []byte(p), 0o644)
}

func allPlayers(conn *dbus.Conn) ([]string, error) {
	var s []string
	if err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s); err != nil {
		return nil, fmt.Errorf("ListNames: %w", err)
	}

	var ret []string
	for _, v := range s {
		if !strings.HasPrefix(v, "org.mpris.MediaPlayer2.") {
			continue
		}
		ret = append(ret, v)
	}

	return ret, nil
}

//go:embed media-remote.tmpl
var tmplRaw string
var tmpl *template.Template

//go:embed media-remote.js
var js string

func init() {
	var err error
	tmpl, err = template.New("name").Parse(tmplRaw)
	if err != nil {
		panic(err)
	}
}

func isDBUSNotFound(err error) bool {
	var derr dbus.Error
	return errors.As(err, &derr) &&
		derr.Name == "org.freedesktop.DBus.Error.ServiceUnknown"
}

var (
	statuses  = map[string]bool{}
	statusMux = &sync.Mutex{}
)

func renderOrSelectPlayer(conn *dbus.Conn) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) error {
		// hack to ignore the duplicate request we get from fraidycat
		if r.Header.Get("Accept") == "*/*" {
			fmt.Println("saw derp accept, ignoring")
			return nil
		}
		statusMux.Lock()
		defer statusMux.Unlock()

		rw.Header().Set("Content-Type", "text/html")
		if r.Method == "POST" {
			// resume paused stuff
			for p, s := range statuses {
				obj := conn.Object(p, "/org/mpris/MediaPlayer2")
				meth := "org.mpris.MediaPlayer2.Player.Play"
				if !s {
					meth = "org.mpris.MediaPlayer2.Player.Pause"
				}
				call := obj.Call(meth, 0)
				if err := call.Err; err != nil {
					if isDBUSNotFound(err) {
						continue
					}
					return err
				}
			}
			if err := r.ParseMultipartForm(0); err != nil {
				return err
			}
			p := r.Form.Get("player")
			if p == "" {
				return nil
			}

			if err := storeSelectedPlayer(p); err != nil {
				return err
			}

			go func() {
				time.Sleep(3 * time.Second)
				fmt.Fprintln(os.Stderr, "selected", p, "shutting it down.")
				os.Exit(0)
			}()

			return nil
		}

		all, err := allPlayers(conn)
		if err != nil {
			return err
		}

		for k := range statuses {
			delete(statuses, k)
		}
		for _, p := range all {
			obj := conn.Object(p, "/org/mpris/MediaPlayer2")
			v, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
			if err != nil {
				if isDBUSNotFound(err) {
					continue
				}
				return err
			}
			statuses[p] = v.Value().(string) == "Playing"

			if statuses[p] {
				call := obj.Call("org.mpris.MediaPlayer2.Player.Pause", 0)
				if call.Err != nil {
					if isDBUSNotFound(err) {
						delete(statuses, p)
						continue
					}
					return call.Err
				}
			}
		}

		return tmpl.Execute(rw, all)
	})
}

func pausePlayer(conn *dbus.Conn) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) error {
		if err := r.ParseMultipartForm(0); err != nil {
			return err
		}
		p := r.Form.Get("player")
		if p == "" {
			return nil
		}

		obj := conn.Object(p, "/org/mpris/MediaPlayer2")
		call := obj.Call("org.mpris.MediaPlayer2.Player.Pause", 0)
		if isDBUSNotFound(call.Err) {
			rw.WriteHeader(404)
			return nil
		}
		return call.Err
	})
}

func playPlayer(conn *dbus.Conn) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) error {
		if err := r.ParseMultipartForm(0); err != nil {
			return err
		}
		p := r.Form.Get("player")
		if p == "" {
			return nil
		}

		obj := conn.Object(p, "/org/mpris/MediaPlayer2")
		call := obj.Call("org.mpris.MediaPlayer2.Player.Play", 0)
		if isDBUSNotFound(call.Err) {
			rw.WriteHeader(404)
			return nil
		}
		return call.Err
	})
}

func selectPlayer(conn *dbus.Conn) error {
	mux := http.NewServeMux()
	mux.Handle("/select-player", renderOrSelectPlayer(conn))
	mux.Handle("/play", playPlayer(conn))
	mux.Handle("/pause", pausePlayer(conn))
	mux.Handle("/media-remote.js", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "application/javascript")
		io.Copy(rw, strings.NewReader(js))
	}))
	h := middleware.Adapt(mux, middleware.Log(os.Stdout))
	s := http.Server{Handler: h}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}
	fmt.Fprintln(os.Stderr, "starting on", listener.Addr())
	go func() {
		if err := exec.Command("xdg-open", "http://"+listener.Addr().String()+"/select-player").Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	return s.Serve(listener)
}
