package main

import (
	"bufio"
	"fmt"
	"log"
	"net/mail"
	"os"
)

// AddrspecToTabs reads emails (`"Foo Bar" <foo@bar.com>`) and produces addrbook
// format (`Foo Bar	foo@bar.com`)
func AddrspecToTabs(args []string) {
	if len(args) > 1 && args[1] == "-h" {
		fmt.Println("reads emails (`\"Foo Bar\" <foo@bar.com>`) and produces addrbook",
			"format (`Foo Bar\tfoo@bar.com`)")
		return
	}
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
