package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// FiltersService handles filter operations for the Jira API.
type FiltersService struct {
	client *Client
}

// Filter represents a Jira filter.
type Filter struct {
	Self             string           `json:"self,omitempty"`
	ID               string           `json:"id,omitempty"`
	Name             string           `json:"name,omitempty"`
	Description      string           `json:"description,omitempty"`
	Owner            *User            `json:"owner,omitempty"`
	JQL              string           `json:"jql,omitempty"`
	ViewURL          string           `json:"viewUrl,omitempty"`
	SearchURL        string           `json:"searchUrl,omitempty"`
	Favourite        bool             `json:"favourite,omitempty"`
	FavouritedCount  int64            `json:"favouritedCount,omitempty"`
	SharePermissions []*SharePermission `json:"sharePermissions,omitempty"`
	EditPermissions  []*SharePermission `json:"editPermissions,omitempty"`
	Subscriptions    []*FilterSubscription `json:"subscriptions,omitempty"`
	Expand           string           `json:"expand,omitempty"`
}

// SharePermission represents a filter share permission.
type SharePermission struct {
	ID      int64    `json:"id,omitempty"`
	Type    string   `json:"type,omitempty"`
	Project *Project `json:"project,omitempty"`
	Role    *ProjectRole `json:"role,omitempty"`
	Group   *Group   `json:"group,omitempty"`
	User    *User    `json:"user,omitempty"`
}

// FilterSubscription represents a filter subscription.
type FilterSubscription struct {
	ID    int64  `json:"id,omitempty"`
	User  *User  `json:"user,omitempty"`
	Group *Group `json:"group,omitempty"`
}

// FilterCreateRequest represents a request to create a filter.
type FilterCreateRequest struct {
	Name             string           `json:"name"`
	Description      string           `json:"description,omitempty"`
	JQL              string           `json:"jql,omitempty"`
	Favourite        bool             `json:"favourite,omitempty"`
	SharePermissions []*SharePermission `json:"sharePermissions,omitempty"`
	EditPermissions  []*SharePermission `json:"editPermissions,omitempty"`
}

// FilterUpdateRequest represents a request to update a filter.
type FilterUpdateRequest struct {
	Name             string           `json:"name,omitempty"`
	Description      string           `json:"description,omitempty"`
	JQL              string           `json:"jql,omitempty"`
	Favourite        bool             `json:"favourite,omitempty"`
	SharePermissions []*SharePermission `json:"sharePermissions,omitempty"`
	EditPermissions  []*SharePermission `json:"editPermissions,omitempty"`
}

// Create creates a new filter.
func (s *FiltersService) Create(ctx context.Context, filter *FilterCreateRequest, expand []string, overrideSharePermissions bool) (*Filter, *Response, error) {
	u := "/rest/api/3/filter"

	params := url.Values{}
	if len(expand) > 0 {
		params.Set("expand", strings.Join(expand, ","))
	}
	if overrideSharePermissions {
		params.Set("overrideSharePermissions", "true")
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, filter)
	if err != nil {
		return nil, nil, err
	}

	result := new(Filter)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetOptions specifies options for getting a filter.
type FilterGetOptions struct {
	Expand []string `url:"expand,omitempty"`
	OverrideSharePermissions bool `url:"overrideSharePermissions,omitempty"`
}

// Get returns a filter by ID.
func (s *FiltersService) Get(ctx context.Context, filterID int64, opts *FilterGetOptions) (*Filter, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/filter/%d", filterID)

	if opts != nil {
		params := url.Values{}
		if len(opts.Expand) > 0 {
			params.Set("expand", strings.Join(opts.Expand, ","))
		}
		if opts.OverrideSharePermissions {
			params.Set("overrideSharePermissions", "true")
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	filter := new(Filter)
	resp, err := s.client.Do(req, filter)
	if err != nil {
		return nil, resp, err
	}

	return filter, resp, nil
}

// Update updates a filter.
func (s *FiltersService) Update(ctx context.Context, filterID int64, filter *FilterUpdateRequest, expand []string, overrideSharePermissions bool) (*Filter, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/filter/%d", filterID)

	params := url.Values{}
	if len(expand) > 0 {
		params.Set("expand", strings.Join(expand, ","))
	}
	if overrideSharePermissions {
		params.Set("overrideSharePermissions", "true")
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, filter)
	if err != nil {
		return nil, nil, err
	}

	result := new(Filter)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a filter.
func (s *FiltersService) Delete(ctx context.Context, filterID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/filter/%d", filterID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ListMyFiltersOptions specifies options for listing my filters.
type ListMyFiltersOptions struct {
	Expand             []string `url:"expand,omitempty"`
	IncludeFavourites  bool     `url:"includeFavourites,omitempty"`
}

// ListMy returns filters owned by the current user.
func (s *FiltersService) ListMy(ctx context.Context, opts *ListMyFiltersOptions) ([]*Filter, *Response, error) {
	u := "/rest/api/3/filter/my"

	if opts != nil {
		params := url.Values{}
		if len(opts.Expand) > 0 {
			params.Set("expand", strings.Join(opts.Expand, ","))
		}
		if opts.IncludeFavourites {
			params.Set("includeFavourites", "true")
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var filters []*Filter
	resp, err := s.client.Do(req, &filters)
	if err != nil {
		return nil, resp, err
	}

	return filters, resp, nil
}

// SearchFiltersOptions specifies options for searching filters.
type SearchFiltersOptions struct {
	FilterName    string   `url:"filterName,omitempty"`
	AccountID     string   `url:"accountId,omitempty"`
	Owner         string   `url:"owner,omitempty"`
	GroupName     string   `url:"groupname,omitempty"`
	GroupID       string   `url:"groupId,omitempty"`
	ProjectID     int64    `url:"projectId,omitempty"`
	IDs           []int64  `url:"id,omitempty"`
	OrderBy       string   `url:"orderBy,omitempty"`
	StartAt       int      `url:"startAt,omitempty"`
	MaxResults    int      `url:"maxResults,omitempty"`
	Expand        []string `url:"expand,omitempty"`
	OverrideSharePermissions bool `url:"overrideSharePermissions,omitempty"`
}

// SearchFiltersResult represents a paginated list of filters.
type SearchFiltersResult struct {
	Self       string    `json:"self,omitempty"`
	NextPage   string    `json:"nextPage,omitempty"`
	MaxResults int       `json:"maxResults,omitempty"`
	StartAt    int       `json:"startAt,omitempty"`
	Total      int       `json:"total,omitempty"`
	IsLast     bool      `json:"isLast,omitempty"`
	Values     []*Filter `json:"values,omitempty"`
}

// Search searches for filters.
func (s *FiltersService) Search(ctx context.Context, opts *SearchFiltersOptions) (*SearchFiltersResult, *Response, error) {
	u := "/rest/api/3/filter/search"

	if opts != nil {
		params := url.Values{}
		if opts.FilterName != "" {
			params.Set("filterName", opts.FilterName)
		}
		if opts.AccountID != "" {
			params.Set("accountId", opts.AccountID)
		}
		if opts.Owner != "" {
			params.Set("owner", opts.Owner)
		}
		if opts.GroupName != "" {
			params.Set("groupname", opts.GroupName)
		}
		if opts.GroupID != "" {
			params.Set("groupId", opts.GroupID)
		}
		if opts.ProjectID > 0 {
			params.Set("projectId", strconv.FormatInt(opts.ProjectID, 10))
		}
		for _, id := range opts.IDs {
			params.Add("id", strconv.FormatInt(id, 10))
		}
		if opts.OrderBy != "" {
			params.Set("orderBy", opts.OrderBy)
		}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		if len(opts.Expand) > 0 {
			params.Set("expand", strings.Join(opts.Expand, ","))
		}
		if opts.OverrideSharePermissions {
			params.Set("overrideSharePermissions", "true")
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SearchFiltersResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListFavourite returns favourite filters.
func (s *FiltersService) ListFavourite(ctx context.Context, expand []string) ([]*Filter, *Response, error) {
	u := "/rest/api/3/filter/favourite"

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var filters []*Filter
	resp, err := s.client.Do(req, &filters)
	if err != nil {
		return nil, resp, err
	}

	return filters, resp, nil
}

// SetFavourite adds a filter to favourites.
func (s *FiltersService) SetFavourite(ctx context.Context, filterID int64, expand []string) (*Filter, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/filter/%d/favourite", filterID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, nil)
	if err != nil {
		return nil, nil, err
	}

	filter := new(Filter)
	resp, err := s.client.Do(req, filter)
	if err != nil {
		return nil, resp, err
	}

	return filter, resp, nil
}

// RemoveFavourite removes a filter from favourites.
func (s *FiltersService) RemoveFavourite(ctx context.Context, filterID int64, expand []string) (*Filter, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/filter/%d/favourite", filterID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, nil, err
	}

	filter := new(Filter)
	resp, err := s.client.Do(req, filter)
	if err != nil {
		return nil, resp, err
	}

	return filter, resp, nil
}

// GetDefaultShareScope returns the default share scope.
func (s *FiltersService) GetDefaultShareScope(ctx context.Context) (*DefaultShareScope, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/filter/defaultShareScope", nil)
	if err != nil {
		return nil, nil, err
	}

	scope := new(DefaultShareScope)
	resp, err := s.client.Do(req, scope)
	if err != nil {
		return nil, resp, err
	}

	return scope, resp, nil
}

// DefaultShareScope represents the default share scope.
type DefaultShareScope struct {
	Scope string `json:"scope,omitempty"`
}

// SetDefaultShareScope sets the default share scope.
func (s *FiltersService) SetDefaultShareScope(ctx context.Context, scope string) (*DefaultShareScope, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, "/rest/api/3/filter/defaultShareScope", &DefaultShareScope{Scope: scope})
	if err != nil {
		return nil, nil, err
	}

	result := new(DefaultShareScope)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSharePermissions returns share permissions for a filter.
func (s *FiltersService) GetSharePermissions(ctx context.Context, filterID int64) ([]*SharePermission, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/filter/%d/permission", filterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var permissions []*SharePermission
	resp, err := s.client.Do(req, &permissions)
	if err != nil {
		return nil, resp, err
	}

	return permissions, resp, nil
}

// SharePermissionRequest represents a request to add a share permission.
type SharePermissionRequest struct {
	Type           string `json:"type"`
	ProjectID      string `json:"projectId,omitempty"`
	GroupName      string `json:"groupname,omitempty"`
	GroupID        string `json:"groupId,omitempty"`
	ProjectRoleID  string `json:"projectRoleId,omitempty"`
	AccountID      string `json:"accountId,omitempty"`
	Rights         int    `json:"rights,omitempty"`
}

// AddSharePermission adds a share permission to a filter.
func (s *FiltersService) AddSharePermission(ctx context.Context, filterID int64, permission *SharePermissionRequest) ([]*SharePermission, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/filter/%d/permission", filterID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, permission)
	if err != nil {
		return nil, nil, err
	}

	var permissions []*SharePermission
	resp, err := s.client.Do(req, &permissions)
	if err != nil {
		return nil, resp, err
	}

	return permissions, resp, nil
}

// GetSharePermission returns a specific share permission for a filter.
func (s *FiltersService) GetSharePermission(ctx context.Context, filterID, permissionID int64) (*SharePermission, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/filter/%d/permission/%d", filterID, permissionID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	permission := new(SharePermission)
	resp, err := s.client.Do(req, permission)
	if err != nil {
		return nil, resp, err
	}

	return permission, resp, nil
}

// DeleteSharePermission removes a share permission from a filter.
func (s *FiltersService) DeleteSharePermission(ctx context.Context, filterID, permissionID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/filter/%d/permission/%d", filterID, permissionID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ChangeOwner changes the owner of a filter.
func (s *FiltersService) ChangeOwner(ctx context.Context, filterID int64, accountID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/filter/%d/owner", filterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, map[string]string{"accountId": accountID})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
