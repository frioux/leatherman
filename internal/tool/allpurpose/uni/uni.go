package uni

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/runenames"
)

func Describe(args []string, _ io.Reader) error {
	if len(args) < 2 {
		fmt.Printf("Usage: %s <string>\n", args[0])
		return nil
	}

	for i, arg := range args[1:] {
		if i != 0 {
			fmt.Println()
		}
		for _, r := range arg {
			fmt.Println(describe(r))
		}
	}

	return nil
}

func describe(r rune) string {
	name := runenames.Name(r)

	var t []string
	if unicode.IsControl(r) {
		t = append(t, "control")
	}
	if unicode.IsDigit(r) {
		t = append(t, "digit")
	}
	if unicode.IsGraphic(r) {
		t = append(t, "graphic")
	}
	if unicode.IsLetter(r) {
		t = append(t, "letter")
	}
	if unicode.IsLower(r) {
		t = append(t, "lower")
	}
	if unicode.IsMark(r) {
		t = append(t, "mark")
	}
	if unicode.IsNumber(r) {
		t = append(t, "number")
	}
	if unicode.IsPrint(r) {
		t = append(t, "printable")
	}
	if unicode.IsPunct(r) {
		t = append(t, "punct")
	}
	if unicode.IsSpace(r) {
		t = append(t, "space")
	}
	if unicode.IsSymbol(r) {
		t = append(t, "symbol")
	}
	if unicode.IsTitle(r) {
		t = append(t, "title")
	}
	if unicode.IsUpper(r) {
		t = append(t, "upper")
	}
	return fmt.Sprintf("%q @ %d aka %s ( %s )", r, r, name, strings.Join(t, " | "))
}
