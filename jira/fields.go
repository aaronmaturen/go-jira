package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// FieldsService handles field operations for the Jira API.
type FieldsService struct {
	client *Client
}

// Field represents a Jira field.
type Field struct {
	ID             string        `json:"id,omitempty"`
	Key            string        `json:"key,omitempty"`
	Name           string        `json:"name,omitempty"`
	Custom         bool          `json:"custom,omitempty"`
	Orderable      bool          `json:"orderable,omitempty"`
	Navigable      bool          `json:"navigable,omitempty"`
	Searchable     bool          `json:"searchable,omitempty"`
	ClauseNames    []string      `json:"clauseNames,omitempty"`
	Schema         *FieldSchema  `json:"schema,omitempty"`
	Scope          *FieldScope   `json:"scope,omitempty"`
	Description    string        `json:"description,omitempty"`
	IsLocked       bool          `json:"isLocked,omitempty"`
	SearcherKey    string        `json:"searcherKey,omitempty"`
	UntranslatedName string      `json:"untranslatedName,omitempty"`
}

// FieldSchema represents the schema of a field.
type FieldSchema struct {
	Type     string `json:"type,omitempty"`
	Items    string `json:"items,omitempty"`
	System   string `json:"system,omitempty"`
	Custom   string `json:"custom,omitempty"`
	CustomID int64  `json:"customId,omitempty"`
}

// FieldScope represents the scope of a field.
type FieldScope struct {
	Type    string   `json:"type,omitempty"`
	Project *Project `json:"project,omitempty"`
}

// List returns all fields.
func (s *FieldsService) List(ctx context.Context) ([]*Field, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/field", nil)
	if err != nil {
		return nil, nil, err
	}

	var fields []*Field
	resp, err := s.client.Do(req, &fields)
	if err != nil {
		return nil, resp, err
	}

	return fields, resp, nil
}

// FieldCreateRequest represents a request to create a custom field.
type FieldCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"` // e.g., "com.atlassian.jira.plugin.system.customfieldtypes:textfield"
	SearcherKey string `json:"searcherKey,omitempty"`
}

// Create creates a custom field.
func (s *FieldsService) Create(ctx context.Context, field *FieldCreateRequest) (*Field, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/field", field)
	if err != nil {
		return nil, nil, err
	}

	result := new(Field)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// FieldUpdateRequest represents a request to update a custom field.
type FieldUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SearcherKey string `json:"searcherKey,omitempty"`
}

// Update updates a custom field.
func (s *FieldsService) Update(ctx context.Context, fieldID string, field *FieldUpdateRequest) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s", fieldID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, field)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Delete removes a custom field.
func (s *FieldsService) Delete(ctx context.Context, fieldID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s", fieldID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// FieldListResult represents a paginated list of fields.
type FieldListResult struct {
	Self       string   `json:"self,omitempty"`
	NextPage   string   `json:"nextPage,omitempty"`
	MaxResults int      `json:"maxResults,omitempty"`
	StartAt    int      `json:"startAt,omitempty"`
	Total      int      `json:"total,omitempty"`
	IsLast     bool     `json:"isLast,omitempty"`
	Values     []*Field `json:"values,omitempty"`
}

// SearchOptions specifies options for searching fields.
type FieldSearchOptions struct {
	StartAt        int      `url:"startAt,omitempty"`
	MaxResults     int      `url:"maxResults,omitempty"`
	Type           []string `url:"type,omitempty"`
	ID             []string `url:"id,omitempty"`
	Query          string   `url:"query,omitempty"`
	OrderBy        string   `url:"orderBy,omitempty"`
	Expand         []string `url:"expand,omitempty"`
}

// Search searches for fields with pagination.
func (s *FieldsService) Search(ctx context.Context, opts *FieldSearchOptions) (*FieldListResult, *Response, error) {
	u := "/rest/api/3/field/search"

	if opts != nil {
		params := url.Values{}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		for _, t := range opts.Type {
			params.Add("type", t)
		}
		for _, id := range opts.ID {
			params.Add("id", id)
		}
		if opts.Query != "" {
			params.Set("query", opts.Query)
		}
		if opts.OrderBy != "" {
			params.Set("orderBy", opts.OrderBy)
		}
		if len(opts.Expand) > 0 {
			params.Set("expand", strings.Join(opts.Expand, ","))
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(FieldListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Trash moves a custom field to the trash.
func (s *FieldsService) Trash(ctx context.Context, fieldID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/trash", fieldID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Restore restores a custom field from the trash.
func (s *FieldsService) Restore(ctx context.Context, fieldID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/restore", fieldID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Context represents a custom field context.
type FieldContext struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	IsGlobalContext bool  `json:"isGlobalContext,omitempty"`
	IsAnyIssueType bool   `json:"isAnyIssueType,omitempty"`
}

// ContextListResult represents a paginated list of field contexts.
type ContextListResult struct {
	Self       string          `json:"self,omitempty"`
	NextPage   string          `json:"nextPage,omitempty"`
	MaxResults int             `json:"maxResults,omitempty"`
	StartAt    int             `json:"startAt,omitempty"`
	Total      int             `json:"total,omitempty"`
	IsLast     bool            `json:"isLast,omitempty"`
	Values     []*FieldContext `json:"values,omitempty"`
}

// ListContexts returns contexts for a custom field.
func (s *FieldsService) ListContexts(ctx context.Context, fieldID string, startAt, maxResults int, isAnyIssueType, isGlobalContext bool, contextID []int64) (*ContextListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/context", fieldID)

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	if isAnyIssueType {
		params.Set("isAnyIssueType", "true")
	}
	if isGlobalContext {
		params.Set("isGlobalContext", "true")
	}
	for _, id := range contextID {
		params.Add("contextId", strconv.FormatInt(id, 10))
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(ContextListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ContextCreateRequest represents a request to create a field context.
type ContextCreateRequest struct {
	Name           string   `json:"name"`
	Description    string   `json:"description,omitempty"`
	ProjectIDs     []string `json:"projectIds,omitempty"`
	IssueTypeIDs   []string `json:"issueTypeIds,omitempty"`
}

// CreateContext creates a context for a custom field.
func (s *FieldsService) CreateContext(ctx context.Context, fieldID string, contextReq *ContextCreateRequest) (*FieldContext, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/context", fieldID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, contextReq)
	if err != nil {
		return nil, nil, err
	}

	result := new(FieldContext)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateContext updates a context for a custom field.
func (s *FieldsService) UpdateContext(ctx context.Context, fieldID string, contextID int64, name, description string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/context/%d", fieldID, contextID)

	body := map[string]string{}
	if name != "" {
		body["name"] = name
	}
	if description != "" {
		body["description"] = description
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteContext removes a context for a custom field.
func (s *FieldsService) DeleteContext(ctx context.Context, fieldID string, contextID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/context/%d", fieldID, contextID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// FieldOption represents a custom field option.
type FieldOption struct {
	ID       string `json:"id,omitempty"`
	Value    string `json:"value,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

// OptionsListResult represents a paginated list of field options.
type OptionsListResult struct {
	Self       string         `json:"self,omitempty"`
	NextPage   string         `json:"nextPage,omitempty"`
	MaxResults int            `json:"maxResults,omitempty"`
	StartAt    int            `json:"startAt,omitempty"`
	Total      int            `json:"total,omitempty"`
	IsLast     bool           `json:"isLast,omitempty"`
	Values     []*FieldOption `json:"values,omitempty"`
}

// ListContextOptions returns options for a field context.
func (s *FieldsService) ListContextOptions(ctx context.Context, fieldID string, contextID int64, startAt, maxResults int, optionID []int64) (*OptionsListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/context/%d/option", fieldID, contextID)

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	for _, id := range optionID {
		params.Add("optionId", strconv.FormatInt(id, 10))
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(OptionsListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// OptionCreateRequest represents a request to create options.
type OptionCreateRequest struct {
	Options []*FieldOptionInput `json:"options"`
}

// FieldOptionInput represents input for creating/updating an option.
type FieldOptionInput struct {
	Value    string `json:"value"`
	Disabled bool   `json:"disabled,omitempty"`
	OptionID string `json:"optionId,omitempty"` // For updates
}

// CreateContextOptions creates options for a field context.
func (s *FieldsService) CreateContextOptions(ctx context.Context, fieldID string, contextID int64, options []*FieldOptionInput) (*OptionsListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/context/%d/option", fieldID, contextID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, &OptionCreateRequest{Options: options})
	if err != nil {
		return nil, nil, err
	}

	result := new(OptionsListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateContextOptions updates options for a field context.
func (s *FieldsService) UpdateContextOptions(ctx context.Context, fieldID string, contextID int64, options []*FieldOptionInput) (*OptionsListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/context/%d/option", fieldID, contextID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, &OptionCreateRequest{Options: options})
	if err != nil {
		return nil, nil, err
	}

	result := new(OptionsListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteContextOption removes an option from a field context.
func (s *FieldsService) DeleteContextOption(ctx context.Context, fieldID string, contextID, optionID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/context/%d/option/%d", fieldID, contextID, optionID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
