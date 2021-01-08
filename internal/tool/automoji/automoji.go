package automoji

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/frioux/leatherman/internal/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
)

var registry = prometheus.NewRegistry()

func mustRegister(cs ...prometheus.Collector) {
	registry.MustRegister(cs...)
	prometheus.MustRegister(cs...)
}

/*
Run comments to discord and reacts to all messages with vaguely related emoji.

The following env vars should be set:

 * LM_DROPBOX_TOKEN should be set to load a responses.json.
 * LM_BOT_LUA_PATH should be set to the location of lua to process emoji data within dropbox.
 * LM_DISCORD_TOKEN should be set for this to actually function.

Here's an example of lua code that works for this:

	if es:messagematches("cronos") then
		es:addrequired("👶")
		es:addrequired("🥘")
	end

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

All of the following are thin veneers atop
[github.com/hackebrot/turtle](https://github.com/hackebrot/turtle):

 * `turtle.findbyname("skull")` // turtleemoji
 * `turtle.findbychar("💀")` // turtleemoji
 * `turtleemoji#name()` // string
 * `turtleemoji#category()` // string
 * `turtleemoji#char()` // string
 * `turtleemoji#haskeyword("keyword")` // bool

Command: auto-emote
*/
func Run(args []string, _ io.Reader) error {
	var dbCl dropbox.Client
	if p := os.Getenv("LM_BOT_LUA_PATH"); p != "" {
		if strings.HasPrefix(p, "file://") {
			p = strings.TrimPrefix(p, "file://")
			b, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			luaC = string(b)
		} else {
			var err error
			dbCl, err = dropbox.NewClient(dropbox.Client{Token: os.Getenv("LM_DROPBOX_TOKEN")})
			if err != nil {
				return err
			}
			luaMu.Lock()
			luaC, err = loadLua(dbCl, p)
			if err != nil {
				luaMu.Unlock()
				return err
			}
			luaMu.Unlock()
		}
	}
	if len(args) > 1 {
		for _, arg := range args[1:] {
			es := newEmojiSet(arg)

			if err := luaEval(es, luaC); err != nil {
				return err
			}

			fmt.Println(es.all(0))
		}
		return nil
	}

	if p := os.Getenv("LM_BOT_LUA_PATH"); p != "" {
		responsesChanged := make(chan struct{})
		go func() {
			for range responsesChanged {
				var err error
				luaMu.Lock()

				luaC, err = loadLua(dbCl, p)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					luaMu.Unlock()
					continue
				}
				fmt.Fprintf(os.Stderr, "updated lua (%d bytes)\n", len(luaC))
				luaMu.Unlock()
			}
		}()
		go dbCl.Longpoll(context.Background(), filepath.Dir(p), responsesChanged)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	rand.Seed(time.Now().UnixNano())
	token := os.Getenv("LM_DISCORD_TOKEN")
	if token == "" {
		return errors.New("set LM_DISCORD_TOKEN to use auto-emote")
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(emojiAdd)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions)

	if err := dg.Open(); err != nil {
		return err
	}

	x := make(chan bool)

	<-x

	return nil
}

var maxes = map[int]int{
	0: 1,
	1: 1,
	2: 1,
	3: 1,
	4: 2,
	5: 10,
}

var reactTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "automoji_react_total",
	Help: "counter incremented each time a message is reacted to",
}, []string{"max"})

func init() {
	mustRegister(reactTotal)
}

func react(s *discordgo.Session, channelID, messageID string, es *emojiSet) {
	max := maxes[rand.Intn(6)]
	reactTotal.WithLabelValues(strconv.Itoa(max)).Inc()
	for i, e := range es.all(max) {
		// the 20 here is to limit to possibly fewer than were returned
		if i == 20 {
			break
		}
		if err := s.MessageReactionAdd(channelID, messageID, e); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

var messageReactionAddTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "automoji_message_reaction_add_total",
	Help: "counter incremented for each message reaction add",
}, []string{"react"})

func init() {
	mustRegister(messageReactionAddTotal)
}

func emojiAdd(s *discordgo.Session, a *discordgo.MessageReactionAdd) {
	if a.Emoji.Name != "bot" {
		messageReactionAddTotal.WithLabelValues("no").Inc()
		return
	}

	messageReactionAddTotal.WithLabelValues("yes").Inc()

	m, err := s.ChannelMessage(a.ChannelID, a.MessageID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	es := newEmojiSet(m.Content)
	if err := luaEval(es, luaC); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	react(s, a.ChannelID, a.MessageID, es)
}

var messageCreateTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "automoji_message_create_total",
	Help: "counter incremented for each message create",
}, []string{"react"})

func init() {
	mustRegister(messageCreateTotal)
}

var luaC string
var luaMu = &sync.Mutex{}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Message.Content == "||hidden knowledge||" {
		messageCreateTotal.WithLabelValues("hidden_knowledge").Inc()
		showMetrics(s, m.ChannelID)
		return
	}

	if m.Message.Content == "||version||" {
		messageCreateTotal.WithLabelValues("version").Inc()
		if _, err := s.ChannelMessageSend(m.ChannelID, version.Version); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	es := newEmojiSet(m.Message.Content)
	if err := luaEval(es, luaC); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	lucky := rand.Intn(100) == 0

	if m == nil || m.Message == nil {
		messageCreateTotal.WithLabelValues("wtf").Inc()
		return
	}

	if !lucky && len(es.required) == 0 {
		messageCreateTotal.WithLabelValues("unlucky").Inc()
		return
	}

	if len(es.required) == 0 {
		messageCreateTotal.WithLabelValues("lucky").Inc()
		es.required = append(es.required, "🎰")
	}

	react(s, m.ChannelID, m.ID, es)
}

func showMetrics(s *discordgo.Session, channelID string) {
	ms, err := registry.Gather()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	buf := &bytes.Buffer{}
	buf.Write([]byte("```\n"))
	for _, m := range ms {
		expfmt.MetricFamilyToText(buf, m)
	}
	buf.Write([]byte("```\n"))

	if _, err := s.ChannelMessageSend(channelID, buf.String()); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
