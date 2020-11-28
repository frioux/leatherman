package selfupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/frioux/leatherman/internal/version"
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

	go func() {
		rand.Seed(time.Now().UnixNano() & int64(os.Getpid()) & int64(os.Getppid()))
		for {
			if os.Getenv("LM_GH_TOKEN") == "" {
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

func doUpdate(url string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var (
		curp = os.Args[0]
		newp = os.Args[0] + ".new"
		oldp = os.Args[0] + ".old"
	)

	f, err := os.Create(newp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't make file to update: %s\n", err)
		return
	}
	resp, err := lmhttp.Get(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't download update: %s\n", err)
		return
	}

	xzr, err := xz.NewReader(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't decompress update: %s\n", err)
		return
	}

	// download to os.Args[0] + ".new"
	if _, err := io.Copy(f, xzr); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't write update: %s\n", err)
		return
	}

	if err := os.Chmod(newp, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't chmod update: %s\n", err)
		return
	}

	// move os.Args[0] to + ".old"
	if err := os.Rename(curp, oldp); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't rename original file: %s\n", err)
		return
	}

	// move .new
	if err := os.Rename(newp, curp); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't rename new file: %s\n", err)
		return
	}

	// new file is copied so anything else that fails shouldn't keep this
	// from happening
	defer func() {
		fmt.Fprintln(os.Stderr, "new version downloaded, exiting to get restarted")
		os.Exit(0)
	}()

	// remove .old
	if err := os.Remove(os.Args[0] + ".old"); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't remove old file: %s\n", err)
		return
	}
}

func checkUpdate() string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req, err := lmhttp.NewRequest(ctx, "GET", "https://api.github.com/repos/frioux/leatherman/releases/latest", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating request: %s\n", err)
		return ""
	}

	if e := os.Getenv("LM_GH_TOKEN"); e != "" {
		req.Header.Set("Authorization", e)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finding latest leatherman: %s\n", err)
		return ""
	}
	defer resp.Body.Close()

	var found struct {
		Assets []struct {
			Name               string
			BrowserDownloadURL string `json:"browser_download_url"`
		}
	}
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&found); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing json for latest leatherman: %s\n", err)
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
