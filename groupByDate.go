package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	arg "github.com/alexflint/go-arg"
	"github.com/jeffjen/datefmt"
)

var cmdArgs struct {
	InLayout    string `arg:"-i"`
	GroupLayout string `arg:"-g"`
	OutLayout   string `arg:"-o"`
	ByDay       bool
}

func parseArgs(args []string) error {
	cmdArgs.InLayout = time.RFC3339Nano
	cmdArgs.OutLayout = time.RFC3339Nano
	cmdArgs.GroupLayout = "2006-01-02"
	p, err := arg.NewParser(arg.Config{}, &cmdArgs)
	if err != nil {
		return err
	}
	err = p.Parse(args)
	if err == arg.ErrHelp {
		p.WriteHelp(os.Stdout)
		os.Exit(0)
	}
	if err != nil {
		return err
	}

	return nil
}

// GroupByDate takes dates on stdin in format -i, will group them by format -g,
// and write them in format -o.
func GroupByDate(args []string) {
	err := parseArgs(args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "group-by-date: %v\n", err)
		os.Exit(1)
	}
	groupByDate(os.Stdin, os.Stdout)
}

func parseDate(format, input string) (time.Time, error) {
	if strings.ContainsRune(format, '%') {
		return datefmt.Strptime(format, input)
	}
	return time.Parse(format, input)
}

func formatDate(format string, date time.Time) (string, error) {
	if strings.ContainsRune(format, '%') {
		return datefmt.Strftime(format, date)
	}
	return date.Format(format), nil
}

func groupByDate(i io.Reader, o io.Writer) {
	in := csv.NewReader(i)
	ret := map[string]int{}

	for {
		record, err := in.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't parse line: %v\n", err)
			continue
		}
		date, err := parseDate(cmdArgs.InLayout, string(record[0]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't parse date: %v\n", err)
			continue
		}
		dateStr, err := formatDate(cmdArgs.GroupLayout, date)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't format date: %v\n", err)
			continue
		}
		ret[dateStr]++
	}

	out := csv.NewWriter(o)
	for dateStr, val := range ret {
		date, err := parseDate(cmdArgs.GroupLayout, dateStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't parse internally rendered date: %v\n", err)
			continue
		}
		outDate, err := formatDate(cmdArgs.OutLayout, date)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Coudln't format date for output: %v\n", err)
			continue
		}
		out.Write([]string{outDate, strconv.Itoa(val)})
	}
	out.Flush()
	if err := out.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush: %v\n", err)
		os.Exit(1)
	}
}
