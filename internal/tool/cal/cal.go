package cal

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/frioux/leatherman/pkg/ics"
)

/*
Today shows today's agenda.

Command: today
*/
func Today(_ []string, _ io.Reader) error {
	resp, err := http.Get("https://calendar.google.com/calendar/ical/frew%40ziprecruiter.com/private-386d64fed4b2d31e8b462c7af52c68e2/basic.ics")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't http.Get: %s\n", err)
		os.Exit(1)
	}

	cal, err := ics.Parse(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't ics.ParseCalendar: %s", err)
		os.Exit(1)
	}

	// now := time.Date(2019, 6, 10, 8, 12, 12, 0, time.Local)
	// tonight := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
	for _, e := range cal.Events {
		// if e.Start.After(now) && e.End.Before(tonight) {
		fmt.Println(e.Summary, e.Start, e.End)
		// }
		// fmt.Println(e.GetProperty(ics.ComponentPropertyDtStart))
	}

	return nil
}
