Addrs sorts the addresses passed on stdin (in the mutt addrbook format, see
`addrspec-to-tabs`) and sorts them based on how recently they were used, from
the glob passed on the arguments.  The tool exists so that you can create an
address list either with an export tool (like `goobook`), a subset of your sent
addresses, or whatever else, and then you can sort it based on your sent mail
folder.

``` bash
$ <someaddrs.txt addrs "$HOME/mail/gmail.sent/cur/*" >sortedaddrs.txt
```
