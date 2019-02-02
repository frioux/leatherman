package personality

import (
	"math/rand"
)

var responses = []string{
	"station",
	"got em.",
	"👍",
	"ack",
	"10-4",
	"wilco",
	"aye aye cap'm'",
}

// Ack returns a random string meaning "yes"; you should rand.Seed() before
// calling this.
func Ack() string {
	const offset = 100
	res := rand.Intn(offset + len(responses))
	if res > offset {
		return responses[res-offset]
	}

	return "Aight"
}
