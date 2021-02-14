Serves a directory over http, automatically refreshing when files change.

It takes an optional dir to serve, the default is `.`.

```bash
$ srv ~
Serving /home/frew on [::]:21873
```

You can pass -port if you care to choose the listen port.

It will set up file watchers and trigger page reloads (via SSE,) this
functionality can be disabled with -no-autoreload.

```bash
$ srv -port 8080 -no-autoreload ~
Serving /home/frew on [::]:8080
```
