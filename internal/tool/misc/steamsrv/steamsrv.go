package steamsrv

import (
	"bufio"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/frioux/leatherman/internal/lmhttp"
)

//go:embed templates/*
var templateFS embed.FS

var templates = template.Must(template.New("tmpl").ParseFS(templateFS, "templates/*"))

type appIDs struct {
	mu  sync.Mutex
	raw map[int]string

	LastLoad time.Time
}

func (a *appIDs) Autoload() {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	for {
		if err := a.Load(context.Background()); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// reload some time between 0 and 24h from now
		time.Sleep(time.Duration(rnd.Float32() * float32(time.Hour*24)))
	}
}

func (a *appIDs) Load(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if time.Now().Sub(a.LastLoad) < time.Hour*24 {
		return nil
	}

	if a.raw == nil {
		a.raw = map[int]string{}
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	resp, err := lmhttp.Get(ctx, "https://api.steampowered.com/ISteamApps/GetAppList/v2/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("non-200 status")
	}

	var remoteData struct {
		Applist struct {
			Apps []struct {
				AppID int
				Name  string
			}
		}
	}

	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&remoteData); err != nil {
		return err
	}

	for k := range a.raw {
		delete(a.raw, k)
	}

	for _, app := range remoteData.Applist.Apps {
		a.raw[app.AppID] = app.Name
	}

	a.LastLoad = time.Now()

	return nil
}

func (a *appIDs) App(appid int) string {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.raw[appid]
}

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

	a := &appIDs{}
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

// shotPattern decomposes steam screenshots filenames
//                                     1             2                   3          4     5    6     7     8      9
//                                     ?            appid               year      month  day   hour  min  sec     i
//                                    760           319630              2021       04    15    21    05    01     1
var shotPattern = regexp.MustCompile(`([^/]+)/remote/([^/]+)/screenshots/(\d\d\d\d)(\d\d)(\d\d)(\d\d)(\d\d)(\d\d)_(\d+).jpg`)

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

func fileToSteamAppStates(a *appIDs, f fs.File) ([]steamAppState, error) {
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

func logHandler(logFS fs.FS, a *appIDs) http.Handler {
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

func screenshotHandler(a *appIDs, fss fs.FS) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) error {
		filenames, err := fs.Glob(fss, "*/remote/*/screenshots/*.jpg")
		if err != nil {
			return err
		}
		appid := r.URL.Query().Get("appid")

		type Image struct {
			Name, Thumbnail string

			Date time.Time
		}
		var data struct {
			Name, Title string

			Apps   map[int]string
			Images []Image
		}
		data.Apps = map[int]string{}

		for _, filename := range filenames {
			m := shotPattern.FindStringSubmatch(filename)
			if len(m) == 0 {
				fmt.Fprintf(os.Stderr, "path didn't match pattern: %s\n", filename)
				continue
			}

			seenAppID := mustAtoi(m[2])
			data.Apps[seenAppID] = a.App(seenAppID)
			date := time.Date(
				// year         month                       day
				mustAtoi(m[3]), time.Month(mustAtoi(m[4])), mustAtoi(m[5]),
				// hour         minute          second
				mustAtoi(m[6]), mustAtoi(m[7]), mustAtoi(m[8]),
				// ns timezone
				0, time.Local)

			if m[2] == appid {
				thumbnail := fmt.Sprintf("%s/remote/%s/screenshots/thumbnails/%s%s%s%s%s%s_%s.jpg", m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9])

				if f, err := fss.Open(thumbnail); err == nil {
					f.Close()
				} else {
					thumbnail = ""
				}
				data.Images = append(data.Images, Image{
					Name:      filename,
					Thumbnail: thumbnail,
					Date:      date,
				})
			}
		}

		sort.Slice(data.Images, func(i, j int) bool { return data.Images[i].Date.Before(data.Images[j].Date) })
		rw.Header().Add("Content-Type", "text/html")
		if appid == "" {
			data.Title = "Steam Apps"
			if err := templates.ExecuteTemplate(rw, "apps.html", data); err != nil {
				return err
			}
		} else {
			data.Name = a.App(mustAtoi(appid))
			data.Title = data.Name + " Screenshots"
			if err := templates.ExecuteTemplate(rw, "shots.html", data); err != nil {
				return err
			}
		}

		return nil
	})
}

func steamHandler(logFS, shotFS fs.FS, a *appIDs) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/files/", lmhttp.TrimHandlerPrefix("/files", http.FileServer(http.FS(shotFS))))
	mux.Handle("/", screenshotHandler(a, shotFS))
	mux.Handle("/log/", logHandler(logFS, a))
	mux.Handle("/favicon/", faviconHandler('🎮'))

	return mux
}

func faviconHandler(favicon rune) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprintf(rw, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y=".9em" font-size="90">%c</text></svg>`, favicon)
	})
}
