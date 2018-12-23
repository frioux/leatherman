package groupbydate

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	arg "github.com/alexflint/go-arg"
	"github.com/pkg/errors"
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

// Run takes dates on stdin in format -i, will group them by format -g,
// and write them in format -o.
func Run(args []string, stdin io.Reader) error {
	err := parseArgs(args[1:])
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
