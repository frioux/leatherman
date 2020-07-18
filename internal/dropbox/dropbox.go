package dropbox

import (
	"errors"
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
