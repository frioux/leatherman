## Debugging

In an effort to make debugging simpler, I've created three ways to see what
`leatherman` is doing:

### Tracing

`LMTRACE=$somefile` will write an execution trace to `$somefile`; look at that with `go tool trace $somefile`

Since so many of the tools are short lived my assumption is that the execution
trace will be the most useful.

### Profiling

`LMPROF=$somefile` will write a cpu profile to `$somefile`; look at that with `go tool pprof -http localhost:10123 $somefile`

If you have a long running tool, the pprof http endpoint is exposed on
`localhost:6060/debug/pprof` but picks a random port if that port is in use; the
port can be overridden by setting `LMHTTPPROF=$someport`.

## Auxiliary Tools

Some tools are annoying to have in the main leatherman tool.  Maybe they pull
in deps that are huge or need cgo, but in any case I try to keep them separate.
For now, these tools are under `leatherman/cmd` and must be built and run
separately.  At some point I may come up with a policy around building or naming these,
but for now they are simply extra tools.
