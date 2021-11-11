package twilio_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/frioux/leatherman/internal/twilio"
)

func TestCheckMAC(t *testing.T) {
	// Example from https://www.twilio.com/docs/usage/security#notes
	ok, err := twilio.CheckMAC([]byte("12345"), []byte("https://mycompany.com/myapp.php?foo=1&bar=2"), &http.Request{
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

	if !ok {
		t.Error("ok should be true")
	}

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
