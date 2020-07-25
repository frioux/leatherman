package notes

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"
)

// parseNow reads markdown and returns html.  The main difference from normal
// markdown is that a section titled ## 2020-02-02 ## will get special
// rendering treatment.
//
// A json header is discarded for now.
func parseNow(r io.Reader, when time.Time) ([]byte, error) {
	desiredHeader := "## " + when.Format("2006-01-02") + " ##"
	ret := &strings.Builder{}

	a, err := readArticle(r)
	if err != nil {
		return nil, fmt.Errorf("readArticle: %w", err)
	}

	var inToday, wroteAddItem bool
	s := bufio.NewScanner(bytes.NewReader(a.Body))
	for s.Scan() {
		line := s.Text()

		switch {
		case !inToday && line == desiredHeader:
			inToday = true
		case inToday && strings.HasPrefix(line, "## "):
			ret.WriteString(`<form action="/add-item" method="POST"><input type="input" name="item"><button>Add Item</button></form>`)
			ret.WriteString("\n\n")
			wroteAddItem = true
			inToday = false
		case inToday && strings.HasPrefix(line, " * "):
			md := md5.Sum([]byte(line))
			linkable := hex.EncodeToString(md[:])
			line += ` <form action="/toggle" method="POST"><input type="hidden" name="chunk" value="` + linkable + `"><button>Toggle</button></form>`
		}

		ret.WriteString(line)
		ret.WriteRune('\n')
	}

	if !wroteAddItem {
		ret.WriteString(`<form action="/add-item" method="POST"><input type="input" name="item"><button>Add Item</button></form>`)
		ret.WriteString("\n\n")
	}

	return []byte(ret.String()), nil
}

// toggleNow will mark a list item done (surround with ~~'s) if it's in the
// section for when and it's md5sum matches sum.  If the item has already been
// done, this function will mark it undone.
func toggleNow(r io.Reader, when time.Time, sum string) ([]byte, error) {
	desiredHeader := "## " + when.Format("2006-01-02") + " ##"
	ret := &strings.Builder{}

	var inToday bool
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		switch {
		case !inToday && line == desiredHeader:
			inToday = true
		case inToday && strings.HasPrefix(line, "## "):
			inToday = false
		case inToday && strings.HasPrefix(line, " * "):
			md := md5.Sum([]byte(line))
			linkable := hex.EncodeToString(md[:])
			if sum == linkable {
				if strings.HasPrefix(line, " * ~~") && strings.HasSuffix(line, "~~") { // already done, undo
					line = " * " + line[5:len(line)-2]
				} else { // not done, mark done
					line = " * ~~" + line[3:] + "~~"
				}
			}
		}

		ret.WriteString(line)
		ret.WriteRune('\n')
	}

	return []byte(ret.String()), nil
}
