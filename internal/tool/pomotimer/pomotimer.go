package pomotimer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

const clear = "\r\x1b[J"

// Run starts a timer for 25m or the duration expressed in the first
// argument.
func Run(args []string, stdin io.Reader) error {
	duration := 25 * time.Minute
	if len(args) > 1 {
		var err error
		duration, err = time.ParseDuration(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "pomotimer: couldn't parse duration: %s", err)
			os.Exit(1)
		}
	}

	initialSeconds := int(duration.Seconds())

	setProcessName("PT" + formatTime(initialSeconds))
	fmt.Print("[p]ause [r]eset abort[!]\n\n")

	// disable input buffering
	err := exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	if err != nil {
		return errors.Wrap(err, "couldn't disable input buffering")
	}
	// do not display entered characters on the screen
	err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	if err != nil {
		return errors.Wrap(err, "couldn't hide input")
	}

	// restore the echoing state when exiting
	defer func() {
		_ = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
		_ = exec.Command("tmux", "setw", "automatic-rename", "on").Run()
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	kb := make(chan rune)
	go kbChan(kb, stdin)

	running := true
	secondsRemaining := initialSeconds
LOOP:
	for {
		if secondsRemaining < 0 {
			fmt.Println(clear + "Take a break!\a")
			break
		}
		select {
		case <-ticker.C:
			if !running {
				continue
			}
			if secondsRemaining%30 == 0 {
				setProcessName("PT" + formatTime(secondsRemaining))
			}
			fmt.Print(clear+formatTime(secondsRemaining), " remaining")
			secondsRemaining--
		case key := <-kb:
			if key == 'p' {
				running = !running
			} else if key == 'r' {
				secondsRemaining = initialSeconds
			} else if key == '!' {
				fmt.Println(clear + "aborting timer!")
				break LOOP
			}
		}
	}

	return nil
}

// XXX messy
func kbChan(keys chan<- rune, stdin io.Reader) {
	var b = make([]byte, 1)
	for {
		_, err := stdin.Read(b)
		if err != nil {
			break
		}
		keys <- rune(b[0])
	}
}

func formatTime(s int) string {
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}
