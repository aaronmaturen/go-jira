package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ProjectRolesService handles project role operations for the Jira API.
type ProjectRolesService struct {
	client *Client
}

// ListAll returns all project roles.
func (s *ProjectRolesService) ListAll(ctx context.Context) ([]*ProjectRole, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/role", nil)
	if err != nil {
		return nil, nil, err
	}

	var roles []*ProjectRole
	resp, err := s.client.Do(req, &roles)
	if err != nil {
		return nil, resp, err
	}

	return roles, resp, nil
}

// Get returns a project role by ID.
func (s *ProjectRolesService) Get(ctx context.Context, roleID int64) (*ProjectRole, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/role/%d", roleID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	role := new(ProjectRole)
	resp, err := s.client.Do(req, role)
	if err != nil {
		return nil, resp, err
	}

	return role, resp, nil
}

// ProjectRoleCreateRequest represents a request to create a project role.
type ProjectRoleCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Create creates a project role.
func (s *ProjectRolesService) Create(ctx context.Context, role *ProjectRoleCreateRequest) (*ProjectRole, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/role", role)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectRole)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ProjectRoleUpdateRequest represents a request to fully update a project role.
type ProjectRoleUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Update fully updates a project role.
func (s *ProjectRolesService) Update(ctx context.Context, roleID int64, role *ProjectRoleUpdateRequest) (*ProjectRole, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/role/%d", roleID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, role)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectRole)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// PartialUpdate partially updates a project role.
func (s *ProjectRolesService) PartialUpdate(ctx context.Context, roleID int64, name, description string) (*ProjectRole, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/role/%d", roleID)

	body := map[string]string{}
	if name != "" {
		body["name"] = name
	}
	if description != "" {
		body["description"] = description
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectRole)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a project role.
func (s *ProjectRolesService) Delete(ctx context.Context, roleID int64, swap int64) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/role/%d", roleID)

	if swap > 0 {
		u = fmt.Sprintf("%s?swap=%d", u, swap)
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ListForProject returns project roles for a project.
func (s *ProjectRolesService) ListForProject(ctx context.Context, projectIDOrKey string) (map[string]string, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/role", projectIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result map[string]string
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetForProject returns a project role for a project.
func (s *ProjectRolesService) GetForProject(ctx context.Context, projectIDOrKey string, roleID int64, excludeInactiveUsers bool) (*ProjectRole, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/role/%d", projectIDOrKey, roleID)

	if excludeInactiveUsers {
		u = fmt.Sprintf("%s?excludeInactiveUsers=true", u)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	role := new(ProjectRole)
	resp, err := s.client.Do(req, role)
	if err != nil {
		return nil, resp, err
	}

	return role, resp, nil
}

// GetRoleDetails returns project role details for a project.
func (s *ProjectRolesService) GetRoleDetails(ctx context.Context, projectIDOrKey string, currentMember, excludeConnectAddons bool, roleIDs []int64) ([]*ProjectRole, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/roledetails", projectIDOrKey)

	params := url.Values{}
	if currentMember {
		params.Set("currentMember", "true")
	}
	if excludeConnectAddons {
		params.Set("excludeConnectAddons", "true")
	}
	for _, id := range roleIDs {
		params.Add("id", strconv.FormatInt(id, 10))
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var roles []*ProjectRole
	resp, err := s.client.Do(req, &roles)
	if err != nil {
		return nil, resp, err
	}

	return roles, resp, nil
}

// ActorRequest represents actors to add/set for a role.
type ActorRequest struct {
	User  []string `json:"user,omitempty"`
	Group []string `json:"group,omitempty"`
}

// SetActors sets actors for a project role.
func (s *ProjectRolesService) SetActors(ctx context.Context, projectIDOrKey string, roleID int64, actors *ActorRequest) (*ProjectRole, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/role/%d", projectIDOrKey, roleID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, map[string]interface{}{
		"categorisedActors": actors,
	})
	if err != nil {
		return nil, nil, err
	}

	role := new(ProjectRole)
	resp, err := s.client.Do(req, role)
	if err != nil {
		return nil, resp, err
	}

	return role, resp, nil
}

// AddActors adds actors to a project role.
func (s *ProjectRolesService) AddActors(ctx context.Context, projectIDOrKey string, roleID int64, actors *ActorRequest) (*ProjectRole, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/role/%d", projectIDOrKey, roleID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, actors)
	if err != nil {
		return nil, nil, err
	}

	role := new(ProjectRole)
	resp, err := s.client.Do(req, role)
	if err != nil {
		return nil, resp, err
	}

	return role, resp, nil
}

// RemoveActor removes an actor from a project role.
func (s *ProjectRolesService) RemoveActor(ctx context.Context, projectIDOrKey string, roleID int64, user, group string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/role/%d", projectIDOrKey, roleID)

	params := url.Values{}
	if user != "" {
		params.Set("user", user)
	}
	if group != "" {
		params.Set("group", group)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetDefaultActors returns default actors for a project role.
func (s *ProjectRolesService) GetDefaultActors(ctx context.Context, roleID int64) (*ProjectRole, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/role/%d/actors", roleID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	role := new(ProjectRole)
	resp, err := s.client.Do(req, role)
	if err != nil {
		return nil, resp, err
	}

	return role, resp, nil
}

// AddDefaultActors adds default actors to a project role.
func (s *ProjectRolesService) AddDefaultActors(ctx context.Context, roleID int64, users, groups []string) (*ProjectRole, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/role/%d/actors", roleID)

	body := map[string][]string{}
	if len(users) > 0 {
		body["user"] = users
	}
	if len(groups) > 0 {
		body["group"] = groups
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}

	role := new(ProjectRole)
	resp, err := s.client.Do(req, role)
	if err != nil {
		return nil, resp, err
	}

	return role, resp, nil
}

// RemoveDefaultActor removes a default actor from a project role.
func (s *ProjectRolesService) RemoveDefaultActor(ctx context.Context, roleID int64, user, group string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/role/%d/actors", roleID)

	params := url.Values{}
	if user != "" {
		params.Set("user", user)
	}
	if group != "" {
		params.Set("group", group)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
