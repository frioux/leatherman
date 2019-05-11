package sshquote

import (
	"fmt"
	"io"
	"os"

	"github.com/frioux/leatherman/pkg/shellquote"
	"golang.org/x/xerrors"
)

/*
Run prepares the arguments such that they can be pasted into an ssh command
safely.  For example:

```bash
$ ssh-quote ls 'f*f'
'sh -c '\''ls '\''\'"''"'f*f'\''\'\'
```

Command: ssh-quote
*/
func Run(args []string, _ io.Reader) error {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [some tokens to quote]\n", args[0])
		os.Exit(1)
	}
	quoted, err := shellquote.Quote(args[1:])
	if err != nil {
		return xerrors.Errorf("Couldn't quote input: %w", err)
	}
	double, _ := shellquote.Quote([]string{"sh", "-c", quoted})
	triple, _ := shellquote.Quote([]string{double})
	fmt.Println(triple)

	return nil
}
