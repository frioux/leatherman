package bamboo

import (
	"io"
	"os"
)

// ExportDirectory will write the JSON extracted from bamboohr to stdout.
func ExportDirectory([]string, io.Reader) error {
	c := newClient()
	if err := c.auth(); err != nil {
		return err
	}

	return c.directory(os.Stdout)
}

// ExportOrgChart will write the JSON extracted from the bamboohr org chart
// to stdout.
func ExportOrgChart([]string, io.Reader) error {
	c := newClient()
	if err := c.auth(); err != nil {
		return err
	}

	return c.tree(os.Stdout)
}
