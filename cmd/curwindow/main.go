package main

// curwindow prints the name of the currently selected window

import (
	"fmt"
	"os"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	X, err := xgbutil.NewConn()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	w, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		return fmt.Errorf("coudln't ActiveWindowGet: %s", err)
	}
	name, err := ewmh.WmNameGet(X, w)
	if err != nil {
		return fmt.Errorf("coudln't WmNameGet: %s", err)
	}
	fmt.Println(name)

	return nil
}
