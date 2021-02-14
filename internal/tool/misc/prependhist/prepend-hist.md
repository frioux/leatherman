Combines a history file and stdin.

Prints out deduplicated lines from the history file in reverse order and
then prints out the lines from STDIN, filtering out what's already been printed.

```bash
$ echo "1\n2\n3" > eg_history
$ echo "1\n5 | prepend-hist eg_history
3
2
1
5
```
