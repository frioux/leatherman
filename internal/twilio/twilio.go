package twilio

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"sort"

	"github.com/pkg/errors"
)

func CheckMAC(key, url []byte, r *http.Request) (bool, error) {
	var buf bytes.Buffer

	buf.Write(url)

	keys := make(sort.StringSlice, 0, len(r.Form))
	for k := range r.Form {
		keys = append(keys, k)
	}

	keys.Sort()

	for _, k := range keys {
		buf.WriteString(k)
		for _, v := range r.Form[k] {
			buf.WriteString(v)
		}
	}

	mac := hmac.New(sha1.New, key)
	mac.Write(buf.Bytes())
	expectedMAC := mac.Sum(nil)
	messageMAC, err := base64.StdEncoding.DecodeString(r.Header.Get("X-Twilio-Signature"))
	if err != nil {
		return false, errors.Wrap(err, "base64.Decode")
	}
	return hmac.Equal(messageMAC, expectedMAC), nil
}
