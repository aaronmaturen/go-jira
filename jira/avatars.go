package jira

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// AvatarsService handles avatar operations for the Jira API.
type AvatarsService struct {
	client *Client
}

// Avatar represents an avatar.
type Avatar struct {
	ID       string `json:"id,omitempty"`
	Owner    string `json:"owner,omitempty"`
	IsSystemAvatar bool `json:"isSystemAvatar,omitempty"`
	IsSelected bool   `json:"isSelected,omitempty"`
	IsDeletable bool  `json:"isDeletable,omitempty"`
	FileName   string `json:"fileName,omitempty"`
	URLs       map[string]string `json:"urls,omitempty"`
}

// Avatars represents a collection of avatars.
type Avatars struct {
	System []*Avatar `json:"system,omitempty"`
	Custom []*Avatar `json:"custom,omitempty"`
}

// GetSystemAvatars returns system avatars by type.
func (s *AvatarsService) GetSystemAvatars(ctx context.Context, avatarType string) (*Avatars, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/avatar/%s/system", avatarType)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(Avatars)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjectAvatars returns avatars for a project.
func (s *AvatarsService) GetProjectAvatars(ctx context.Context, projectIDOrKey string) (*Avatars, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/avatars", projectIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(Avatars)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SetProjectAvatar sets the avatar for a project.
func (s *AvatarsService) SetProjectAvatar(ctx context.Context, projectIDOrKey string, avatarID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/avatar", projectIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, map[string]string{"id": avatarID})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteProjectAvatar deletes an avatar from a project.
func (s *AvatarsService) DeleteProjectAvatar(ctx context.Context, projectIDOrKey string, avatarID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/avatar/%d", projectIDOrKey, avatarID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// LoadProjectAvatar loads a custom avatar for a project.
func (s *AvatarsService) LoadProjectAvatar(ctx context.Context, projectIDOrKey string, x, y, size int, data []byte) (*Avatar, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/avatar2?x=%d&y=%d&size=%d", projectIDOrKey, x, y, size)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.client.baseURL.String()+u, bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "image/png")
	req.Header.Set("X-Atlassian-Token", "no-check")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", s.client.UserAgent)

	if s.client.auth != nil {
		s.client.auth.Apply(req)
	}

	avatar := new(Avatar)
	resp, err := s.client.Do(req, avatar)
	if err != nil {
		return nil, resp, err
	}

	return avatar, resp, nil
}

// GetIssueTypeAvatars returns avatars for an issue type.
func (s *AvatarsService) GetIssueTypeAvatars(ctx context.Context, issueTypeID string) (*Avatars, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetype/%s/avatars", issueTypeID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(Avatars)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// LoadIssueTypeAvatar loads a custom avatar for an issue type.
func (s *AvatarsService) LoadIssueTypeAvatar(ctx context.Context, issueTypeID string, x, y, size int, data []byte) (*Avatar, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetype/%s/avatar2?x=%d&y=%d&size=%d", issueTypeID, x, y, size)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.client.baseURL.String()+u, bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "image/png")
	req.Header.Set("X-Atlassian-Token", "no-check")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", s.client.UserAgent)

	if s.client.auth != nil {
		s.client.auth.Apply(req)
	}

	avatar := new(Avatar)
	resp, err := s.client.Do(req, avatar)
	if err != nil {
		return nil, resp, err
	}

	return avatar, resp, nil
}

// GetUniversalAvatar returns a universal avatar image.
func (s *AvatarsService) GetUniversalAvatar(ctx context.Context, avatarType string, owningType string, ownerID string, avatarID int64, size string) (io.ReadCloser, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/universal_avatar/view/type/%s/owner/%s", avatarType, ownerID)

	params := ""
	if avatarID > 0 {
		params = fmt.Sprintf("?id=%d", avatarID)
	}
	if size != "" {
		if params == "" {
			params = "?"
		} else {
			params += "&"
		}
		params += fmt.Sprintf("size=%s", size)
	}
	u += params

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return resp.Body, newResponse(resp), nil
}
