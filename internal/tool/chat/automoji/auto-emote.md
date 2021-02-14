Reacts to discord messages with vaguely related emoji.

The following env vars should be set:

 * LM_DROPBOX_TOKEN should be set to load a responses.json.
 * LM_BOT_LUA_PATH should be set to the location of lua to process emoji data within dropbox.
 * LM_DISCORD_TOKEN should be set for this to actually function.

Here's an example of lua code that works for this:

```lua
function add(w)
        es:addoptional(w:char())
end

function addtoken(tok)
        if tok == "" then
                return
        end

        bn = turtleemoji.findbyname(tok)
        if not bn == nil then
                add(bn)
        end

        for i, te in ipairs(turtleemoji.searchbycategory(tok)) do
                add(te)
        end

        for i, te in ipairs(turtleemoji.searchbykeyword(tok)) do
                add(te)
        end

        if not es:optionallen() == 0 then
                return
        end

        -- this returns so many results for basic words that we only
        -- use it if nothing else found anything for the whole message
        for i, te in ipairs(turtleemoji.search(tok)) do
                add(te)
        end
end

for i, w in ipairs(es:words()) do
        addtoken(w)
end
```

The lua code has a global var called `es` (for emoji set) and an imported
package called `turtleemoji`.  `es` is how you access the current message,
currently added emoji, etc.  Here are the methods on `es`:

#### `es:optional()` // table of string to bool

Returns a copy of the optional emoji.  Modifications of the table will not
affect the final result; other methods should be used for modification.

#### `es:addoptional("💀")`

Adds an emoji to randomly include in the reaction.

#### `es:hasoptional("💀")` // bool

Returns true of the passed emoji is in the list of optional emoji to include
(at random) on the reaction.

#### `es:removeoptional("💀")`

Remove the passed emoji from the optionally included emoji.

#### `es:required()` // table of required emoji

Returns a copy of the required emoji.  Modifications of the table will not
affect the final result; other methods should be used for modification.

#### `es:hasrequired("💀")` // bool

Returns true if the passed emoji is going to be included in the reaction.

#### `es:addrequired("💀")`

Add an emoji to the reaction.

#### `es:removerequired("💀")`

Remove an emoji that is going to be included in the reaction.

#### `es:message()` // string

Returns the message that triggered the reaction.

#### `es:messagematches("regexp")` // bool

True if the message matches the passed regex.
[Docs for regex syntax are here](https://golang.org/pkg/regexp/syntax/).

#### `es:words()` // table of tokenized words

Returns a copy of the tokenized words.  Tokenization of words happens on all
non-alpha characters and the message is lowerecased.

#### `es:hasword("word")` // bool

True if the word is included in the message.

#### `es:len()` // int

Total count of items in optional and required.

All of the following are thin veneers atop
[github.com/hackebrot/turtle](https://github.com/hackebrot/turtle):

 * `turtle.findbyname("skull")` // turtleemoji
 * `turtle.findbychar("💀")` // turtleemoji
 * `turtle.searchbycategory("people")` // table of turtleemoji
 * `turtle.searchbykeyword("animal")` // table of turtleemoji
 * `turtle.search("foo")` // table of turtleemoji
 * `turtleemoji#name()` // string
 * `turtleemoji#category()` // string
 * `turtleemoji#char()` // string
 * `turtleemoji#haskeyword("keyword")` // bool
