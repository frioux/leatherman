package status

import (
	"fmt"
	"os/exec"
	"sync"
)

func manageLight(camMu, soundMu *sync.Mutex, c *cam, s *sound) error {
	camMu.Lock()
	defer camMu.Unlock()

	soundMu.Lock()
	defer soundMu.Unlock()

	var red, green, blue int
	if c.value {
		green = 255
	}
	if s.value {
		red = 255
	}

	colorSpec := fmt.Sprintf("--rgb=%d,%d,%d", red, green, blue)
	return exec.Command("blink1-tool", colorSpec).Run()
}
