// Package personality returns pseudorandom responses for the lulz.
// If you don't call rand.Seed() with something sensible it won't even be
// pseudorandom.
package personality

import (
	"math/rand"
)

var acks = []string{
	"station",
	"got em.",
	"👍",
	"ack",
	"10-4",
	"wilco",
	"aye aye cap'm'",
}

// Ack returns a string meaning "yes"
func Ack() string {
	const offset = 100
	res := rand.Intn(offset + len(acks))
	if res > offset {
		return acks[res-offset]
	}

	return "Aight"
}

var errs = []string{
	"COMPTER FAIL",
	"Shucks Howdy! 🤠",
	"FAIL🐳",
}

// Err returns a string meaning something went wrong
func Err() string {
	return errs[rand.Intn(len(errs))]
}
