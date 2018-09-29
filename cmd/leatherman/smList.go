package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"

	"github.com/frioux/leatherman/pkg/sweetmarias"
)

// SMList prints a line of JSON for each Sweet Maria's coffee.
func SMList(_ []string, _ io.Reader) error {
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
			c, _ := sweetmarias.LoadCoffee(url)
			if err != nil {
				fmt.Fprintln(os.Stderr, errors.Wrap(err, "sweetmarias.AllCoffees"))
			}
			wg.Done()
			<-tokens
			_ = e.Encode(c)
			if err != nil {
				fmt.Fprintln(os.Stderr, errors.Wrap(err, "json.Encode"))
			}
		}()
	}

	wg.Wait()

	return nil
}
