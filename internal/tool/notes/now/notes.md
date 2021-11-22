Serves my notes from disk or Dropbox.

When notes change (either on disk or in dropbox) the change is noticed quickly
and should appear within seconds of the change being made.  If a note is open
in a browser window, the window will be reloaded when the change is noticed, so
old information should almost never be visible.

Notes should be stored in a single directory and named `$somenote.md`.

## Running

When notes are loaded from disk, run notes as follows:

```bash
LM_NOTES_PATH=/some/directory notes -listen :8080
```

`-listen` is optional; by default notes will pick a random port to listen on.

When notes are loaded from dropbox, you must also set `LM_DROPBOX_TOKEN` to a
Dropbox token, and the PATH should not include a `/` prefix.

### systemd.service

The following is a working systemd unit you can use to set up a service:

```
[Unit]
Description=Notes Server, port 5000

[Service]
Environment=LM_NOTES_PATH=notes/content/posts LM_DROPBOX_TOKEN=yyy 'LM_GH_TOKEN=Bearer xxx'
ExecStart=/home/pi/leatherman notes -listen 0:5000
Restart=always
StartLimitBurst=0

[Install]
WantedBy=default.target
```

You can put it at either `/etc/systemd/system/notes.service` or
`~/.config/systemd/user/notes.service`.

Then do one of these:

```bash
$ systemctl --user daemon-reload
$ systemctl --user enable notes
$ systemctl --user start notes
```

```bash
$ systemctl daemon-reload
$ systemctl enable notes
$ systemctl start notes
```

## Note Format

This is what a note looks like, called something like `funny.md`:

        {
          "tags": [ "public", "reference" ],
          "title": "Funny",
        }

         * [What's the quickest way to make someone feel uncomfortable using only one sentence?](https://www.reddit.com/r/AskReddit/comments/1ivjae/whats_the_quickest_way_to_make_someone_feel/)
         * [Parents of Reddit, what is the creepiest thing your young child has ever said to you?](https://www.reddit.com/r/AskReddit/comments/1d2v7i/parents_of_reddit_what_is_the_creepiest_thing/)
         * [Two Medieval Monks Invent Bestiaries - The Toast](http://the-toast.net/2015/04/01/two-medieval-monks-invent-bestiaries/)
         * [Monologue: I’m Comic Sans, Asshole - McSweeney’s Internet Tendency](https://www.mcsweeneys.net/articles/im-comic-sans-asshole)

        ## Videos

         * [Woodford Reserve Mint Julep - YouTube](https://www.youtube.com/watch?feature=player_embedded&v=Nk57WmewiRA)
         * [The Baby Bullet - YouTube](https://www.youtube.com/watch?v=n5Gn8jt55LQ)
         * [Relax With Sonic - YouTube](https://www.youtube.com/watch?v=Y8rt69ztDjA)
         * [Aussie Shames A Pair Of Kangaroos For Brawling In His Yard - Digg](http://digg.com/video/guy-shames-roos)

        ## GIFs

         * [kickflip](https://i.imgur.com/NmKZKSB.gif)
         * [brie decay](https://gfycat.com/requiredallazurevasesponge)
         * [Friday](https://i.imgur.com/nSD6W.gif)
         * [Peekaboo!](https://i.imgur.com/LsPKOwj.jpg)
         * [Red Pandas](http://jennipoos.tumblr.com/post/39955516867/storiesinsilentnights-i-still-fucking-love)
         * [All about timing](https://i.imgur.com/CBIaS78.gifv)
         * [Burger Alignment](https://i.imgur.com/AtKtw86.png)

The top is a variant of JSON that allows trailing commas.  That is metadata
about the note.

## Database

The metadata discussed above gets inserted into a database that can be
referenced from a note.  The schema is:

```sql
CREATE TABLE articles (
        title,
        url,
        filename,
        reviewed_on NULLABLE,
        review_by NULLABLE,
        body
);
CREATE TABLE article_tag ( id, tag );
CREATE VIEW _ ( id, title, url, filename, body, reviewed_on, review_by, tag) AS
        SELECT a.rowid, title, url, filename, body, reviewed_on, review_by, tag
        FROM articles a
        JOIN article_tag at ON a.rowid = at.id;
```

The SQL can be used within a note to render links (or whatever) from other
notes.  Here's an example from my `recipes.md`:

```
## Dinner

{{range (q "SELECT title, url FROM _ WHERE tag = 'dinner'") }}
 * [{{.title}}]({{.url}})
{{- end}}
 * [30 Minute Indian Pumpkin Butter Chickpeas. - Half Baked Harvest](https://www.halfbakedharvest.com/30-minute-indian-pumpkin-butter-chickpeas/)
 * [Alton brown chili](https://www.foodnetwork.com/recipes/alton-brown/pressure-cooker-chili-recipe-1942714) - recommended by Melinda Baustian
```

That finds all notes with a tag of `dinner` and puts links to them in the list.

## Lua

Notes can include lua scripts that are executed within the server.  The
intention is that these scripts can be used to hook into the rendering of a
note, as well as being able to modify a note directly.

To embed a lua script in a page use a fenced code block with type `mdlua`, so
something like:

     ```mdlua
     function x(rw, r)
        rw:write("boo!")
     end
     ```

You'd call the above by going to `/$note?lua=x`.

The currently implemented API is a work in progress and will definitely change.
Here's a vague rundown of how it looks for now:

### fs

 * `fs:writefile(path, contents)`

### goquery

 * `goquery.newdocumentfromstring(content)` -> `goqueryselection`
 * `goqueryselection:attr(name)` -> `stringvalue`
 * `goqueryselection:each(function(i, goqueryselection) end)` -> `goqueryselection`
 * `goqueryselection:find(cssselector)` -> `goqueryselection`
 * `goqueryselection:text()` -> `string`

### http

 * `http.get(url)` -> `contentstring`
 * `http.multiget({url1, url2})` -> `{url1="content1", url2="content2"}`
 * `responsewriter:write(string)`
 * `responsewriter:writeheader(statuscode)`
 * `request:url()` -> url
 * `url:query()` -> values
 * `values:get(param)` -> `val`
 * `url.parse(str)` -> url

### regexp

 * `regexp.compile(string)` -> regexp
 * `regexp:findallstringsubmatch(bigstring)` -> `{{matched, capture1, cap2}, {...}}`
 * `regexp:replaceallstringfunc(bigstring, function(found) return "replacement" end)` -> `newbigstring`

### notes specific

 * `notes.readarticlefromfs(fs, path)` -> `article`
 * `article:rawcontents()` string of the full note
