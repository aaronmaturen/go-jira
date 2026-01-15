package jira

import (
	"context"
	"fmt"
	"net/http"
)

// WatchersService handles watcher operations for the Jira API.
type WatchersService struct {
	client *Client
}

// Watchers represents the watchers of an issue.
type Watchers struct {
	Self       string  `json:"self,omitempty"`
	IsWatching bool    `json:"isWatching,omitempty"`
	WatchCount int     `json:"watchCount,omitempty"`
	Watchers   []*User `json:"watchers,omitempty"`
}

// Get returns the watchers for an issue.
func (s *WatchersService) Get(ctx context.Context, issueIDOrKey string) (*Watchers, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/watchers", issueIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	watchers := new(Watchers)
	resp, err := s.client.Do(req, watchers)
	if err != nil {
		return nil, resp, err
	}

	return watchers, resp, nil
}

// Add adds a watcher to an issue.
func (s *WatchersService) Add(ctx context.Context, issueIDOrKey, accountID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/watchers", issueIDOrKey)

	// The API expects the account ID as a JSON string, not an object
	req, err := s.client.NewRequest(ctx, http.MethodPost, u, accountID)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Remove removes a watcher from an issue.
func (s *WatchersService) Remove(ctx context.Context, issueIDOrKey, accountID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/watchers?accountId=%s", issueIDOrKey, accountID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// BulkWatchersResult represents the result of bulk watching operations.
type BulkWatchersResult struct {
	Errors   []string `json:"errors,omitempty"`
	Success  []string `json:"success,omitempty"`
}

// BulkAdd adds a watcher to multiple issues.
func (s *WatchersService) BulkAdd(ctx context.Context, issueKeys []string, accountID string) (*BulkWatchersResult, *Response, error) {
	body := map[string]interface{}{
		"issueIds":  issueKeys,
		"accountId": accountID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/issue/watching", body)
	if err != nil {
		return nil, nil, err
	}

	result := new(BulkWatchersResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// BulkRemove removes a watcher from multiple issues.
func (s *WatchersService) BulkRemove(ctx context.Context, issueKeys []string, accountID string) (*BulkWatchersResult, *Response, error) {
	body := map[string]interface{}{
		"issueIds":  issueKeys,
		"accountId": accountID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, "/rest/api/3/issue/watching", body)
	if err != nil {
		return nil, nil, err
	}

	result := new(BulkWatchersResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
