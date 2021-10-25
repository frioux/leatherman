package now

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sort"
	"sync"
	"time"

	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/frioux/leatherman/internal/notes"
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

func peerVersion(ctx context.Context, hostname string) (ret option) {
	defer func() {
		if ret.error != nil {
			ret.error = fmt.Errorf("%s: %w", hostname, ret.error)
		} else {
			ret.value = fmt.Sprintf("%s: %s", hostname, ret.value)
		}
	}()
	resp, err := lmhttp.Get(ctx, "http://"+hostname+":8081/version")
	if err != nil {
		return option{error: err}
	}
	defer resp.Body.Close()

	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		return option{value: s.Text()} // silly, return first line
	}

	return option{error: errors.New("???")}
}

func allVersions(ctx context.Context) []option {
	peers, err := loadPeers()
	if err != nil {
		return []option{{error: fmt.Errorf("couldn't load peers: %w", err)}}
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(peers))

	ret := make([]option, len(peers))

	for i, p := range peers {
		go func(i int, p string) {
			defer wg.Done()
			ret[i] = peerVersion(ctx, p)
		}(i, p)
	}

	wg.Wait()

	return ret
}

func handlerSup(z *notes.Zine) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		rw.Header().Set("Cache-Control", "no-cache")
		ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
		defer cancel()

		v := supVars{HTMLVars: &HTMLVars{Zine: z}}
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := lmhttp.Get(ctx, "http://retropie:8081/retropie")
			if err != nil {
				// whyyyyy
				v.retroPie.error = err
				return
			}
			defer resp.Body.Close()

			var decodeable struct{ Game string }
			dec := json.NewDecoder(resp.Body)
			v.retroPie.error = dec.Decode(&decodeable)
			v.retroPie.value = decodeable.Game
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := lmhttp.Get(ctx, "http://steamos:8081/steambox")
			if err != nil {
				v.steamOS.error = err
				return
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			v.steamOS.value, v.steamOS.error = string(b), err
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := lmhttp.Get(ctx, "http://pi400:8081/x11title")
			if err != nil {
				v.pi400.error = err
				return
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				v.pi400.error = err
			} else {
				v.pi400.value = `<a href="http://pi400:8081/x11shot">` + string(b) + `</a>`
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			v.Versions = allVersions(ctx)
		}()

		wg.Wait()

		return tpl.ExecuteTemplate(rw, "sup.html", v)
	})
}
