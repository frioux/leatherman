Reads CSV on stdin and writes Markdown on stdout

```bash
$ echo "foo,bar\n1,2\n3,4" | csv2md
foo | bar
 --- | ---
 1 | 2
 3 | 4
```
