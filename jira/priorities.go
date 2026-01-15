package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// PrioritiesService handles priority operations for the Jira API.
type PrioritiesService struct {
	client *Client
}

// List returns all priorities.
func (s *PrioritiesService) List(ctx context.Context) ([]*Priority, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/priority", nil)
	if err != nil {
		return nil, nil, err
	}

	var priorities []*Priority
	resp, err := s.client.Do(req, &priorities)
	if err != nil {
		return nil, resp, err
	}

	return priorities, resp, nil
}

// Get returns a priority by ID.
func (s *PrioritiesService) Get(ctx context.Context, priorityID string) (*Priority, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/priority/%s", priorityID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	priority := new(Priority)
	resp, err := s.client.Do(req, priority)
	if err != nil {
		return nil, resp, err
	}

	return priority, resp, nil
}

// PriorityListResult represents a paginated list of priorities.
type PriorityListResult struct {
	Self       string      `json:"self,omitempty"`
	NextPage   string      `json:"nextPage,omitempty"`
	MaxResults int         `json:"maxResults,omitempty"`
	StartAt    int         `json:"startAt,omitempty"`
	Total      int         `json:"total,omitempty"`
	IsLast     bool        `json:"isLast,omitempty"`
	Values     []*Priority `json:"values,omitempty"`
}

// Search searches for priorities with pagination.
func (s *PrioritiesService) Search(ctx context.Context, startAt, maxResults int, ids []string, projectIDs []string, onlyDefault bool) (*PriorityListResult, *Response, error) {
	u := "/rest/api/3/priority/search"

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
	for _, pid := range projectIDs {
		params.Add("projectId", pid)
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

	result := new(PriorityListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// PriorityCreateRequest represents a request to create a priority.
type PriorityCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"iconUrl,omitempty"`
	StatusColor string `json:"statusColor,omitempty"`
}

// PriorityCreateResponse represents the response from creating a priority.
type PriorityCreateResponse struct {
	ID string `json:"id,omitempty"`
}

// Create creates a new priority.
func (s *PrioritiesService) Create(ctx context.Context, priority *PriorityCreateRequest) (*PriorityCreateResponse, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/priority", priority)
	if err != nil {
		return nil, nil, err
	}

	result := new(PriorityCreateResponse)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// PriorityUpdateRequest represents a request to update a priority.
type PriorityUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"iconUrl,omitempty"`
	StatusColor string `json:"statusColor,omitempty"`
}

// Update updates a priority.
func (s *PrioritiesService) Update(ctx context.Context, priorityID string, priority *PriorityUpdateRequest) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/priority/%s", priorityID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, priority)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Delete removes a priority.
func (s *PrioritiesService) Delete(ctx context.Context, priorityID string, newPriorityID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/priority/%s", priorityID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SetDefault sets the default priority.
func (s *PrioritiesService) SetDefault(ctx context.Context, priorityID string) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, "/rest/api/3/priority/default", map[string]string{"id": priorityID})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Move changes the order of priorities.
func (s *PrioritiesService) Move(ctx context.Context, ids []string, position string, after string) (*Response, error) {
	body := map[string]interface{}{
		"ids": ids,
	}
	if position != "" {
		body["position"] = position
	}
	if after != "" {
		body["after"] = after
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, "/rest/api/3/priority/move", body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// PriorityScheme represents a priority scheme.
type PriorityScheme struct {
	ID                string      `json:"id,omitempty"`
	Self              string      `json:"self,omitempty"`
	Name              string      `json:"name,omitempty"`
	Description       string      `json:"description,omitempty"`
	DefaultPriorityID string      `json:"defaultPriorityId,omitempty"`
	Priorities        []*Priority `json:"priorities,omitempty"`
	ProjectIDs        []string    `json:"projectIds,omitempty"`
	IsDefault         bool        `json:"isDefault,omitempty"`
}

// PrioritySchemeListResult represents a paginated list of priority schemes.
type PrioritySchemeListResult struct {
	Self       string            `json:"self,omitempty"`
	NextPage   string            `json:"nextPage,omitempty"`
	MaxResults int               `json:"maxResults,omitempty"`
	StartAt    int               `json:"startAt,omitempty"`
	Total      int               `json:"total,omitempty"`
	IsLast     bool              `json:"isLast,omitempty"`
	Values     []*PriorityScheme `json:"values,omitempty"`
}

// ListSchemes returns all priority schemes.
func (s *PrioritiesService) ListSchemes(ctx context.Context, startAt, maxResults int, ids []int64, onlyDefault bool, expand string) (*PrioritySchemeListResult, *Response, error) {
	u := "/rest/api/3/priorityscheme"

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	for _, id := range ids {
		params.Add("id", strconv.FormatInt(id, 10))
	}
	if onlyDefault {
		params.Set("onlyDefault", "true")
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

	result := new(PrioritySchemeListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetScheme returns a priority scheme by ID.
func (s *PrioritiesService) GetScheme(ctx context.Context, schemeID string, expand string) (*PriorityScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/priorityscheme/%s", schemeID)

	if expand != "" {
		u = fmt.Sprintf("%s?expand=%s", u, expand)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	scheme := new(PriorityScheme)
	resp, err := s.client.Do(req, scheme)
	if err != nil {
		return nil, resp, err
	}

	return scheme, resp, nil
}

// PrioritySchemeCreateRequest represents a request to create a priority scheme.
type PrioritySchemeCreateRequest struct {
	Name              string   `json:"name"`
	Description       string   `json:"description,omitempty"`
	DefaultPriorityID int64    `json:"defaultPriorityId,omitempty"`
	PriorityIDs       []int64  `json:"priorityIds,omitempty"`
	ProjectIDs        []int64  `json:"projectIds,omitempty"`
	Mappings          map[string]string `json:"mappings,omitempty"`
}

// CreateScheme creates a priority scheme.
func (s *PrioritiesService) CreateScheme(ctx context.Context, scheme *PrioritySchemeCreateRequest) (*PriorityScheme, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/priorityscheme", scheme)
	if err != nil {
		return nil, nil, err
	}

	result := new(PriorityScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// PrioritySchemeUpdateRequest represents a request to update a priority scheme.
type PrioritySchemeUpdateRequest struct {
	Name              string            `json:"name,omitempty"`
	Description       string            `json:"description,omitempty"`
	DefaultPriorityID int64             `json:"defaultPriorityId,omitempty"`
	PriorityIDs       []int64           `json:"priorityIds,omitempty"`
	Mappings          map[string]string `json:"mappings,omitempty"`
}

// UpdateScheme updates a priority scheme.
func (s *PrioritiesService) UpdateScheme(ctx context.Context, schemeID string, scheme *PrioritySchemeUpdateRequest) (*PriorityScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/priorityscheme/%s", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, scheme)
	if err != nil {
		return nil, nil, err
	}

	result := new(PriorityScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteScheme removes a priority scheme.
func (s *PrioritiesService) DeleteScheme(ctx context.Context, schemeID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/priorityscheme/%s", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ListSchemeProjects returns projects using a priority scheme.
func (s *PrioritiesService) ListSchemeProjects(ctx context.Context, schemeID string, startAt, maxResults int) (*ProjectListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/priorityscheme/%s/project", schemeID)

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// AssignSchemeToProjects assigns a priority scheme to projects.
func (s *PrioritiesService) AssignSchemeToProjects(ctx context.Context, schemeID string, projectIDs []int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/priorityscheme/%s/project", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, map[string][]int64{"projectIds": projectIDs})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UnassignSchemeFromProjects removes a priority scheme from projects.
func (s *PrioritiesService) UnassignSchemeFromProjects(ctx context.Context, schemeID string, projectIDs []int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/priorityscheme/%s/project", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, map[string][]int64{"projectIds": projectIDs})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
