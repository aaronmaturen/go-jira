package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// LabelsService handles label operations for the Jira API.
type LabelsService struct {
	client *Client
}

// LabelsListResult represents a paginated list of labels.
type LabelsListResult struct {
	Self       string   `json:"self,omitempty"`
	NextPage   string   `json:"nextPage,omitempty"`
	MaxResults int      `json:"maxResults,omitempty"`
	StartAt    int      `json:"startAt,omitempty"`
	Total      int      `json:"total,omitempty"`
	IsLast     bool     `json:"isLast,omitempty"`
	Values     []string `json:"values,omitempty"`
}

// List returns all labels.
func (s *LabelsService) List(ctx context.Context, startAt, maxResults int) (*LabelsListResult, *Response, error) {
	u := "/rest/api/3/label"

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(LabelsListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
