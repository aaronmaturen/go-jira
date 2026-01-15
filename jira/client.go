// Package jira provides a Go client for the Jira Cloud REST API v3.
//
// # Authentication
//
// The client supports API token authentication for Jira Cloud:
//
//	client := jira.NewClient(
//		"https://yoursite.atlassian.net",
//		jira.WithBasicAuth("email@example.com", "your-api-token"),
//	)
//
// # Usage
//
// Create a client and use the service endpoints:
//
//	client := jira.NewClient("https://yoursite.atlassian.net",
//		jira.WithBasicAuth(email, token),
//	)
//
//	// Get an issue
//	issue, _, err := client.Issues.Get(ctx, "PROJ-123", nil)
//
//	// Search with JQL
//	results, _, err := client.Search.Do(ctx, "project = PROJ", nil)
//
//	// List projects
//	projects, _, err := client.Projects.List(ctx, nil)
package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default Jira Cloud API base URL.
	// Users should replace this with their own instance URL.
	DefaultBaseURL = "https://your-domain.atlassian.net"

	// APIVersion is the Jira REST API version this client targets.
	APIVersion = "3"

	// UserAgent is the default user agent string.
	UserAgent = "go-jira/1.0"
)

// Client manages communication with the Jira API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	baseURL *url.URL

	// User agent used when communicating with the API.
	UserAgent string

	// Authentication method
	auth Authenticator

	// Services for different API groups
	Issues           *IssuesService
	Search           *SearchService
	Projects         *ProjectsService
	Users            *UsersService
	Groups           *GroupsService
	Filters          *FiltersService
	Dashboards       *DashboardsService
	IssueTypes       *IssueTypesService
	Priorities       *PrioritiesService
	Resolutions      *ResolutionsService
	Statuses         *StatusesService
	Components       *ComponentsService
	Versions         *VersionsService
	IssueLinks       *IssueLinksService
	IssueLinkTypes   *IssueLinkTypesService
	Attachments      *AttachmentsService
	Comments         *CommentsService
	Worklogs         *WorklogsService
	Watchers         *WatchersService
	Votes            *VotesService
	Fields           *FieldsService
	Screens          *ScreensService
	Workflows        *WorkflowsService
	WorkflowSchemes  *WorkflowSchemesService
	Permissions      *PermissionsService
	ProjectRoles     *ProjectRolesService
	Labels           *LabelsService
	ServerInfo       *ServerInfoService
	Myself           *MyselfService
	ApplicationRoles *ApplicationRolesService
	AuditRecords     *AuditRecordsService
	Avatars          *AvatarsService
	JQL              *JQLService
}

// Authenticator is the interface for authentication methods.
type Authenticator interface {
	// Apply adds authentication to the request.
	Apply(req *http.Request)
}

// BasicAuth implements basic authentication with email and API token.
type BasicAuth struct {
	Email    string
	APIToken string
}

// Apply adds basic auth header to the request.
func (a *BasicAuth) Apply(req *http.Request) {
	req.SetBasicAuth(a.Email, a.APIToken)
}

// BearerAuth implements bearer token authentication.
type BearerAuth struct {
	Token string
}

// Apply adds bearer token header to the request.
func (a *BearerAuth) Apply(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+a.Token)
}

// ClientOption configures the Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.client = httpClient
	}
}

// WithBasicAuth sets basic authentication with email and API token.
func WithBasicAuth(email, apiToken string) ClientOption {
	return func(c *Client) {
		c.auth = &BasicAuth{Email: email, APIToken: apiToken}
	}
}

// WithBearerToken sets bearer token authentication.
func WithBearerToken(token string) ClientOption {
	return func(c *Client) {
		c.auth = &BearerAuth{Token: token}
	}
}

// WithUserAgent sets a custom user agent string.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.UserAgent = userAgent
	}
}

// NewClient returns a new Jira API client.
func NewClient(baseURL string, opts ...ClientOption) (*Client, error) {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	c := &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:   parsedURL,
		UserAgent: UserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	// Initialize services
	c.Issues = &IssuesService{client: c}
	c.Search = &SearchService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.Users = &UsersService{client: c}
	c.Groups = &GroupsService{client: c}
	c.Filters = &FiltersService{client: c}
	c.Dashboards = &DashboardsService{client: c}
	c.IssueTypes = &IssueTypesService{client: c}
	c.Priorities = &PrioritiesService{client: c}
	c.Resolutions = &ResolutionsService{client: c}
	c.Statuses = &StatusesService{client: c}
	c.Components = &ComponentsService{client: c}
	c.Versions = &VersionsService{client: c}
	c.IssueLinks = &IssueLinksService{client: c}
	c.IssueLinkTypes = &IssueLinkTypesService{client: c}
	c.Attachments = &AttachmentsService{client: c}
	c.Comments = &CommentsService{client: c}
	c.Worklogs = &WorklogsService{client: c}
	c.Watchers = &WatchersService{client: c}
	c.Votes = &VotesService{client: c}
	c.Fields = &FieldsService{client: c}
	c.Screens = &ScreensService{client: c}
	c.Workflows = &WorkflowsService{client: c}
	c.WorkflowSchemes = &WorkflowSchemesService{client: c}
	c.Permissions = &PermissionsService{client: c}
	c.ProjectRoles = &ProjectRolesService{client: c}
	c.Labels = &LabelsService{client: c}
	c.ServerInfo = &ServerInfoService{client: c}
	c.Myself = &MyselfService{client: c}
	c.ApplicationRoles = &ApplicationRolesService{client: c}
	c.AuditRecords = &AuditRecordsService{client: c}
	c.Avatars = &AvatarsService{client: c}
	c.JQL = &JQLService{client: c}

	return c, nil
}

// Response represents an API response.
type Response struct {
	*http.Response

	// For paginated responses
	StartAt    int
	MaxResults int
	Total      int
}

// newResponse creates a new Response from an http.Response.
func newResponse(r *http.Response) *Response {
	return &Response{Response: r}
}

// ErrorResponse represents an error response from the Jira API.
type ErrorResponse struct {
	Response      *http.Response    `json:"-"`
	ErrorMessages []string          `json:"errorMessages,omitempty"`
	Errors        map[string]string `json:"errors,omitempty"`
}

// Error implements the error interface.
func (e *ErrorResponse) Error() string {
	if len(e.ErrorMessages) > 0 {
		return fmt.Sprintf("%s %s: %d %s",
			e.Response.Request.Method,
			e.Response.Request.URL,
			e.Response.StatusCode,
			strings.Join(e.ErrorMessages, ", "))
	}
	if len(e.Errors) > 0 {
		var msgs []string
		for k, v := range e.Errors {
			msgs = append(msgs, fmt.Sprintf("%s: %s", k, v))
		}
		return fmt.Sprintf("%s %s: %d %s",
			e.Response.Request.Method,
			e.Response.Request.URL,
			e.Response.StatusCode,
			strings.Join(msgs, ", "))
	}
	return fmt.Sprintf("%s %s: %d",
		e.Response.Request.Method,
		e.Response.Request.URL,
		e.Response.StatusCode)
}

// NewRequest creates an API request.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	// Ensure the URL starts with the API path
	if !strings.HasPrefix(urlStr, "/") {
		urlStr = "/" + urlStr
	}

	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	if c.auth != nil {
		c.auth.Apply(req)
	}

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	if err := checkResponse(resp); err != nil {
		return response, err
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
		if err != nil && err != io.EOF {
			return response, err
		}
	}

	return response, nil
}

// checkResponse checks the API response for errors.
func checkResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool { return &v }

// Int returns a pointer to the given int value.
func Int(v int) *int { return &v }

// String returns a pointer to the given string value.
func String(v string) *string { return &v }
