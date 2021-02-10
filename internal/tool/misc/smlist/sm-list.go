package smlist

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/frioux/leatherman/pkg/sweetmarias"
)

func Run(_ []string, _ io.Reader) error {
	wg := sync.WaitGroup{}

	tokens := make(chan struct{}, 10)

	// 15s for the index
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	urls, err := sweetmarias.AllCoffees(ctx)
	if err != nil {
		return fmt.Errorf("sweetmarias.AllCoffees: %w", err)
	}

	coffees := make([]sweetmarias.Coffee, len(urls))

	// 5m for all the details
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	for i, url := range urls {
		wg.Add(1)
		tokens <- struct{}{}
		i, url := i, url
		go func() {
			defer func() { wg.Done(); <-tokens }()
			var err error
			coffees[i], err = sweetmarias.LoadCoffee(ctx, url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "sweetmarias.LoadCoffee: %s\n", err)
				return
			}
		}()
	}

	wg.Wait()

	sort.Slice(coffees, func(i, j int) bool { return coffees[i].URL < coffees[j].URL })

	e := json.NewEncoder(os.Stdout)
	for _, c := range coffees {
		if err := e.Encode(c); err != nil {
			return fmt.Errorf("json.Encode: %w", err)
		}
	}

	return nil
}
