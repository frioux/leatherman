Run debounces input from stdin to stdout

The default lockout time is one second, you can override that with the
`--lockoutTime` argument.  By default the trailing edge triggers output, so
output is emitted after there is no input for the lockout time.  You can change
this behavior by passing the `--leadingEdge` flag.
