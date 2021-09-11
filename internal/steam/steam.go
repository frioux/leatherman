package steam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/frioux/leatherman/internal/lmhttp"
)

// AppIDs automatically maintains a map of steam appids to steam names.
type AppIDs struct {
	mu  sync.Mutex
	raw map[int]string

	LastLoad time.Time
}

// Autoload repopulates the internal map about daily.
func (a *AppIDs) Autoload() {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	for {
		if err := a.Load(context.Background()); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// reload some time between 0 and 24h from now
		time.Sleep(time.Duration(rnd.Float32() * float32(time.Hour*24)))
	}
}

func (a *AppIDs) Load(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if time.Now().Sub(a.LastLoad) < time.Hour*24 {
		return nil
	}

	if a.raw == nil {
		a.raw = map[int]string{}
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	resp, err := lmhttp.Get(ctx, "https://api.steampowered.com/ISteamApps/GetAppList/v2/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("non-200 status")
	}

	var remoteData struct {
		Applist struct {
			Apps []struct {
				AppID int
				Name  string
			}
		}
	}

	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&remoteData); err != nil {
		return err
	}

	for k := range a.raw {
		delete(a.raw, k)
	}

	for _, app := range remoteData.Applist.Apps {
		a.raw[app.AppID] = app.Name
	}

	a.LastLoad = time.Now()

	return nil
}

func (a *AppIDs) App(appid int) string {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.raw[appid]
}
