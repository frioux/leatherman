package main

import (
	"fmt"
	"github.com/frioux/netrc"
	"os"
	"os/user"
	"path/filepath"
)

// NetrcPassword prints passsword for passed machine and login
func NetrcPassword(args []string) {
	usr, err := user.Current()

	if len(args) != 3 {
		fmt.Println("Usage:\n\tnetrc-password $machine $login")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Couldn't get current user:", err)
		os.Exit(-1)
	}

	n, err := netrc.Parse(filepath.Join(usr.HomeDir, ".netrc"))

	if err != nil {
		fmt.Println("Couldn't parse netrc", err)
		os.Exit(-2)
	}

	login := n.MachineAndLogin(args[1], args[2])
	if login == nil {
		fmt.Println("Couldn't find login for", args[1], "and", args[2])
		os.Exit(2)
	}

	fmt.Println(login.Get("password"))
}
