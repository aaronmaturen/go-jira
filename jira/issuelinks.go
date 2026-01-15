package jira

import (
	"context"
	"fmt"
	"net/http"
)

// IssueLinksService handles issue link operations for the Jira API.
type IssueLinksService struct {
	client *Client
}

// IssueLink represents a link between two issues.
type IssueLink struct {
	ID           string         `json:"id,omitempty"`
	Self         string         `json:"self,omitempty"`
	Type         *IssueLinkType `json:"type,omitempty"`
	InwardIssue  *LinkedIssue   `json:"inwardIssue,omitempty"`
	OutwardIssue *LinkedIssue   `json:"outwardIssue,omitempty"`
}

// LinkedIssue represents a linked issue reference.
type LinkedIssue struct {
	ID     string           `json:"id,omitempty"`
	Key    string           `json:"key,omitempty"`
	Self   string           `json:"self,omitempty"`
	Fields *LinkedIssueFields `json:"fields,omitempty"`
}

// LinkedIssueFields represents the fields of a linked issue.
type LinkedIssueFields struct {
	Summary    string     `json:"summary,omitempty"`
	Status     *Status    `json:"status,omitempty"`
	Priority   *Priority  `json:"priority,omitempty"`
	IssueType  *IssueType `json:"issuetype,omitempty"`
}

// Get returns an issue link by ID.
func (s *IssueLinksService) Get(ctx context.Context, linkID string) (*IssueLink, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issueLink/%s", linkID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	link := new(IssueLink)
	resp, err := s.client.Do(req, link)
	if err != nil {
		return nil, resp, err
	}

	return link, resp, nil
}

// IssueLinkCreateRequest represents a request to create an issue link.
type IssueLinkCreateRequest struct {
	Type         *IssueLinkTypeRef `json:"type"`
	InwardIssue  *IssueRef         `json:"inwardIssue"`
	OutwardIssue *IssueRef         `json:"outwardIssue"`
	Comment      *CommentRef       `json:"comment,omitempty"`
}

// IssueLinkTypeRef represents a reference to an issue link type.
type IssueLinkTypeRef struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// IssueRef represents a reference to an issue.
type IssueRef struct {
	ID  string `json:"id,omitempty"`
	Key string `json:"key,omitempty"`
}

// CommentRef represents a comment reference for link creation.
type CommentRef struct {
	Body       interface{} `json:"body,omitempty"` // Can be string or ADF
	Visibility *Visibility `json:"visibility,omitempty"`
}

// Create creates an issue link.
func (s *IssueLinksService) Create(ctx context.Context, link *IssueLinkCreateRequest) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/issueLink", link)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Delete removes an issue link.
func (s *IssueLinksService) Delete(ctx context.Context, linkID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issueLink/%s", linkID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
