package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ScreensService handles screen operations for the Jira API.
type ScreensService struct {
	client *Client
}

// Screen represents a Jira screen.
type Screen struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Scope       *Scope `json:"scope,omitempty"`
}

// ScreenListResult represents a paginated list of screens.
type ScreenListResult struct {
	Self       string    `json:"self,omitempty"`
	NextPage   string    `json:"nextPage,omitempty"`
	MaxResults int       `json:"maxResults,omitempty"`
	StartAt    int       `json:"startAt,omitempty"`
	Total      int       `json:"total,omitempty"`
	IsLast     bool      `json:"isLast,omitempty"`
	Values     []*Screen `json:"values,omitempty"`
}

// ListOptions specifies options for listing screens.
type ScreenListOptions struct {
	StartAt    int      `url:"startAt,omitempty"`
	MaxResults int      `url:"maxResults,omitempty"`
	IDs        []int64  `url:"id,omitempty"`
	QueryString string  `url:"queryString,omitempty"`
	Scope      []string `url:"scope,omitempty"`
	OrderBy    string   `url:"orderBy,omitempty"`
}

// List returns screens with pagination.
func (s *ScreensService) List(ctx context.Context, opts *ScreenListOptions) (*ScreenListResult, *Response, error) {
	u := "/rest/api/3/screens"

	if opts != nil {
		params := url.Values{}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		for _, id := range opts.IDs {
			params.Add("id", strconv.FormatInt(id, 10))
		}
		if opts.QueryString != "" {
			params.Set("queryString", opts.QueryString)
		}
		for _, scope := range opts.Scope {
			params.Add("scope", scope)
		}
		if opts.OrderBy != "" {
			params.Set("orderBy", opts.OrderBy)
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(ScreenListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ScreenCreateRequest represents a request to create a screen.
type ScreenCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Create creates a screen.
func (s *ScreensService) Create(ctx context.Context, screen *ScreenCreateRequest) (*Screen, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/screens", screen)
	if err != nil {
		return nil, nil, err
	}

	result := new(Screen)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ScreenUpdateRequest represents a request to update a screen.
type ScreenUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Update updates a screen.
func (s *ScreensService) Update(ctx context.Context, screenID int64, screen *ScreenUpdateRequest) (*Screen, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d", screenID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, screen)
	if err != nil {
		return nil, nil, err
	}

	result := new(Screen)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a screen.
func (s *ScreensService) Delete(ctx context.Context, screenID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d", screenID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ScreenTab represents a tab on a screen.
type ScreenTab struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ListTabs returns tabs for a screen.
func (s *ScreensService) ListTabs(ctx context.Context, screenID int64, projectKey string) ([]*ScreenTab, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d/tabs", screenID)

	if projectKey != "" {
		u = fmt.Sprintf("%s?projectKey=%s", u, projectKey)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tabs []*ScreenTab
	resp, err := s.client.Do(req, &tabs)
	if err != nil {
		return nil, resp, err
	}

	return tabs, resp, nil
}

// CreateTab creates a tab on a screen.
func (s *ScreensService) CreateTab(ctx context.Context, screenID int64, name string) (*ScreenTab, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d/tabs", screenID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, map[string]string{"name": name})
	if err != nil {
		return nil, nil, err
	}

	tab := new(ScreenTab)
	resp, err := s.client.Do(req, tab)
	if err != nil {
		return nil, resp, err
	}

	return tab, resp, nil
}

// UpdateTab updates a tab on a screen.
func (s *ScreensService) UpdateTab(ctx context.Context, screenID, tabID int64, name string) (*ScreenTab, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d", screenID, tabID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, map[string]string{"name": name})
	if err != nil {
		return nil, nil, err
	}

	tab := new(ScreenTab)
	resp, err := s.client.Do(req, tab)
	if err != nil {
		return nil, resp, err
	}

	return tab, resp, nil
}

// DeleteTab removes a tab from a screen.
func (s *ScreensService) DeleteTab(ctx context.Context, screenID, tabID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d", screenID, tabID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// MoveTab moves a tab on a screen.
func (s *ScreensService) MoveTab(ctx context.Context, screenID, tabID int64, pos int) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d/move/%d", screenID, tabID, pos)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ScreenTabField represents a field on a screen tab.
type ScreenTabField struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ListTabFields returns fields on a screen tab.
func (s *ScreensService) ListTabFields(ctx context.Context, screenID, tabID int64, projectKey string) ([]*ScreenTabField, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d/fields", screenID, tabID)

	if projectKey != "" {
		u = fmt.Sprintf("%s?projectKey=%s", u, projectKey)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var fields []*ScreenTabField
	resp, err := s.client.Do(req, &fields)
	if err != nil {
		return nil, resp, err
	}

	return fields, resp, nil
}

// AddTabField adds a field to a screen tab.
func (s *ScreensService) AddTabField(ctx context.Context, screenID, tabID int64, fieldID string) (*ScreenTabField, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d/fields", screenID, tabID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, map[string]string{"fieldId": fieldID})
	if err != nil {
		return nil, nil, err
	}

	field := new(ScreenTabField)
	resp, err := s.client.Do(req, field)
	if err != nil {
		return nil, resp, err
	}

	return field, resp, nil
}

// RemoveTabField removes a field from a screen tab.
func (s *ScreensService) RemoveTabField(ctx context.Context, screenID, tabID int64, fieldID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d/fields/%s", screenID, tabID, fieldID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// MoveTabField moves a field on a screen tab.
func (s *ScreensService) MoveTabField(ctx context.Context, screenID, tabID int64, fieldID string, after string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/screens/%d/tabs/%d/fields/%s/move", screenID, tabID, fieldID)

	body := map[string]string{}
	if after != "" {
		body["after"] = after
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ScreenScheme represents a screen scheme.
type ScreenScheme struct {
	ID          int64              `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Screens     *ScreenSchemeScreens `json:"screens,omitempty"`
}

// ScreenSchemeScreens represents screen mappings in a scheme.
type ScreenSchemeScreens struct {
	Default int64 `json:"default,omitempty"`
	View    int64 `json:"view,omitempty"`
	Edit    int64 `json:"edit,omitempty"`
	Create  int64 `json:"create,omitempty"`
}

// ScreenSchemeListResult represents a paginated list of screen schemes.
type ScreenSchemeListResult struct {
	Self       string          `json:"self,omitempty"`
	NextPage   string          `json:"nextPage,omitempty"`
	MaxResults int             `json:"maxResults,omitempty"`
	StartAt    int             `json:"startAt,omitempty"`
	Total      int             `json:"total,omitempty"`
	IsLast     bool            `json:"isLast,omitempty"`
	Values     []*ScreenScheme `json:"values,omitempty"`
}

// ListSchemes returns screen schemes.
func (s *ScreensService) ListSchemes(ctx context.Context, startAt, maxResults int, ids []int64, expand string, queryString string, orderBy string) (*ScreenSchemeListResult, *Response, error) {
	u := "/rest/api/3/screenscheme"

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
	if orderBy != "" {
		params.Set("orderBy", orderBy)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(ScreenSchemeListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ScreenSchemeCreateRequest represents a request to create a screen scheme.
type ScreenSchemeCreateRequest struct {
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	Screens     *ScreenSchemeScreens `json:"screens"`
}

// CreateScheme creates a screen scheme.
func (s *ScreensService) CreateScheme(ctx context.Context, scheme *ScreenSchemeCreateRequest) (*ScreenScheme, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/screenscheme", scheme)
	if err != nil {
		return nil, nil, err
	}

	result := new(ScreenScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateScheme updates a screen scheme.
func (s *ScreensService) UpdateScheme(ctx context.Context, schemeID int64, name, description string, screens *ScreenSchemeScreens) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/screenscheme/%d", schemeID)

	body := map[string]interface{}{}
	if name != "" {
		body["name"] = name
	}
	if description != "" {
		body["description"] = description
	}
	if screens != nil {
		body["screens"] = screens
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteScheme removes a screen scheme.
func (s *ScreensService) DeleteScheme(ctx context.Context, schemeID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/screenscheme/%d", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// FieldScreensResult represents screens for a field.
type FieldScreensResult struct {
	Self       string         `json:"self,omitempty"`
	NextPage   string         `json:"nextPage,omitempty"`
	MaxResults int            `json:"maxResults,omitempty"`
	StartAt    int            `json:"startAt,omitempty"`
	Total      int            `json:"total,omitempty"`
	IsLast     bool           `json:"isLast,omitempty"`
	Values     []*FieldScreen `json:"values,omitempty"`
}

// FieldScreen represents a screen containing a field.
type FieldScreen struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetFieldScreens returns screens for a field.
func (s *ScreensService) GetFieldScreens(ctx context.Context, fieldID string, startAt, maxResults int, expand []string) (*FieldScreensResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/field/%s/screens", fieldID)

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	if len(expand) > 0 {
		params.Set("expand", strings.Join(expand, ","))
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(FieldScreensResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
