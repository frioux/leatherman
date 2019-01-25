package groupbydate

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
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
		return errors.Wrap(err, "flags.Parse")
	}

	return nil
}

// Run takes dates on stdin in format -i, will group them by format -g,
// and write them in format -o.
func Run(args []string, stdin io.Reader) error {
	err := parseArgs(args)
	if err != nil {
		return errors.Wrap(err, "Couldn't parse args")
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
			return errors.Wrap(err, "csv.Write")
		}
	}
	out.Flush()
	if err := out.Error(); err != nil {
		return errors.Wrap(err, "Failed to flush")
	}

	return nil
}
