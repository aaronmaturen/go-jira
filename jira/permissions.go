package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// PermissionsService handles permission operations for the Jira API.
type PermissionsService struct {
	client *Client
}

// Permission represents a Jira permission.
type Permission struct {
	ID              string `json:"id,omitempty"`
	Key             string `json:"key,omitempty"`
	Name            string `json:"name,omitempty"`
	Type            string `json:"type,omitempty"`
	Description     string `json:"description,omitempty"`
	HavePermission  bool   `json:"havePermission,omitempty"`
	DeprecatedKey   bool   `json:"deprecatedKey,omitempty"`
}

// PermissionsResult represents a set of permissions.
type PermissionsResult struct {
	Permissions map[string]*Permission `json:"permissions,omitempty"`
}

// ListAll returns all permissions.
func (s *PermissionsService) ListAll(ctx context.Context) (*PermissionsResult, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/permissions", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// MyPermissionsOptions specifies options for getting my permissions.
type MyPermissionsOptions struct {
	ProjectKey  string `url:"projectKey,omitempty"`
	ProjectID   string `url:"projectId,omitempty"`
	IssueKey    string `url:"issueKey,omitempty"`
	IssueID     string `url:"issueId,omitempty"`
	Permissions string `url:"permissions,omitempty"`
	ProjectUUID string `url:"projectUuid,omitempty"`
	ProjectConfigurationUUID string `url:"projectConfigurationUuid,omitempty"`
}

// GetMyPermissions returns permissions for the current user.
func (s *PermissionsService) GetMyPermissions(ctx context.Context, opts *MyPermissionsOptions) (*PermissionsResult, *Response, error) {
	u := "/rest/api/3/mypermissions"

	if opts != nil {
		params := url.Values{}
		if opts.ProjectKey != "" {
			params.Set("projectKey", opts.ProjectKey)
		}
		if opts.ProjectID != "" {
			params.Set("projectId", opts.ProjectID)
		}
		if opts.IssueKey != "" {
			params.Set("issueKey", opts.IssueKey)
		}
		if opts.IssueID != "" {
			params.Set("issueId", opts.IssueID)
		}
		if opts.Permissions != "" {
			params.Set("permissions", opts.Permissions)
		}
		if opts.ProjectUUID != "" {
			params.Set("projectUuid", opts.ProjectUUID)
		}
		if opts.ProjectConfigurationUUID != "" {
			params.Set("projectConfigurationUuid", opts.ProjectConfigurationUUID)
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// BulkPermissionsRequest represents a request to check bulk permissions.
type BulkPermissionsRequest struct {
	ProjectPermissions []*BulkProjectPermission `json:"projectPermissions,omitempty"`
	GlobalPermissions  []string                 `json:"globalPermissions,omitempty"`
	AccountID          string                   `json:"accountId,omitempty"`
}

// BulkProjectPermission represents a project permission check.
type BulkProjectPermission struct {
	Issues      []int64  `json:"issues,omitempty"`
	Projects    []int64  `json:"projects,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// BulkPermissionsResult represents the result of bulk permission checks.
type BulkPermissionsResult struct {
	ProjectPermissions []*BulkProjectPermissionGrant `json:"projectPermissions,omitempty"`
	GlobalPermissions  []string                      `json:"globalPermissions,omitempty"`
}

// BulkProjectPermissionGrant represents granted project permissions.
type BulkProjectPermissionGrant struct {
	Permission string  `json:"permission,omitempty"`
	Issues     []int64 `json:"issues,omitempty"`
	Projects   []int64 `json:"projects,omitempty"`
}

// CheckBulk checks bulk permissions.
func (s *PermissionsService) CheckBulk(ctx context.Context, req *BulkPermissionsRequest) (*BulkPermissionsResult, *Response, error) {
	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/permissions/check", req)
	if err != nil {
		return nil, nil, err
	}

	result := new(BulkPermissionsResult)
	resp, err := s.client.Do(httpReq, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// PermittedProjectsResult represents projects where user has permission.
type PermittedProjectsResult struct {
	Projects []*Project `json:"projects,omitempty"`
}

// GetPermittedProjects returns projects where user has specified permission.
func (s *PermissionsService) GetPermittedProjects(ctx context.Context, permissions []string) (*PermittedProjectsResult, *Response, error) {
	u := "/rest/api/3/permissions/project"

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, map[string][]string{"permissions": permissions})
	if err != nil {
		return nil, nil, err
	}

	result := new(PermittedProjectsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// PermissionScheme represents a permission scheme.
type PermissionScheme struct {
	ID          int64                      `json:"id,omitempty"`
	Self        string                     `json:"self,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Description string                     `json:"description,omitempty"`
	Scope       *Scope                     `json:"scope,omitempty"`
	Permissions []*PermissionGrant         `json:"permissions,omitempty"`
	Expand      string                     `json:"expand,omitempty"`
}

// PermissionGrant represents a permission grant in a scheme.
type PermissionGrant struct {
	ID         int64                     `json:"id,omitempty"`
	Self       string                    `json:"self,omitempty"`
	Holder     *PermissionHolder         `json:"holder,omitempty"`
	Permission string                    `json:"permission,omitempty"`
}

// PermissionHolder represents a holder of a permission.
type PermissionHolder struct {
	Type      string `json:"type,omitempty"`
	Parameter string `json:"parameter,omitempty"`
	Expand    string `json:"expand,omitempty"`
	Value     string `json:"value,omitempty"`
}

// PermissionSchemeListResult represents a list of permission schemes.
type PermissionSchemeListResult struct {
	PermissionSchemes []*PermissionScheme `json:"permissionSchemes,omitempty"`
}

// ListSchemes returns all permission schemes.
func (s *PermissionsService) ListSchemes(ctx context.Context, expand []string) (*PermissionSchemeListResult, *Response, error) {
	u := "/rest/api/3/permissionscheme"

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionSchemeListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetScheme returns a permission scheme.
func (s *PermissionsService) GetScheme(ctx context.Context, schemeID int64, expand []string) (*PermissionScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/permissionscheme/%d", schemeID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	scheme := new(PermissionScheme)
	resp, err := s.client.Do(req, scheme)
	if err != nil {
		return nil, resp, err
	}

	return scheme, resp, nil
}

// PermissionSchemeCreateRequest represents a request to create a permission scheme.
type PermissionSchemeCreateRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Permissions []*PermissionGrantInput `json:"permissions,omitempty"`
	Scope       *Scope             `json:"scope,omitempty"`
}

// PermissionGrantInput represents input for creating a permission grant.
type PermissionGrantInput struct {
	Holder     *PermissionHolder `json:"holder,omitempty"`
	Permission string            `json:"permission,omitempty"`
}

// CreateScheme creates a permission scheme.
func (s *PermissionsService) CreateScheme(ctx context.Context, scheme *PermissionSchemeCreateRequest, expand []string) (*PermissionScheme, *Response, error) {
	u := "/rest/api/3/permissionscheme"

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, scheme)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateScheme updates a permission scheme.
func (s *PermissionsService) UpdateScheme(ctx context.Context, schemeID int64, scheme *PermissionSchemeCreateRequest, expand []string) (*PermissionScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/permissionscheme/%d", schemeID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, scheme)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionScheme)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteScheme removes a permission scheme.
func (s *PermissionsService) DeleteScheme(ctx context.Context, schemeID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/permissionscheme/%d", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetSchemeGrants returns grants for a permission scheme.
func (s *PermissionsService) GetSchemeGrants(ctx context.Context, schemeID int64, expand []string) ([]*PermissionGrant, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/permissionscheme/%d/permission", schemeID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Permissions []*PermissionGrant `json:"permissions"`
	}
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Permissions, resp, nil
}

// CreateSchemeGrant creates a permission grant in a scheme.
func (s *PermissionsService) CreateSchemeGrant(ctx context.Context, schemeID int64, grant *PermissionGrantInput, expand []string) (*PermissionGrant, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/permissionscheme/%d/permission", schemeID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, grant)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionGrant)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSchemeGrant returns a permission grant from a scheme.
func (s *PermissionsService) GetSchemeGrant(ctx context.Context, schemeID, grantID int64, expand []string) (*PermissionGrant, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/permissionscheme/%d/permission/%d", schemeID, grantID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	grant := new(PermissionGrant)
	resp, err := s.client.Do(req, grant)
	if err != nil {
		return nil, resp, err
	}

	return grant, resp, nil
}

// DeleteSchemeGrant removes a permission grant from a scheme.
func (s *PermissionsService) DeleteSchemeGrant(ctx context.Context, schemeID, grantID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/permissionscheme/%d/permission/%d", schemeID, grantID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ProjectPermissionScheme represents a project's permission scheme.
type ProjectPermissionScheme struct {
	ID          int64  `json:"id,omitempty"`
	Self        string `json:"self,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// GetProjectScheme returns the permission scheme for a project.
func (s *PermissionsService) GetProjectScheme(ctx context.Context, projectKeyOrID string, expand []string) (*ProjectPermissionScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/permissionscheme", projectKeyOrID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	scheme := new(ProjectPermissionScheme)
	resp, err := s.client.Do(req, scheme)
	if err != nil {
		return nil, resp, err
	}

	return scheme, resp, nil
}

// AssignProjectScheme assigns a permission scheme to a project.
func (s *PermissionsService) AssignProjectScheme(ctx context.Context, projectKeyOrID string, schemeID int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/permissionscheme", projectKeyOrID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, map[string]int64{"id": schemeID})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetSecurityLevelsForProject returns the issue security levels for a project.
func (s *PermissionsService) GetSecurityLevelsForProject(ctx context.Context, projectKeyOrID string) (*ProjectIssueSecurityLevels, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/securitylevel", projectKeyOrID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	levels := new(ProjectIssueSecurityLevels)
	resp, err := s.client.Do(req, levels)
	if err != nil {
		return nil, resp, err
	}

	return levels, resp, nil
}

// IssueSecurityScheme represents an issue security scheme.
type IssueSecurityScheme struct {
	Self              string `json:"self,omitempty"`
	ID                int64  `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	DefaultSecurityLevelID int64 `json:"defaultSecurityLevelId,omitempty"`
}

// IssueSecuritySchemeListResult represents a list of issue security schemes.
type IssueSecuritySchemeListResult struct {
	Self       string                 `json:"self,omitempty"`
	NextPage   string                 `json:"nextPage,omitempty"`
	MaxResults int                    `json:"maxResults,omitempty"`
	StartAt    int                    `json:"startAt,omitempty"`
	Total      int                    `json:"total,omitempty"`
	IsLast     bool                   `json:"isLast,omitempty"`
	Values     []*IssueSecurityScheme `json:"values,omitempty"`
}

// ListSecuritySchemes returns all issue security schemes.
func (s *PermissionsService) ListSecuritySchemes(ctx context.Context, startAt, maxResults int, ids []int64, projectID string) (*IssueSecuritySchemeListResult, *Response, error) {
	u := "/rest/api/3/issuesecurityschemes"

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
	if projectID != "" {
		params.Set("projectId", projectID)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssueSecuritySchemeListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSecurityScheme returns an issue security scheme.
func (s *PermissionsService) GetSecurityScheme(ctx context.Context, schemeID int64) (*IssueSecurityScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issuesecurityschemes/%d", schemeID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	scheme := new(IssueSecurityScheme)
	resp, err := s.client.Do(req, scheme)
	if err != nil {
		return nil, resp, err
	}

	return scheme, resp, nil
}
