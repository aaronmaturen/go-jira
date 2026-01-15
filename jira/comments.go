package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CommentsService handles comment operations for the Jira API.
type CommentsService struct {
	client *Client
}

// CommentListResult represents a paginated list of comments.
type CommentListResult struct {
	StartAt    int        `json:"startAt,omitempty"`
	MaxResults int        `json:"maxResults,omitempty"`
	Total      int        `json:"total,omitempty"`
	Comments   []*Comment `json:"comments,omitempty"`
}

// ListIssueComments returns comments for an issue.
func (s *CommentsService) ListIssueComments(ctx context.Context, issueIDOrKey string, startAt, maxResults int, orderBy string, expand []string) (*CommentListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/comment", issueIDOrKey)

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	if orderBy != "" {
		params.Set("orderBy", orderBy)
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

	result := new(CommentListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Get returns a comment by ID.
func (s *CommentsService) Get(ctx context.Context, issueIDOrKey, commentID string, expand []string) (*Comment, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueIDOrKey, commentID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	comment := new(Comment)
	resp, err := s.client.Do(req, comment)
	if err != nil {
		return nil, resp, err
	}

	return comment, resp, nil
}

// CommentCreateRequest represents a request to create a comment.
type CommentCreateRequest struct {
	Body       interface{} `json:"body"` // Can be string or ADF document
	Visibility *Visibility `json:"visibility,omitempty"`
}

// Visibility represents comment visibility settings.
type Visibility struct {
	Type       string `json:"type,omitempty"` // group, role
	Value      string `json:"value,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

// Add adds a comment to an issue.
func (s *CommentsService) Add(ctx context.Context, issueIDOrKey string, comment *CommentCreateRequest, expand []string) (*Comment, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/comment", issueIDOrKey)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, comment)
	if err != nil {
		return nil, nil, err
	}

	result := new(Comment)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CommentUpdateRequest represents a request to update a comment.
type CommentUpdateRequest struct {
	Body       interface{} `json:"body,omitempty"` // Can be string or ADF document
	Visibility *Visibility `json:"visibility,omitempty"`
}

// Update updates a comment.
func (s *CommentsService) Update(ctx context.Context, issueIDOrKey, commentID string, comment *CommentUpdateRequest, notifyUsers bool, overrideEditableFlag bool, expand []string) (*Comment, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueIDOrKey, commentID)

	params := url.Values{}
	if !notifyUsers {
		params.Set("notifyUsers", "false")
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

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, comment)
	if err != nil {
		return nil, nil, err
	}

	result := new(Comment)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a comment.
func (s *CommentsService) Delete(ctx context.Context, issueIDOrKey, commentID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueIDOrKey, commentID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetCommentsByIDs returns comments by their IDs.
type GetCommentsByIDsRequest struct {
	IDs []int64 `json:"ids"`
}

// GetCommentsByIDsResult represents the result of getting comments by IDs.
type GetCommentsByIDsResult struct {
	MaxResults int        `json:"maxResults,omitempty"`
	StartAt    int        `json:"startAt,omitempty"`
	Total      int        `json:"total,omitempty"`
	IsLast     bool       `json:"isLast,omitempty"`
	Values     []*Comment `json:"values,omitempty"`
}

// GetByIDs returns comments by their IDs.
func (s *CommentsService) GetByIDs(ctx context.Context, ids []int64, expand []string) (*GetCommentsByIDsResult, *Response, error) {
	u := "/rest/api/3/comment/list"

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, &GetCommentsByIDsRequest{IDs: ids})
	if err != nil {
		return nil, nil, err
	}

	result := new(GetCommentsByIDsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CommentProperty represents a comment property.
type CommentProperty struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// GetPropertyKeys returns property keys for a comment.
func (s *CommentsService) GetPropertyKeys(ctx context.Context, issueIDOrKey, commentID string) ([]string, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/comment/%s/properties", commentID)

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

// GetProperty returns a comment property.
func (s *CommentsService) GetProperty(ctx context.Context, issueIDOrKey, commentID, propertyKey string) (*CommentProperty, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/comment/%s/properties/%s", commentID, propertyKey)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	prop := new(CommentProperty)
	resp, err := s.client.Do(req, prop)
	if err != nil {
		return nil, resp, err
	}

	return prop, resp, nil
}

// SetProperty sets a comment property.
func (s *CommentsService) SetProperty(ctx context.Context, issueIDOrKey, commentID, propertyKey string, value interface{}) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/comment/%s/properties/%s", commentID, propertyKey)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, value)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteProperty deletes a comment property.
func (s *CommentsService) DeleteProperty(ctx context.Context, issueIDOrKey, commentID, propertyKey string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/comment/%s/properties/%s", commentID, propertyKey)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
