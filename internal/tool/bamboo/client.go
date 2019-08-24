package bamboo

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

type client struct {
	authURL, dirURL, treeURL string

	user, password string

	b *browser.Browser
}

func newClient(user, password string) client {
	return client{
		authURL: "https://ziprecruiter1.bamboohr.com/login.php",
		dirURL:  "https://ziprecruiter1.bamboohr.com/employee_directory/ajax/get_directory_info",
		treeURL: "https://ziprecruiter1.bamboohr.com/employees/orgchart.php?pin",

		user:     user,
		password: password,
	}
}

func (c *client) auth() error {
	ua := surf.NewBrowser()
	ua.SetUserAgent(lmhttp.UserAgent)
	err := ua.Open(c.authURL)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	fm, err := ua.Form("form")
	if err != nil {
		return fmt.Errorf("auth: %s", err)
	}

	err = fm.Input("username", c.user)
	if err != nil {
		return fmt.Errorf("fm.Input: %w", err)
	}
	err = fm.Input("password", c.password)
	if err != nil {
		return fmt.Errorf("fm.Input: %w", err)
	}

	err = fm.Submit()
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}
	c.b = ua

	return nil
}

func (c *client) directory(w io.Writer) error {
	if err := c.b.Open(c.dirURL); err != nil {
		return err
	}

	if _, err := c.b.Download(w); err != nil {
		return err
	}

	return nil
}

func (c *client) tree(w io.Writer) error {
	if err := c.b.Open(c.treeURL); err != nil {
		return fmt.Errorf("export-bamboohr-tree: %w", err)
	}
	buff := &bytes.Buffer{}

	if _, err := c.b.Download(buff); err != nil {
		return fmt.Errorf("export-bamboohr-tree: %w", err)
	}

	reader := bufio.NewReader(strings.NewReader(buff.String()))
	re := regexp.MustCompile("json = (.*);")

	var err error

	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if strings.Contains(line, "json = ") {
			if m := re.FindStringSubmatch(line); len(m) > 0 {
				_, err := fmt.Fprint(w, m[1])
				return err
			}
		}
	}

	return errors.New("export-bamboohr-tree: couldn't find json")
}
