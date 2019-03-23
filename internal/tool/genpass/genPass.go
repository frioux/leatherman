package genpass

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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
			return errors.Wrap(err, "couldn't parse "+args[2])
		}
	}

	return run(os.Stdout, os.Stderr, pass, cost)
}

func run(w, wErr io.Writer, password string, cost int) error {
	t0 := time.Now()
	out, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return errors.Wrap(err, "couldn't hash password")
	}
	fmt.Fprintln(w, string(out))
	fmt.Fprintf(wErr, "%0.2fs elapsed\n", time.Since(t0).Seconds())

	return nil
}
