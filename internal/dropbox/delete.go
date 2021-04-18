package dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Delete a file or folder.  Official docs are at
// https://www.dropbox.com/developers/documentation/http/documentation#files-delete
func (cl Client) Delete(path string) error {
	buf := &bytes.Buffer{}
	e := json.NewEncoder(buf)
	if err := e.Encode(struct {
		Path string `json:"path"`
	}{path}); err != nil {
		return fmt.Errorf("dropbox.Client.Delete: %w", err)
	}
	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/delete_v2", buf)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cl.Token)
	req.Header.Set("Content-Type", "application/json")

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
