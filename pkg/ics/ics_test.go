package ics

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	type row struct {
		name string
		in   string
		out  parsedLine
		err  error
	}

	tests := []row{
		{
			name: "simple",
			in:   "FOO:1",
			out:  parsedLine{name: "FOO", value: "1"},
		},
		{
			name: "one param",
			in:   "FOO;P=1:1",
			out:  parsedLine{name: "FOO", value: "1", params: []string{"P", "1"}},
		},
		{
			name: "two params",
			in:   "FOO;P=1;A=2:1",
			out:  parsedLine{name: "FOO", value: "1", params: []string{"P", "1", "A", "2"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l, err := parseLine(test.in)
			if test.err != nil {
				assert.Equal(t, test.err, err)
				return
			}
			if err != nil {
				t.Fatalf("parseLine: %s", err)
			}
			assert.Equal(t, test.out, l)
		})
	}
}

func TestParseDate(t *testing.T) {
	type row struct {
		name, in string
		err      error
		t        time.Time
	}
	tests := []row{
		{
			name: "plain",
			in:   "DTSTART:20190617T100000Z",
			t:    time.Date(2019, 6, 17, 10, 0, 0, 0, time.Local),
		},
		{
			name: "datezone",
			in:   "DTSTART;TZID=America/Los_Angeles:20190617T100000",
			t:    time.Date(2019, 6, 17, 10, 0, 0, 0, time.Local),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l, err := parseLine(test.in)
			if err != nil {
				t.Fatalf("parseLine: %s", err)
			}
			d, err := parseDate(l)
			if err != nil {
				t.Fatalf("parseLine: %s", err)
			}
			assert.Equal(t, test.t, d)
		})
	}
}

func TestParse(t *testing.T) {
	r, err := os.Open("./testdata/basic.ics")
	if err != nil {
		t.Fatalf("Couldn't os.Open: %s", err)
	}
	c, err := Parse(r)
	if err != nil {
		t.Fatalf("Couldn't Parse: %s", err)
	}
	assert.Equal(t, Calendar{}, c)
}
