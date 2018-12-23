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

func auth() (*browser.Browser, error) {
	ua := surf.NewBrowser()
	err := ua.Open("https://ziprecruiter1.bamboohr.com/login.php")
	if err != nil {
		return nil, errors.Wrap(err, "auth")
	}

	fm, err := ua.Form("form")
	if err != nil {
		return nil, fmt.Errorf("auth: %s", err)
	}

	err = fm.Input("username", os.Getenv("BAMBOO_USER"))
	if err != nil {
		return nil, errors.Wrap(err, "fm.Input")
	}
	err = fm.Input("password", os.Getenv("BAMBOO_PASSWORD"))
	if err != nil {
		return nil, errors.Wrap(err, "fm.Input")
	}

	err = fm.Submit()
	if err != nil {
		return nil, errors.Wrap(err, "auth")
	}

	return ua, nil
}

const bamboo = "https://ziprecruiter1.bamboohr.com/employee_directory/ajax/get_directory_info"

// ExportDirectory will write the JSON extracted from bamboohr to stdout.
func ExportDirectory([]string, io.Reader) error {
	ua, err := auth()
	if err != nil {
		fmt.Fprintf(os.Stderr, "export-bamboohr: %s\n", err)
		os.Exit(1)
	}

	err = ua.Open(bamboo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "export-bamboohr: %s\n", err)
		os.Exit(1)
	}
	_, err = ua.Download(os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "export-bamboohr: %s\n", err)
		os.Exit(1)
	}

	return nil
}

const tree = "https://ziprecruiter1.bamboohr.com/employees/orgchart.php?pin"

// ExportOrgChart will write the JSON extracted from the bamboohr org chart
// to stdout.
func ExportOrgChart([]string, io.Reader) error {
	ua, err := auth()
	if err != nil {
		return errors.Wrap(err, "export-bamboohr-tree")
	}

	err = ua.Open(tree)
	if err != nil {
		return errors.Wrap(err, "export-bamboohr-tree")
	}
	buff := bytes.NewBuffer([]byte{})

	_, err = ua.Download(buff)
	if err != nil {
		return errors.Wrap(err, "export-bamboohr-tree")
	}

	reader := bufio.NewReader(strings.NewReader(buff.String()))
	re := regexp.MustCompile("json = (.*);")

	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "export-bamboohr-tree: %s\n", err)
			os.Exit(1)
		}
		if strings.Contains(line, "json = ") {
			if m := re.FindStringSubmatch(line); len(m) > 0 {
				fmt.Print(m[1])
				return nil
			}
		}
	}

	return errors.New("export-bamboohr-tree: couldn't find json")
}
