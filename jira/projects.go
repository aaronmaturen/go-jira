package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ProjectsService handles project operations for the Jira API.
type ProjectsService struct {
	client *Client
}

// ProjectListOptions specifies optional parameters for listing projects.
type ProjectListOptions struct {
	// StartAt index of the first item to return.
	StartAt int `url:"startAt,omitempty"`

	// MaxResults maximum number of items to return.
	MaxResults int `url:"maxResults,omitempty"`

	// OrderBy field to order by.
	OrderBy string `url:"orderBy,omitempty"`

	// ID filter by project IDs.
	ID []int64 `url:"id,omitempty"`

	// Keys filter by project keys.
	Keys []string `url:"keys,omitempty"`

	// Query search query.
	Query string `url:"query,omitempty"`

	// TypeKey filter by project type key.
	TypeKey string `url:"typeKey,omitempty"`

	// CategoryID filter by project category ID.
	CategoryID int64 `url:"categoryId,omitempty"`

	// Action filter by projects user has permission to perform action on.
	Action string `url:"action,omitempty"`

	// Expand additional fields to include.
	Expand []string `url:"expand,omitempty"`

	// Status filter by project status (live, archived, deleted).
	Status []string `url:"status,omitempty"`

	// Properties filter by project properties.
	Properties []string `url:"properties,omitempty"`

	// PropertyQuery JQL-like query for properties.
	PropertyQuery string `url:"propertyQuery,omitempty"`
}

// ProjectListResult represents a paginated list of projects.
type ProjectListResult struct {
	Self       string     `json:"self,omitempty"`
	NextPage   string     `json:"nextPage,omitempty"`
	MaxResults int        `json:"maxResults,omitempty"`
	StartAt    int        `json:"startAt,omitempty"`
	Total      int        `json:"total,omitempty"`
	IsLast     bool       `json:"isLast,omitempty"`
	Values     []*Project `json:"values,omitempty"`
}

// List returns a paginated list of projects.
func (s *ProjectsService) List(ctx context.Context, opts *ProjectListOptions) (*ProjectListResult, *Response, error) {
	u := "/rest/api/3/project/search"

	if opts != nil {
		params := url.Values{}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		if opts.OrderBy != "" {
			params.Set("orderBy", opts.OrderBy)
		}
		if opts.Query != "" {
			params.Set("query", opts.Query)
		}
		if opts.TypeKey != "" {
			params.Set("typeKey", opts.TypeKey)
		}
		if opts.CategoryID > 0 {
			params.Set("categoryId", strconv.FormatInt(opts.CategoryID, 10))
		}
		if opts.Action != "" {
			params.Set("action", opts.Action)
		}
		if len(opts.Expand) > 0 {
			params.Set("expand", strings.Join(opts.Expand, ","))
		}
		if len(opts.Status) > 0 {
			for _, st := range opts.Status {
				params.Add("status", st)
			}
		}
		if len(opts.Keys) > 0 {
			for _, k := range opts.Keys {
				params.Add("keys", k)
			}
		}
		if len(opts.ID) > 0 {
			for _, id := range opts.ID {
				params.Add("id", strconv.FormatInt(id, 10))
			}
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
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

// GetProjectOptions specifies options for getting a project.
type GetProjectOptions struct {
	Expand     []string `url:"expand,omitempty"`
	Properties []string `url:"properties,omitempty"`
}

// Get returns a project by its ID or key.
func (s *ProjectsService) Get(ctx context.Context, projectIDOrKey string, opts *GetProjectOptions) (*Project, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s", projectIDOrKey)

	if opts != nil {
		params := url.Values{}
		if len(opts.Expand) > 0 {
			params.Set("expand", strings.Join(opts.Expand, ","))
		}
		if len(opts.Properties) > 0 {
			for _, p := range opts.Properties {
				params.Add("properties", p)
			}
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, nil
}

// ProjectCreateRequest represents a request to create a project.
type ProjectCreateRequest struct {
	Key                      string `json:"key"`
	Name                     string `json:"name"`
	Description              string `json:"description,omitempty"`
	Lead                     string `json:"lead,omitempty"`
	LeadAccountID            string `json:"leadAccountId,omitempty"`
	URL                      string `json:"url,omitempty"`
	AssigneeType             string `json:"assigneeType,omitempty"`
	AvatarID                 int64  `json:"avatarId,omitempty"`
	IssueSecurityScheme      int64  `json:"issueSecurityScheme,omitempty"`
	PermissionScheme         int64  `json:"permissionScheme,omitempty"`
	NotificationScheme       int64  `json:"notificationScheme,omitempty"`
	CategoryID               int64  `json:"categoryId,omitempty"`
	ProjectTypeKey           string `json:"projectTypeKey,omitempty"`
	ProjectTemplateKey       string `json:"projectTemplateKey,omitempty"`
	WorkflowScheme           int64  `json:"workflowScheme,omitempty"`
	IssueTypeScreenScheme    int64  `json:"issueTypeScreenScheme,omitempty"`
	IssueTypeScheme          int64  `json:"issueTypeScheme,omitempty"`
	FieldConfigurationScheme int64  `json:"fieldConfigurationScheme,omitempty"`
}

// ProjectCreateResponse represents the response from creating a project.
type ProjectCreateResponse struct {
	Self string `json:"self,omitempty"`
	ID   int64  `json:"id,omitempty"`
	Key  string `json:"key,omitempty"`
}

// Create creates a new project.
func (s *ProjectsService) Create(ctx context.Context, project *ProjectCreateRequest) (*ProjectCreateResponse, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/project", project)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectCreateResponse)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ProjectUpdateRequest represents a request to update a project.
type ProjectUpdateRequest struct {
	Key                 string `json:"key,omitempty"`
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	Lead                string `json:"lead,omitempty"`
	LeadAccountID       string `json:"leadAccountId,omitempty"`
	URL                 string `json:"url,omitempty"`
	AssigneeType        string `json:"assigneeType,omitempty"`
	AvatarID            int64  `json:"avatarId,omitempty"`
	IssueSecurityScheme int64  `json:"issueSecurityScheme,omitempty"`
	PermissionScheme    int64  `json:"permissionScheme,omitempty"`
	NotificationScheme  int64  `json:"notificationScheme,omitempty"`
	CategoryID          int64  `json:"categoryId,omitempty"`
}

// Update updates a project.
func (s *ProjectsService) Update(ctx context.Context, projectIDOrKey string, project *ProjectUpdateRequest) (*Project, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s", projectIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, project)
	if err != nil {
		return nil, nil, err
	}

	result := new(Project)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete deletes a project.
func (s *ProjectsService) Delete(ctx context.Context, projectIDOrKey string, enableUndo bool) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s", projectIDOrKey)

	if enableUndo {
		u = fmt.Sprintf("%s?enableUndo=true", u)
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Archive archives a project.
func (s *ProjectsService) Archive(ctx context.Context, projectIDOrKey string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/archive", projectIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Restore restores an archived or deleted project.
func (s *ProjectsService) Restore(ctx context.Context, projectIDOrKey string) (*Project, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/restore", projectIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, nil
}

// GetStatuses returns all statuses for a project.
func (s *ProjectsService) GetStatuses(ctx context.Context, projectIDOrKey string) ([]*IssueTypeWithStatuses, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/statuses", projectIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var statuses []*IssueTypeWithStatuses
	resp, err := s.client.Do(req, &statuses)
	if err != nil {
		return nil, resp, err
	}

	return statuses, resp, nil
}

// IssueTypeWithStatuses represents an issue type with its associated statuses.
type IssueTypeWithStatuses struct {
	Self     string    `json:"self,omitempty"`
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Subtask  bool      `json:"subtask,omitempty"`
	Statuses []*Status `json:"statuses,omitempty"`
}

// GetHierarchy returns the issue type hierarchy for a project.
func (s *ProjectsService) GetHierarchy(ctx context.Context, projectID int64) (*ProjectIssueTypeHierarchy, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%d/hierarchy", projectID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	hierarchy := new(ProjectIssueTypeHierarchy)
	resp, err := s.client.Do(req, hierarchy)
	if err != nil {
		return nil, resp, err
	}

	return hierarchy, resp, nil
}

// ProjectIssueTypeHierarchy represents the issue type hierarchy for a project.
type ProjectIssueTypeHierarchy struct {
	ProjectID int64             `json:"projectId,omitempty"`
	Hierarchy []*HierarchyLevel `json:"hierarchy,omitempty"`
}

// HierarchyLevel represents a level in the issue type hierarchy.
type HierarchyLevel struct {
	EntityID   string       `json:"entityId,omitempty"`
	Level      int          `json:"level,omitempty"`
	Name       string       `json:"name,omitempty"`
	IssueTypes []*IssueType `json:"issueTypes,omitempty"`
}

// GetNotificationScheme returns the notification scheme for a project.
func (s *ProjectsService) GetNotificationScheme(ctx context.Context, projectKeyOrID string, expand []string) (*NotificationScheme, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/notificationscheme", projectKeyOrID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	scheme := new(NotificationScheme)
	resp, err := s.client.Do(req, scheme)
	if err != nil {
		return nil, resp, err
	}

	return scheme, resp, nil
}

// NotificationScheme represents a notification scheme.
type NotificationScheme struct {
	Self                     string                     `json:"self,omitempty"`
	ID                       int64                      `json:"id,omitempty"`
	Name                     string                     `json:"name,omitempty"`
	Description              string                     `json:"description,omitempty"`
	NotificationSchemeEvents []*NotificationSchemeEvent `json:"notificationSchemeEvents,omitempty"`
	Expand                   string                     `json:"expand,omitempty"`
}

// NotificationSchemeEvent represents an event in a notification scheme.
type NotificationSchemeEvent struct {
	Event         *NotificationEvent   `json:"event,omitempty"`
	Notifications []*EventNotification `json:"notifications,omitempty"`
}

// NotificationEvent represents a notification event.
type NotificationEvent struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// EventNotification represents a notification for an event.
type EventNotification struct {
	ID               int64        `json:"id,omitempty"`
	NotificationType string       `json:"notificationType,omitempty"`
	Parameter        string       `json:"parameter,omitempty"`
	Expand           string       `json:"expand,omitempty"`
	Group            *Group       `json:"group,omitempty"`
	Field            *Field       `json:"field,omitempty"`
	EmailAddress     string       `json:"emailAddress,omitempty"`
	ProjectRole      *ProjectRole `json:"projectRole,omitempty"`
	User             *User        `json:"user,omitempty"`
}

// ProjectRole represents a project role.
type ProjectRole struct {
	Self        string       `json:"self,omitempty"`
	ID          int64        `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Actors      []*RoleActor `json:"actors,omitempty"`
	Scope       *Scope       `json:"scope,omitempty"`
	Admin       bool         `json:"admin,omitempty"`
	Default     bool         `json:"default,omitempty"`
}

// RoleActor represents an actor in a project role.
type RoleActor struct {
	ID          int64       `json:"id,omitempty"`
	DisplayName string      `json:"displayName,omitempty"`
	Type        string      `json:"type,omitempty"`
	Name        string      `json:"name,omitempty"`
	AvatarURL   string      `json:"avatarUrl,omitempty"`
	ActorUser   *ActorUser  `json:"actorUser,omitempty"`
	ActorGroup  *ActorGroup `json:"actorGroup,omitempty"`
}

// ActorUser represents a user actor.
type ActorUser struct {
	AccountID string `json:"accountId,omitempty"`
}

// ActorGroup represents a group actor.
type ActorGroup struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	GroupID     string `json:"groupId,omitempty"`
}

// Scope represents a project scope.
type Scope struct {
	Type    string   `json:"type,omitempty"`
	Project *Project `json:"project,omitempty"`
}

// ListRecent returns recently accessed projects.
func (s *ProjectsService) ListRecent(ctx context.Context, expand []string) ([]*Project, *Response, error) {
	u := "/rest/api/3/project/recent"

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var projects []*Project
	resp, err := s.client.Do(req, &projects)
	if err != nil {
		return nil, resp, err
	}

	return projects, resp, nil
}

// GetSecurityLevels returns all issue security levels for a project.
func (s *ProjectsService) GetSecurityLevels(ctx context.Context, projectKeyOrID string) (*ProjectIssueSecurityLevels, *Response, error) {
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

// ProjectIssueSecurityLevels represents issue security levels for a project.
type ProjectIssueSecurityLevels struct {
	Levels []*SecurityLevel `json:"levels,omitempty"`
}
