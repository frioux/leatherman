package automoji

import (
	"fmt"
	"testing"
)

func benchmark() {
	fmt.Print("Alpha")
	r := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := newEmojiSet("a b c d e f g h i j k l m n o p q r s t u v w x y z")
			if err != nil {
				panic(err)
			}
		}
	})
	fmt.Println(r)

	fmt.Print("OneWord")
	r = testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := newEmojiSet("hello!")
			if err != nil {
				panic(err)
			}
		}
	})
	fmt.Println(r)

	fmt.Print("NeilGen")
	r = testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := newEmojiSet("gear grip strength hiccup bleed garbage tourist wriggle miscarriage crash trait feedback application relative prince hilarious matrix reserve velvet account good trick invite attractive disorder period drawer harm monk land cower governor knowledge pedestrian payment sniff beautiful nominate color possession width facility embryo thick refer wind moon mutter battle prove")
			if err != nil {
				panic(err)
			}
		}
	})
	fmt.Println(r)
}
