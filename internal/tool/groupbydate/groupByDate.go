package groupbydate

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/frioux/leatherman/pkg/datefmt"
	"golang.org/x/xerrors"
)

var cmdArgs struct {
	InLayout    string
	GroupLayout string
	OutLayout   string
}

func parseArgs(args []string) error {
	flags := flag.NewFlagSet("group-by-date", flag.ExitOnError)

	flags.StringVar(&cmdArgs.InLayout, "in", time.RFC3339Nano, "format to parse when reading input")
	flags.StringVar(&cmdArgs.OutLayout, "out", time.RFC3339Nano, "format to write when writing output")
	flags.StringVar(&cmdArgs.GroupLayout, "group", "2006-01-02", "format to group by internally")

	err := flags.Parse(args[1:])
	if err != nil {
		return xerrors.Errorf("flags.Parse: %w", err)
	}

	return nil
}

func parseDate(format, input string) (time.Time, error) {
	if strings.ContainsRune(format, '%') {
		format = datefmt.TranslateFormat(format)
	}
	return time.Parse(format, input)
}

func formatDate(format string, date time.Time) (string, error) {
	if strings.ContainsRune(format, '%') {
		format = datefmt.TranslateFormat(format)
	}
	return date.Format(format), nil
}

/*
Run creates time series data by counting lines and grouping them by a given date
format.  takes dates on stdin in format -i, will group them by format -g, and
write them in format -o.

Command: group-by-date
*/
func Run(args []string, stdin io.Reader) error {
	err := parseArgs(args)
	if err != nil {
		return xerrors.Errorf("Couldn't parse args: %w", err)
	}
	return groupByDate(stdin, os.Stdout)
}

func groupByDate(i io.Reader, o io.Writer) error {
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
		date, err := parseDate(cmdArgs.InLayout, record[0])
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
		err = out.Write([]string{outDate, strconv.Itoa(val)})
		if err != nil {
			return xerrors.Errorf("csv.Write: %w", err)
		}
	}
	out.Flush()
	if err := out.Error(); err != nil {
		return xerrors.Errorf("Failed to flush: %w", err)
	}

	return nil
}
