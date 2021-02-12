Reads CSV on stdin and writes JSON on stdout

```bash
$ echo "foo,bar\n1,2\n3,4" | csv2json
{"bar":"2","foo":"1"}
{"bar":"4","foo":"3"}
```

First line of input is the header, and thus the keys of the JSON.
