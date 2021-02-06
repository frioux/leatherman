package email

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/mail"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type addr struct {
	addr, name string
	score      float64
}

func (a addr) String() string { return fmt.Sprintf("%s\t%s", a.addr, a.name) }

// allAddrs returns all email addresses an email was sent to (To, Cc, and Bcc)
func allAddrs(email *mail.Message) []*mail.Address {
	addrs := []*mail.Address{}
	headers := []string{"To", "Cc", "Bcc"}
	for _, x := range headers {
		if email.Header.Get(x) != "" {
			// Emails tend to be messy, ignore error
			iAddrs, _ := email.Header.AddressList(x)
			addrs = append(addrs, iAddrs...)
		}
	}

	return addrs
}

// buildFrecencyMapFromGlob calls buildFrecencyMap passing emails found from the
// passed glob
func (score frecencyMap) scoreFromGlob(glob string) {
	matches, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal("couldn't get glob", err)
	}

	for _, path := range matches {
		file, err := os.Open(path)
		if err != nil {
			log.Println("Coudln't open email", path, err)
		}
		email, err := mail.ReadMessage(file)
		file.Close()
		if err != nil {
			log.Println("Coudln't parse email", path, err)
			continue
		}
		score.scoreEmail(email, time.Now())
	}
}

type frecencyMap map[string]*addr

// math.Log(2) / 30
const lambda = 0.02310490601866484364

func newFrecencyMap() frecencyMap { return map[string]*addr{} }

// buildFrecencyMap returns a map of addresses, scored based on how recently
// they were mailed to.  See
// https://wiki.mozilla.org/User:Jesse/NewFrecency#Proposed_new_definition
func (score frecencyMap) scoreEmail(email *mail.Message, now time.Time) {
	for _, addr := range allAddrs(email) {
		if val, ok := score[strings.ToLower(addr.Address)]; ok {
			time, err := email.Header.Date()
			if err != nil {
				log.Println("Couldn't read date", err)
				continue
			}
			age := now.Sub(time).Hours() / 24

			val.score += math.Exp(-lambda * age)
		}
	}
}

// buildAddrMap returns a map of address and content, based on os.Stdin
func buildAddrMap(reader io.Reader) frecencyMap {
	scanner := bufio.NewScanner(reader)

	ret := newFrecencyMap()
	for scanner.Scan() {
		z := strings.SplitN(scanner.Text(), "\t", 2)
		if len(z) < 2 {
			continue
		}
		if _, ok := ret[z[0]]; ok {
			continue
		}
		ret[z[0]] = &addr{
			addr: z[0],
			name: z[1],
		}
	}

	return ret
}

// sortAddrMap sorts the addrs arg based on the values in the score arg;
// leftover values are printed in alphabetical order.
func sortAddrMap(score frecencyMap) []addr {
	addrs := make([]addr, len(score))

	i := 0
	for _, addr := range score {
		addrs[i] = *addr
		i++
	}

	// sort keys based on score
	sort.Slice(
		addrs,
		func(i, j int) bool { return addrs[i].score > addrs[j].score },
	)

	return addrs
}

/*
Addrs sorts the addresses passed on stdin (in the mutt addrbook format, see
`addrspec-to-tabs`) and sorts them based on how recently they were used, from
the glob passed on the arguments.  The tool exists so that you can create an
address list either with an export tool (like `goobook`), a subset of your sent
addresses, or whatever else, and then you can sort it based on your sent mail
folder.

``` bash
$ <someaddrs.txt addrs "$HOME/mail/gmail.sent/cur/*" >sortedaddrs.txt
```
*/
func Addrs(args []string, stdin io.Reader) error {
	if len(args) != 2 {
		return errors.New("please pass a glob")
	}

	addrMap := buildAddrMap(stdin)
	addrMap.scoreFromGlob(args[1])

	addrs := sortAddrMap(addrMap)

	// first line is blank
	fmt.Println()
	for _, line := range addrs {
		fmt.Println(line)
	}

	return nil
}
