package main

import (
	"bufio"
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

// allAddrs returns all email addresses an email was sent to (To, Cc, and Bcc)
func allAddrs(email *mail.Message) []*mail.Address {
	addrs := []*mail.Address{}
	headers := []string{"To", "Cc", "Bcc"}
	for _, x := range headers {
		if email.Header.Get(x) != "" {
			iAddrs, err := email.Header.AddressList(x)
			// Emails tend to be messy, just move on.
			if err != nil {
				continue
			}
			addrs = append(addrs, iAddrs...)
		}
	}

	return addrs
}

// buildFrecencyMapFromGlob calls buildFrecencyMap passing emails found from the
// passed glob
func buildFrecencyMapFromGlob(glob string) map[string]float64 {
	matches, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal("couldn't get glob", err)
	}

	emailChan := make(chan *mail.Message)
	go func() {
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
			emailChan <- email
		}
		close(emailChan)
	}()

	return buildFrecencyMap(emailChan, time.Now())
}

// math.Log(2) / 30
const lambda = 0.02310490601866484364

// buildFrecencyMap returns a map of addresses, scored based on how recently
// they were mailed to.  See
// https://wiki.mozilla.org/User:Jesse/NewFrecency#Proposed_new_definition
func buildFrecencyMap(emails chan *mail.Message, now time.Time) map[string]float64 {
	score := map[string]float64{}

	for email := range emails {
		for _, addr := range allAddrs(email) {
			time, err := email.Header.Date()
			if err != nil {
				log.Println("Couldn't read date", err)
				continue
			}
			age := now.Sub(time).Hours() / 24

			score[strings.ToLower(addr.Address)] += math.Exp(-lambda * age)
		}
	}
	return score
}

// buildAddrMap returns a map of address and content, based on os.Stdin
func buildAddrMap(reader io.Reader) map[string]string {
	scanner := bufio.NewScanner(reader)

	ret := map[string]string{}
	for scanner.Scan() {
		z := strings.SplitN(scanner.Text(), "\t", 2)
		if len(z) < 2 {
			continue
		}
		if _, ok := ret[z[0]]; ok {
			continue
		}
		ret[z[0]] = z[1]
	}

	return ret
}

// sortAddrMap sorts the addrs arg based on the values in the score arg;
// leftover values are printed in alphabetical order.
func sortAddrMap(score map[string]float64, addrs map[string]string) []string {
	// map of addresses that have been scored
	scored := map[string]string{}
	// keys list, for sorting based on score later
	scoredKeys := []string{}
	for key := range score {
		var ok bool
		scored[key], ok = addrs[key]
		if ok {
			delete(addrs, key)
			scoredKeys = append(scoredKeys, key)
		}
	}

	// sort keys based on score
	sort.Slice(
		sort.StringSlice(scoredKeys),
		func(i, j int) bool { return score[scoredKeys[i]] > score[scoredKeys[j]] },
	)

	ret := []string{}

	for _, key := range scoredKeys {
		ret = append(ret, key+"\t"+scored[key])
	}

	// sort remaining addrs based on keys
	addrKeys := []string{}
	for key := range addrs {
		addrKeys = append(addrKeys, key)
	}
	sort.Sort(sort.StringSlice(addrKeys))
	for _, key := range addrKeys {
		ret = append(ret, key+"\t"+addrs[key])
	}
	return ret
}

// Addrs sorts the addresses passed on stdin based on how recently they were
// used, based on the glob passed on the arguments.
func Addrs(args []string) {
	if len(args) != 2 {
		log.Fatal("Please pass a glob")
	}

	addrs := sortAddrMap(
		buildFrecencyMapFromGlob(args[1]), buildAddrMap(os.Stdin))

	// first line is blank
	fmt.Println()
	for _, line := range addrs {
		fmt.Println(line)
	}
}
