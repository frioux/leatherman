Steamsrv renders steam screenshots and the steam log (what games you played and when) over http.

## systemd.service

The following is a working systemd unit you can use to set up a service:

```
[Unit]
Description=Steam Server, port 8080

[Service]
Environment='LM_GH_TOKEN=Bearer xxx'
ExecStart=/home/pi/leatherman steamsrv -screenshot-prefix %h/.local/share/Steam/userdata/1234324321
Restart=always
StartLimitBurst=0

[Install]
WantedBy=default.target
```

You can put it at either `/etc/systemd/system/steamsrv.service` or
`~/.config/systemd/user/steamsrv.service`.

Then do one of these:

```bash
$ systemctl --user daemon-reload
$ systemctl --user enable steamsrv
$ systemctl --user start steamsrv
```

```bash
$ systemctl daemon-reload
$ systemctl enable steamsrv
$ systemctl start steamsrv
```
