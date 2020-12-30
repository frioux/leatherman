package notes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
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

	cmd := exec.Command("tailscale status")
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

	return ret, nil
}

func allVersions() []byte {
	peers, err := loadPeers()
	if err != nil {
		return []byte("couldn't load peers: " + err.Error())
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(peers))

	buf := &bytes.Buffer{}

	for _, p := range peers {
		go func(p string) {

		}(p)

	}

	return buf.Bytes()
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

	wg.Wait()

	fmt.Fprintf(rw, prelude, "now: sup")
	fmt.Fprintf(rw, "retropie: %s<br>", rpi.Game)
	fmt.Fprintf(rw, "steamos: %s<br>", steamos)
	fmt.Fprintf(rw, "pi400: %s<br>", pi400)

	return nil
}
