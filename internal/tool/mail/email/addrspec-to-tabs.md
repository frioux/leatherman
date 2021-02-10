ToTabs converts email addresses from the standard format (`"Hello Friend"
<foo@bar>`) from stdin to the mutt address book format, ie tab separated fields,
on stdout.

``` bash
$ <someaddrs.txt addrs "$HOME/mail/gmail.sent/cur/*" | addrspec-to-tabs >addrbook.txt
```

This tool ignores the comment because, after actually auditing my addressbook,
most comments are incorrectly recognized by all tools. (for example:
`<5555555555@vzw.com> (555) 555-5555` should not have a comment of `(555)`.)

It exists to be combined with `addrs` and mutt.
