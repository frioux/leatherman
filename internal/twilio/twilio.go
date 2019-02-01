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

func GenerateMAC(key, url []byte, r *http.Request) []byte {
	var buf bytes.Buffer

	buf.Write(url)

	keys := make(sort.StringSlice, 0, len(r.PostForm))
	for k := range r.PostForm {
		keys = append(keys, k)
	}

	keys.Sort()

	for _, k := range keys {
		buf.WriteString(k)
		for _, v := range r.PostForm[k] {
			buf.WriteString(v)
		}
	}

	mac := hmac.New(sha1.New, key)
	mac.Write(buf.Bytes())
	return mac.Sum(nil)
}

func CheckMAC(key, url []byte, r *http.Request) (bool, error) {
	expectedMAC := GenerateMAC(key, url, r)
	messageMAC, err := base64.StdEncoding.DecodeString(r.Header.Get("X-Twilio-Signature"))
	if err != nil {
		return false, errors.Wrap(err, "base64.Decode")
	}
	return hmac.Equal(messageMAC, expectedMAC), nil
}
