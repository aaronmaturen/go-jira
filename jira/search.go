package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// SearchService handles search operations for the Jira API.
type SearchService struct {
	client *Client
}

// SearchOptions specifies optional parameters for search requests.
type SearchOptions struct {
	// Fields to return for each issue.
	Fields []string `url:"fields,omitempty"`

	// Expand additional information.
	Expand []string `url:"expand,omitempty"`

	// Properties to return for each issue.
	Properties []string `url:"properties,omitempty"`

	// FieldsByKeys whether fields are referenced by keys rather than IDs.
	FieldsByKeys bool `url:"fieldsByKeys,omitempty"`

	// MaxResults maximum number of results to return.
	MaxResults int `url:"maxResults,omitempty"`

	// NextPageToken for pagination (Jira API v3 style).
	NextPageToken string `url:"nextPageToken,omitempty"`

	// StartAt index of the first result to return (legacy pagination).
	StartAt int `url:"startAt,omitempty"`

	// ValidateQuery level of JQL query validation.
	ValidateQuery string `url:"validateQuery,omitempty"`
}

// SearchResult represents the result of a search query.
type SearchResult struct {
	Expand        string   `json:"expand,omitempty"`
	StartAt       int      `json:"startAt,omitempty"`
	MaxResults    int      `json:"maxResults,omitempty"`
	Total         int      `json:"total,omitempty"`
	Issues        []*Issue `json:"issues,omitempty"`
	WarningMessages []string `json:"warningMessages,omitempty"`
	Names         map[string]string `json:"names,omitempty"`
	Schema        map[string]interface{} `json:"schema,omitempty"`
	NextPageToken string   `json:"nextPageToken,omitempty"`
}

// SearchRequest represents a POST search request body.
type SearchRequest struct {
	JQL           string   `json:"jql,omitempty"`
	StartAt       int      `json:"startAt,omitempty"`
	MaxResults    int      `json:"maxResults,omitempty"`
	Fields        []string `json:"fields,omitempty"`
	Expand        []string `json:"expand,omitempty"`
	Properties    []string `json:"properties,omitempty"`
	FieldsByKeys  bool     `json:"fieldsByKeys,omitempty"`
	ValidateQuery string   `json:"validateQuery,omitempty"`
}

// Do performs a JQL search using the new v3 endpoint.
// This is the recommended search method for Jira Cloud.
func (s *SearchService) Do(ctx context.Context, jql string, opts *SearchOptions) (*SearchResult, *Response, error) {
	u := "/rest/api/3/search/jql"

	params := url.Values{}
	params.Set("jql", jql)

	if opts != nil {
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		if opts.NextPageToken != "" {
			params.Set("nextPageToken", opts.NextPageToken)
		}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if len(opts.Fields) > 0 {
			for _, f := range opts.Fields {
				params.Add("fields", f)
			}
		}
		if len(opts.Expand) > 0 {
			for _, e := range opts.Expand {
				params.Add("expand", e)
			}
		}
		if opts.ValidateQuery != "" {
			params.Set("validateQuery", opts.ValidateQuery)
		}
		if opts.FieldsByKeys {
			params.Set("fieldsByKeys", "true")
		}
	}

	u = fmt.Sprintf("%s?%s", u, params.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SearchResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DoPost performs a JQL search using POST method.
// Use this for complex queries that might exceed URL length limits.
func (s *SearchService) DoPost(ctx context.Context, searchReq *SearchRequest) (*SearchResult, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/search/jql", searchReq)
	if err != nil {
		return nil, nil, err
	}

	result := new(SearchResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Legacy performs a search using the legacy /rest/api/3/search endpoint.
// Deprecated: Use Do() instead which uses the new /rest/api/3/search/jql endpoint.
func (s *SearchService) Legacy(ctx context.Context, jql string, opts *SearchOptions) (*SearchResult, *Response, error) {
	u := "/rest/api/3/search"

	params := url.Values{}
	params.Set("jql", jql)

	if opts != nil {
		if opts.MaxResults > 0 {
			params.Set("maxResults", strconv.Itoa(opts.MaxResults))
		}
		if opts.StartAt > 0 {
			params.Set("startAt", strconv.Itoa(opts.StartAt))
		}
		if len(opts.Fields) > 0 {
			for _, f := range opts.Fields {
				params.Add("fields", f)
			}
		}
		if len(opts.Expand) > 0 {
			for _, e := range opts.Expand {
				params.Add("expand", e)
			}
		}
	}

	u = fmt.Sprintf("%s?%s", u, params.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SearchResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// PickerSuggestions represents issue picker suggestions.
type PickerSuggestions struct {
	Sections []*PickerSection `json:"sections,omitempty"`
}

// PickerSection represents a section in issue picker suggestions.
type PickerSection struct {
	Label      string             `json:"label,omitempty"`
	Sub        string             `json:"sub,omitempty"`
	ID         string             `json:"id,omitempty"`
	Msg        string             `json:"msg,omitempty"`
	Issues     []*PickerIssue     `json:"issues,omitempty"`
}

// PickerIssue represents an issue in picker suggestions.
type PickerIssue struct {
	Key            string `json:"key,omitempty"`
	KeyHTML        string `json:"keyHtml,omitempty"`
	Img            string `json:"img,omitempty"`
	Summary        string `json:"summary,omitempty"`
	SummaryText    string `json:"summaryText,omitempty"`
}

// PickerOptions specifies options for the issue picker.
type PickerOptions struct {
	Query           string `url:"query,omitempty"`
	CurrentJQL      string `url:"currentJQL,omitempty"`
	CurrentIssueKey string `url:"currentIssueKey,omitempty"`
	CurrentProjectID string `url:"currentProjectId,omitempty"`
	ShowSubTasks    bool   `url:"showSubTasks,omitempty"`
	ShowSubTaskParent bool `url:"showSubTaskParent,omitempty"`
}

// Picker returns issue picker suggestions.
func (s *SearchService) Picker(ctx context.Context, opts *PickerOptions) (*PickerSuggestions, *Response, error) {
	u := "/rest/api/3/issue/picker"

	if opts != nil {
		params := url.Values{}
		if opts.Query != "" {
			params.Set("query", opts.Query)
		}
		if opts.CurrentJQL != "" {
			params.Set("currentJQL", opts.CurrentJQL)
		}
		if opts.CurrentIssueKey != "" {
			params.Set("currentIssueKey", opts.CurrentIssueKey)
		}
		if opts.CurrentProjectID != "" {
			params.Set("currentProjectId", opts.CurrentProjectID)
		}
		if opts.ShowSubTasks {
			params.Set("showSubTasks", "true")
		}
		if opts.ShowSubTaskParent {
			params.Set("showSubTaskParent", "true")
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(PickerSuggestions)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// MatchRequest represents a request to check issues against JQL.
type MatchRequest struct {
	IssueIDs []int64  `json:"issueIds,omitempty"`
	JQLs     []string `json:"jqls,omitempty"`
}

// MatchResult represents the result of matching issues against JQL.
type MatchResult struct {
	Matches []*MatchEntry `json:"matches,omitempty"`
}

// MatchEntry represents a single match result.
type MatchEntry struct {
	MatchedIssues []int64  `json:"matchedIssues,omitempty"`
	Errors        []string `json:"errors,omitempty"`
}

// Match checks whether one or more issues match one or more JQL queries.
func (s *SearchService) Match(ctx context.Context, req *MatchRequest) (*MatchResult, *Response, error) {
	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/jql/match", req)
	if err != nil {
		return nil, nil, err
	}

	result := new(MatchResult)
	resp, err := s.client.Do(httpReq, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
