package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// UsersService handles user operations for the Jira API.
type UsersService struct {
	client *Client
}

// UserSearchOptions specifies options for searching users.
type UserSearchOptions struct {
	// Query string to search for in user properties.
	Query string `url:"query,omitempty"`

	// Username deprecated, use query instead.
	Username string `url:"username,omitempty"`

	// AccountID filter by account ID.
	AccountID string `url:"accountId,omitempty"`

	// StartAt index of the first item to return.
	StartAt int `url:"startAt,omitempty"`

	// MaxResults maximum number of items to return.
	MaxResults int `url:"maxResults,omitempty"`

	// Property filter by user property key and value.
	Property string `url:"property,omitempty"`
}

// Search finds users matching the query.
func (s *UsersService) Search(ctx context.Context, opts *UserSearchOptions) ([]*User, *Response, error) {
	u := "/rest/api/3/user/search"

	if opts != nil {
		params := url.Values{}
		if opts.Query != "" {
			params.Set("query", opts.Query)
		}
		if opts.AccountID != "" {
			params.Set("accountId", opts.AccountID)
		}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		if opts.Property != "" {
			params.Set("property", opts.Property)
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

// Get returns a user by account ID.
func (s *UsersService) Get(ctx context.Context, accountID string, expand []string) (*User, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/user?accountId=%s", accountID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s&expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

// UserCreateRequest represents a request to create a user.
type UserCreateRequest struct {
	EmailAddress string   `json:"emailAddress"`
	DisplayName  string   `json:"displayName,omitempty"`
	Name         string   `json:"name,omitempty"`
	Password     string   `json:"password,omitempty"`
	Products     []string `json:"products,omitempty"`
}

// Create creates a new user.
func (s *UsersService) Create(ctx context.Context, user *UserCreateRequest) (*User, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/user", user)
	if err != nil {
		return nil, nil, err
	}

	result := new(User)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a user.
func (s *UsersService) Delete(ctx context.Context, accountID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/user?accountId=%s", accountID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// BulkGetOptions specifies options for bulk getting users.
type BulkGetOptions struct {
	AccountIDs []string `url:"accountId,omitempty"`
	StartAt    int      `url:"startAt,omitempty"`
	MaxResults int      `url:"maxResults,omitempty"`
}

// BulkGetResult represents the result of a bulk get operation.
type BulkGetResult struct {
	Self       string  `json:"self,omitempty"`
	NextPage   string  `json:"nextPage,omitempty"`
	MaxResults int     `json:"maxResults,omitempty"`
	StartAt    int     `json:"startAt,omitempty"`
	Total      int     `json:"total,omitempty"`
	IsLast     bool    `json:"isLast,omitempty"`
	Values     []*User `json:"values,omitempty"`
}

// BulkGet returns multiple users by their account IDs.
func (s *UsersService) BulkGet(ctx context.Context, opts *BulkGetOptions) (*BulkGetResult, *Response, error) {
	u := "/rest/api/3/user/bulk"

	if opts != nil {
		params := url.Values{}
		for _, id := range opts.AccountIDs {
			params.Add("accountId", id)
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

	result := new(BulkGetResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// AccountIDMigrationRequest represents a request to get account IDs.
type AccountIDMigrationRequest struct {
	QueryStrings []string `json:"queryStrings"`
}

// AccountIDMigrationResult represents account ID migration results.
type AccountIDMigrationResult struct {
	Key       string `json:"key,omitempty"`
	AccountID string `json:"accountId,omitempty"`
}

// BulkGetMigration returns account IDs for users given their keys or usernames.
func (s *UsersService) BulkGetMigration(ctx context.Context, startAt, maxResults int, userKeys, usernames []string) ([]*AccountIDMigrationResult, *Response, error) {
	u := "/rest/api/3/user/bulk/migration"

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	for _, k := range userKeys {
		params.Add("key", k)
	}
	for _, u := range usernames {
		params.Add("username", u)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var results []*AccountIDMigrationResult
	resp, err := s.client.Do(req, &results)
	if err != nil {
		return nil, resp, err
	}

	return results, resp, nil
}

// GetDefaultColumns returns the default issue table columns for a user.
func (s *UsersService) GetDefaultColumns(ctx context.Context, accountID string) ([]*ColumnItem, *Response, error) {
	u := "/rest/api/3/user/columns"
	if accountID != "" {
		u = fmt.Sprintf("%s?accountId=%s", u, accountID)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var columns []*ColumnItem
	resp, err := s.client.Do(req, &columns)
	if err != nil {
		return nil, resp, err
	}

	return columns, resp, nil
}

// ColumnItem represents a column in the issue navigator.
type ColumnItem struct {
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}

// SetDefaultColumns sets the default issue table columns for a user.
func (s *UsersService) SetDefaultColumns(ctx context.Context, accountID string, columns []string) (*Response, error) {
	u := "/rest/api/3/user/columns"
	if accountID != "" {
		u = fmt.Sprintf("%s?accountId=%s", u, accountID)
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, map[string][]string{"columns": columns})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ResetDefaultColumns resets the default issue table columns for a user to the system default.
func (s *UsersService) ResetDefaultColumns(ctx context.Context, accountID string) (*Response, error) {
	u := "/rest/api/3/user/columns"
	if accountID != "" {
		u = fmt.Sprintf("%s?accountId=%s", u, accountID)
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetGroups returns the groups a user belongs to.
func (s *UsersService) GetGroups(ctx context.Context, accountID string) ([]*GroupName, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/user/groups?accountId=%s", accountID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var groups []*GroupName
	resp, err := s.client.Do(req, &groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}

// GroupName represents a group name.
type GroupName struct {
	Name   string `json:"name,omitempty"`
	Self   string `json:"self,omitempty"`
	GroupID string `json:"groupId,omitempty"`
}

// FindAssignableOptions specifies options for finding assignable users.
type FindAssignableOptions struct {
	Query              string `url:"query,omitempty"`
	SessionID          string `url:"sessionId,omitempty"`
	Username           string `url:"username,omitempty"`
	AccountID          string `url:"accountId,omitempty"`
	Project            string `url:"project,omitempty"`
	IssueKey           string `url:"issueKey,omitempty"`
	StartAt            int    `url:"startAt,omitempty"`
	MaxResults         int    `url:"maxResults,omitempty"`
	ActionDescriptorID int    `url:"actionDescriptorId,omitempty"`
	Recommend          bool   `url:"recommend,omitempty"`
}

// FindAssignableUsers finds users that can be assigned to an issue.
func (s *UsersService) FindAssignableUsers(ctx context.Context, opts *FindAssignableOptions) ([]*User, *Response, error) {
	u := "/rest/api/3/user/assignable/search"

	if opts != nil {
		params := url.Values{}
		if opts.Query != "" {
			params.Set("query", opts.Query)
		}
		if opts.Project != "" {
			params.Set("project", opts.Project)
		}
		if opts.IssueKey != "" {
			params.Set("issueKey", opts.IssueKey)
		}
		if opts.AccountID != "" {
			params.Set("accountId", opts.AccountID)
		}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		if opts.ActionDescriptorID > 0 {
			params.Set("actionDescriptorId", strconv.Itoa(opts.ActionDescriptorID))
		}
		if opts.Recommend {
			params.Set("recommend", "true")
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

// FindAssignableMultiProject finds users that can be assigned to issues in multiple projects.
func (s *UsersService) FindAssignableMultiProject(ctx context.Context, projectKeys []string, query string, startAt, maxResults int) ([]*User, *Response, error) {
	u := "/rest/api/3/user/assignable/multiProjectSearch"

	params := url.Values{}
	params.Set("projectKeys", strings.Join(projectKeys, ","))
	if query != "" {
		params.Set("query", query)
	}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}

	u = fmt.Sprintf("%s?%s", u, params.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

// FindUsersWithPermissions finds users that have specific permissions.
func (s *UsersService) FindUsersWithPermissions(ctx context.Context, permissions []string, query, accountID string, issueKey, projectKey string, startAt, maxResults int) ([]*User, *Response, error) {
	u := "/rest/api/3/user/permission/search"

	params := url.Values{}
	params.Set("permissions", strings.Join(permissions, ","))
	if query != "" {
		params.Set("query", query)
	}
	if accountID != "" {
		params.Set("accountId", accountID)
	}
	if issueKey != "" {
		params.Set("issueKey", issueKey)
	}
	if projectKey != "" {
		params.Set("projectKey", projectKey)
	}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}

	u = fmt.Sprintf("%s?%s", u, params.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

// FindUsersForPicker finds users for the user picker.
func (s *UsersService) FindUsersForPicker(ctx context.Context, query string, maxResults int, showAvatar bool, exclude, excludeAccountIDs, avatarSize string, excludeConnectUsers bool) (*UserPickerResult, *Response, error) {
	u := "/rest/api/3/user/picker"

	params := url.Values{}
	params.Set("query", query)
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	if showAvatar {
		params.Set("showAvatar", "true")
	}
	if exclude != "" {
		params.Set("exclude", exclude)
	}
	if excludeAccountIDs != "" {
		params.Set("excludeAccountIds", excludeAccountIDs)
	}
	if avatarSize != "" {
		params.Set("avatarSize", avatarSize)
	}
	if excludeConnectUsers {
		params.Set("excludeConnectUsers", "true")
	}

	u = fmt.Sprintf("%s?%s", u, params.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(UserPickerResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UserPickerResult represents the result of a user picker search.
type UserPickerResult struct {
	Users  []*UserPickerUser `json:"users,omitempty"`
	Total  int               `json:"total,omitempty"`
	Header string            `json:"header,omitempty"`
}

// UserPickerUser represents a user in the picker results.
type UserPickerUser struct {
	AccountID   string `json:"accountId,omitempty"`
	AccountType string `json:"accountType,omitempty"`
	HTML        string `json:"html,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
}

// GetAllUsersOptions specifies options for getting all users.
type GetAllUsersOptions struct {
	StartAt    int `url:"startAt,omitempty"`
	MaxResults int `url:"maxResults,omitempty"`
}

// GetAllUsers returns all users visible to the current user.
func (s *UsersService) GetAllUsers(ctx context.Context, opts *GetAllUsersOptions) ([]*User, *Response, error) {
	u := "/rest/api/3/users/search"

	if opts != nil {
		params := url.Values{}
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

	var users []*User
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

// GetEmail returns the email address for a user.
func (s *UsersService) GetEmail(ctx context.Context, accountID string) (*UserEmail, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/user/email?accountId=%s", accountID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	email := new(UserEmail)
	resp, err := s.client.Do(req, email)
	if err != nil {
		return nil, resp, err
	}

	return email, resp, nil
}

// UserEmail represents a user's email address.
type UserEmail struct {
	AccountID string `json:"accountId,omitempty"`
	Email     string `json:"email,omitempty"`
}

// BulkGetEmail returns email addresses for multiple users.
func (s *UsersService) BulkGetEmail(ctx context.Context, accountIDs []string) (*UserEmailList, *Response, error) {
	u := "/rest/api/3/user/email/bulk"

	params := url.Values{}
	for _, id := range accountIDs {
		params.Add("accountId", id)
	}

	u = fmt.Sprintf("%s?%s", u, params.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(UserEmailList)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UserEmailList represents a list of user emails.
type UserEmailList struct {
	Emails map[string]string `json:"emails,omitempty"`
}
