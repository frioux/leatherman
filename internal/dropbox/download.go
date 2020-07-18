package dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func encodeDownloadParams(path string) (string, error) {
	buf := &bytes.Buffer{}

	e := json.NewEncoder(buf)
	err := e.Encode(struct {
		Path string `json:"path"`
	}{path})
	if err != nil {
		return "", fmt.Errorf("json.Encode: %w", err)
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// Download a file
func (cl Client) Download(path string) (io.Reader, error) {
	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/download", &bytes.Buffer{})
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cl.Token)
	apiArg, err := encodeDownloadParams(path)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Dropbox-API-Arg", apiArg)

	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do: %w", err)
	}

	if resp.StatusCode > 399 {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, resp.Body); err != nil {
			return nil, fmt.Errorf("io.Copy: %w", err)
		}
		return nil, errors.New(buf.String())
	}

	return resp.Body, nil
}
