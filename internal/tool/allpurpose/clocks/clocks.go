package clocks

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
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
	switch cmpDates(relativeThere, relativeHere) {
	case 0:
		day = "today"
	case 1:
		day = "tomorrow"
	case -1:
		day = "yesterday"
	}
	// I can't figure out why I need two tabs at the end or why the final column
	// isn't right aligned :(
	return l + "\t" + day + "\t" + relativeThere.Format("15:04\t3:04 PM") + "\t\t" + offsetStr
}

/*
Run shows my personal, digital, wall of clocks.  Pass one or more timezone names
to choose which timezones are shown.

```bash
clocks Africa/Johannesburg America/Los_Angeles Europe/Copenhagen
```
*/
func Run(args []string, _ io.Reader) error {
	if len(args) > 1 && args[1] == "-h" {
		fmt.Println("my personal, digital, wall of clocks")
		return nil
	}

	now := time.Now().In(time.Local)

	zones := []string{"Local", "America/Los_Angeles", "America/Chicago", "America/New_York", "Asia/Jerusalem", "UTC"}
	if len(args) > 1 {
		zones = args[1:]
	}
	run(now, zones, os.Stdout)

	return nil
}

func run(now time.Time, zones []string, out io.Writer) {
	w := tabwriter.NewWriter(out, 0, 8, 2, ' ', tabwriter.AlignRight)
	for _, tz := range zones {
		fmt.Fprintln(w, t(now, tz))
	}
	w.Flush()
}
