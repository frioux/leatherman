package bamboo

import (
	"errors"
	"io"
	"os"
)

func ExportDirectory([]string, io.Reader) error {
	k, d := os.Getenv("BAMBOO_APIKEY"), os.Getenv("BAMBOO_COMPANY_DOMAIN")
	if k == "" || d == "" {
		return errors.New("BAMBOO_APIKEY and BAMBOO_COMPANY_DOMAIN are required")
	}
	c := newClient(k, d)

	return c.directory(os.Stdout)
}
