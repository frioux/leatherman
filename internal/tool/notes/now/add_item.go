package now

import (
	"bufio"
	"io"
	"strings"
	"time"
)

// addItem will add item to the end of today's list, and create today's list if
// it doesn't already exist.
func addItem(r io.Reader, when time.Time, item string) ([]byte, error) {
	desiredHeader := "## " + when.Format("2006-01-02") + " ##"
	ret := &strings.Builder{}

	var inToday, inList, addedItem bool
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		switch {
		case !addedItem && !inToday && strings.HasPrefix(line, "## ") && strings.HasSuffix(line, " ##") && line < desiredHeader: // We found a previous day, stop searching and make a new day:
			ret.WriteString(desiredHeader + "\n\n * " + item + "\n\n")
			addedItem = true
		case !inToday && line == desiredHeader:
			inToday = true
		case inToday && !inList && strings.HasPrefix(line, " * "):
			inList = true
		case inToday && strings.HasPrefix(line, "## "):
			inToday = false
		case inToday && strings.HasPrefix(line, " * ") && !addedItem:
			foundItem := strings.TrimPrefix(line, " * ")
			if foundItem == item {
				addedItem = true
			}
		case inToday && !addedItem && inList && line == "":
			ret.WriteString(" * " + item + "\n")
			addedItem = true
		}

		ret.WriteString(line)
		ret.WriteRune('\n')
	}

	return []byte(ret.String()), nil
}
