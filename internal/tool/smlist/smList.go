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

/*
Run lists all of the available [Sweet Maria's](https://www.sweetmarias.com/) coffees
as json documents per line.  Here's how you might see the top ten coffees by
rating:

```bash
$ sm-list | jq -r '[.Score, .Title, .URL ] | @tsv' | sort -n | tail -10
```
Command: sm-list */
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
