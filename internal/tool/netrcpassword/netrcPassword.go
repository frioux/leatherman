package netrcpassword

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/frioux/leatherman/pkg/netrc"
	"github.com/pkg/errors"
)

// Run prints passsword for passed machine and login
func Run(args []string, _ io.Reader) error {
	if len(args) != 3 {
		fmt.Println("Usage:\n\tnetrc-password $machine $login")
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "Couldn't get current user")
	}

	password, err := run(filepath.Join(usr.HomeDir, ".netrc"), args[1], args[2])
	if err != nil {
		return errors.Wrap(err, "Couldn't load password")
	}

	fmt.Println(password)

	return nil
}

func run(path, machine, user string) (string, error) {
	n, err := netrc.Parse(path)
	if err != nil {
		return "", errors.Wrap(err, "Couldn't parse netrc")
	}

	login, ok := n.MachineAndLogin(machine, user)
	if !ok {
		return "", errors.New("Couldn't find login for " + user + "@" + machine)
	}

	return login.Password, nil
}
