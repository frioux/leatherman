package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func meeting() bool {
	c := exec.Command("xdotool", "search", "--name", "Meet")
	o, err := c.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "xdotool seach --name Meet: %s\n", err)
		return false
	}

	return len(o) > 0
}

func main() {
	for {
		var red, green, blue int
		if meeting() {
			green = 255
		}
		if sound() {
			red = 255
		}

		colorSpec := fmt.Sprintf("--rgb=%d,%d,%d", red, green, blue)
		err := exec.Command("blink1-tool", colorSpec).Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "blink1-tool %s: %s\n", colorSpec, err)
		}
		time.Sleep(time.Second)
	}
}
