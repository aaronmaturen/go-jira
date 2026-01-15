package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ComponentsService handles component operations for the Jira API.
type ComponentsService struct {
	client *Client
}

// Get returns a component by ID.
func (s *ComponentsService) Get(ctx context.Context, componentID string) (*Component, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/component/%s", componentID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	component := new(Component)
	resp, err := s.client.Do(req, component)
	if err != nil {
		return nil, resp, err
	}

	return component, resp, nil
}

// ComponentCreateRequest represents a request to create a component.
type ComponentCreateRequest struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description,omitempty"`
	LeadAccountID string                 `json:"leadAccountId,omitempty"`
	LeadUserName  string                 `json:"leadUserName,omitempty"`
	AssigneeType  string                 `json:"assigneeType,omitempty"` // PROJECT_DEFAULT, COMPONENT_LEAD, PROJECT_LEAD, UNASSIGNED
	Project       string                 `json:"project,omitempty"`
	ProjectID     int64                  `json:"projectId,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Create creates a new component.
func (s *ComponentsService) Create(ctx context.Context, component *ComponentCreateRequest) (*Component, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/component", component)
	if err != nil {
		return nil, nil, err
	}

	result := new(Component)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ComponentUpdateRequest represents a request to update a component.
type ComponentUpdateRequest struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	LeadAccountID string `json:"leadAccountId,omitempty"`
	LeadUserName  string `json:"leadUserName,omitempty"`
	AssigneeType  string `json:"assigneeType,omitempty"`
}

// Update updates a component.
func (s *ComponentsService) Update(ctx context.Context, componentID string, component *ComponentUpdateRequest) (*Component, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/component/%s", componentID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, component)
	if err != nil {
		return nil, nil, err
	}

	result := new(Component)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a component.
func (s *ComponentsService) Delete(ctx context.Context, componentID string, moveIssuesTo string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/component/%s", componentID)

	if moveIssuesTo != "" {
		u = fmt.Sprintf("%s?moveIssuesTo=%s", u, moveIssuesTo)
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ComponentIssueCount represents the issue count for a component.
type ComponentIssueCount struct {
	Self       string `json:"self,omitempty"`
	IssueCount int    `json:"issueCount,omitempty"`
}

// GetIssueCount returns the issue count for a component.
func (s *ComponentsService) GetIssueCount(ctx context.Context, componentID string) (*ComponentIssueCount, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/component/%s/relatedIssueCounts", componentID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	count := new(ComponentIssueCount)
	resp, err := s.client.Do(req, count)
	if err != nil {
		return nil, resp, err
	}

	return count, resp, nil
}

// ComponentListResult represents a paginated list of components.
type ComponentListResult struct {
	Self       string       `json:"self,omitempty"`
	NextPage   string       `json:"nextPage,omitempty"`
	MaxResults int          `json:"maxResults,omitempty"`
	StartAt    int          `json:"startAt,omitempty"`
	Total      int          `json:"total,omitempty"`
	IsLast     bool         `json:"isLast,omitempty"`
	Values     []*Component `json:"values,omitempty"`
}

// ListProjectComponents returns components for a project.
func (s *ComponentsService) ListProjectComponents(ctx context.Context, projectIDOrKey string, startAt, maxResults int, orderBy, query string) (*ComponentListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/component", projectIDOrKey)

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	if orderBy != "" {
		params.Set("orderBy", orderBy)
	}
	if query != "" {
		params.Set("query", query)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(ComponentListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListAllProjectComponents returns all components for a project (non-paginated).
func (s *ComponentsService) ListAllProjectComponents(ctx context.Context, projectIDOrKey string) ([]*Component, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/components", projectIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var components []*Component
	resp, err := s.client.Do(req, &components)
	if err != nil {
		return nil, resp, err
	}

	return components, resp, nil
}

// FindComponentsUsedByUser returns components the user has permission for.
func (s *ComponentsService) FindComponentsUsedByUser(ctx context.Context, projectIDs []int64, componentIDs []int64, startAt, maxResults int) (*ComponentListResult, *Response, error) {
	u := "/rest/api/3/component"

	params := url.Values{}
	for _, pid := range projectIDs {
		params.Add("projectIds", strconv.FormatInt(pid, 10))
	}
	for _, cid := range componentIDs {
		params.Add("componentIds", strconv.FormatInt(cid, 10))
	}
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

	result := new(ComponentListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
