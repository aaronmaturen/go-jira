package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GroupsService handles group operations for the Jira API.
type GroupsService struct {
	client *Client
}

// Group represents a Jira group.
type Group struct {
	Name    string  `json:"name,omitempty"`
	GroupID string  `json:"groupId,omitempty"`
	Self    string  `json:"self,omitempty"`
	Users   *Users  `json:"users,omitempty"`
	Expand  string  `json:"expand,omitempty"`
}

// Users represents a paginated list of users in a group.
type Users struct {
	Size       int     `json:"size,omitempty"`
	Items      []*User `json:"items,omitempty"`
	MaxResults int     `json:"max-results,omitempty"`
	StartIndex int     `json:"start-index,omitempty"`
	EndIndex   int     `json:"end-index,omitempty"`
}

// GroupCreateRequest represents a request to create a group.
type GroupCreateRequest struct {
	Name string `json:"name"`
}

// Create creates a new group.
func (s *GroupsService) Create(ctx context.Context, name string) (*Group, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/group", &GroupCreateRequest{Name: name})
	if err != nil {
		return nil, nil, err
	}

	group := new(Group)
	resp, err := s.client.Do(req, group)
	if err != nil {
		return nil, resp, err
	}

	return group, resp, nil
}

// Delete removes a group.
func (s *GroupsService) Delete(ctx context.Context, groupName string, swapGroup string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/group?groupname=%s", url.QueryEscape(groupName))

	if swapGroup != "" {
		u = fmt.Sprintf("%s&swapGroup=%s", u, url.QueryEscape(swapGroup))
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Get returns a group by name.
func (s *GroupsService) Get(ctx context.Context, groupName string, expand []string) (*Group, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/group?groupname=%s", url.QueryEscape(groupName))

	if len(expand) > 0 {
		for _, e := range expand {
			u = fmt.Sprintf("%s&expand=%s", u, e)
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	group := new(Group)
	resp, err := s.client.Do(req, group)
	if err != nil {
		return nil, resp, err
	}

	return group, resp, nil
}

// GroupBulkResult represents a paginated list of groups.
type GroupBulkResult struct {
	Self       string   `json:"self,omitempty"`
	NextPage   string   `json:"nextPage,omitempty"`
	MaxResults int      `json:"maxResults,omitempty"`
	StartAt    int      `json:"startAt,omitempty"`
	Total      int      `json:"total,omitempty"`
	IsLast     bool     `json:"isLast,omitempty"`
	Values     []*Group `json:"values,omitempty"`
}

// BulkGetOptions specifies options for bulk getting groups.
type GroupBulkGetOptions struct {
	StartAt    int      `url:"startAt,omitempty"`
	MaxResults int      `url:"maxResults,omitempty"`
	GroupIDs   []string `url:"groupId,omitempty"`
	GroupNames []string `url:"groupName,omitempty"`
}

// BulkGet returns multiple groups.
func (s *GroupsService) BulkGet(ctx context.Context, opts *GroupBulkGetOptions) (*GroupBulkResult, *Response, error) {
	u := "/rest/api/3/group/bulk"

	if opts != nil {
		params := url.Values{}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		for _, id := range opts.GroupIDs {
			params.Add("groupId", id)
		}
		for _, name := range opts.GroupNames {
			params.Add("groupName", name)
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(GroupBulkResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetMembersOptions specifies options for getting group members.
type GetMembersOptions struct {
	IncludeInactiveUsers bool `url:"includeInactiveUsers,omitempty"`
	StartAt              int  `url:"startAt,omitempty"`
	MaxResults           int  `url:"maxResults,omitempty"`
}

// GroupMembersResult represents a paginated list of group members.
type GroupMembersResult struct {
	Self       string  `json:"self,omitempty"`
	NextPage   string  `json:"nextPage,omitempty"`
	MaxResults int     `json:"maxResults,omitempty"`
	StartAt    int     `json:"startAt,omitempty"`
	Total      int     `json:"total,omitempty"`
	IsLast     bool    `json:"isLast,omitempty"`
	Values     []*User `json:"values,omitempty"`
}

// GetMembers returns members of a group.
func (s *GroupsService) GetMembers(ctx context.Context, groupName string, opts *GetMembersOptions) (*GroupMembersResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/group/member?groupname=%s", url.QueryEscape(groupName))

	if opts != nil {
		if opts.IncludeInactiveUsers {
			u = fmt.Sprintf("%s&includeInactiveUsers=true", u)
		}
		if opts.StartAt > 0 {
			u = fmt.Sprintf("%s&startAt=%d", u, opts.StartAt)
		}
		if opts.MaxResults > 0 {
			u = fmt.Sprintf("%s&maxResults=%d", u, opts.MaxResults)
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(GroupMembersResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// AddUserRequest represents a request to add a user to a group.
type AddUserRequest struct {
	AccountID string `json:"accountId,omitempty"`
	Name      string `json:"name,omitempty"`
}

// AddUser adds a user to a group.
func (s *GroupsService) AddUser(ctx context.Context, groupName, accountID string) (*Group, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/group/user?groupname=%s", url.QueryEscape(groupName))

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, &AddUserRequest{AccountID: accountID})
	if err != nil {
		return nil, nil, err
	}

	group := new(Group)
	resp, err := s.client.Do(req, group)
	if err != nil {
		return nil, resp, err
	}

	return group, resp, nil
}

// RemoveUser removes a user from a group.
func (s *GroupsService) RemoveUser(ctx context.Context, groupName, accountID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/group/user?groupname=%s&accountId=%s", url.QueryEscape(groupName), accountID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// FindGroupsOptions specifies options for finding groups.
type FindGroupsOptions struct {
	AccountID         string   `url:"accountId,omitempty"`
	Query             string   `url:"query,omitempty"`
	Exclude           []string `url:"exclude,omitempty"`
	ExcludeID         []string `url:"excludeId,omitempty"`
	MaxResults        int      `url:"maxResults,omitempty"`
	CaseInsensitive   bool     `url:"caseInsensitive,omitempty"`
}

// FoundGroups represents the result of finding groups.
type FoundGroups struct {
	Header string         `json:"header,omitempty"`
	Total  int            `json:"total,omitempty"`
	Groups []*GroupSuggestion `json:"groups,omitempty"`
}

// GroupSuggestion represents a group suggestion.
type GroupSuggestion struct {
	Name    string       `json:"name,omitempty"`
	HTML    string       `json:"html,omitempty"`
	Labels  []*GroupLabel `json:"labels,omitempty"`
	GroupID string       `json:"groupId,omitempty"`
}

// GroupLabel represents a label on a group suggestion.
type GroupLabel struct {
	Text  string `json:"text,omitempty"`
	Title string `json:"title,omitempty"`
	Type  string `json:"type,omitempty"`
}

// Find finds groups matching the query (used for pickers).
func (s *GroupsService) Find(ctx context.Context, opts *FindGroupsOptions) (*FoundGroups, *Response, error) {
	u := "/rest/api/3/groups/picker"

	if opts != nil {
		params := url.Values{}
		if opts.AccountID != "" {
			params.Set("accountId", opts.AccountID)
		}
		if opts.Query != "" {
			params.Set("query", opts.Query)
		}
		for _, e := range opts.Exclude {
			params.Add("exclude", e)
		}
		for _, e := range opts.ExcludeID {
			params.Add("excludeId", e)
		}
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		if opts.CaseInsensitive {
			params.Set("caseInsensitive", "true")
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(FoundGroups)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
