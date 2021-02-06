package bamboo

import (
	"io"
	"os"
)

/*
ExportOrgChart exports company org chart as JSON.
*/
func ExportOrgChart([]string, io.Reader) error {
	c := newClient(os.Getenv("BAMBOO_USER"), os.Getenv("BAMBOO_PASSWORD"))
	if err := c.auth(); err != nil {
		return err
	}

	return c.tree(os.Stdout)
}
