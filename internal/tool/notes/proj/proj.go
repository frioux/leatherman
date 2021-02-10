package proj

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
)

var proj, vimSessions, notes, smartcd string

func init() {
	u, err := user.Current()
	if err != nil {
		panic("Couldn't get current user: " + err.Error())
	}
	vimSessions = u.HomeDir + "/.vvar/sessions"
	notes = u.HomeDir + "/code/notes/content/posts"
	smartcd = u.HomeDir + "/.smartcd/scripts"

	proj = os.Getenv("PROJ")
}

func Proj(args []string, _ io.Reader) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s init | vim | note", args[0])
	}
	switch args[1] {
	case "init":
		return initialize(args[1:])
	case "vim":
		return vim()
	case "note":
		return errors.New("nyi")
	default:
		return errors.New("unknown subcommand " + args[1])
	}

}

// XXX would be nice to make this use exec instead; I know go doens't
// technically support that but I also know it does actually work.
func vim() error {
	if proj == "" {
		return errors.New("cannot infer session without PROJ set")
	}

	vim := exec.Command("vim", "-S", vimSessions+"/"+proj)
	vim.Stdin = os.Stdin
	vim.Stdout = os.Stdout
	vim.Stderr = os.Stderr
	return vim.Run()
}
