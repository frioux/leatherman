package twilio

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckMAC(t *testing.T) {
	// Example from https://www.twilio.com/docs/usage/security#notes
	ok, err := CheckMAC([]byte("12345"), &http.Request{
		URL: &url.URL{
			Scheme:   "https",
			Host:     "mycompany.com",
			Path:     "/myapp.php",
			RawQuery: "foo=1&bar=2",
		},
		Header: http.Header(map[string][]string{
			"X-Twilio-Signature": {
				"0/KCTR6DLpKmkAf8muzZqo1nDgQ=",
			},
			"Content-Type": {
				"application/x-www-form-urlencoded",
			},
		}),
		PostForm: url.Values(map[string][]string{
			"CallSid": {
				"CA1234567890ABCDE",
			},
			"Caller": {
				"+12349013030",
			},
			"Digits": {
				"1234",
			},
			"From": {
				"+12349013030",
			},
			"To": {
				"+18005551212",
			},
		}),
	})

	assert.True(t, ok)
	assert.NoError(t, err)
}
