package now

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/frioux/leatherman/internal/dropbox"
)

var errARGone = errors.New("auto-reload channel closed")

func doReload(ch <-chan struct{}, dir string, generation *chan bool) error {
	var timeout <-chan time.Time

	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return errARGone
			}
			timeout = time.After(time.Second)
		case <-timeout:
			close(*generation)
			*generation = make(chan bool)
		}

	}
}

func autoReload(db dropbox.Client, h http.Handler, generation *chan bool, dir string) (handler http.Handler, err error) {
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
			case <-*generation:
				fmt.Fprintf(rw, "data: Message: reload!!!\n\n")
				f.Flush()
			case <-r.Context().Done():
				// client went away
			}
			return
		} else if r.URL.Path == "/_force_reload" {
			rw.Header().Set("Cache-Control", "no-cache")

			close(*generation)
			*generation = make(chan bool)
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

			res := brw.Result()
			defer res.Body.Close()

			// Copy headers back out
			for h := range res.Header {
				rw.Header().Set(h, res.Header.Get(h))
			}

			rw.Header().Del("Content-Length")

			// When the file is not found it's usually that we
			// caught a file event before the file was recreated,
			// so we turn the 404 page into text/html so the
			// reloader JS gets injected.
			if res.StatusCode == 404 {
				rw.Header().Set("Content-Type", "text/html")
			}

			rw.WriteHeader(res.StatusCode)

			// Copy body
			if _, err := io.Copy(rw, res.Body); err != nil {
				fmt.Fprintf(os.Stderr, "error writing body: %s\n", err)
			}

			const js = `<script>
			function sleep(n) {
			  n = 1000*n;
			  return new Promise(done => {
			    setTimeout(() => {
			      done();
			    }, n);
			  });
			}

			let sinking = false;

			// Disable the reload when we navigate away from the page.
			//
			// This event happens when we try to download files in Firefox, but
			// the reload triggered on the navigated-away-from page which is still
			// in the tab behind the download prompt actually causes the download
			// window to close.
			window.addEventListener('beforeunload', function(event) {
			  sinking = true;
			});
			const evtSource = new EventSource("/_reload");
			evtSource.onerror = async function(event) {
			  if (!sinking && event.target.readyState == EventSource.CLOSED) {
			    // the server went away, poll till it's back, then reload.
                            let i = 0;
                            while(true) {
                              try {
                                await fetch('/');
                                location.reload();
                                break;
                              } catch(e) {
                                await sleep(Math.random() * i**2);
	                        if (i < 7) { // ~a minute
	                          i++;
	                        }
                              }
                            }
			    return;
			  }
			  console.log(event);
			};

			evtSource.onmessage = function(event) { location.reload() }
			</script>`

			if mt, _, _ := mime.ParseMediaType(rw.Header().Get("Content-Type")); mt == "text/html" {
				fmt.Fprint(rw, js)
			}
		}
	}), nil
}
