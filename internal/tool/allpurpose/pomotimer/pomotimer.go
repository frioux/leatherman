package pomotimer

import (
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"time"
)

const clear = "\r\x1b[J"

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

	// disable input buffering
	err := exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	if err != nil {
		return fmt.Errorf("couldn't disable input buffering: %w", err)
	}
	// do not display entered characters on the screen
	err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	if err != nil {
		return fmt.Errorf("couldn't hide input: %w", err)
	}

	// restore the echoing state when exiting
	defer func() {
		_ = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
		setProcessName("")
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	deadline := time.Now().Add(duration)
	remaining := deadline.Sub(time.Now())

	fmt.Print("[p]ause [r]eset abort[!]\n\n")
	setProcessName("PT" + formatTime(remaining+time.Second))
	fmt.Print(clear+formatTime(remaining+time.Second), " remaining")

	kb := make(chan rune)
	go kbChan(kb, stdin)

	running := true
LOOP:
	for {
		remaining = deadline.Sub(time.Now())
		select {
		case <-ticker.C:
			if time.Now().After(deadline) {
				fmt.Println(clear + "Take a break!\a")
				break LOOP
			}
			if int(math.Round(remaining.Seconds()))%30 == 1 {
				setProcessName("PT" + formatTime(remaining))
			}
			fmt.Print(clear+formatTime(remaining), " remaining")
		case key := <-kb:
			switch key {
			case 'p':
				{
					if running {
						ticker.Stop()
					} else {
						ticker = time.NewTicker(time.Second)
						deadline = time.Now().Add(remaining)
					}
					running = !running
				}
			case 'r':
				deadline = time.Now().Add(duration)
			case '!', 'q':
				{
					fmt.Println(clear + "aborting timer!")
					break LOOP
				}
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

func formatTime(r time.Duration) string {
	s := int(math.Round(r.Seconds())) - 1
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

func setProcessName(name string) {
	fmt.Print("\033k" + name + "\033\\")
}
