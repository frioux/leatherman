Run watches one or more directories (before the `--`) and runs a script when
events in those directories occur.

```bash
minotaur -include-args -include internal -ignore yaml \
   ~/code/leatherman -- \
   go test ~/code/leatherman/...
```

If the `-include-args` flag is set, the script receives the events as
arguments, so you can exit early if only irrelevant files changed.

The arguments are of the form `$event\t$filename`; for example `CREATE	x.pl`.
As far as I know the valid events are;

 * `CHMOD`
 * `CREATE`
 * `REMOVE`
 * `RENAME`
 * `WRITE`

The events are deduplicated and also debounced, so your script will never fire
more often than once a second.  If events are happening every half second the
debouncing will cause the script to never run.

The underlying library supports emitting multiple events in a single line (ie
`CREATE|CHMOD`) though I've not seen that in Linux.

`minotaur` reëmits all output (both stderr and stdout) of the passed script to
standard out, so you could make a script like this to experiment with the
events with timestamps:

```bash
#!/bin/sh

for x in "$@"; do
   echo "$x"
done | ts
```

You can do all kinds of interesting things in the script, for example you could
verify that the events deserve a restart, then restart a service, then block till
the service can serve traffic, then restart some other related service.

The `-include` and `-ignore` arguments are optional; by default `-include` is
empty, so matches everything, and `-ignore` matches `.git`.  You can also pass
`-verbose` to include output about minotaur itself, like which directories it's
watching.

The flag `-no-run-at-start` will not the the script until there are any events.

The flag `-report` will decorate output with a text wrapper to clarify when the
script is run.
