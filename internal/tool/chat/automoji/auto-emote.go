package automoji

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/frioux/leatherman/internal/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	lua "github.com/yuin/gopher-lua"
)

var registry = prometheus.NewRegistry()

func mustRegister(cs ...prometheus.Collector) {
	registry.MustRegister(cs...)
	prometheus.MustRegister(cs...)
}

var luaC string

func Run(args []string, _ io.Reader) error {
	fs := flag.NewFlagSet("automoji", flag.ContinueOnError)
	var bench bool
	fs.BoolVar(&bench, "bench", false, "run standard benchmarks against the current lua code")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	var dbCl dropbox.Client
	if p := os.Getenv("LM_BOT_LUA_PATH"); p != "" {
		var err error
		dbCl, err = dropbox.NewClient(dropbox.Client{Token: os.Getenv("LM_DROPBOX_TOKEN")})
		if err != nil {
			return err
		}
		if err := loadLua(dbCl, p); err != nil {
			return err
		}
	}
	if len(fs.Args()) > 0 {
		for _, arg := range fs.Args() {
			t0 := time.Now()
			es, err := newEmojiSet(arg)
			if err != nil {
				return err
			}

			optional := make([]string, 0, len(es.optional))
			for o := range es.optional {
				optional = append(optional, o)
			}

			fmt.Println("lua time", time.Now().Sub(t0))
			fmt.Println("required", es.required)
			sort.Strings(optional)
			fmt.Println("optional", optional)
		}
		return nil
	}

	if p := os.Getenv("LM_BOT_LUA_PATH"); p != "" {
		responsesChanged := make(chan []dropbox.Metadata)
		go func() {
			for range responsesChanged {
				if err := loadLua(dbCl, p); err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				fmt.Fprintf(os.Stderr, "updated lua (%d bytes)\n", len(luaC))
			}
		}()
		go dbCl.Longpoll(context.Background(), filepath.Dir(p), responsesChanged)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

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

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsDirectMessages | discordgo.IntentsDirectMessageReactions)

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
	if (a.GuildID != "" && a.Emoji.Name != "bot") || (a.GuildID == "" && a.Emoji.Name != "🤖") {
		messageReactionAddTotal.WithLabelValues("no").Inc()
		return
	}

	messageReactionAddTotal.WithLabelValues("yes").Inc()

	m, err := s.ChannelMessage(a.ChannelID, a.MessageID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	es, err := newEmojiSet(m.Content)
	if err != nil {
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

var luaFn *lua.FunctionProto
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

	es, err := newEmojiSet(m.Message.Content)
	if err != nil {
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
