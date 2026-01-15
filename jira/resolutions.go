package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ResolutionsService handles resolution operations for the Jira API.
type ResolutionsService struct {
	client *Client
}

// List returns all resolutions.
func (s *ResolutionsService) List(ctx context.Context) ([]*Resolution, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/resolution", nil)
	if err != nil {
		return nil, nil, err
	}

	var resolutions []*Resolution
	resp, err := s.client.Do(req, &resolutions)
	if err != nil {
		return nil, resp, err
	}

	return resolutions, resp, nil
}

// Get returns a resolution by ID.
func (s *ResolutionsService) Get(ctx context.Context, resolutionID string) (*Resolution, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/resolution/%s", resolutionID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	resolution := new(Resolution)
	resp, err := s.client.Do(req, resolution)
	if err != nil {
		return nil, resp, err
	}

	return resolution, resp, nil
}

// ResolutionListResult represents a paginated list of resolutions.
type ResolutionListResult struct {
	Self       string        `json:"self,omitempty"`
	NextPage   string        `json:"nextPage,omitempty"`
	MaxResults int           `json:"maxResults,omitempty"`
	StartAt    int           `json:"startAt,omitempty"`
	Total      int           `json:"total,omitempty"`
	IsLast     bool          `json:"isLast,omitempty"`
	Values     []*Resolution `json:"values,omitempty"`
}

// Search searches for resolutions with pagination.
func (s *ResolutionsService) Search(ctx context.Context, startAt, maxResults int, ids []string, onlyDefault bool) (*ResolutionListResult, *Response, error) {
	u := "/rest/api/3/resolution/search"

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	for _, id := range ids {
		params.Add("id", id)
	}
	if onlyDefault {
		params.Set("onlyDefault", "true")
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(ResolutionListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ResolutionCreateRequest represents a request to create a resolution.
type ResolutionCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ResolutionCreateResponse represents the response from creating a resolution.
type ResolutionCreateResponse struct {
	ID string `json:"id,omitempty"`
}

// Create creates a new resolution.
func (s *ResolutionsService) Create(ctx context.Context, resolution *ResolutionCreateRequest) (*ResolutionCreateResponse, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/resolution", resolution)
	if err != nil {
		return nil, nil, err
	}

	result := new(ResolutionCreateResponse)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ResolutionUpdateRequest represents a request to update a resolution.
type ResolutionUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Update updates a resolution.
func (s *ResolutionsService) Update(ctx context.Context, resolutionID string, resolution *ResolutionUpdateRequest) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/resolution/%s", resolutionID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, resolution)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Delete removes a resolution.
func (s *ResolutionsService) Delete(ctx context.Context, resolutionID string, replaceWith string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/resolution/%s", resolutionID)

	if replaceWith != "" {
		u = fmt.Sprintf("%s?replaceWith=%s", u, replaceWith)
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SetDefault sets the default resolution.
func (s *ResolutionsService) SetDefault(ctx context.Context, resolutionID string) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, "/rest/api/3/resolution/default", map[string]string{"id": resolutionID})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Move changes the order of resolutions.
func (s *ResolutionsService) Move(ctx context.Context, ids []string, position string, after string) (*Response, error) {
	body := map[string]interface{}{
		"ids": ids,
	}
	if position != "" {
		body["position"] = position
	}
	if after != "" {
		body["after"] = after
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, "/rest/api/3/resolution/move", body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
