package steamsrv

import (
	"bufio"
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/frioux/leatherman/internal/lmfav"
	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/frioux/leatherman/internal/steam"
)

//go:embed templates/*
var templateFS embed.FS

var templates = template.Must(template.New("tmpl").ParseFS(templateFS, "templates/*"))

func Serve(args []string, _ io.Reader) error {
	var (
		port       int
		shotPrefix string
		logPrefix  string
	)

	u, err := user.Current()
	if err != nil {
		return err
	}

	fs := flag.NewFlagSet("steamsrv", flag.ContinueOnError)
	fs.IntVar(&port, "port", 8080, "port to listen on; default is 8080")
	fs.StringVar(&shotPrefix, "screenshot-prefix", "", "screenshot path to serve, required, should probably be something like ~/.local/share/Steam/userdata/?*")
	fs.StringVar(&logPrefix, "log-prefix", u.HomeDir+"/.local/share/Steam/logs", "location of steam logs")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if shotPrefix == "" {
		return errors.New("path-prefix is required")
	}

	a := &steam.AppIDs{}
	if err := a.Load(context.Background()); err != nil {
		return err
	}

	go a.Autoload()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}
	s := http.Server{Handler: steamHandler(os.DirFS(logPrefix), os.DirFS(shotPrefix), a)}

	return s.Serve(listener)
}

func mustAtoi(s string) int {
	ret, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return ret
}

type steamAppState struct{ Name, Date, State string }

// lineMatcher matches lines about game state changes
//                                     1              2                    3
//                                     date           appid                state
var lineMatcher = regexp.MustCompile(`^\[(.*?)\] AppID (\d+) state changed : (.*)$`)

func fileToSteamAppStates(a *steam.AppIDs, f fs.File) ([]steamAppState, error) {
	ret := []steamAppState{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		m := lineMatcher.FindStringSubmatch(s.Text())
		if len(m) == 0 {
			continue
		}

		for _, state := range strings.Split(m[3], ",") {
			if state == "" {
				continue
			}
			ret = append(ret, steamAppState{
				State: state,
				Name:  a.App(mustAtoi(m[2])),
				Date:  m[1],
			})
		}
	}

	if s.Err() != nil {
		return nil, s.Err()
	}

	return ret, nil
}

func logHandler(logFS fs.FS, a *steam.AppIDs) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) error {
		var data struct {
			Logs  []steamAppState
			Title string
		}
		data.Title = "Steam Logs"

		f, err := logFS.Open("content_log.previous.txt")
		if err != nil {
			return err
		}
		defer f.Close()

		ss, err := fileToSteamAppStates(a, f)
		if err != nil {
			return err
		}

		for _, s := range ss {
			if s.State == "App Running" {
				data.Logs = append(data.Logs, s)
			}
		}

		f, err = logFS.Open("content_log.txt")
		if err != nil {
			return err
		}
		defer f.Close()

		ss, err = fileToSteamAppStates(a, f)
		if err != nil {
			return err
		}

		for _, s := range ss {
			if s.State == "App Running" {
				data.Logs = append(data.Logs, s)
			}
		}

		rw.Header().Add("Content-Type", "text/html")
		if err := templates.ExecuteTemplate(rw, "logs.html", data); err != nil {
			return err
		}

		return nil
	})
}

func screenshotHandler(a *steam.AppIDs, fss fs.FS) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) error {
		s := screenshots{fss: fss}
		var appid int
		if a := r.URL.Query().Get("appid"); a != "" {
			appid = mustAtoi(a)
		}

		var data struct {
			Name, Title string

			Apps   map[int]string
			Images []screenshot
		}
		data.Apps = map[int]string{}
		var err error
		data.Images, err = s.Screenshots(appid)
		if err != nil {
			return err
		}

		rw.Header().Add("Content-Type", "text/html")
		if appid == 0 {
			data.Title = "Steam Apps"
			for _, image := range data.Images {
				data.Apps[image.AppID] = a.App(image.AppID)
			}
			if err := templates.ExecuteTemplate(rw, "apps.html", data); err != nil {
				return err
			}
		} else {
			data.Name = a.App(appid)
			data.Title = data.Name + " Screenshots"
			if err := templates.ExecuteTemplate(rw, "shots.html", data); err != nil {
				return err
			}
		}

		return nil
	})
}

func steamHandler(logFS, shotFS fs.FS, a *steam.AppIDs) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/files/", lmhttp.TrimHandlerPrefix("/files", http.FileServer(http.FS(shotFS))))
	mux.Handle("/", screenshotHandler(a, shotFS))
	mux.Handle("/log/", logHandler(logFS, a))
	mux.Handle("/favicon/", lmfav.Emoji('🎮'))

	return mux
}

