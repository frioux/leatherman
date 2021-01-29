package dropbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ListFolderParams docs are at
// https://www.dropbox.com/developers/documentation/http/documentation#files-list_folder
type ListFolderParams struct {
	Path string `json:"path"`

	Recursive                       bool   `json:"recursive"`
	IncludeMediaInfo                bool   `json:"include_media_info"`
	IncludeDeleted                  bool   `json:"include_deleted"`
	IncludeHasExplicitSharedMembers bool   `json:"include_has_explicit_shared_members"`
	IncludeMountedFolders           bool   `json:"include_mounted_folders"`
	IncludeNonDownloadableFiles     bool   `json:"include_non_downloadable_files"`
	Limit                           uint32 `json:"limit,omitempty"`
}

// Metadata is defined at
// https://www.dropbox.com/developers/documentation/http/documentation#files-list_folder
type Metadata struct {
	Tag            string `json:".tag"`
	Name           string `json:"name"`
	ID             string `json:"id"`
	ClientModified string `json:"client_modified"` // could be time.Time
	ServerModified string `json:"server_modified"` // could be time.Time
	Rev            string `json:"rev"`
	Size           int    `json:"size"`
	PathLower      string `json:"path_lower"`
	PathDisplay    string `json:"path_display"`
	SharingInfo    struct {
		ParentSharedFolderID string `json:"parent_shared_folder_id"`
		ModifiedBy           string `json:"modified_by"`
		ReadOnly             bool   `json:"read_only"`
		TraverseOnly         bool   `json:"traverse_only"`
		NoAccess             bool   `json:"no_access"`
	} `json:"sharing_info"`
	IsDownloadable bool `json:"is_downloadable"`
	PropertyGroups []struct {
		TemplateID string `json:"template_id"`
		Fields     []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}
	} `json:"property_groups"`
	HasExplicitSharedMembers bool   `json:"has_explicit_shared_members"`
	ContentHash              string `json:"content_hash"`
	FileLockInfo             struct {
		IsLockholder   bool   `json:"is_lockholder"`
		LockholderName string `json:"lockholder_name"`
		Created        string `json:"created"` // could be time.Time
	} `json:"file_lock_info"`
}

type ListFolderResult struct {
	Entries []Metadata `json:"entries"`
	Cursor  string     `json:"cursor"`
	HasMore bool       `json:"has_more"`
}

func (cl Client) ListFolder(p ListFolderParams) (ListFolderResult, error) {
	body := &bytes.Buffer{}

	e := json.NewEncoder(body)

	if err := e.Encode(p); err != nil {
		return ListFolderResult{}, err
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/list_folder", body)
	if err != nil {
		return ListFolderResult{}, fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cl.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cl.Do(req)
	if err != nil {
		return ListFolderResult{}, fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	var ret ListFolderResult
	d := json.NewDecoder(resp.Body)

	if err := d.Decode(&ret); err != nil {
		return ListFolderResult{}, err
	}

	return ret, nil
}

func (cl Client) ListFolderContinue(cursor string) (ListFolderResult, error) {
	body := &bytes.Buffer{}

	e := json.NewEncoder(body)

	if err := e.Encode(struct {
		Cursor string `json:"cursor"`
	}{cursor}); err != nil {
		return ListFolderResult{}, err
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/list_folder/continue", body)
	if err != nil {
		return ListFolderResult{}, fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cl.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cl.Do(req)
	if err != nil {
		return ListFolderResult{}, fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	var ret ListFolderResult
	d := json.NewDecoder(resp.Body)

	if err := d.Decode(&ret); err != nil {
		return ListFolderResult{}, err
	}

	return ret, nil
}

func (cl Client) ListFolderLongPoll(ctx context.Context, cursor string, timeout int) (bool, int, error) {
	body := &bytes.Buffer{}

	// match default of a missing value, for our own timeout calculation
	// later.
	if timeout == 0 {
		timeout = 30
	}

	e := json.NewEncoder(body)

	if err := e.Encode(struct {
		Cursor  string `json:"cursor"`
		Timeout int    `json:"timeout"`
	}{cursor, timeout}); err != nil {
		return false, 0, err
	}

	// up to 90s added by dropbox to avoid thundering herd.  We add an
	// extra 1s grace.
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeout+90+1))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "https://notify.dropboxapi.com/2/files/list_folder/longpoll", body)
	if err != nil {
		return false, 0, fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := cl.Do(req)
	if err != nil {
		return false, 0, fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	var ret struct {
		Changes bool `json:"changes"`
		Backoff int  `json:"backoff"`
	}
	d := json.NewDecoder(resp.Body)

	if err := d.Decode(&ret); err != nil {
		return false, 0, err
	}

	return ret.Changes, ret.Backoff, nil
}
