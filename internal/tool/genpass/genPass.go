package genpass

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

// Run bcrypts the first argument with the second argument rounds.
func Run(args []string, _ io.Reader) error {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s $password [$cost]\n", args[0])
		os.Exit(1)
	}
	pass := args[1]
	var cost int
	var err error

	if len(args) == 3 {
		cost, err = strconv.Atoi(args[2])
		if err != nil {
			return xerrors.Errorf("couldn't parse %s: %w", args[2], err)
		}
	}

	return run(os.Stdout, os.Stderr, pass, cost)
}

func run(w, wErr io.Writer, password string, cost int) error {
	t0 := time.Now()
	out, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return xerrors.Errorf("couldn't hash password: %w", err)
	}
	fmt.Fprintln(w, string(out))
	fmt.Fprintf(wErr, "%0.2fs elapsed\n", time.Since(t0).Seconds())

	return nil
}
