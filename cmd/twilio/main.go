package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/frioux/amygdala/internal/twilio"
)

var endpoint, auth, message, from string

func init() {
	flag.StringVar(&endpoint, "endpoint", "http://localhost:8080/twilio", "endpoint to post to")
	flag.StringVar(&auth, "auth", "xyzzy", "auth token")
	flag.StringVar(&message, "message", "this is a test", "message to submit to amygdala")
	flag.StringVar(&from, "from", "+15555555555", "cell message is from")
}

func main() {
	flag.Parse()

	vals := url.Values{
		"Body": {message},
		"From": {from},
	}
	r := strings.NewReader(vals.Encode())
	req, err := http.NewRequest("POST", endpoint, r)
	if err != nil {
		fmt.Printf("Failed building request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.ParseForm()
	sig := twilio.GenerateMAC([]byte(auth), []byte(endpoint), req)
	encodedSignature := base64.StdEncoding.EncodeToString(sig)

	req.Header.Set("X-Twilio-Signature", encodedSignature)
	r.Reset(vals.Encode())

	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		fmt.Printf("Failed submitting to amygdala: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		fmt.Printf("Failed load bodu: %s\n", err)
		os.Exit(1)
	}
}
