package jira

import (
	"context"
	"fmt"
	"net/http"
)

// IssueLinkTypesService handles issue link type operations for the Jira API.
type IssueLinkTypesService struct {
	client *Client
}

// IssueLinkType represents a type of link between issues.
type IssueLinkType struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Inward  string `json:"inward,omitempty"`
	Outward string `json:"outward,omitempty"`
	Self    string `json:"self,omitempty"`
}

// IssueLinkTypesResult represents a list of issue link types.
type IssueLinkTypesResult struct {
	IssueLinkTypes []*IssueLinkType `json:"issueLinkTypes,omitempty"`
}

// List returns all issue link types.
func (s *IssueLinkTypesService) List(ctx context.Context) (*IssueLinkTypesResult, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/issueLinkType", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssueLinkTypesResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Get returns an issue link type by ID.
func (s *IssueLinkTypesService) Get(ctx context.Context, linkTypeID string) (*IssueLinkType, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issueLinkType/%s", linkTypeID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	linkType := new(IssueLinkType)
	resp, err := s.client.Do(req, linkType)
	if err != nil {
		return nil, resp, err
	}

	return linkType, resp, nil
}

// IssueLinkTypeCreateRequest represents a request to create an issue link type.
type IssueLinkTypeCreateRequest struct {
	Name    string `json:"name"`
	Inward  string `json:"inward"`
	Outward string `json:"outward"`
}

// Create creates a new issue link type.
func (s *IssueLinkTypesService) Create(ctx context.Context, linkType *IssueLinkTypeCreateRequest) (*IssueLinkType, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/issueLinkType", linkType)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssueLinkType)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// IssueLinkTypeUpdateRequest represents a request to update an issue link type.
type IssueLinkTypeUpdateRequest struct {
	Name    string `json:"name,omitempty"`
	Inward  string `json:"inward,omitempty"`
	Outward string `json:"outward,omitempty"`
}

// Update updates an issue link type.
func (s *IssueLinkTypesService) Update(ctx context.Context, linkTypeID string, linkType *IssueLinkTypeUpdateRequest) (*IssueLinkType, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issueLinkType/%s", linkTypeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, linkType)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssueLinkType)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes an issue link type.
func (s *IssueLinkTypesService) Delete(ctx context.Context, linkTypeID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issueLinkType/%s", linkTypeID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
