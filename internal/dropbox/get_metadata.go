package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetMetadataParams maps to the parameters to get-metadata, documented here:
// https://www.dropbox.com/developers/documentation/http/documentation#files-get_metadata
type GetMetadataParams struct {
	Path string `json:"path"`

	IncludeMediaInfo                bool `json:"include_media_info,omitempty"`
	IncludeDeleted                  bool `json:"include_deleted,omitempty"`
	IncludeHasExplicitSharedMembers bool `json:"include_has_explicit_shared_members,omitempty"`

	IncludePropertyGroups []string `json:"include_property_groups,omitempty"`
}

// GetMetadata gets metadata for a file or directory.
func (cl Client) GetMetadata(p GetMetadataParams) (Metadata, error) {
	body := &bytes.Buffer{}

	e := json.NewEncoder(body)

	if err := e.Encode(p); err != nil {
		return Metadata{}, fmt.Errorf("dropbox.Client.GetMetadata: json.Encode: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/get_metadata", body)
	if err != nil {
		return Metadata{}, fmt.Errorf("dropbox.Client.GetMetadata: http.NewRequest: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cl.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cl.Do(req)
	if err != nil {
		return Metadata{}, fmt.Errorf("dropbox.Client.GetMetadata: http.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	var ret Metadata
	d := json.NewDecoder(resp.Body)

	if err := d.Decode(&ret); err != nil {
		return Metadata{}, fmt.Errorf("dropbox.Client.GetMetadata: json.Decode: %w", err)
	}

	return ret, nil

}
