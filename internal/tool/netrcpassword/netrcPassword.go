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

	n, err := netrc.Parse(filepath.Join(usr.HomeDir, ".netrc"))
	if err != nil {
		return errors.Wrap(err, "Couldn't parse netrc")
	}

	login := n.MachineAndLogin(args[1], args[2])
	if login == nil {
		return errors.New("Couldn't find login for " + args[2] + "@" + args[1])
	}

	fmt.Println(login.Password)

	return nil
}
