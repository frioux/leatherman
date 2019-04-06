package bamboo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
	"github.com/pkg/errors"
)

type client struct {
	authURL, dirURL, treeURL string

	b *browser.Browser
}

func newClient() client {
	return client{
		authURL: "https://ziprecruiter1.bamboohr.com/login.php",
		dirURL:  "https://ziprecruiter1.bamboohr.com/employee_directory/ajax/get_directory_info",
		treeURL: "https://ziprecruiter1.bamboohr.com/employees/orgchart.php?pin",
	}
}

func (c *client) auth() error {
	ua := surf.NewBrowser()
	err := ua.Open(c.authURL)
	if err != nil {
		return errors.Wrap(err, "auth")
	}

	fm, err := ua.Form("form")
	if err != nil {
		return fmt.Errorf("auth: %s", err)
	}

	err = fm.Input("username", os.Getenv("BAMBOO_USER"))
	if err != nil {
		return errors.Wrap(err, "fm.Input")
	}
	err = fm.Input("password", os.Getenv("BAMBOO_PASSWORD"))
	if err != nil {
		return errors.Wrap(err, "fm.Input")
	}

	err = fm.Submit()
	if err != nil {
		return errors.Wrap(err, "auth")
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
		return errors.Wrap(err, "export-bamboohr-tree")
	}
	buff := &bytes.Buffer{}

	if _, err := c.b.Download(buff); err != nil {
		return errors.Wrap(err, "export-bamboohr-tree")
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
		fmt.Printf("%q\n", line)
		if strings.Contains(line, "json = ") {
			if m := re.FindStringSubmatch(line); len(m) > 0 {
				_, err := fmt.Fprint(w, m[1])
				return err
			}
		}
	}

	return errors.New("export-bamboohr-tree: couldn't find json")
}
