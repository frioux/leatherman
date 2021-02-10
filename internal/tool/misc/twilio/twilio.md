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
