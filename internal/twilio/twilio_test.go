package twilio

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckMAC(t *testing.T) {
	// Example from https://www.twilio.com/docs/usage/security#notes
	ok, err := CheckMAC([]byte("12345"), []byte("https://mycompany.com/myapp.php?foo=1&bar=2"),
		&http.Request{
			Form: url.Values(map[string][]string{
				"CallSid": []string{"CA1234567890ABCDE"},
				"Caller":  []string{"+12349013030"},
				"Digits":  []string{"1234"},
				"From":    []string{"+12349013030"},
				"To":      []string{"+18005551212"},
			}),
			Header: http.Header(map[string][]string{
				"X-Twilio-Signature": []string{"0/KCTR6DLpKmkAf8muzZqo1nDgQ="},
			}),
		})

	assert.True(t, ok)
	assert.NoError(t, err)
}
