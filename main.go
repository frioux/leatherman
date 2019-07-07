package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/frioux/amygdala/internal/log"
	"github.com/frioux/amygdala/internal/middleware"
	"github.com/frioux/amygdala/internal/notes"
	"github.com/frioux/amygdala/internal/twilio"
	"golang.org/x/xerrors"
)

var (
	dropboxAccessToken, myCell string
)

var twilioAuthToken, twilioURL []byte

func init() {
	rand.Seed(time.Now().UnixNano())

	dropboxAccessToken = os.Getenv("DROPBOX_ACCESS_TOKEN")
	if dropboxAccessToken == "" {
		panic("DROPBOX_ACCESS_TOKEN is missing")
	}

	myCell = os.Getenv("MY_CELL")
	if myCell == "" {
		myCell = "+15555555555"
	}

	twilioAuthToken = []byte(os.Getenv("TWILIO_AUTH_TOKEN"))
	if len(twilioAuthToken) == 0 {
		twilioAuthToken = []byte("xyzzy")
	}

	twilioURL = []byte(os.Getenv("TWILIO_URL"))
	if len(twilioURL) == 0 {
		twilioURL = []byte("http://localhost:8080/twilio")
	}
}

var port int

var version string

func init() {
	flag.IntVar(&port, "port", 8080, "port to listen on")

	if version == "" {
		version = "unknown"
	}
}

func main() {
	flag.Parse()
	cl := &http.Client{}

	http.Handle("/version", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("content-type", "text/plain")

		bi, ok := debug.ReadBuildInfo()
		if !ok {
			rw.WriteHeader(500)
		}

		fmt.Fprintln(rw, "version:", version)

		for _, dep := range bi.Deps {
			fmt.Printf("%s@%s (%s)\n", dep.Path, dep.Version, dep.Sum)
			if dep.Replace != nil {
				r := dep.Replace
				fmt.Printf("   replaced by %s@%s (%s)\n", r.Path, r.Version, r.Sum)
			}
		}
	}))

	http.Handle("/twilio", middleware.Adapt(receiveSMS(cl, dropboxAccessToken),
		middleware.Log(os.Stdout),
	))

	log.Err(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	os.Exit(1)
}

// receiveSMS handles https://www.twilio.com/docs/sms/twiml
func receiveSMS(cl *http.Client, tok string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			io.WriteString(rw, "Couldn't Parse Form")
			log.Err(xerrors.Errorf("http.Request.ParseForm: %w", err))
			return
		}

		if ok, err := twilio.CheckMAC(twilioAuthToken, twilioURL, r); err != nil || !ok {
			rw.WriteHeader(403)
			if err != nil {
				log.Err(xerrors.Errorf("twilio.CheckMAC: %w", err))
			}
			return
		}

		if r.Form.Get("From") != myCell {
			rw.WriteHeader(http.StatusForbidden)
			io.WriteString(rw, "Wrong Cell\n")
			return
		}

		message := r.Form.Get("Body")
		media, _ := twilio.ExtractMedia(r.Form)

		if message == "" && len(media) == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			io.WriteString(rw, "No Message\n")
			return
		}

		resp, err := notes.Dispatch(cl, tok, message, media)
		if err != nil {
			// normally it's a really bad idea to use other values if the error is
			// non-nil, but care has been taken to propogate cheeky responses even
			// in that situation.
			//
			// Note that the cheeky values won't work unless we return a 200 OK.
			io.WriteString(rw, resp+"\n")
			log.Err(xerrors.Errorf("notes.Dispatch: %w", err))
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain")

		io.WriteString(rw, resp+"\n")
	}
}
