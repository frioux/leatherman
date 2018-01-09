package main

import (
	"bufio"
	"fmt"
	"log"
	"net/mail"
	"os"
)

func AddrspecToTabs() {
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
