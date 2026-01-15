package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// VersionsService handles version operations for the Jira API.
type VersionsService struct {
	client *Client
}

// Version represents a Jira project version.
type Version struct {
	Self                      string    `json:"self,omitempty"`
	ID                        string    `json:"id,omitempty"`
	Name                      string    `json:"name,omitempty"`
	Description               string    `json:"description,omitempty"`
	Archived                  bool      `json:"archived,omitempty"`
	Released                  bool      `json:"released,omitempty"`
	StartDate                 string    `json:"startDate,omitempty"`
	ReleaseDate               string    `json:"releaseDate,omitempty"`
	UserStartDate             string    `json:"userStartDate,omitempty"`
	UserReleaseDate           string    `json:"userReleaseDate,omitempty"`
	ProjectID                 int64     `json:"projectId,omitempty"`
	Project                   string    `json:"project,omitempty"`
	Overdue                   bool      `json:"overdue,omitempty"`
	Operations                []*VersionOperation `json:"operations,omitempty"`
	IssuesStatusForFixVersion *IssuesStatusForVersion `json:"issuesStatusForFixVersion,omitempty"`
}

// VersionOperation represents an operation available on a version.
type VersionOperation struct {
	ID           string `json:"id,omitempty"`
	StyleClass   string `json:"styleClass,omitempty"`
	Label        string `json:"label,omitempty"`
	Href         string `json:"href,omitempty"`
	Weight       int    `json:"weight,omitempty"`
}

// IssuesStatusForVersion represents issue status counts for a version.
type IssuesStatusForVersion struct {
	Unmapped   int `json:"unmapped,omitempty"`
	ToDo       int `json:"toDo,omitempty"`
	InProgress int `json:"inProgress,omitempty"`
	Done       int `json:"done,omitempty"`
}

// Get returns a version by ID.
func (s *VersionsService) Get(ctx context.Context, versionID string, expand []string) (*Version, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/version/%s", versionID)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	version := new(Version)
	resp, err := s.client.Do(req, version)
	if err != nil {
		return nil, resp, err
	}

	return version, resp, nil
}

// VersionCreateRequest represents a request to create a version.
type VersionCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ProjectID   int64  `json:"projectId,omitempty"`
	Project     string `json:"project,omitempty"`
	Archived    bool   `json:"archived,omitempty"`
	Released    bool   `json:"released,omitempty"`
	StartDate   string `json:"startDate,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	MoveUnfixedIssuesTo string `json:"moveUnfixedIssuesTo,omitempty"`
}

// Create creates a new version.
func (s *VersionsService) Create(ctx context.Context, version *VersionCreateRequest) (*Version, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/version", version)
	if err != nil {
		return nil, nil, err
	}

	result := new(Version)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// VersionUpdateRequest represents a request to update a version.
type VersionUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Archived    bool   `json:"archived,omitempty"`
	Released    bool   `json:"released,omitempty"`
	StartDate   string `json:"startDate,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	MoveUnfixedIssuesTo string `json:"moveUnfixedIssuesTo,omitempty"`
}

// Update updates a version.
func (s *VersionsService) Update(ctx context.Context, versionID string, version *VersionUpdateRequest) (*Version, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/version/%s", versionID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, version)
	if err != nil {
		return nil, nil, err
	}

	result := new(Version)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a version.
func (s *VersionsService) Delete(ctx context.Context, versionID string, moveFixIssuesTo, moveAffectedIssuesTo string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/version/%s", versionID)

	params := url.Values{}
	if moveFixIssuesTo != "" {
		params.Set("moveFixIssuesTo", moveFixIssuesTo)
	}
	if moveAffectedIssuesTo != "" {
		params.Set("moveAffectedIssuesTo", moveAffectedIssuesTo)
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

// DeleteAndReplace deletes a version and replaces it in issues.
type DeleteAndReplaceRequest struct {
	MoveFixIssuesTo      int64 `json:"moveFixIssuesTo,omitempty"`
	MoveAffectedIssuesTo int64 `json:"moveAffectedIssuesTo,omitempty"`
	CustomFieldReplacementList []*CustomFieldReplacement `json:"customFieldReplacementList,omitempty"`
}

// CustomFieldReplacement represents a custom field replacement.
type CustomFieldReplacement struct {
	CustomFieldID int64 `json:"customFieldId,omitempty"`
	MoveTo        int64 `json:"moveTo,omitempty"`
}

// DeleteAndReplace deletes a version and replaces it in issues.
func (s *VersionsService) DeleteAndReplace(ctx context.Context, versionID string, request *DeleteAndReplaceRequest) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/version/%s/removeAndSwap", versionID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, request)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Merge merges one version into another.
func (s *VersionsService) Merge(ctx context.Context, versionID, moveIssuesTo string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/version/%s/mergeto/%s", versionID, moveIssuesTo)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Move changes the position of a version.
type VersionMoveRequest struct {
	After    string `json:"after,omitempty"`
	Position string `json:"position,omitempty"` // Earlier, Later, First, Last
}

// Move changes the position of a version.
func (s *VersionsService) Move(ctx context.Context, versionID string, request *VersionMoveRequest) (*Version, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/version/%s/move", versionID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, request)
	if err != nil {
		return nil, nil, err
	}

	version := new(Version)
	resp, err := s.client.Do(req, version)
	if err != nil {
		return nil, resp, err
	}

	return version, resp, nil
}

// VersionIssueCounts represents issue counts for a version.
type VersionIssueCounts struct {
	Self                                     string `json:"self,omitempty"`
	IssuesFixedCount                         int    `json:"issuesFixedCount,omitempty"`
	IssuesAffectedCount                      int    `json:"issuesAffectedCount,omitempty"`
	IssueCountWithCustomFieldsShowingVersion int    `json:"issueCountWithCustomFieldsShowingVersion,omitempty"`
	CustomFieldUsage                         []*VersionUsageInCustomField `json:"customFieldUsage,omitempty"`
}

// VersionUsageInCustomField represents version usage in a custom field.
type VersionUsageInCustomField struct {
	FieldName              string `json:"fieldName,omitempty"`
	CustomFieldID          int64  `json:"customFieldId,omitempty"`
	IssueCountWithVersionInCustomField int `json:"issueCountWithVersionInCustomField,omitempty"`
}

// GetIssueCounts returns issue counts for a version.
func (s *VersionsService) GetIssueCounts(ctx context.Context, versionID string) (*VersionIssueCounts, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/version/%s/relatedIssueCounts", versionID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	counts := new(VersionIssueCounts)
	resp, err := s.client.Do(req, counts)
	if err != nil {
		return nil, resp, err
	}

	return counts, resp, nil
}

// VersionUnresolvedIssueCounts represents unresolved issue counts.
type VersionUnresolvedIssueCounts struct {
	Self                  string `json:"self,omitempty"`
	IssuesUnresolvedCount int    `json:"issuesUnresolvedCount,omitempty"`
	IssuesCount           int    `json:"issuesCount,omitempty"`
}

// GetUnresolvedIssueCounts returns unresolved issue counts for a version.
func (s *VersionsService) GetUnresolvedIssueCounts(ctx context.Context, versionID string) (*VersionUnresolvedIssueCounts, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/version/%s/unresolvedIssueCount", versionID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	counts := new(VersionUnresolvedIssueCounts)
	resp, err := s.client.Do(req, counts)
	if err != nil {
		return nil, resp, err
	}

	return counts, resp, nil
}

// VersionListResult represents a paginated list of versions.
type VersionListResult struct {
	Self       string     `json:"self,omitempty"`
	NextPage   string     `json:"nextPage,omitempty"`
	MaxResults int        `json:"maxResults,omitempty"`
	StartAt    int        `json:"startAt,omitempty"`
	Total      int        `json:"total,omitempty"`
	IsLast     bool       `json:"isLast,omitempty"`
	Values     []*Version `json:"values,omitempty"`
}

// ListProjectVersions returns versions for a project.
func (s *VersionsService) ListProjectVersions(ctx context.Context, projectIDOrKey string, startAt, maxResults int, orderBy, query, status string, expand []string) (*VersionListResult, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/version", projectIDOrKey)

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
	if query != "" {
		params.Set("query", query)
	}
	if status != "" {
		params.Set("status", status)
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

	result := new(VersionListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListAllProjectVersions returns all versions for a project (non-paginated).
func (s *VersionsService) ListAllProjectVersions(ctx context.Context, projectIDOrKey string, expand []string) ([]*Version, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/project/%s/versions", projectIDOrKey)

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var versions []*Version
	resp, err := s.client.Do(req, &versions)
	if err != nil {
		return nil, resp, err
	}

	return versions, resp, nil
}

// ParseDate parses a Jira date string.
func (v *Version) ParseStartDate() (time.Time, error) {
	if v.StartDate == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", v.StartDate)
}

// ParseReleaseDate parses a Jira release date string.
func (v *Version) ParseReleaseDate() (time.Time, error) {
	if v.ReleaseDate == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", v.ReleaseDate)
}
