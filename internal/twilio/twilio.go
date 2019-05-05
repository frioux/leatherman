package twilio

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"golang.org/x/xerrors"
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
		return false, xerrors.Errorf("base64.Decode: %w", err)
	}
	return hmac.Equal(messageMAC, expectedMAC), nil
}

type Media struct {
	ContentType, URL string
}

func ExtractMedia(f url.Values) ([]Media, error) {
	numMedia := f.Get("NumMedia")
	if numMedia == "" {
		return nil, nil
	}

	n, err := strconv.Atoi(numMedia)
	if err != nil {
		return nil, xerrors.Errorf("Couldn't parse NumMedia: %w", err)
	}

	ret := make([]Media, n)

	for i := 0; i < n; i++ {
		ret[i].URL = f.Get(fmt.Sprintf("MediaUrl%d", i))
		ret[i].ContentType = f.Get(fmt.Sprintf("MediaContentType%d", i))
	}

	return ret, nil
}
