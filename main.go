package main

import (
	"bufio"
	"fmt"
	"log"
	"net/mail"
	"os"
	"path"
	"strconv"
	"time"
)

var dispatch map[string]func()

func main() {
	which := path.Base(os.Args[0])

	dispatch = map[string]func(){
		"addrspec-to-tabs": addrspecToTabs,
		"clocks":           clocks,
	}

	if which == "main" {
		if len(os.Args) > 1 {
			which = os.Args[1]
		}
	}

	fn, ok := dispatch[which]
	if !ok {
		help()
	}
	fn()
}

func help() {
	str := "Tools:\n"
	for k := range dispatch {
		str += " * " + k + "\n"
	}
	fmt.Println(str)
	os.Exit(1)
}

func addrspecToTabs() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		list := scanner.Text()
		emails, err := mail.ParseAddressList(list)
		if err != nil {
			log.Print(err, list)
		}

		for _, v := range emails {
			fmt.Println(v.Address + "\t" + v.Name + "\t")
		}
	}
}

var now time.Time = time.Now().In(time.Local)

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

func t(l string) string {
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

func clocks() {
	fmt.Println("here : " + t("Local"))
	fmt.Println("L.A. : " + t("America/Los_Angeles"))
	fmt.Println("MS/TX: " + t("America/Chicago"))
	fmt.Println("rjbs : " + t("America/New_York"))
	fmt.Println("riba : " + t("Europe/Berlin"))
	fmt.Println("seo  : " + t("Asia/Jerusalem"))
	fmt.Println("UTC  : " + t("UTC"))
}
