package smlist

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"

	"github.com/frioux/leatherman/pkg/sweetmarias"
)

// Run prints a line of JSON for each Sweet Maria's coffee.
func Run(_ []string, _ io.Reader) error {
	wg := sync.WaitGroup{}

	tokens := make(chan struct{}, 10)

	coffees, err := sweetmarias.AllCoffees()
	if err != nil {
		return errors.Wrap(err, "sweetmarias.AllCoffees")
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
				fmt.Fprintln(os.Stderr, errors.Wrap(err, "sweetmarias.LoadCoffee"))
				return
			}
			err = e.Encode(c)
			if err != nil {
				fmt.Fprintln(os.Stderr, errors.Wrap(err, "json.Encode"))
				return
			}
		}()
	}

	wg.Wait()

	return nil
}
