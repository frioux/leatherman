package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func cam() bool {
	c := exec.Command("lsof", "/dev/video0")
	o, err := c.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lsof /dev/video*: %s\n", err)
		if eErr, ok := err.(*exec.ExitError); ok {
			fmt.Fprintf(os.Stderr, "stderr: %s\n", eErr.Stderr)
		}
		return false
	}

	return bytes.Contains(o, []byte("mem"))
}

func main() {
	for {
		var red, green, blue int
		if cam() {
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
