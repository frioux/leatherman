package srv

import (
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

const js = `<script>
const evtSource = new EventSource("/_reload");
evtSource.onmessage = function(event) { location.reload() }
</script>`

func autoReload(h http.Handler, dir string) (http.Handler, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("fsnotify.NewWatcher: %w", err)
	}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
	if err != nil {
		return nil, fmt.Errorf("filepath.Walk: %w", err)
	}

	var timeout <-chan time.Time
	generation := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					stat, err := os.Stat(event.Name)
					if err != nil {
						if os.IsNotExist(err) {
							continue
						}
						fmt.Fprintf(os.Stderr, "Couldn't stat created thing: %s\n", err)
					} else if stat.IsDir() {
						err := watcher.Add(event.Name)
						if err != nil {
							fmt.Fprintf(os.Stderr, "failed to watch %s: %s\n", event.Name, err)
						}
					}
				}

				timeout = time.After(time.Second)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			case <-timeout:
				close(generation)
				generation = make(chan bool)
			}

		}
	}()

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		f, ok := rw.(http.Flusher)
		if !ok {
			http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		if r.URL.Path == "/_reload" {
			rw.Header().Set("Cache-Control", "no-cache")
			rw.Header().Set("Content-Type", "text/event-stream")

			<-generation
			fmt.Fprintf(rw, "data: Message: reload!!!\n\n")
			f.Flush()
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

			if mt, _, _ := mime.ParseMediaType(res.Header.Get("Content-Type")); mt == "text/html" {
				fmt.Fprint(rw, js)
			}
		}
	}), nil
}
