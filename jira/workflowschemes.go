package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// WorkflowSchemesService handles workflow scheme operations for the Jira API.
type WorkflowSchemesService struct {
	client *Client
}

// WorkflowScheme represents a workflow scheme.
type WorkflowScheme struct {
	ID                  int64              `json:"id,omitempty"`
	Name                string             `json:"name,omitempty"`
	Description         string             `json:"description,omitempty"`
	DefaultWorkflow     string             `json:"defaultWorkflow,omitempty"`
	IssueTypeMappings   map[string]string  `json:"issueTypeMappings,omitempty"`
	OriginalDefaultWorkflow string         `json:"originalDefaultWorkflow,omitempty"`
	OriginalIssueTypeMappings map[string]string `json:"originalIssueTypeMappings,omitempty"`
	Draft               bool               `json:"draft,omitempty"`
	LastModifiedUser    *User              `json:"lastModifiedUser,omitempty"`
	LastModified        string             `json:"lastModified,omitempty"`
	Self                string             `json:"self,omitempty"`
	UpdateDraftIfNeeded bool               `json:"updateDraftIfNeeded,omitempty"`
	IssueTypes          map[string]*IssueType `json:"issueTypes,omitempty"`
}

// WorkflowSchemeListResult represents a paginated list of workflow schemes.
type WorkflowSchemeListResult struct {
	Self       string            `json:"self,omitempty"`
	NextPage   string            `json:"nextPage,omitempty"`
	MaxResults int               `json:"maxResults,omitempty"`
	StartAt    int               `json:"startAt,omitempty"`
	Total      int               `json:"total,omitempty"`
	IsLast     bool              `json:"isLast,omitempty"`
	Values     []*WorkflowScheme `json:"values,omitempty"`
}

// List returns all workflow schemes.
func (s *WorkflowSchemesService) List(ctx context.Context, startAt, maxResults int, expand string) (*WorkflowSchemeListResult, *Response, error) {
	u := "/rest/api/3/workflowscheme"

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
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

	result := new(WorkflowSchemeListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Get returns a workflow scheme by ID.
func (s *WorkflowSchemesService) Get(ctx context.Context, schemeID int64, returnDraftIfExists bool) (*WorkflowScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d", schemeID)

	if returnDraftIfExists {
		u = fmt.Sprintf("%s?returnDraftIfExists=true", u)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	scheme := new(WorkflowScheme)
	resp, err := s.client.Do(req, scheme)
	if err != nil {
		return nil, resp, err
	}

	return scheme, resp, nil
}

// WorkflowSchemeCreateRequest represents a request to create a workflow scheme.
type WorkflowSchemeCreateRequest struct {
	Name                string            `json:"name"`
	Description         string            `json:"description,omitempty"`
	DefaultWorkflow     string            `json:"defaultWorkflow,omitempty"`
	IssueTypeMappings   map[string]string `json:"issueTypeMappings,omitempty"`
}

// Create creates a workflow scheme.
func (s *WorkflowSchemesService) Create(ctx context.Context, scheme *WorkflowSchemeCreateRequest) (*WorkflowScheme, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/workflowscheme", scheme)
	if err != nil {
		return nil, nil, err
	}

	result := new(WorkflowScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// WorkflowSchemeUpdateRequest represents a request to update a workflow scheme.
type WorkflowSchemeUpdateRequest struct {
	Name                string            `json:"name,omitempty"`
	Description         string            `json:"description,omitempty"`
	DefaultWorkflow     string            `json:"defaultWorkflow,omitempty"`
	IssueTypeMappings   map[string]string `json:"issueTypeMappings,omitempty"`
	UpdateDraftIfNeeded bool              `json:"updateDraftIfNeeded,omitempty"`
}

// Update updates a workflow scheme.
func (s *WorkflowSchemesService) Update(ctx context.Context, schemeID int64, scheme *WorkflowSchemeUpdateRequest) (*WorkflowScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, scheme)
	if err != nil {
		return nil, nil, err
	}

	result := new(WorkflowScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a workflow scheme.
func (s *WorkflowSchemesService) Delete(ctx context.Context, schemeID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetDefault returns the default workflow for a scheme.
func (s *WorkflowSchemesService) GetDefault(ctx context.Context, schemeID int64, returnDraftIfExists bool) (*DefaultWorkflow, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/default", schemeID)

	if returnDraftIfExists {
		u = fmt.Sprintf("%s?returnDraftIfExists=true", u)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(DefaultWorkflow)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DefaultWorkflow represents the default workflow for a scheme.
type DefaultWorkflow struct {
	Workflow            string `json:"workflow,omitempty"`
	UpdateDraftIfNeeded bool   `json:"updateDraftIfNeeded,omitempty"`
}

// SetDefault sets the default workflow for a scheme.
func (s *WorkflowSchemesService) SetDefault(ctx context.Context, schemeID int64, workflow string, updateDraftIfNeeded bool) (*WorkflowScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/default", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, &DefaultWorkflow{
		Workflow:            workflow,
		UpdateDraftIfNeeded: updateDraftIfNeeded,
	})
	if err != nil {
		return nil, nil, err
	}

	result := new(WorkflowScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteDefault removes the default workflow from a scheme.
func (s *WorkflowSchemesService) DeleteDefault(ctx context.Context, schemeID int64, updateDraftIfNeeded bool) (*WorkflowScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/default", schemeID)

	if updateDraftIfNeeded {
		u = fmt.Sprintf("%s?updateDraftIfNeeded=true", u)
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(WorkflowScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// IssueTypeWorkflow represents an issue type workflow mapping.
type IssueTypeWorkflow struct {
	IssueType           string `json:"issueType,omitempty"`
	Workflow            string `json:"workflow,omitempty"`
	UpdateDraftIfNeeded bool   `json:"updateDraftIfNeeded,omitempty"`
}

// GetIssueTypeMapping returns the workflow for an issue type in a scheme.
func (s *WorkflowSchemesService) GetIssueTypeMapping(ctx context.Context, schemeID int64, issueType string, returnDraftIfExists bool) (*IssueTypeWorkflow, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/issuetype/%s", schemeID, issueType)

	if returnDraftIfExists {
		u = fmt.Sprintf("%s?returnDraftIfExists=true", u)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssueTypeWorkflow)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SetIssueTypeMapping sets the workflow for an issue type in a scheme.
func (s *WorkflowSchemesService) SetIssueTypeMapping(ctx context.Context, schemeID int64, issueType, workflow string, updateDraftIfNeeded bool) (*WorkflowScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/issuetype/%s", schemeID, issueType)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, &IssueTypeWorkflow{
		IssueType:           issueType,
		Workflow:            workflow,
		UpdateDraftIfNeeded: updateDraftIfNeeded,
	})
	if err != nil {
		return nil, nil, err
	}

	result := new(WorkflowScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteIssueTypeMapping removes the workflow for an issue type from a scheme.
func (s *WorkflowSchemesService) DeleteIssueTypeMapping(ctx context.Context, schemeID int64, issueType string, updateDraftIfNeeded bool) (*WorkflowScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/issuetype/%s", schemeID, issueType)

	if updateDraftIfNeeded {
		u = fmt.Sprintf("%s?updateDraftIfNeeded=true", u)
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(WorkflowScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetDraft returns the draft of a workflow scheme.
func (s *WorkflowSchemesService) GetDraft(ctx context.Context, schemeID int64) (*WorkflowScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/draft", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	draft := new(WorkflowScheme)
	resp, err := s.client.Do(req, draft)
	if err != nil {
		return nil, resp, err
	}

	return draft, resp, nil
}

// CreateDraft creates a draft of a workflow scheme.
func (s *WorkflowSchemesService) CreateDraft(ctx context.Context, schemeID int64) (*WorkflowScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/createdraft", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, nil, err
	}

	draft := new(WorkflowScheme)
	resp, err := s.client.Do(req, draft)
	if err != nil {
		return nil, resp, err
	}

	return draft, resp, nil
}

// UpdateDraft updates the draft of a workflow scheme.
func (s *WorkflowSchemesService) UpdateDraft(ctx context.Context, schemeID int64, scheme *WorkflowSchemeUpdateRequest) (*WorkflowScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/draft", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, scheme)
	if err != nil {
		return nil, nil, err
	}

	result := new(WorkflowScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteDraft removes the draft of a workflow scheme.
func (s *WorkflowSchemesService) DeleteDraft(ctx context.Context, schemeID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/draft", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// PublishDraft publishes the draft of a workflow scheme.
func (s *WorkflowSchemesService) PublishDraft(ctx context.Context, schemeID int64, statusMappings map[string]string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflowscheme/%d/draft/publish", schemeID)

	body := map[string]interface{}{}
	if statusMappings != nil {
		body["statusMappings"] = statusMappings
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
