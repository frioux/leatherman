package email

import (
	"bufio"
	"fmt"
	"io"
	"net/mail"
	"regexp"
	"time"
)

var dateRe = regexp.MustCompile(`^Date:\s+(.*)\s*$`)

func Render(args []string, stdin io.Reader) error {
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		line := scanner.Text()

		match := dateRe.FindSubmatch([]byte(line))
		if match == nil {
			fmt.Println(line)
		} else {
			date, err := mail.ParseDate(string(match[1]))

			if err == nil && date.Location() != time.Local {
				fmt.Println("Local-Date: " + date.In(time.Local).Format("Mon, 02 Jan 2006 15:04:05"))
			}
			fmt.Println(line)
		}
	}

	return nil
}
