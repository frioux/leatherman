package backlight

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"golang.org/x/xerrors"
)

const path = "/sys/class/backlight/intel_backlight"

// Run modifies backlight brightness, assuming first arg is a percent.
func Run(args []string, _ io.Reader) error {
	err := os.Chdir(path)
	if err != nil {
		return xerrors.Errorf("Couldn't chdir: %w", err)
	}

	if len(args) != 2 {
		return xerrors.Errorf("Usage: %s <change-as-integer-percent>", args[0])
	}

	change, err := strconv.Atoi(args[1])
	if err != nil {
		return xerrors.Errorf("Couldn't parse arg: %w", err)
	}

	return run(change)
}

func run(change int) error {
	max, err := getMaxBrightness()
	if err != nil {
		return xerrors.Errorf("Couldn't getMaxBrightness: %w", err)
	}

	cur, err := getCurBrightness()
	if err != nil {
		return xerrors.Errorf("getCurBrightness: %w", err)
	}

	var toWrite = change*max/100 + cur
	fmt.Fprintf(os.Stderr, "%d = %d*%d/100 + %d\n", toWrite, change, max, cur)
	if toWrite < 0 {
		toWrite = 0
	}
	if toWrite > max {
		toWrite = max
	}

	file, err := os.OpenFile("./brightness", os.O_RDWR, 0)
	if err != nil {
		return xerrors.Errorf("Couldn't open brightness for writing: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Setting brightness to %d\n", toWrite)

	_, err = file.WriteString(fmt.Sprintf("%d\n", toWrite))
	if err != nil {
		return xerrors.Errorf("file.WriteString: %w", err)
	}
	err = file.Close()
	if err != nil {
		return xerrors.Errorf("Couldn't write brightness: %w", err)
	}

	return nil
}

func getMaxBrightness() (int, error) {
	file, err := os.Open("./max_brightness")
	if err != nil {
		return 0, xerrors.Errorf("couldn't open max_brightness: %s", err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	line, err := r.ReadSlice('\n')
	if err != nil {
		return 0, xerrors.Errorf("couldn't read line: %s", err)
	}

	i, err := strconv.Atoi(string(line[:len(line)-1]))
	if err != nil {
		return 0, xerrors.Errorf("couldn't parse line: %s", err)
	}

	return i, nil
}

func getCurBrightness() (int, error) {
	file, err := os.Open("./brightness")
	if err != nil {
		return 0, xerrors.Errorf("couldn't open brightness: %s", err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	line, err := r.ReadSlice('\n')
	if err != nil {
		return 0, xerrors.Errorf("couldn't read line: %s", err)
	}

	i, err := strconv.Atoi(string(line[:len(line)-1]))
	if err != nil {
		return 0, xerrors.Errorf("couldn't parse line: %s", err)
	}

	return i, nil
}
