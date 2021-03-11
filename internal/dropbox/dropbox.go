package dropbox

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Client gives access to the Dropbox API
type Client struct {
	Token string
	*http.Client
}

// NewClient returns a fully created Client value
func NewClient(cl Client) (Client, error) {
	if cl.Token == "" {
		return Client{}, errors.New("Token is required")
	}

	if cl.Client == nil {
		cl.Client = &http.Client{}
	}

	return cl, nil
}

func (cl Client) handleError(resp *http.Response) error {
	if resp.StatusCode > 399 {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, resp.Body); err != nil {
			return fmt.Errorf("io.Copy: %w", err)
		}
		return errors.New(buf.String())
	}
	return nil
}
