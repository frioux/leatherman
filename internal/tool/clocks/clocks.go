package clocks

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func cmpDates(there, here time.Time) int8 {
	tDate := there.Truncate(time.Duration(24) * time.Hour)
	hDate := here.Truncate(time.Duration(24) * time.Hour)
	if tDate == hDate {
		return 0
	} else if tDate.Before(hDate) {
		return -1
	} else {
		return 1
	}
}

func t(now time.Time, l string) string {
	loc, err := time.LoadLocation(l)
	if err != nil {
		log.Fatal(err)
	}
	thereNow := now.In(loc)

	relativeHere := time.Date(now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		time.UTC,
	)
	relativeThere := time.Date(thereNow.Year(), thereNow.Month(), thereNow.Day(),
		thereNow.Hour(), thereNow.Minute(), thereNow.Second(), thereNow.Nanosecond(),
		time.UTC,
	)

	offset := relativeThere.Sub(relativeHere).Hours()

	offsetStr := strconv.FormatFloat(offset, 'g', -1, 64)
	if offset >= 0 {
		offsetStr = "+" + offsetStr
	}

	day := "wtf"
	switch cmpDates(relativeHere, relativeThere) {
	case 0:
		day = "today     "
	case 1:
		day = "tomorrow  "
	case -1:
		day = "yesterday "
	}
	return day + " \t" + relativeThere.Format("15:04\t03:04 PM") + "\t" + offsetStr
}

// Run shows my personal, digital, wall of clocks.
func Run(args []string, _ io.Reader) error {
	if len(args) > 1 && args[1] == "-h" {
		fmt.Println("my personal, digital, wall of clocks")
		return nil
	}

	now := time.Now().In(time.Local)

	run(now, os.Stdout)

	return nil
}

func run(now time.Time, w io.Writer) {
	fmt.Fprintln(w, "here : "+t(now, "Local"))
	fmt.Fprintln(w, "L.A. : "+t(now, "America/Los_Angeles"))
	fmt.Fprintln(w, "MS/TX: "+t(now, "America/Chicago"))
	fmt.Fprintln(w, "rjbs : "+t(now, "America/New_York"))
	fmt.Fprintln(w, "riba : "+t(now, "Europe/Berlin"))
	fmt.Fprintln(w, "seo  : "+t(now, "Asia/Jerusalem"))
	fmt.Fprintln(w, "UTC  : "+t(now, "UTC"))
}
