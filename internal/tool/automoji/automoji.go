package automoji

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
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

Command: auto-emote
*/
func Run(args []string, _ io.Reader) error {
	if len(args) > 1 {
		for _, arg := range args[1:] {
			fmt.Println(newEmojiSet(arg).all(0))
		}
		return nil
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
		if i == max || i == 20 {
			break
		}
		s.MessageReactionAdd(channelID, messageID, e)
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

	react(s, a.ChannelID, a.MessageID, newEmojiSet(m.Content))
}

var messageCreateTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "automoji_message_create_total",
	Help: "counter incremented for each message create",
}, []string{"react"})

func init() {
	mustRegister(messageCreateTotal)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Message.Content == "||hidden knowledge||" {
		messageCreateTotal.WithLabelValues("hidden_knowledge").Inc()
		showMetrics(s, m.ChannelID)
		return
	}

	if m.Message.Content == "||version||" {
		messageCreateTotal.WithLabelValues("version").Inc()
		s.ChannelMessageSend(m.ChannelID, version.Version)
		return
	}

	es := newEmojiSet(m.Message.Content)

	if strings.Contains(m.Message.Content, "did no back reading") ||
		strings.Contains(m.Message.Content, "have no back scroll") ||
		strings.Contains(m.Message.Content, "have no scroll back") ||
		strings.Contains(m.Message.Content, "have no scrollback") {
		es.required = append(es.required, "costanza")
	}

	lucky := rand.Intn(100) == 0

	if m == nil || m.Message == nil {
		messageCreateTotal.WithLabelValues("wtf").Inc()
		return
	}

	if !lucky {
		messageCreateTotal.WithLabelValues("unlucky").Inc()
		return
	}

	messageCreateTotal.WithLabelValues("lucky").Inc()
	es.required = append(es.required, "🎰")

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

	s.ChannelMessageSend(channelID, buf.String())
}
