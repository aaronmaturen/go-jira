package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// IssueTypesService handles issue type operations for the Jira API.
type IssueTypesService struct {
	client *Client
}

// List returns all issue types.
func (s *IssueTypesService) List(ctx context.Context) ([]*IssueType, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/issuetype", nil)
	if err != nil {
		return nil, nil, err
	}

	var issueTypes []*IssueType
	resp, err := s.client.Do(req, &issueTypes)
	if err != nil {
		return nil, resp, err
	}

	return issueTypes, resp, nil
}

// Get returns an issue type by ID.
func (s *IssueTypesService) Get(ctx context.Context, issueTypeID string) (*IssueType, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetype/%s", issueTypeID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	issueType := new(IssueType)
	resp, err := s.client.Do(req, issueType)
	if err != nil {
		return nil, resp, err
	}

	return issueType, resp, nil
}

// IssueTypeCreateRequest represents a request to create an issue type.
type IssueTypeCreateRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Type           string `json:"type,omitempty"` // "standard" or "subtask"
	HierarchyLevel int    `json:"hierarchyLevel,omitempty"`
}

// Create creates a new issue type.
func (s *IssueTypesService) Create(ctx context.Context, issueType *IssueTypeCreateRequest) (*IssueType, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/issuetype", issueType)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssueType)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// IssueTypeUpdateRequest represents a request to update an issue type.
type IssueTypeUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	AvatarID    int64  `json:"avatarId,omitempty"`
}

// Update updates an issue type.
func (s *IssueTypesService) Update(ctx context.Context, issueTypeID string, issueType *IssueTypeUpdateRequest) (*IssueType, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetype/%s", issueTypeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, issueType)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssueType)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes an issue type.
func (s *IssueTypesService) Delete(ctx context.Context, issueTypeID string, alternativeIssueTypeID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetype/%s", issueTypeID)

	if alternativeIssueTypeID != "" {
		u = fmt.Sprintf("%s?alternativeIssueTypeId=%s", u, alternativeIssueTypeID)
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetAlternatives returns alternative issue types for when an issue type is deleted.
func (s *IssueTypesService) GetAlternatives(ctx context.Context, issueTypeID string) ([]*IssueType, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetype/%s/alternatives", issueTypeID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var issueTypes []*IssueType
	resp, err := s.client.Do(req, &issueTypes)
	if err != nil {
		return nil, resp, err
	}

	return issueTypes, resp, nil
}

// LoadAvatar loads an avatar for an issue type.
func (s *IssueTypesService) LoadAvatar(ctx context.Context, issueTypeID string, x, y, size int, filename string, avatarData []byte) (*Avatar, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetype/%s/avatar2?x=%d&y=%d&size=%d", issueTypeID, x, y, size)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "image/png")
	req.Header.Set("X-Atlassian-Token", "no-check")

	avatar := new(Avatar)
	resp, err := s.client.Do(req, avatar)
	if err != nil {
		return nil, resp, err
	}

	return avatar, resp, nil
}

// IssueTypeScheme represents an issue type scheme.
type IssueTypeScheme struct {
	ID                 string   `json:"id,omitempty"`
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty"`
	DefaultIssueTypeID string   `json:"defaultIssueTypeId,omitempty"`
	IsDefault          bool     `json:"isDefault,omitempty"`
}

// IssueTypeSchemeListResult represents a paginated list of issue type schemes.
type IssueTypeSchemeListResult struct {
	Self       string             `json:"self,omitempty"`
	NextPage   string             `json:"nextPage,omitempty"`
	MaxResults int                `json:"maxResults,omitempty"`
	StartAt    int                `json:"startAt,omitempty"`
	Total      int                `json:"total,omitempty"`
	IsLast     bool               `json:"isLast,omitempty"`
	Values     []*IssueTypeScheme `json:"values,omitempty"`
}

// ListSchemes returns all issue type schemes.
func (s *IssueTypesService) ListSchemes(ctx context.Context, startAt, maxResults int, ids []int64, expand string, queryString string) (*IssueTypeSchemeListResult, *Response, error) {
	u := "/rest/api/3/issuetypescheme"

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
	if expand != "" {
		params.Set("expand", expand)
	}
	if queryString != "" {
		params.Set("queryString", queryString)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssueTypeSchemeListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// IssueTypeSchemeCreateRequest represents a request to create an issue type scheme.
type IssueTypeSchemeCreateRequest struct {
	Name               string   `json:"name"`
	Description        string   `json:"description,omitempty"`
	DefaultIssueTypeID string   `json:"defaultIssueTypeId,omitempty"`
	IssueTypeIDs       []string `json:"issueTypeIds"`
}

// IssueTypeSchemeCreateResponse represents the response from creating an issue type scheme.
type IssueTypeSchemeCreateResponse struct {
	IssueTypeSchemeID string `json:"issueTypeSchemeId,omitempty"`
}

// CreateScheme creates an issue type scheme.
func (s *IssueTypesService) CreateScheme(ctx context.Context, scheme *IssueTypeSchemeCreateRequest) (*IssueTypeSchemeCreateResponse, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/issuetypescheme", scheme)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssueTypeSchemeCreateResponse)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// IssueTypeSchemeUpdateRequest represents a request to update an issue type scheme.
type IssueTypeSchemeUpdateRequest struct {
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	DefaultIssueTypeID string `json:"defaultIssueTypeId,omitempty"`
}

// UpdateScheme updates an issue type scheme.
func (s *IssueTypesService) UpdateScheme(ctx context.Context, schemeID int64, scheme *IssueTypeSchemeUpdateRequest) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetypescheme/%d", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, scheme)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteScheme removes an issue type scheme.
func (s *IssueTypesService) DeleteScheme(ctx context.Context, schemeID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetypescheme/%d", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// AddIssueTypesToScheme adds issue types to a scheme.
func (s *IssueTypesService) AddIssueTypesToScheme(ctx context.Context, schemeID int64, issueTypeIDs []string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetypescheme/%d/issuetype", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, map[string][]string{"issueTypeIds": issueTypeIDs})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RemoveIssueTypeFromScheme removes an issue type from a scheme.
func (s *IssueTypesService) RemoveIssueTypeFromScheme(ctx context.Context, schemeID int64, issueTypeID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetypescheme/%d/issuetype/%s", schemeID, issueTypeID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ReorderIssueTypesInScheme changes the order of issue types in a scheme.
func (s *IssueTypesService) ReorderIssueTypesInScheme(ctx context.Context, schemeID int64, issueTypeIDs []string, position string, after string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetypescheme/%d/issuetype/move", schemeID)

	body := map[string]interface{}{
		"issueTypeIds": issueTypeIDs,
	}
	if position != "" {
		body["position"] = position
	}
	if after != "" {
		body["after"] = after
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ListProjectsForScheme returns projects using an issue type scheme.
func (s *IssueTypesService) ListProjectsForScheme(ctx context.Context, schemeID int64, startAt, maxResults int) (*ProjectListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetypescheme/%d/project", schemeID)

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

// AssignSchemeToProject assigns an issue type scheme to a project.
func (s *IssueTypesService) AssignSchemeToProject(ctx context.Context, schemeID int64, projectID string) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, "/rest/api/3/issuetypescheme/project", map[string]interface{}{
		"issueTypeSchemeId": strconv.FormatInt(schemeID, 10),
		"projectId":         projectID,
	})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetIssueTypesForProject returns issue types for a project.
func (s *IssueTypesService) GetIssueTypesForProject(ctx context.Context, projectID int64, level int) ([]*IssueType, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuetype/project?projectId=%d", projectID)

	if level > 0 {
		u = fmt.Sprintf("%s&level=%d", u, level)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var issueTypes []*IssueType
	resp, err := s.client.Do(req, &issueTypes)
	if err != nil {
		return nil, resp, err
	}

	return issueTypes, resp, nil
}
