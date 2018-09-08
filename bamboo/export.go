package bamboo // import "github.com/frioux/leatherman/bamboo"

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
)

func auth() (*browser.Browser, error) {
	ua := surf.NewBrowser()
	err := ua.Open("https://ziprecruiter1.bamboohr.com/login.php")
	if err != nil {
		return nil, fmt.Errorf("auth: %s", err)
	}

	fm, err := ua.Form("form")
	if err != nil {
		return nil, fmt.Errorf("auth: %s", err)
	}

	fm.Input("username", os.Getenv("BAMBOO_USER"))
	fm.Input("password", os.Getenv("BAMBOO_PASSWORD"))

	if err := fm.Submit(); err != nil {
		return nil, fmt.Errorf("auth: %s", err)
	}

	return ua, nil
}

const bamboo = "https://ziprecruiter1.bamboohr.com/employee_directory/ajax/get_directory_info"

// ExportDirectory will write the JSON extracted from bamboohr to stdout.
func ExportDirectory([]string, io.Reader) {
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
}

const tree = "https://ziprecruiter1.bamboohr.com/employees/orgchart.php?pin"

// ExportOrgChart will write the JSON extracted from the bamboohr org chart
// to stdout.
func ExportOrgChart([]string, io.Reader) {
	ua, err := auth()
	if err != nil {
		fmt.Fprintf(os.Stderr, "export-bamboohr-tree: %s\n", err)
		os.Exit(1)
	}

	err = ua.Open(tree)
	if err != nil {
		fmt.Fprintf(os.Stderr, "export-bamboohr-tree: %s\n", err)
		os.Exit(1)
	}
	buff := bytes.NewBuffer([]byte{})

	_, err = ua.Download(buff)
	if err != nil {
		fmt.Fprintf(os.Stderr, "export-bamboohr-tree: %s\n", err)
		os.Exit(1)
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
				return
			}
		}
	}

	fmt.Fprint(os.Stderr, "export-bamboohr-tree: couldn't find json\n")
	os.Exit(1)
}
