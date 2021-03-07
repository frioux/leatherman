`media-remote` control media players on Linux.

The intention is that this tool is bound to media keys, so a play button on
your keyboard runs `media-remote -play`, etc.  The `-select-player` feature
is the main reason to use this tool.

The following flags do the obvious thing:

 * `-play`
 * `-pause`
 * `-play-pause`
 * `-next`
 * `-prev`

You can use the following to call methods not defined in this tool:

 * `-raw` `<method-to-call>`

Finally, the `-select-player` flag will start up a web UI that will allow the
user to select a different player than the default to control with the media
keys.  Try it out by playing a youtube video in both firefox and chrome.
`-select-player` should show two links, allowing you to select firefox or
chrome to bind media keys to.
