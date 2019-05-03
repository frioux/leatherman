package smlist

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/frioux/leatherman/pkg/sweetmarias"
	"golang.org/x/xerrors"
)

// Run prints a line of JSON for each Sweet Maria's coffee.
func Run(_ []string, _ io.Reader) error {
	wg := sync.WaitGroup{}

	tokens := make(chan struct{}, 10)

	coffees, err := sweetmarias.AllCoffees()
	if err != nil {
		return xerrors.Errorf("sweetmarias.AllCoffees: %w", err)
	}

	e := json.NewEncoder(os.Stdout)

	for _, url := range coffees {
		wg.Add(1)
		tokens <- struct{}{}
		url := url
		go func() {
			defer func() { wg.Done(); <-tokens }()
			c, err := sweetmarias.LoadCoffee(url)
			if err != nil {
				fmt.Fprintln(os.Stderr, xerrors.Errorf("sweetmarias.LoadCoffee: %w", err))
				return
			}
			err = e.Encode(c)
			if err != nil {
				fmt.Fprintln(os.Stderr, xerrors.Errorf("json.Encode: %w", err))
				return
			}
		}()
	}

	wg.Wait()

	return nil
}
