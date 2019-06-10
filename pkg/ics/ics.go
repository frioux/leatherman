package ics

import (
	"bufio"
	"errors"
	"io"
	"net/textproto"
	"regexp"
	"strconv"
	"time"

	"github.com/frioux/leatherman/pkg/timeutil"
)

type Calendar struct {
	Events []Event
}

type params []string

func (p params) Param(key string) (string, bool) {
	for i := 0; i < len(p); i += 2 {
		if p[i] == key {
			return p[i+1], true
		}
	}
	return "", false
}

type Event struct {
	Start, End time.Time
	Summary    string
}

type parsedEvent map[string][]parsedLine

func (p parsedEvent) simpleValue(key string) (string, bool) {
	lines, ok := p[key]
	if !ok || len(lines) == 0 {
		return "", false
	}

	return lines[0].value, true
}

// func addMonth(t time.Time) time.Time {
// 	t.
// }

func recur(e Event, p params) (func() (Event, bool), error) {
	// https://tools.ietf.org/html/rfc5545#page-40
	freq, ok := p.Param("FREQ")
	if !ok {
		return nil, errors.New("FREQ is required")
	}

	// see https://tools.ietf.org/html/rfc5545#page-41

	interval := 1
	intervalString, ok := p.Param("INTERVAL")
	if ok {
		var err error
		interval, err = strconv.Atoi(intervalString)
		if err != nil {
			interval = 1
			// XXX warn error???
		}
	}

	until := time.Now().AddDate(1, 0, 0) // set a sensible end, just in case.
	untilString, ok := p.Param("UNTIL")
	if ok {
		if len(untilString) == 8 {
			t, err := time.ParseInLocation("20060102", untilString, time.Local)
			if err != nil {
				// XXX
			} else {
				until = t
			}
		} else {
			t, err := time.ParseInLocation("20060102T150405Z", untilString, time.Local)
			if err != nil {
				// XXX
			} else {
				until = t
			}
		}
	}
	if _, ok := p.Param("COUNT"); ok {
		return nil, errors.New("COUNT not supported")
	}
	if _, ok := p.Param("BYSECOND"); ok {
		return nil, errors.New("BYSECOND not supported")
	}
	if _, ok := p.Param("BYDAY"); ok {
		return nil, errors.New("BYDAY not supported")
	}

	// see https://tools.ietf.org/html/rfc5545#page-42

	if _, ok := p.Param("BYMONTHDAY"); ok {
		return nil, errors.New("BYMONTHDAY not supported")
	}
	if _, ok := p.Param("BYYEARDAY"); ok {
		return nil, errors.New("BYYEARDAY not supported")
	}
	if _, ok := p.Param("BYWEEKNO"); ok {
		return nil, errors.New("BYWEEKNO not supported")
	}
	if _, ok := p.Param("BYMONTH"); ok {
		return nil, errors.New("BYMONTH not supported")
	}
	if _, ok := p.Param("WKST"); ok {
		return nil, errors.New("WKST not supported")
	}

	// see https://tools.ietf.org/html/rfc5545#page-43

	// bySetPosition, ok := p.Param("BYSETPOS")
	switch freq {
	case "WEEKLY":
		return func() (Event, bool) {
			if e.Start.After(until) {
				return Event{}, false
			}
			weekdayName, ok := p.Param("BYDAY")
			if !ok {
				return Event{}, false // errors.New("need BYDAY for now")
			}
			var weekday time.Weekday
			switch weekdayName {
			case "SU":
				weekday = time.Sunday
			case "MO":
				weekday = time.Monday
			case "TU":
				weekday = time.Tuesday
			case "WE":
				weekday = time.Wednesday
			case "TH":
				weekday = time.Thursday
			case "FR":
				weekday = time.Friday
			case "SA":
				weekday = time.Saturday
			default: // ignoring numberic prefixes and multiple values for now
			}
			for i := 0; i < interval; i++ {
				if e.Start.Weekday() == weekday {
					e.Start = e.Start.AddDate(0, 0, 7)
				} else {
					e.Start = timeutil.JumpTo(e.Start, weekday)
				}
			}
			return e, true
		}, nil
	// case "MONTHLY":
	// 	return func() (Event, bool) {
	// 		return e, true
	// 	}, nil
	// case "YEARLY":
	// 	return func() (Event, bool) {
	// 		e.Start = e.Start.AddDate(1, 0, 0)
	// 		return e, true
	// 	}, nil
	default:
		return nil, errors.New("unsupported FREQ: " + freq)
	}

}

func (p parsedEvent) Event() ([]Event, error) {
	summary, ok := p.simpleValue("SUMMARY")
	if !ok {
		return nil, errors.New("No SUMMARY")
	}
	r := Event{Summary: summary}

	start, err := parseDate(p["DTSTART"][0])
	if err != nil {
		return nil, err
	}
	r.Start = start

	end, err := parseDate(p["DTEND"][0])
	if err != nil {
		return nil, err
	}
	r.End = end

	ret := []Event{r}

	// see https://tools.ietf.org/html/rfc5545#page-123
	// and https://tools.ietf.org/html/rfc5545#section-3.3.10
	rrule, ok := p.simpleValue("RRULE")
	if !ok {
		return ret, nil
	}
	gen, err := recur(r, parseParams(rrule))
	if err != nil {
		return nil, err
	}
	for {
		t, ok := gen()
		if !ok {
			break
		}
		ret = append(ret, t)
	}

	return ret, nil
}

type parsedLine struct {
	name, value string
	params
}

var lineMatcher = regexp.MustCompile("^(?i)([A-Z0-9-]+)(;[^:]+)?(?::(.*))?$")
var paramMatcher = regexp.MustCompile("(?i)([A-Z0-9-]+)=([^:;.]+)")

func parseParams(s string) params {
	var params params
	p := paramMatcher.FindAllStringSubmatch(s, -1)
	for _, param := range p {
		params = append(params, param[1:]...)
	}

	return params
}

func parseLine(s string) (parsedLine, error) {
	m := lineMatcher.FindStringSubmatch(s)
	if len(m) != 4 {
		return parsedLine{}, errors.New("oh no")
	}

	return parsedLine{
		name:   m[1],
		params: parseParams(m[2]),
		value:  m[3],
	}, nil
}

func parseDate(l parsedLine) (time.Time, error) {
	if locName, ok := l.Param("TZID"); ok {
		loc, err := time.LoadLocation(locName)
		t, err := time.ParseInLocation("20060102T150405", l.value, loc)
		if err != nil {
			return t, err
		}

		return t.In(time.Local), nil
	}

	if value, ok := l.Param("VALUE"); ok && value == "DATE" {
		return time.ParseInLocation("20060102", l.value, time.Local)
	}

	return time.ParseInLocation("20060102T150405Z", l.value, time.Local)
}

func Parse(r io.Reader) (Calendar, error) {
	ret := Calendar{}

	var event parsedEvent
	var inEvent bool
	tr := textproto.NewReader(bufio.NewReader(r))
	for {
		line, err := tr.ReadContinuedLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if line == "BEGIN:VEVENT" {
			inEvent = true
			event = make(map[string][]parsedLine)
			continue
		}
		if line == "END:VEVENT" {
			inEvent = false
			events, err := event.Event()
			if err != nil {
				return Calendar{}, err
			}
			ret.Events = append(ret.Events, events...)
			continue
		}
		if !inEvent {
			continue
		}
		l, err := parseLine(line)
		if err != nil {
			return Calendar{}, err
		}

		event[l.name] = append(event[l.name], l)
		// switch l.name {
		// case "SUMMARY":
		// 	event.Summary = l.value
		// case "DTSTART":
		// 	t, err := parseDate(l)
		// 	if err != nil {
		// 		return Calendar{}, err
		// 	}
		// 	event.Start = t
		// case "DTEND":
		// 	t, err := parseDate(l)
		// 	if err != nil {
		// 		return Calendar{}, err
		// 	}
		// 	event.End = t
		// }
	}

	return ret, nil
}
