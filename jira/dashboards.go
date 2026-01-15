package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// DashboardsService handles dashboard operations for the Jira API.
type DashboardsService struct {
	client *Client
}

// Dashboard represents a Jira dashboard.
type Dashboard struct {
	Self             string           `json:"self,omitempty"`
	ID               string           `json:"id,omitempty"`
	IsFavourite      bool             `json:"isFavourite,omitempty"`
	Name             string           `json:"name,omitempty"`
	Description      string           `json:"description,omitempty"`
	Owner            *User            `json:"owner,omitempty"`
	Popularity       int              `json:"popularity,omitempty"`
	Rank             int              `json:"rank,omitempty"`
	SharePermissions []*SharePermission `json:"sharePermissions,omitempty"`
	EditPermissions  []*SharePermission `json:"editPermissions,omitempty"`
	View             string           `json:"view,omitempty"`
	IsWritable       bool             `json:"isWritable,omitempty"`
	SystemDashboard  bool             `json:"systemDashboard,omitempty"`
}

// DashboardListResult represents a paginated list of dashboards.
type DashboardListResult struct {
	StartAt    int          `json:"startAt,omitempty"`
	MaxResults int          `json:"maxResults,omitempty"`
	Total      int          `json:"total,omitempty"`
	Prev       string       `json:"prev,omitempty"`
	Next       string       `json:"next,omitempty"`
	Dashboards []*Dashboard `json:"dashboards,omitempty"`
}

// ListDashboardsOptions specifies options for listing dashboards.
type ListDashboardsOptions struct {
	Filter     string `url:"filter,omitempty"`
	StartAt    int    `url:"startAt,omitempty"`
	MaxResults int    `url:"maxResults,omitempty"`
}

// List returns all dashboards.
func (s *DashboardsService) List(ctx context.Context, opts *ListDashboardsOptions) (*DashboardListResult, *Response, error) {
	u := "/rest/api/3/dashboard"

	if opts != nil {
		params := url.Values{}
		if opts.Filter != "" {
			params.Set("filter", opts.Filter)
		}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(DashboardListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchDashboardsOptions specifies options for searching dashboards.
type SearchDashboardsOptions struct {
	DashboardName            string   `url:"dashboardName,omitempty"`
	AccountID                string   `url:"accountId,omitempty"`
	Owner                    string   `url:"owner,omitempty"`
	Groupname                string   `url:"groupname,omitempty"`
	GroupID                  string   `url:"groupId,omitempty"`
	ProjectID                int64    `url:"projectId,omitempty"`
	OrderBy                  string   `url:"orderBy,omitempty"`
	StartAt                  int      `url:"startAt,omitempty"`
	MaxResults               int      `url:"maxResults,omitempty"`
	Status                   string   `url:"status,omitempty"`
	Expand                   []string `url:"expand,omitempty"`
}

// SearchDashboardsResult represents a paginated list of dashboards from search.
type SearchDashboardsResult struct {
	Self       string       `json:"self,omitempty"`
	NextPage   string       `json:"nextPage,omitempty"`
	MaxResults int          `json:"maxResults,omitempty"`
	StartAt    int          `json:"startAt,omitempty"`
	Total      int          `json:"total,omitempty"`
	IsLast     bool         `json:"isLast,omitempty"`
	Values     []*Dashboard `json:"values,omitempty"`
}

// Search searches for dashboards.
func (s *DashboardsService) Search(ctx context.Context, opts *SearchDashboardsOptions) (*SearchDashboardsResult, *Response, error) {
	u := "/rest/api/3/dashboard/search"

	if opts != nil {
		params := url.Values{}
		if opts.DashboardName != "" {
			params.Set("dashboardName", opts.DashboardName)
		}
		if opts.AccountID != "" {
			params.Set("accountId", opts.AccountID)
		}
		if opts.Owner != "" {
			params.Set("owner", opts.Owner)
		}
		if opts.Groupname != "" {
			params.Set("groupname", opts.Groupname)
		}
		if opts.GroupID != "" {
			params.Set("groupId", opts.GroupID)
		}
		if opts.ProjectID > 0 {
			params.Set("projectId", strconv.FormatInt(opts.ProjectID, 10))
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
		if opts.Status != "" {
			params.Set("status", opts.Status)
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

	result := new(SearchDashboardsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Get returns a dashboard by ID.
func (s *DashboardsService) Get(ctx context.Context, dashboardID string) (*Dashboard, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/dashboard/%s", dashboardID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	dashboard := new(Dashboard)
	resp, err := s.client.Do(req, dashboard)
	if err != nil {
		return nil, resp, err
	}

	return dashboard, resp, nil
}

// DashboardCreateRequest represents a request to create a dashboard.
type DashboardCreateRequest struct {
	Name             string           `json:"name"`
	Description      string           `json:"description,omitempty"`
	SharePermissions []*SharePermission `json:"sharePermissions,omitempty"`
	EditPermissions  []*SharePermission `json:"editPermissions,omitempty"`
}

// Create creates a new dashboard.
func (s *DashboardsService) Create(ctx context.Context, dashboard *DashboardCreateRequest) (*Dashboard, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/dashboard", dashboard)
	if err != nil {
		return nil, nil, err
	}

	result := new(Dashboard)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DashboardUpdateRequest represents a request to update a dashboard.
type DashboardUpdateRequest struct {
	Name             string           `json:"name,omitempty"`
	Description      string           `json:"description,omitempty"`
	SharePermissions []*SharePermission `json:"sharePermissions,omitempty"`
	EditPermissions  []*SharePermission `json:"editPermissions,omitempty"`
}

// Update updates a dashboard.
func (s *DashboardsService) Update(ctx context.Context, dashboardID string, dashboard *DashboardUpdateRequest) (*Dashboard, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/dashboard/%s", dashboardID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, dashboard)
	if err != nil {
		return nil, nil, err
	}

	result := new(Dashboard)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a dashboard.
func (s *DashboardsService) Delete(ctx context.Context, dashboardID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/dashboard/%s", dashboardID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Copy copies a dashboard.
func (s *DashboardsService) Copy(ctx context.Context, dashboardID string, dashboard *DashboardCreateRequest) (*Dashboard, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/dashboard/%s/copy", dashboardID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, dashboard)
	if err != nil {
		return nil, nil, err
	}

	result := new(Dashboard)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DashboardGadget represents a gadget on a dashboard.
type DashboardGadget struct {
	ID                 int64           `json:"id,omitempty"`
	ModuleKey          string          `json:"moduleKey,omitempty"`
	URI                string          `json:"uri,omitempty"`
	Color              string          `json:"color,omitempty"`
	Position           *GadgetPosition `json:"position,omitempty"`
	Title              string          `json:"title,omitempty"`
}

// GadgetPosition represents the position of a gadget.
type GadgetPosition struct {
	Row    int `json:"row,omitempty"`
	Column int `json:"column,omitempty"`
}

// GadgetListResult represents a list of gadgets.
type GadgetListResult struct {
	Gadgets []*DashboardGadget `json:"gadgets,omitempty"`
}

// ListGadgets returns all gadgets on a dashboard.
func (s *DashboardsService) ListGadgets(ctx context.Context, dashboardID string, moduleKey, uri, gadgetID string) (*GadgetListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/dashboard/%s/gadget", dashboardID)

	params := url.Values{}
	if moduleKey != "" {
		params.Set("moduleKey", moduleKey)
	}
	if uri != "" {
		params.Set("uri", uri)
	}
	if gadgetID != "" {
		params.Set("gadgetId", gadgetID)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(GadgetListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GadgetCreateRequest represents a request to add a gadget.
type GadgetCreateRequest struct {
	ModuleKey string          `json:"moduleKey,omitempty"`
	URI       string          `json:"uri,omitempty"`
	Color     string          `json:"color,omitempty"`
	Position  *GadgetPosition `json:"position,omitempty"`
	Title     string          `json:"title,omitempty"`
	IgnoreURI bool            `json:"ignoreURI,omitempty"`
}

// AddGadget adds a gadget to a dashboard.
func (s *DashboardsService) AddGadget(ctx context.Context, dashboardID string, gadget *GadgetCreateRequest) (*DashboardGadget, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/dashboard/%s/gadget", dashboardID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, gadget)
	if err != nil {
		return nil, nil, err
	}

	result := new(DashboardGadget)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GadgetUpdateRequest represents a request to update a gadget.
type GadgetUpdateRequest struct {
	Color    string          `json:"color,omitempty"`
	Position *GadgetPosition `json:"position,omitempty"`
	Title    string          `json:"title,omitempty"`
}

// UpdateGadget updates a gadget on a dashboard.
func (s *DashboardsService) UpdateGadget(ctx context.Context, dashboardID string, gadgetID int64, gadget *GadgetUpdateRequest) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/dashboard/%s/gadget/%d", dashboardID, gadgetID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, gadget)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RemoveGadget removes a gadget from a dashboard.
func (s *DashboardsService) RemoveGadget(ctx context.Context, dashboardID string, gadgetID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/dashboard/%s/gadget/%d", dashboardID, gadgetID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// AvailableGadget represents an available gadget.
type AvailableGadget struct {
	ModuleKey string `json:"moduleKey,omitempty"`
	URI       string `json:"uri,omitempty"`
	Title     string `json:"title,omitempty"`
}

// AvailableGadgetsResult represents a list of available gadgets.
type AvailableGadgetsResult struct {
	Gadgets []*AvailableGadget `json:"gadgets,omitempty"`
}

// ListAvailableGadgets returns all available gadgets.
func (s *DashboardsService) ListAvailableGadgets(ctx context.Context) (*AvailableGadgetsResult, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/dashboard/gadgets", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(AvailableGadgetsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// BulkEdit edits multiple dashboards at once.
func (s *DashboardsService) BulkEdit(ctx context.Context, action string, dashboardIDs []string, changeOwnerAccountID string, sharePermissions []*SharePermission, extendAdminPermissions bool) (*BulkEditResult, *Response, error) {
	u := "/rest/api/3/dashboard/bulk/edit"

	body := map[string]interface{}{
		"action":       action,
		"selectedDashboardIds": dashboardIDs,
	}
	if changeOwnerAccountID != "" {
		body["newOwner"] = map[string]string{"accountId": changeOwnerAccountID}
	}
	if sharePermissions != nil {
		body["sharePermissions"] = sharePermissions
	}
	if extendAdminPermissions {
		body["extendAdminPermissions"] = true
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, body)
	if err != nil {
		return nil, nil, err
	}

	result := new(BulkEditResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// BulkEditResult represents the result of a bulk edit operation.
type BulkEditResult struct {
	SuccessfulDashboardIDs []string `json:"modifiedDashboards,omitempty"`
	FailedDashboardIDs     []string `json:"notModifiedDashboards,omitempty"`
}
