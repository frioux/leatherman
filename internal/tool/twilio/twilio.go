package twilio

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	coretwilio "github.com/frioux/leatherman/internal/twilio"
)

/*
Twilio allows interacting with a service that recieves callbacks from twilio
for testing.

It takes four arguments:

 * `-endpoint`: the url to hit (`http://localhost:8080/twilio`, for example)
 * `-auth`: the auth token to use
 * `-message`: the message to send
 * `-from`: the phone number the message is from (`+15555555555`, for example)

Run `twilio -help` to see the defaults.

```bash
$ twilio -message "the building is on fire!"
```
*/
func Twilio(args []string, _ io.Reader) error {
	var endpoint, auth, message, from string

	fs := flag.NewFlagSet("twilio", flag.ContinueOnError)

	fs.StringVar(&endpoint, "endpoint", "http://localhost:8080/twilio", "endpoint to post to")
	fs.StringVar(&auth, "auth", "xyzzy", "auth token")
	fs.StringVar(&message, "message", "this is a test", "message to submit to amygdala")
	fs.StringVar(&from, "from", "+15555555555", "cell message is from")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	vals := url.Values{
		"Body": {message},
		"From": {from},
	}
	r := strings.NewReader(vals.Encode())
	req, err := http.NewRequest("POST", endpoint, r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.ParseForm()
	sig := coretwilio.GenerateMAC([]byte(auth), []byte(endpoint), req)
	encodedSignature := base64.StdEncoding.EncodeToString(sig)

	req.Header.Set("X-Twilio-Signature", encodedSignature)
	req.Header.Set("User-Agent", "twilioemu")

	r.Reset(vals.Encode())

	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	for k, v := range resp.Header {
		fmt.Printf("%s: %s\n", k, v[0])
	}
	fmt.Println()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return err
	}

	return nil
}
