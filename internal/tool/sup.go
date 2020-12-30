package notes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/frioux/leatherman/internal/lmhttp"
)

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
