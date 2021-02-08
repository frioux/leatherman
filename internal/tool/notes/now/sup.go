package now

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sort"
	"sync"
	"time"

	"github.com/frioux/leatherman/internal/lmhttp"
)

func loadPeers() ([]string, error) {
	var tsStatus struct {
		Peer map[string]struct {
			DNSName string
		}
	}

	cmd := exec.Command("tailscale", "status", "-json")
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &tsStatus); err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(tsStatus.Peer))
	for _, p := range tsStatus.Peer {
		ret = append(ret, p.DNSName)
	}

	sort.Strings(ret)

	return ret, nil
}

func peerVersion(ctx context.Context, hostname string) (ret string) {
	defer func() { ret = hostname + ": " + ret }()
	resp, err := lmhttp.Get(ctx, "http://"+hostname+":8081/version")
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		return s.Text() // silly, return first line
	}

	return "???"
}

func allVersions(ctx context.Context) []string {
	peers, err := loadPeers()
	if err != nil {
		return []string{"couldn't load peers: " + err.Error()}
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(peers))

	ret := make([]string, len(peers))

	for i, p := range peers {
		go func(i int, p string) {
			defer wg.Done()
			ret[i] = peerVersion(ctx, p)
		}(i, p)
	}

	wg.Wait()

	return ret
}

func sup(rw http.ResponseWriter, req *http.Request) error {
	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	var rpi struct{ Game string }
	go func() {
		defer wg.Done()
		resp, err := lmhttp.Get(ctx, "http://retropie:8081/retropie")
		if err != nil {
			// whyyyyy
			rpi.Game = err.Error()
			return
		}
		defer resp.Body.Close()

		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&rpi); err != nil {
			// ugh wtf
			rpi.Game = "ERR: " + err.Error()
		}
	}()

	wg.Add(1)
	var steamos []byte
	go func() {
		defer wg.Done()
		resp, err := lmhttp.Get(ctx, "http://steamos:8081/steambox")
		if err != nil {
			// I don't like it
			steamos = []byte(err.Error())
			return
		}
		defer resp.Body.Close()

		steamos, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			// I should have thought this through more carefully
			steamos = []byte("ERR: " + err.Error())
		}
	}()

	wg.Add(1)
	var pi400 []byte
	go func() {
		defer wg.Done()
		resp, err := lmhttp.Get(ctx, "http://pi400:8081/x11title")
		if err != nil {
			// I don't like it
			pi400 = []byte(err.Error())
			return
		}
		defer resp.Body.Close()

		pi400, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			// I should have thought this through more carefully
			pi400 = []byte("ERR: " + err.Error())
		} else {
			pi400 = []byte(`<a href="http://pi400:8081/x11shot">` + string(pi400) + `</a>`)
		}
	}()

	wg.Add(1)
	var versions []string
	go func() {
		defer wg.Done()
		versions = allVersions(ctx)
	}()

	wg.Wait()

	fmt.Fprintf(rw, prelude, "now: sup")
	fmt.Fprintf(rw, "retropie: %s<br>\n", rpi.Game)
	fmt.Fprintf(rw, "steamos: %s<br>\n", steamos)
	fmt.Fprintf(rw, "pi400: %s<br>\n", pi400)
	fmt.Fprintln(rw, "\n<br>versions:<br><ul>")
	for _, v := range versions {
		fmt.Fprintf(rw, "<li>%s</li>\n", v)
	}
	fmt.Fprintln(rw, "</ul>")

	return nil
}
