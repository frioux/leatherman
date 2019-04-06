package netrc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"unicode"
)

// ErrInvalidNetrc means there was an error parsing the netrc file
var ErrInvalidNetrc = errors.New("Invalid netrc")

// Netrc file
type Netrc []Login

// Login from the netrc file
type Login struct {
	IsDefault bool

	Name, Login, Password, Account, Macdef string
}

// IsZero tells whether you got a real Login or an (effectively) nil Login
func (l Login) IsZero() bool {
	return !l.IsDefault && l.Name == "" && l.Login == "" && l.Password == "" &&
		l.Account == "" && l.Macdef == ""
}

// Parse the netrc file at the given path
// It returns a Netrc instance
func Parse(path string) (Netrc, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	netrc, err := parse(lex(file))
	if err != nil {
		return nil, err
	}
	return netrc, nil
}

// Machine gets a login by machine name
func (n Netrc) Machine(name string) Login {
	for _, m := range n {
		if m.Name == name {
			return m
		}
	}
	return Login{}
}

// MachineAndLogin gets a login by machine name and login name
func (n Netrc) MachineAndLogin(name, login string) Login {
	for _, m := range n {
		if m.Name == name && m.Login == login {
			return m
		}
	}
	return Login{}
}

func lex(file io.Reader) []string {
	commentRe := regexp.MustCompile("\\s*#")
	scanner := bufio.NewScanner(file)
	scanner.Split(func(data []byte, eof bool) (int, []byte, error) {
		if eof && len(data) == 0 {
			return 0, nil, nil
		}
		inWhitespace := unicode.IsSpace(rune(data[0]))
		for i, c := range data {
			if c == '#' {
				// line has a comment
				i = commentRe.FindIndex(data)[0]
				if i == 0 {
					// currently in a comment
					i = bytes.IndexByte(data, '\n')
					if i == -1 {
						// no newline at end
						if !eof {
							return 0, nil, nil
						}
						i = len(data)
					}
					for i < len(data) {
						if !unicode.IsSpace(rune(data[i])) {
							break
						}
						i++
					}
				}
				return i, data[0:i], nil
			}
			if unicode.IsSpace(rune(c)) != inWhitespace {
				return i, data[0:i], nil
			}
		}
		if eof {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	tokens := []string{}
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}
	return tokens
}

func parse(tokens []string) (Netrc, error) {
	n := Netrc([]Login{})
	var machine Login
	for i, token := range tokens {
		// group tokens into machines
		if token == "machine" || token == "default" {
			// start new group
			n = append(n, machine)
			machine = Login{}
			if token == "default" {
				machine.IsDefault = true
				machine.Name = "default"
			} else {
				machine.Name = tokens[i+2]
			}
		}
		switch token {
		case "login":
			machine.Login = tokens[i+2]
		case "password":
			machine.Password = tokens[i+2]
		case "account":
			machine.Account = tokens[i+2]
		case "macdef":
			machine.Macdef = tokens[i+2]
		}
	}
	n = append(n, machine)
	return n, nil
}
