package jira

import (
	"context"
	"fmt"
	"net/http"
)

// VotesService handles vote operations for the Jira API.
type VotesService struct {
	client *Client
}

// Votes represents the votes for an issue.
type Votes struct {
	Self     string  `json:"self,omitempty"`
	Votes    int     `json:"votes,omitempty"`
	HasVoted bool    `json:"hasVoted,omitempty"`
	Voters   []*User `json:"voters,omitempty"`
}

// Get returns the votes for an issue.
func (s *VotesService) Get(ctx context.Context, issueIDOrKey string) (*Votes, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/votes", issueIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	votes := new(Votes)
	resp, err := s.client.Do(req, votes)
	if err != nil {
		return nil, resp, err
	}

	return votes, resp, nil
}

// Add adds a vote to an issue for the current user.
func (s *VotesService) Add(ctx context.Context, issueIDOrKey string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/votes", issueIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Remove removes the current user's vote from an issue.
func (s *VotesService) Remove(ctx context.Context, issueIDOrKey string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/votes", issueIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
