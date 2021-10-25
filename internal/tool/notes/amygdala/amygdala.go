package amygdala

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/frioux/leatherman/internal/log"
	"github.com/frioux/leatherman/internal/middleware"
	"github.com/frioux/leatherman/internal/notes"
	"github.com/frioux/leatherman/internal/twilio"
)

func Amygdala(args []string, _ io.Reader) error {
	var (
		dropboxAccessToken, myCell, version string
		twilioAuthToken, twilioURL          string
		port                                int
	)

	dropboxAccessToken = os.Getenv("LM_DROPBOX_TOKEN")
	if dropboxAccessToken == "" {
		return errors.New("LM_DROPBOX_TOKEN is missing")
	}

	myCell = os.Getenv("LM_MY_CELL")
	if myCell == "" {
		myCell = "+15555555555"
	}

	twilioAuthToken = os.Getenv("LM_TWILIO_TOKEN")
	if len(twilioAuthToken) == 0 {
		twilioAuthToken = "xyzzy"
	}

	twilioURL = os.Getenv("LM_TWILIO_URL")
	if len(twilioURL) == 0 {
		twilioURL = "http://localhost:8080/twilio"
	}

	fs := flag.NewFlagSet("amgydala", flag.ContinueOnError)
	fs.IntVar(&port, "port", 8080, "port to listen on")

	if version == "" {
		version = "unknown"
	}
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	cl := &http.Client{}

	mux := http.NewServeMux()

	mux.Handle("/version", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")
		rw.Header().Set("Cache-Control", "no-cache")

		bi, ok := debug.ReadBuildInfo()
		if !ok {
			rw.WriteHeader(500)
		}

		fmt.Fprintln(rw, "version:", version)

		for _, dep := range bi.Deps {
			fmt.Fprintf(rw, "%s@%s (%s)\n", dep.Path, dep.Version, dep.Sum)
			if dep.Replace != nil {
				r := dep.Replace
				fmt.Fprintf(rw, "   replaced by %s@%s (%s)\n", r.Path, r.Version, r.Sum)
			}
		}
	}))

	mux.Handle("/twilio", receiveSMS(cl, dropboxAccessToken, twilioAuthToken, twilioURL, myCell))

	h := middleware.Adapt(mux, middleware.Log(os.Stdout))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), h)
}

// receiveSMS handles https://www.twilio.com/docs/sms/twiml
func receiveSMS(cl *http.Client, tok, twilioAuthToken, twilioURL, myCell string) http.HandlerFunc {
	rules, err := notes.NewRules(tok)
	if err != nil {
		panic(err)
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Cache-Control", "no-cache")
		if err := r.ParseForm(); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			io.WriteString(rw, "Couldn't Parse Form")
			log.Err(fmt.Errorf("http.Request.ParseForm: %w", err))
			return
		}

		if ok, err := twilio.CheckMAC([]byte(twilioAuthToken), []byte(twilioURL), r); err != nil || !ok {
			rw.WriteHeader(403)
			if err != nil {
				log.Err(fmt.Errorf("twilio.CheckMAC: %w", err))
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

		resp, err := rules.Dispatch(message, media)
		if err != nil {
			// normally it's a really bad idea to use other values if the error is
			// non-nil, but care has been taken to propogate cheeky responses even
			// in that situation.
			//
			// Note that the cheeky values won't work unless we return a 200 OK.
			io.WriteString(rw, resp+"\n")
			log.Err(fmt.Errorf("notes.Dispatch: %w", err))
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain")

		io.WriteString(rw, resp+"\n")
	}
}
