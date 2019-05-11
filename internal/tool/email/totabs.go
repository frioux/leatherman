package email

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/mail"
)

/*
ToTabs converts email addresses from the standard format (`"Hello Friend"
<foo@bar>`) to the mutt (?) address book format, ie tab separated fields.

Note that this version ignores the comment because, after actually auditing my
addressbook, most comments are incorrectly recognized by all tools. (for
example: `<5555555555@vzw.com> (555) 555-5555` should not have a comment of
`(555)`.)

Command: addrspec-to-tabs
*/
func ToTabs(args []string, stdin io.Reader) error {
	if len(args) > 1 && args[1] == "-h" {
		fmt.Println("reads emails (`\"Foo Bar\" <foo@bar.com>`) and produces addrbook",
			"format (`Foo Bar\tfoo@bar.com`)")
		return nil
	}
	scanner := bufio.NewScanner(stdin)
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

	return nil
}
