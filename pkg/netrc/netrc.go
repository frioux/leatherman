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
type Netrc struct {
	Path   string
	logins []*Login
	tokens []string
}

// Login from the netrc file
type Login struct {
	Name      string
	IsDefault bool
	tokens    []string
}

// Parse the netrc file at the given path
// It returns a Netrc instance
func Parse(path string) (Netrc, error) {
	file, err := os.Open(path)
	if err != nil {
		return Netrc{}, err
	}
	netrc, err := parse(lex(file))
	if err != nil {
		return Netrc{}, err
	}
	netrc.Path = path
	return netrc, nil
}

// Machine gets a login by machine name
func (n Netrc) Machine(name string) *Login {
	for _, m := range n.logins {
		if m.Name == name {
			return m
		}
	}
	return nil
}

// MachineAndLogin gets a login by machine name and login name
func (n Netrc) MachineAndLogin(name, login string) *Login {
	for _, m := range n.logins {
		if m.Name == name && m.Get("login") == login {
			return m
		}
	}
	return nil
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
	tokens := make([]string, 0, 100)
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}
	return tokens
}

func parse(tokens []string) (Netrc, error) {
	n := Netrc{}
	n.logins = make([]*Login, 0, 20)
	var machine *Login
	for i, token := range tokens {
		// group tokens into machines
		if token == "machine" || token == "default" {
			// start new group
			machine = &Login{}
			n.logins = append(n.logins, machine)
			if token == "default" {
				machine.IsDefault = true
				machine.Name = "default"
			} else {
				machine.Name = tokens[i+2]
			}
		}
		if machine == nil {
			n.tokens = append(n.tokens, token)
		} else {
			machine.tokens = append(machine.tokens, token)
		}
	}
	return n, nil
}

// Get a property from a machine
func (m *Login) Get(name string) string {
	i := 4
	if m.IsDefault {
		i = 2
	}
	for {
		if i+2 >= len(m.tokens) {
			return ""
		}
		if m.tokens[i] == name {
			return m.tokens[i+2]
		}
		i = i + 4
	}
}
