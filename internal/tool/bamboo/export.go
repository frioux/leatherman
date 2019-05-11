package bamboo

import (
	"io"
	"os"
)

/*
ExportDirectory exports entire company directory as JSON.

Command: export-bamboohr
*/
func ExportDirectory([]string, io.Reader) error {
	c := newClient(os.Getenv("BAMBOO_USER"), os.Getenv("BAMBOO_PASSWORD"))
	if err := c.auth(); err != nil {
		return err
	}

	return c.directory(os.Stdout)
}

/*
ExportOrgChart exports company org chart as JSON.

Command: export-bamboohr-tree
*/
func ExportOrgChart([]string, io.Reader) error {
	c := newClient(os.Getenv("BAMBOO_USER"), os.Getenv("BAMBOO_PASSWORD"))
	if err := c.auth(); err != nil {
		return err
	}

	return c.tree(os.Stdout)
}
