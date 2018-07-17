package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const path = "/sys/class/backlight/intel_backlight"

// Backlight modifies backlight brightness, assuming first arg is a percent.
func Backlight(args []string) {
	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't chdir into %s: %s\n", path, err)
		os.Exit(1)
	}

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <change-as-integer-percent>\n", args[0])
		os.Exit(1)
	}

	change, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Coudln't parse arg: %s", err)
		os.Exit(1)
	}

	max, err := getMaxBrightness()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't getMaxBrightness: %s\n", err)
		os.Exit(1)
	}

	cur, err := getCurBrightness()

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
		fmt.Fprintf(os.Stderr, "Couldn't open brightness for writing: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Setting brightness to %d\n", toWrite)

	file.WriteString(fmt.Sprintf("%d\n", toWrite))
	err = file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't write brightness: %s\n", err)
		os.Exit(1)
	}
}

func getMaxBrightness() (int, error) {
	file, err := os.Open("./max_brightness")
	if err != nil {
		return 0, fmt.Errorf("couldn't open max_brightness: %s", err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	line, err := r.ReadSlice('\n')
	if err != nil {
		return 0, fmt.Errorf("couldn't read line: %s", err)
	}

	i, err := strconv.Atoi(string(line[:len(line)-1]))
	if err != nil {
		return 0, fmt.Errorf("couldn't parse line: %s", err)
	}

	return i, nil
}

func getCurBrightness() (int, error) {
	file, err := os.Open("./brightness")
	if err != nil {
		return 0, fmt.Errorf("couldn't open brightness: %s", err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	line, err := r.ReadSlice('\n')
	if err != nil {
		return 0, fmt.Errorf("couldn't read line: %s", err)
	}

	i, err := strconv.Atoi(string(line[:len(line)-1]))
	if err != nil {
		return 0, fmt.Errorf("couldn't parse line: %s", err)
	}

	return i, nil
}
