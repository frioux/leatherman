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

// UploadParams maps to the parameters to file-upload, documented here:
// https://www.dropbox.com/developers/documentation/http/documentation#files-upload
type UploadParams struct {
	Path       string `json:"path"`
	Autorename bool   `json:"autorename,omitempty"`

	Mode string `json:"mode,omitempty"`
	// Mute           bool   `json:"mute,omitempty"`
	// StrictConflict bool   `json:"strict_conflict,omitempty"`
	// ClientModified string `json:"client_modified,omitempty"` // should be time.Time?
}

func encodeUploadParams(up UploadParams) (string, error) {
	buf := &bytes.Buffer{}

	e := json.NewEncoder(buf)
	err := e.Encode(up)
	if err != nil {
		return "", fmt.Errorf("json.Encode: %w", err)
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// Create writes the body to path.
func (cl Client) Create(up UploadParams, body io.Reader) error {
	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload", body)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cl.Token)
	req.Header.Set("Content-Type", "application/octet-stream")
	apiArg, err := encodeUploadParams(up)
	if err != nil {
		return err
	}
	req.Header.Set("Dropbox-API-Arg", apiArg)

	resp, err := cl.Do(req)
	if err != nil {
		return fmt.Errorf("http.Client.Do: %w", err)
	}

	if resp.StatusCode > 399 {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, resp.Body); err != nil {
			return fmt.Errorf("io.Copy: %w", err)
		}
		return errors.New(buf.String())
	}

	return nil
}

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
