package twilio

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckMAC(t *testing.T) {
	// Example from https://www.twilio.com/docs/usage/security#notes
	ok, err := CheckMAC([]byte("12345"), []byte("https://mycompany.com/myapp.php?foo=1&bar=2"), &http.Request{
		PostForm: url.Values(map[string][]string{
			"CallSid": {"CA1234567890ABCDE"},
			"Caller":  {"+12349013030"},
			"Digits":  {"1234"},
			"From":    {"+12349013030"},
			"To":      {"+18005551212"},
		}),
		Header: http.Header(map[string][]string{
			"X-Twilio-Signature": {"0/KCTR6DLpKmkAf8muzZqo1nDgQ="},
			"Content-Type":       {"application/x-www-form-urlencoded"},
		}),
	})

	assert.True(t, ok)
	assert.NoError(t, err)
}
