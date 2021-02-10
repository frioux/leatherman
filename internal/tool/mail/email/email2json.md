ToJSON produces a JSON representation of an email from a list of globs.  Only
headers are currently supported, patches welcome to support bodies.

```bash
$ email2json '/home/frew/var/mail/mitsi/cur/*' | head -1 | jq .
{
  "Header": {
    "Content-Type": "multipart/alternative; boundary=00163642688b8ef3070464661533",
    "Date": "Thu, 5 Mar 2009 15:45:17 -0600",
    "Delivered-To": "xxx",
    "From": "fREW Schmidt <xxx>",
    "Message-Id": "<fb3648c60903051345o728960f5l6cfb9e1f324bbf50@mail.gmail.com>",
    "Mime-Version": "1.0",
    "Received": "by 10.103.115.8 with HTTP; Thu, 5 Mar 2009 13:45:17 -0800 (PST)",
    "Subject": "STATION",
    "To": "xyzzy@googlegroups.com"
  }
}
```
