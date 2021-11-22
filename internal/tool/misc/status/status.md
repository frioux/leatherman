Serves information about host machine.

Status runs a little web server that surfaces status information related to how
I'm using the machine.  For example, it can say which window is active, what
firefox tabs are loaded, if the screen is locked, etc.  The main benefit of the
tool is that it caches the values returned.

In the background, it interact swith the [blink(1)](http://blink1.thingm.com/).
It turns the light green when I'm in a meeting and red when audio is playing.

### systemd.service

The following is a working systemd unit you can use to set up a service:

```
[Unit]
Description=Status Server, port 8081

[Service]
Environment='LM_GH_TOKEN=Bearer xxx'
ExecStart=/home/pi/leatherman status
Restart=always
StartLimitBurst=0

[Install]
WantedBy=default.target
```

You can put it at either `/etc/systemd/system/status.service` or
`~/.config/systemd/user/status.service`.

Then do one of these:

```bash
$ systemctl --user daemon-reload
$ systemctl --user enable status
$ systemctl --user start status
```

```bash
$ systemctl daemon-reload
$ systemctl enable status
$ systemctl start status
```
