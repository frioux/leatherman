package srv

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

var errARGone = errors.New("auto-reload watcher disappeared")

func doReload(watcher *fsnotify.Watcher, dir string, generation *chan bool) error {
	var timeout <-chan time.Time

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errARGone
			}
			// sink the ship if a root disappears
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				if event.Name == dir {
					return errors.New("deleted root, capsizing")
				}
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				stat, err := os.Stat(event.Name)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					fmt.Fprintf(os.Stderr, "Couldn't stat created thing: %s\n", err)
				} else if stat.IsDir() {
					err := addDir(watcher, event.Name)
					if err != nil {
						fmt.Fprintf(os.Stderr, "failed to watch %s: %s\n", event.Name, err)
					}
				}
			}

			timeout = time.After(time.Second)
		case err, ok := <-watcher.Errors:
			if !ok {
				return errARGone
			}
			fmt.Println("error:", err)
		case <-timeout:
			close(*generation)
			*generation = make(chan bool)
		}

	}
}

func autoReload(h http.Handler, dir string) (handler http.Handler, sinking chan error, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, fmt.Errorf("fsnotify.NewWatcher: %w", err)
	}
	err = addDir(watcher, dir)
	if err != nil {
		return nil, nil, fmt.Errorf("addDir: %w", err)
	}

	generation := make(chan bool)
	reloadErr := make(chan error)
	go func() { reloadErr <- doReload(watcher, dir, &generation) }()

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		f, ok := rw.(http.Flusher)
		if !ok {
			http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		if r.URL.Path == "/_reload" {
			rw.Header().Set("Cache-Control", "no-cache")
			rw.Header().Set("Content-Type", "text/event-stream")

			select {
			case <-generation:
				fmt.Fprintf(rw, "data: Message: reload!!!\n\n")
				f.Flush()
			case <-r.Context().Done():
				// client went away
			}
			return
		} else if r.URL.Path == "/_force_reload" {
			rw.Header().Set("Cache-Control", "no-cache")

			close(generation)
			generation = make(chan bool)
			return
		} else {
			// This is a pretty inefficient way to do this, but
			// it's reliable at least.  Given time and motivation
			// this could be more stream oriented and not buffer
			// the whole response.
			brw := httptest.NewRecorder()

			// Copy headers into buffer
			for h := range rw.Header() {
				brw.Header().Set(h, rw.Header().Get(h))
			}

			// Run handler against buffer
			h.ServeHTTP(brw, r)

			// Copy headers back out
			for h := range brw.Header() {
				rw.Header().Set(h, brw.Header().Get(h))
			}

			rw.Header().Del("Content-Length")

			res := brw.Result()
			defer res.Body.Close()

			rw.WriteHeader(res.StatusCode)

			// Copy body
			if _, err := io.Copy(rw, res.Body); err != nil {
				fmt.Fprintf(os.Stderr, "error writing body: %s\n", err)
			}

			const js = `<script>
			const evtSource = new EventSource("/_reload");
			evtSource.onerror = function(event) {
			  if (event.target.readyState == EventSource.CLOSED) {
			    // refresh page after 2-5s
			    setTimeout(function() { location.reload() }, 2000 + Math.random() * 3000);
			    return;
			  }
			  console.log(event);
			};
			evtSource.onmessage = function(event) { location.reload() }
			</script>`

			if mt, _, _ := mime.ParseMediaType(res.Header.Get("Content-Type")); mt == "text/html" {
				fmt.Fprint(rw, js)
			}
		}
	}), reloadErr, nil
}

func addDir(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if err := watcher.Add(path); err != nil {
			return fmt.Errorf("fsnotify.Watcher.Add: %w", err)
		}
		return nil
	})
}
