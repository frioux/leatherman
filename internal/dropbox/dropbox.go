package dropbox

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"golang.org/x/xerrors"
)

// UploadParams maps to the parameters to file-upload, documented here:
// https://www.dropbox.com/developers/documentation/http/documentation#files-upload
type UploadParams struct {
	Path       string `json:"path"`
	Autorename bool   `json:"autorename,omitempty"`

	// Mode           string `json:"mode,omitempty"`
	// Mute           bool   `json:"mute,omitempty"`
	// StrictConflict bool   `json:"strict_conflict,omitempty"`
	// ClientModified string `json:"client_modified,omitempty"` // should be time.Time?
}

func encodeUploadParams(up UploadParams) (string, error) {
	buf := &bytes.Buffer{}

	e := json.NewEncoder(buf)
	err := e.Encode(up)
	if err != nil {
		return "", xerrors.Errorf("json.Encode: %w", err)
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// Create writes the body to path.
func Create(cl *http.Client, token string, up UploadParams, body io.Reader) error {
	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload", body)
	if err != nil {
		return xerrors.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/octet-stream")
	apiArg, err := encodeUploadParams(up)
	if err != nil {
		return err
	}
	req.Header.Set("Dropbox-API-Arg", apiArg)

	resp, err := cl.Do(req)
	if err != nil {
		return xerrors.Errorf("http.Client.Do: %w", err)
	}

	if resp.StatusCode > 399 {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, resp.Body); err != nil {
			return xerrors.Errorf("io.Copy: %w", err)
		}
		return xerrors.New(buf.String())
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
		return "", xerrors.Errorf("json.Encode: %w", err)
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// Download a file
func Download(cl *http.Client, token, path string) (io.Reader, error) {
	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/download", &bytes.Buffer{})
	if err != nil {
		return nil, xerrors.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	apiArg, err := encodeDownloadParams(path)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Dropbox-API-Arg", apiArg)

	resp, err := cl.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("http.Client.Do: %w", err)
	}

	if resp.StatusCode > 399 {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, resp.Body); err != nil {
			return nil, xerrors.Errorf("io.Copy: %w", err)
		}
		return nil, xerrors.New(buf.String())
	}

	return resp.Body, nil
}
