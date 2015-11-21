# Light Daemon

This is a very rough set of applications that are meant to interact with the
[blink(1)](http://blink1.thingm.com/).  They will likely change wildly until
they support all that they could.

## poller.pl

The main script is `poller.pl`; it reads from stdin and sets the light's current
color.  It supports a single blink1 with a single rgb light currently.  The
format is:

`[rgb][+-]?0-255`

So `r+5` will increase red by 5, `g-6` will decrease green by 6, and `b212` will
set blue to 212.

## pal.pl

This will simply poll slack for the name you pass it, and when the user is
online the green LED will be set to maximum brightness.

## sound.pl

This will set the light to red if sound is going to pulseaudio.

# Wiring It All Together

To wire everything together you can simply do something like:

```
mkfifo color-pipe

>color-pipe ./poller.pl &
./pal.pl   > color-pipe &
./sound.pl > color-pipe &

```
