package main

import (
	"bufio"
	"fmt"
	"net/mail"
	"os"
	"regexp"
	"time"
)

var dateRe = regexp.MustCompile(`^Date:\s+(.*)\s*$`)

// RenderMail reads email on stdin and reproduces it on stdout, with the Date
// header parsed and duplicated as Local-Date containing the date in the local
// timezone.
func RenderMail(args []string) {
	scanner := bufio.NewScanner(os.Stdin)
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
}
