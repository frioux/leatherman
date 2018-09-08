package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// GenPass bcrypts the first argument with the second argument rounds.
func GenPass(args []string, _ io.Reader) error {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s $password [$cost]\n", args[0])
		os.Exit(1)
	}
	pass := args[1]
	cost, err := strconv.Atoi(args[2])
	if err != nil {
		return errors.Wrap(err, "couldn't parse "+args[2])
	}

	t0 := time.Now()
	out, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	if err != nil {
		return errors.Wrap(err, "couldn't hash password")
	}
	fmt.Println(string(out))
	fmt.Fprintf(os.Stderr, "%0.2fs elapsed\n", time.Since(t0).Seconds())

	return nil
}
