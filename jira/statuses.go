package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// StatusesService handles status operations for the Jira API.
type StatusesService struct {
	client *Client
}

// List returns all statuses.
func (s *StatusesService) List(ctx context.Context) ([]*Status, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/status", nil)
	if err != nil {
		return nil, nil, err
	}

	var statuses []*Status
	resp, err := s.client.Do(req, &statuses)
	if err != nil {
		return nil, resp, err
	}

	return statuses, resp, nil
}

// Get returns a status by ID or name.
func (s *StatusesService) Get(ctx context.Context, statusIDOrName string) (*Status, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/status/%s", statusIDOrName)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	status := new(Status)
	resp, err := s.client.Do(req, status)
	if err != nil {
		return nil, resp, err
	}

	return status, resp, nil
}

// StatusListResult represents a paginated list of statuses.
type StatusListResult struct {
	Self       string    `json:"self,omitempty"`
	NextPage   string    `json:"nextPage,omitempty"`
	MaxResults int       `json:"maxResults,omitempty"`
	StartAt    int       `json:"startAt,omitempty"`
	Total      int       `json:"total,omitempty"`
	IsLast     bool      `json:"isLast,omitempty"`
	Values     []*Status `json:"values,omitempty"`
}

// SearchOptions specifies options for searching statuses.
type StatusSearchOptions struct {
	StartAt        int    `url:"startAt,omitempty"`
	MaxResults     int    `url:"maxResults,omitempty"`
	Expand         string `url:"expand,omitempty"`
	ProjectID      string `url:"projectId,omitempty"`
	SearchString   string `url:"searchString,omitempty"`
	StatusCategory string `url:"statusCategory,omitempty"`
}

// Search searches for statuses with pagination.
func (s *StatusesService) Search(ctx context.Context, opts *StatusSearchOptions) (*StatusListResult, *Response, error) {
	u := "/rest/api/3/statuses/search"

	if opts != nil {
		params := url.Values{}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		if opts.Expand != "" {
			params.Set("expand", opts.Expand)
		}
		if opts.ProjectID != "" {
			params.Set("projectId", opts.ProjectID)
		}
		if opts.SearchString != "" {
			params.Set("searchString", opts.SearchString)
		}
		if opts.StatusCategory != "" {
			params.Set("statusCategory", opts.StatusCategory)
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(StatusListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// BulkGet returns multiple statuses by ID.
func (s *StatusesService) BulkGet(ctx context.Context, ids []string, expand string) ([]*Status, *Response, error) {
	u := "/rest/api/3/statuses"

	params := url.Values{}
	if len(ids) > 0 {
		params.Set("id", strings.Join(ids, ","))
	}
	if expand != "" {
		params.Set("expand", expand)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var statuses []*Status
	resp, err := s.client.Do(req, &statuses)
	if err != nil {
		return nil, resp, err
	}

	return statuses, resp, nil
}

// StatusCreateInput represents the input for creating a status.
type StatusCreateInput struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	StatusCategory string `json:"statusCategory"` // TODO, IN_PROGRESS, DONE
}

// StatusCreateRequest represents a request to create statuses.
type StatusCreateRequest struct {
	Statuses []StatusCreateInput `json:"statuses"`
	Scope    *StatusScope        `json:"scope"`
}

// StatusScope represents the scope of a status.
type StatusScope struct {
	Type    string `json:"type"` // PROJECT, GLOBAL
	Project *struct {
		ID string `json:"id,omitempty"`
	} `json:"project,omitempty"`
}

// Create creates new statuses.
func (s *StatusesService) Create(ctx context.Context, request *StatusCreateRequest) ([]*Status, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/statuses", request)
	if err != nil {
		return nil, nil, err
	}

	var statuses []*Status
	resp, err := s.client.Do(req, &statuses)
	if err != nil {
		return nil, resp, err
	}

	return statuses, resp, nil
}

// StatusUpdateInput represents the input for updating a status.
type StatusUpdateInput struct {
	ID             string `json:"id"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	StatusCategory string `json:"statusCategory,omitempty"`
}

// StatusUpdateRequest represents a request to update statuses.
type StatusUpdateRequest struct {
	Statuses []StatusUpdateInput `json:"statuses"`
}

// Update updates statuses.
func (s *StatusesService) Update(ctx context.Context, request *StatusUpdateRequest) ([]*Status, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, "/rest/api/3/statuses", request)
	if err != nil {
		return nil, nil, err
	}

	var statuses []*Status
	resp, err := s.client.Do(req, &statuses)
	if err != nil {
		return nil, resp, err
	}

	return statuses, resp, nil
}

// Delete removes statuses.
func (s *StatusesService) Delete(ctx context.Context, ids []string) (*Response, error) {
	u := "/rest/api/3/statuses"

	if len(ids) > 0 {
		u = fmt.Sprintf("%s?id=%s", u, strings.Join(ids, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// StatusCategory represents a status category.
type StatusCategory struct {
	Self      string `json:"self,omitempty"`
	ID        int64  `json:"id,omitempty"`
	Key       string `json:"key,omitempty"`
	ColorName string `json:"colorName,omitempty"`
	Name      string `json:"name,omitempty"`
}

// ListCategories returns all status categories.
func (s *StatusesService) ListCategories(ctx context.Context) ([]*StatusCategory, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/statuscategory", nil)
	if err != nil {
		return nil, nil, err
	}

	var categories []*StatusCategory
	resp, err := s.client.Do(req, &categories)
	if err != nil {
		return nil, resp, err
	}

	return categories, resp, nil
}

// GetCategory returns a status category by ID or key.
func (s *StatusesService) GetCategory(ctx context.Context, categoryIDOrKey string) (*StatusCategory, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/statuscategory/%s", categoryIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	category := new(StatusCategory)
	resp, err := s.client.Do(req, category)
	if err != nil {
		return nil, resp, err
	}

	return category, resp, nil
}
