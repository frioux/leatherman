package bamboo

import (
	"context"
	"io"
	"errors"
	"net/http"

	"github.com/frioux/leatherman/internal/lmhttp"
)

type client struct {
	apiKey string

	companyDomain string
}

func newClient(apiKey, companyDomain string) client {
	return client{
		apiKey: apiKey,
		companyDomain: companyDomain,
	}
}

func (c *client) prefix() string {
	return "https://api.bamboohr.com/api/gateway.php/" + c.companyDomain
}

func (c *client) directory(w io.Writer) error {
	req, err := lmhttp.NewRequest(context.Background(), "GET", c.prefix() + "/v1/employees/directory", nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.apiKey, "x")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return errors.New(res.Status)
	}

	_, err = io.Copy(w, res.Body)

	return err
}
