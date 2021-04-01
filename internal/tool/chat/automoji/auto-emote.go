package automoji

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	lua "github.com/yuin/gopher-lua"

	"github.com/frioux/leatherman/internal/drawlua"
	"github.com/frioux/leatherman/internal/dropbox"
	"github.com/frioux/leatherman/internal/version"
)

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

func react(s *discordgo.Session, channelID, messageID string, es *emojiSet) {
	max := maxes[rand.Intn(6)]
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

func emojiAdd(s *discordgo.Session, a *discordgo.MessageReactionAdd) {
	if (a.GuildID != "" && a.Emoji.Name != "bot") || (a.GuildID == "" && a.Emoji.Name != "🤖") {
		return
	}

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

var luaFn *lua.FunctionProto
var luaMu = &sync.Mutex{}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Message.Content, "```lua\n") && strings.HasSuffix(m.Message.Content, "\n```") {
		code := m.Message.Content
		code = strings.TrimPrefix(code, "```lua\n")
		code = strings.TrimSuffix(code, "\n```")

		img := image.NewNRGBA(image.Rect(0, 0, 128, 128))

		L := lua.NewState()
		L.OpenLibs()
		L.DoString("coroutine=nil;debug=nil;io=nil;os=nil;string=nil;table=nil")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		L.SetContext(ctx)
		done := drawlua.RegisterDrawFunctions(L, img)
		defer done()

		if err := L.DoString(code); err != nil {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				if err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, "⏱"); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			} else {
				if err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, "⚠"); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
			fmt.Fprintln(os.Stderr, err)
			return
		}

		buf := &bytes.Buffer{}
		if err := png.Encode(buf, img); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		if _, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:        "drawlua.png",
					ContentType: "image/png",
					Reader:      buf,
				},
			},
		}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		return
	}

	if m.Message.Content == "||version||" {
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
		return
	}

	if !lucky && len(es.required) == 0 {
		return
	}

	if len(es.required) == 0 {
		es.required = append(es.required, "🎰")
	}

	react(s, m.ChannelID, m.ID, es)
}
