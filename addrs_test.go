package main

import (
	"github.com/go-test/deep"
	"net/mail"
	"strings"
	"testing"
)

var eml1 = `MIME-Version: 1.0
Date: Fri, 12 Jan 2018 08:13:08 -0800
Subject: Re: Monday Night Math Tutoring
From: fREW Schmidt <frioux@xxy.com>
To: K Yuz <yyz@gmail.com>,
	"foo" <foo@gmail.com>
Cc: bar <bar@gmail.com>, vaz <vaz@gmail.com>,
	buff <buff@gmail.com>
Bcc: bog <bog@gmail.com>
Content-Type: text/plain; charset="UTF-8"

yyz`

var eml2 = `MIME-Version: 1.0
Date: Thu, 11 Jan 2018 08:13:08 -0800
Subject: Re: Monday Night Math Tutoring
From: fREW Schmidt <frioux@xxy.com>
To: <a@gmail.com>, <yYz@gmail.com>
Content-Type: text/plain; charset="UTF-8"

yyz`

func TestAllAddrs(t *testing.T) {
	email, _ := mail.ReadMessage(strings.NewReader(eml1))
	expected, _ := mail.ParseAddressList(
		`"K Yuz" <yyz@gmail.com>, "foo" <foo@gmail.com>, "bar" <bar@gmail.com>,` +
			`"vaz" <vaz@gmail.com>, "buff" <buff@gmail.com>, "bog" <bog@gmail.com>`)

	got := allAddrs(email)

	if diff := deep.Equal(expected, got); diff != nil {
		t.Error(diff)
	}
}

func TestBuildAddrMap(t *testing.T) {
	addrs := "a@foo.com\tMr. A\n" +
		"b@foo.com\tMs. B\n" +
		"d@foo.com\tMrs. D\n" +
		"c@foo.com\tMx. C\n"

	got := buildAddrMap(strings.NewReader(addrs))

	expected := map[string]string{
		"a@foo.com": "Mr. A",
		"b@foo.com": "Ms. B",
		"d@foo.com": "Mrs. D",
		"c@foo.com": "Mx. C"}

	if diff := deep.Equal(expected, got); diff != nil {
		t.Error(diff)
	}
}

func TestSortAddrMap(t *testing.T) {
	got := sortAddrMap(
		map[string]float64{
			"a": 0.3,
			"b": 0.2,
			"c": 0.1},
		map[string]string{"a": "a", "b": "b", "c": "c", "d": "d"})

	expected := []string{"a\ta", "b\tb", "c\tc", "d\td"}

	if diff := deep.Equal(expected, got); diff != nil {
		t.Error(diff)
	}
}
