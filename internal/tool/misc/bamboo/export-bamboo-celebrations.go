package bamboo

import (
	"fmt"
	"io"
	"os"
	"time"
)

func ExportCelebrations([]string, io.Reader) error {
	c := newClient(os.Getenv("BAMBOO_USER"), os.Getenv("BAMBOO_PASSWORD"))
	if err := c.auth(); err != nil {
		return err
	}

	year := fmt.Sprintf("%d", time.Now().Year())
	return c.celebrations(year + "-01-01", year + "-12-31", os.Stdout)
}
