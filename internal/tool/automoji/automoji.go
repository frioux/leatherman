package automoji

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hackebrot/turtle"
)

/*
Run comments to discord and reacts to all messages with vaguely related emoji.

Command: auto-emote
*/
func Run(args []string, _ io.Reader) error {
	if len(args) > 1 {
		for _, arg := range args[1:] {
			fmt.Println(messageToEmoji(arg))
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

func emojiAdd(s *discordgo.Session, a *discordgo.MessageReactionAdd) {
	if a.Emoji.Name != "bot" {
		return
	}

	m, err := s.ChannelMessage(a.ChannelID, a.MessageID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	react(s, a.ChannelID, a.MessageID, messageToEmoji(m.Content))
}

func react(s *discordgo.Session, channelID, messageID string, emoji []string) {
	max := maxes[rand.Intn(6)]
	for i, e := range emoji {
		// 20 is max, so limit to half the total amount
		if i == max {
			break
		}
		s.MessageReactionAdd(channelID, messageID, e)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m == nil || m.Message == nil || rand.Intn(100) != 0 {
		return
	}

	react(s, m.ChannelID, m.ID, messageToEmoji(m.Message.Content))
}

type emojiSet map[string]bool

func (s emojiSet) add(e *turtle.Emoji) {
	if e.Category == "flags" {
		return
	}

	s[e.Char] = true
}

func (s emojiSet) all() []string {
	ret := make([]string, 0, len(s))
	for e := range s {
		ret = append(ret, e)
	}

	return ret
}

func messageToEmoji(m string) []string {
	words := strings.Split(strings.ToLower(m), " ")
	s := emojiSet(make(map[string]bool, len(words)))

	for _, word := range words {
		if word == "" {
			continue
		}

		if e, ok := turtle.Emojis[word]; ok {
			s.add(e)
		}

		if es := turtle.Category(word); es != nil {
			for _, e := range es {
				s.add(e)
			}
		}

		if es := turtle.Keyword(word); es != nil {
			for _, e := range es {
				s.add(e)
			}
		}
	}
	if len(s) == 0 { // since this always finds too much, only use it when nothing is found
		for _, word := range words {
			if es := turtle.Search(word); es != nil {
				for _, e := range es {
					s.add(e)
				}
			}
		}
	}

	return s.all()
}
