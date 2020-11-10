package automoji

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

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

func react(s *discordgo.Session, channelID, messageID string, es *emojiSet) {
	max := maxes[rand.Intn(6)]
	for i, e := range es.all(max) {
		// the 20 here is to limit to possibly fewer than were returned
		if i == max || i == 20 {
			break
		}
		s.MessageReactionAdd(channelID, messageID, e)
	}
}

func emojiAdd(s *discordgo.Session, a *discordgo.MessageReactionAdd) {
	if a.Emoji.Name != "bot" {
		return
	}

	m, err := s.ChannelMessage(a.ChannelID, a.MessageID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	react(s, a.ChannelID, a.MessageID, newEmojiSet(m.Content))
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	es := newEmojiSet(m.Message.Content)

	lucky := rand.Intn(100) == 0

	if m == nil || m.Message == nil || !lucky {
		return
	}

	if lucky {
		es.required = append(es.required, "🎰")
	}

	react(s, m.ChannelID, m.ID, es)
}
