package automoji

import (
	"os"
	"testing"

	"github.com/frioux/leatherman/internal/dropbox"
)

func init() {
	if os.Getenv("LM_BOT_LUA_PATH") == "" {
		return
	}
	var dbCl dropbox.Client
	if t := os.Getenv("LM_DROPBOX_TOKEN"); t != "" {
		var err error
		dbCl, err = dropbox.NewClient(dropbox.Client{Token: os.Getenv("LM_DROPBOX_TOKEN")})
		if err != nil {
			panic(err)
		}
	}
	if err := loadLua(dbCl, os.Getenv("LM_BOT_LUA_PATH")); err != nil {
		panic(err)
	}
}

func BenchmarkAutomojiAlpha(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := newEmojiSet("a b c d e f g h i j k l m n o p q r s t u v w x y z")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkAutomojiOneWord(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := newEmojiSet("hello!")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkAutomojiNeilGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := newEmojiSet("gear grip strength hiccup bleed garbage tourist wriggle miscarriage crash trait feedback application relative prince hilarious matrix reserve velvet account good trick invite attractive disorder period drawer harm monk land cower governor knowledge pedestrian payment sniff beautiful nominate color possession width facility embryo thick refer wind moon mutter battle prove")
		if err != nil {
			panic(err)
		}
	}
}
