package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// WorklogsService handles worklog operations for the Jira API.
type WorklogsService struct {
	client *Client
}

// Worklog represents a Jira worklog entry.
type Worklog struct {
	Self             string      `json:"self,omitempty"`
	ID               string      `json:"id,omitempty"`
	IssueID          string      `json:"issueId,omitempty"`
	Author           *User       `json:"author,omitempty"`
	UpdateAuthor     *User       `json:"updateAuthor,omitempty"`
	Comment          interface{} `json:"comment,omitempty"` // Can be string or ADF
	Created          *Time       `json:"created,omitempty"`
	Updated          *Time       `json:"updated,omitempty"`
	Started          *Time       `json:"started,omitempty"`
	TimeSpent        string      `json:"timeSpent,omitempty"`
	TimeSpentSeconds int64       `json:"timeSpentSeconds,omitempty"`
	Visibility       *Visibility `json:"visibility,omitempty"`
	Properties       []*EntityProperty `json:"properties,omitempty"`
}

// EntityProperty represents an entity property.
type EntityProperty struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// WorklogListResult represents a paginated list of worklogs.
type WorklogListResult struct {
	StartAt    int        `json:"startAt,omitempty"`
	MaxResults int        `json:"maxResults,omitempty"`
	Total      int        `json:"total,omitempty"`
	Worklogs   []*Worklog `json:"worklogs,omitempty"`
}

// ListIssueWorklogs returns worklogs for an issue.
func (s *WorklogsService) ListIssueWorklogs(ctx context.Context, issueIDOrKey string, startAt, maxResults int, startedAfter, startedBefore int64, expand []string) (*WorklogListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueIDOrKey)

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	if startedAfter > 0 {
		params.Set("startedAfter", strconv.FormatInt(startedAfter, 10))
	}
	if startedBefore > 0 {
		params.Set("startedBefore", strconv.FormatInt(startedBefore, 10))
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

	result := new(WorklogListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Get returns a worklog by ID.
func (s *WorklogsService) Get(ctx context.Context, issueIDOrKey, worklogID string, expand []string) (*Worklog, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s", issueIDOrKey, worklogID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	worklog := new(Worklog)
	resp, err := s.client.Do(req, worklog)
	if err != nil {
		return nil, resp, err
	}

	return worklog, resp, nil
}

// WorklogCreateRequest represents a request to create a worklog.
type WorklogCreateRequest struct {
	Comment          interface{}       `json:"comment,omitempty"` // Can be string or ADF
	Started          string            `json:"started,omitempty"` // Format: "2021-01-17T12:34:00.000+0000"
	TimeSpent        string            `json:"timeSpent,omitempty"`
	TimeSpentSeconds int64             `json:"timeSpentSeconds,omitempty"`
	Visibility       *Visibility       `json:"visibility,omitempty"`
	Properties       []*EntityProperty `json:"properties,omitempty"`
}

// Add adds a worklog to an issue.
func (s *WorklogsService) Add(ctx context.Context, issueIDOrKey string, worklog *WorklogCreateRequest, notifyUsers bool, adjustEstimate string, newEstimate string, reduceBy string, overrideEditableFlag bool, expand []string) (*Worklog, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueIDOrKey)

	params := url.Values{}
	if !notifyUsers {
		params.Set("notifyUsers", "false")
	}
	if adjustEstimate != "" {
		params.Set("adjustEstimate", adjustEstimate)
	}
	if newEstimate != "" {
		params.Set("newEstimate", newEstimate)
	}
	if reduceBy != "" {
		params.Set("reduceBy", reduceBy)
	}
	if overrideEditableFlag {
		params.Set("overrideEditableFlag", "true")
	}
	if len(expand) > 0 {
		params.Set("expand", strings.Join(expand, ","))
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, worklog)
	if err != nil {
		return nil, nil, err
	}

	result := new(Worklog)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// WorklogUpdateRequest represents a request to update a worklog.
type WorklogUpdateRequest struct {
	Comment          interface{}       `json:"comment,omitempty"`
	Started          string            `json:"started,omitempty"`
	TimeSpent        string            `json:"timeSpent,omitempty"`
	TimeSpentSeconds int64             `json:"timeSpentSeconds,omitempty"`
	Visibility       *Visibility       `json:"visibility,omitempty"`
	Properties       []*EntityProperty `json:"properties,omitempty"`
}

// Update updates a worklog.
func (s *WorklogsService) Update(ctx context.Context, issueIDOrKey, worklogID string, worklog *WorklogUpdateRequest, notifyUsers bool, adjustEstimate string, newEstimate string, overrideEditableFlag bool, expand []string) (*Worklog, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s", issueIDOrKey, worklogID)

	params := url.Values{}
	if !notifyUsers {
		params.Set("notifyUsers", "false")
	}
	if adjustEstimate != "" {
		params.Set("adjustEstimate", adjustEstimate)
	}
	if newEstimate != "" {
		params.Set("newEstimate", newEstimate)
	}
	if overrideEditableFlag {
		params.Set("overrideEditableFlag", "true")
	}
	if len(expand) > 0 {
		params.Set("expand", strings.Join(expand, ","))
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, worklog)
	if err != nil {
		return nil, nil, err
	}

	result := new(Worklog)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a worklog.
func (s *WorklogsService) Delete(ctx context.Context, issueIDOrKey, worklogID string, notifyUsers bool, adjustEstimate string, newEstimate string, increaseBy string, overrideEditableFlag bool) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s", issueIDOrKey, worklogID)

	params := url.Values{}
	if !notifyUsers {
		params.Set("notifyUsers", "false")
	}
	if adjustEstimate != "" {
		params.Set("adjustEstimate", adjustEstimate)
	}
	if newEstimate != "" {
		params.Set("newEstimate", newEstimate)
	}
	if increaseBy != "" {
		params.Set("increaseBy", increaseBy)
	}
	if overrideEditableFlag {
		params.Set("overrideEditableFlag", "true")
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

// WorklogIDsResult represents a list of worklog IDs.
type WorklogIDsResult struct {
	Values        []WorklogID `json:"values,omitempty"`
	Since         int64       `json:"since,omitempty"`
	Until         int64       `json:"until,omitempty"`
	Self          string      `json:"self,omitempty"`
	NextPage      string      `json:"nextPage,omitempty"`
	LastPage      bool        `json:"lastPage,omitempty"`
}

// WorklogID represents a worklog ID.
type WorklogID struct {
	WorklogID   int64 `json:"worklogId,omitempty"`
	UpdatedTime int64 `json:"updatedTime,omitempty"`
}

// ListUpdated returns worklog IDs updated since a given time.
func (s *WorklogsService) ListUpdated(ctx context.Context, since int64, expand []string) (*WorklogIDsResult, *Response, error) {
	u := "/rest/api/3/worklog/updated"

	params := url.Values{}
	if since > 0 {
		params.Set("since", strconv.FormatInt(since, 10))
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

	result := new(WorklogIDsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListDeleted returns worklog IDs deleted since a given time.
func (s *WorklogsService) ListDeleted(ctx context.Context, since int64) (*WorklogIDsResult, *Response, error) {
	u := "/rest/api/3/worklog/deleted"

	if since > 0 {
		u = fmt.Sprintf("%s?since=%d", u, since)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(WorklogIDsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetByIDs returns worklogs by their IDs.
func (s *WorklogsService) GetByIDs(ctx context.Context, ids []int64, expand []string) ([]*Worklog, *Response, error) {
	u := "/rest/api/3/worklog/list"

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, map[string][]int64{"ids": ids})
	if err != nil {
		return nil, nil, err
	}

	var worklogs []*Worklog
	resp, err := s.client.Do(req, &worklogs)
	if err != nil {
		return nil, resp, err
	}

	return worklogs, resp, nil
}

// GetPropertyKeys returns property keys for a worklog.
func (s *WorklogsService) GetPropertyKeys(ctx context.Context, issueIDOrKey, worklogID string) ([]string, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s/properties", issueIDOrKey, worklogID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Keys []struct {
			Key string `json:"key"`
		} `json:"keys"`
	}
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	keys := make([]string, len(result.Keys))
	for i, k := range result.Keys {
		keys[i] = k.Key
	}

	return keys, resp, nil
}

// GetProperty returns a worklog property.
func (s *WorklogsService) GetProperty(ctx context.Context, issueIDOrKey, worklogID, propertyKey string) (*EntityProperty, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s/properties/%s", issueIDOrKey, worklogID, propertyKey)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	prop := new(EntityProperty)
	resp, err := s.client.Do(req, prop)
	if err != nil {
		return nil, resp, err
	}

	return prop, resp, nil
}

// SetProperty sets a worklog property.
func (s *WorklogsService) SetProperty(ctx context.Context, issueIDOrKey, worklogID, propertyKey string, value interface{}) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s/properties/%s", issueIDOrKey, worklogID, propertyKey)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, value)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteProperty deletes a worklog property.
func (s *WorklogsService) DeleteProperty(ctx context.Context, issueIDOrKey, worklogID, propertyKey string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s/properties/%s", issueIDOrKey, worklogID, propertyKey)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
