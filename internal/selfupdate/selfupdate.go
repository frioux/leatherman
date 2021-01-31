package selfupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/frioux/leatherman/internal/version"
	"github.com/mattn/go-isatty"
	"github.com/ulikunitz/xz"
)

// MaybeUpdate will check for an update and install it immediately.  If
// the LM_GH_TOKEN environment variable is set, it will be treated as a
// github personal access token, instead of anonymous access.
func MaybeUpdate() {
	url := checkUpdate()
	if url == "" {
		return
	}

	doUpdate(url)

}

// AutoUpdate will periodically check for an update (a little over hourly) and
// install it unless the LM_DISABLE_SELFUPDATE env var is set.  If the
// LM_GH_TOKEN environment variable is set this will check more often (closer
// to every minute.)  The Token should have public_repo access.  The token
// should include the `Basic ` prefix.
func AutoUpdate() {
	if os.Getenv("LM_DISABLE_SELFUPDATE") != "" {
		return
	}

	// if there's a tty connected do not do auto-update, since it could
	// restart the process.  Let that happen either to a service or
	// intentionally by the user.
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return
	}

	if isatty.IsTerminal(os.Stdin.Fd()) {
		return
	}

	go func() {
		for {
			if token == "" {
				time.Sleep(time.Duration(rand.Int63n(int64(time.Minute*60))) + time.Minute*30)
			} else {
				time.Sleep(time.Duration(rand.Int63n(int64(time.Second*30))) + time.Second*30)
			}

			MaybeUpdate()
		}

	}()
}

func whichFilename() string {
	switch {
	case runtime.GOARCH == "amd64" && runtime.GOOS == "linux":
		return "leatherman.xz"
	case runtime.GOARCH == "arm" && runtime.GOOS == "linux":
		return "leatherman.arm.xz"
	case runtime.GOARCH == "amd64" && runtime.GOOS == "windows":
		return "leatherman.zip"
	default:
		return ""
	}
}

var mostRecentFailure error

func doUpdate(url string) {
	var err error

	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			mostRecentFailure = err
		}
	}()

	curp, err := os.Executable()
	if err != nil {
		err = fmt.Errorf("couldn't get os.Executable: %w", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var (
		newp = curp + ".new"
		oldp = curp + ".old"
	)

	f, err := os.Create(newp)
	if err != nil {
		err = fmt.Errorf("couldn't make file to update: %w", err)
		return
	}
	resp, err := lmhttp.Get(ctx, url)
	if err != nil {
		err = fmt.Errorf("couldn't download update: %w", err)
		return
	}

	xzr, err := xz.NewReader(resp.Body)
	if err != nil {
		err = fmt.Errorf("couldn't decompress update: %w", err)
		return
	}

	// download to os.Args[0] + ".new"
	_, err = io.Copy(f, xzr)
	if err != nil {
		err = fmt.Errorf("couldn't write update: %w", err)
		return
	}

	err = os.Chmod(newp, 0755)
	if err != nil {
		err = fmt.Errorf("couldn't chmod update: %w", err)
		return
	}

	// move os.Args[0] to + ".old"
	err = os.Rename(curp, oldp)
	if err != nil {
		err = fmt.Errorf("couldn't rename original file: %w", err)
		return
	}

	// move .new
	err = os.Rename(newp, curp)
	if err != nil {
		err = fmt.Errorf("couldn't rename new file: %w", err)
		return
	}

	// new file is copied so anything else that fails shouldn't keep this
	// from happening
	defer func() {
		fmt.Fprintln(os.Stderr, "new version downloaded, exiting to get restarted")
		os.Exit(0)
	}()

	// remove .old
	err = os.Remove(os.Args[0] + ".old")
	if err != nil {
		err = fmt.Errorf("couldn't remove old file: %w", err)
		return
	}
}

var (
	token        = os.Getenv("LM_GH_TOKEN")
	invalidToken bool
)

func checkUpdate() string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var err error
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			mostRecentFailure = err
		}
	}()

	req, err := lmhttp.NewRequest(ctx, "GET", "https://api.github.com/repos/frioux/leatherman/releases/latest", nil)
	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		return ""
	}

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("error finding latest leatherman: %w", err)
		return ""
	}
	defer resp.Body.Close()

	if h := resp.Header.Get("X-RateLimit-Limit"); token != "" && h != "5000" {
		invalidToken = true
		fmt.Fprintf(os.Stderr, "X-RateLimit-Limit wasn't 5000, your auth token might be invalid (was %s); disabling token\n", h)
		b, _ := httputil.DumpResponse(resp, true)
		fmt.Fprintf(os.Stderr, "Full response follows:\n%s", b)
		token = ""
	}

	var found struct {
		Assets []struct {
			Name               string
			BrowserDownloadURL string `json:"browser_download_url"`
		}
	}
	d := json.NewDecoder(resp.Body)
	err = d.Decode(&found)
	if err != nil {
		err = fmt.Errorf("error parsing json for latest leatherman: %w", err)
		return ""
	}

	for _, a := range found.Assets {
		if whichFilename() != a.Name {
			continue
		}

		// same version
		if strings.HasSuffix(a.BrowserDownloadURL, "untagged-"+version.Version+"/"+a.Name) {
			return ""
		}

		return a.BrowserDownloadURL
	}

	return ""
}
