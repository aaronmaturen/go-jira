package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// IssuesService handles communication with the issue related methods of the Jira API.
type IssuesService struct {
	client *Client
}

// Issue represents a Jira issue.
type Issue struct {
	Self           string             `json:"self,omitempty"`
	ID             string             `json:"id,omitempty"`
	Key            string             `json:"key,omitempty"`
	Expand         string             `json:"expand,omitempty"`
	Fields         *IssueFields       `json:"fields,omitempty"`
	Changelog      *Changelog         `json:"changelog,omitempty"`
	Operations     *Operations        `json:"operations,omitempty"`
	Editmeta       *EditMeta          `json:"editmeta,omitempty"`
	Transitions    []*Transition      `json:"transitions,omitempty"`
	Names          map[string]string  `json:"names,omitempty"`
	Schema         map[string]*Schema `json:"schema,omitempty"`
	RenderedFields map[string]any     `json:"renderedFields,omitempty"`
	Properties     []*EntityProperty  `json:"properties,omitempty"`
}

// IssueFields represents the fields of an issue.
type IssueFields struct {
	Summary              string         `json:"summary,omitempty"`
	Description          any            `json:"description,omitempty"` // Can be string or ADF
	Type                 *IssueType     `json:"issuetype,omitempty"`
	Project              *Project       `json:"project,omitempty"`
	Resolution           *Resolution    `json:"resolution,omitempty"`
	Priority             *Priority      `json:"priority,omitempty"`
	Resolutiondate       *Time          `json:"resolutiondate,omitempty"`
	Created              *Time          `json:"created,omitempty"`
	Updated              *Time          `json:"updated,omitempty"`
	DueDate              *Date          `json:"duedate,omitempty"`
	Watches              *Watches       `json:"watches,omitempty"`
	Assignee             *User          `json:"assignee,omitempty"`
	Reporter             *User          `json:"reporter,omitempty"`
	Creator              *User          `json:"creator,omitempty"`
	Votes                *Votes         `json:"votes,omitempty"`
	Labels               []string       `json:"labels,omitempty"`
	Comment              *Comments      `json:"comment,omitempty"`
	Components           []*Component   `json:"components,omitempty"`
	Status               *Status        `json:"status,omitempty"`
	Progress             *Progress      `json:"progress,omitempty"`
	AggregateProgress    *Progress      `json:"aggregateprogress,omitempty"`
	TimeTracking         *TimeTracking  `json:"timetracking,omitempty"`
	TimeSpent            int            `json:"timespent,omitempty"`
	TimeEstimate         int            `json:"timeestimate,omitempty"`
	TimeOriginalEstimate int            `json:"timeoriginalestimate,omitempty"`
	Worklog              *Worklogs      `json:"worklog,omitempty"`
	IssueLinks           []*IssueLink   `json:"issuelinks,omitempty"`
	Attachment           []*Attachment  `json:"attachment,omitempty"`
	Subtasks             []*Issue       `json:"subtasks,omitempty"`
	Parent               *Issue         `json:"parent,omitempty"`
	FixVersions          []*Version     `json:"fixVersions,omitempty"`
	AffectsVersions      []*Version     `json:"versions,omitempty"`
	Environment          any            `json:"environment,omitempty"` // Can be string or ADF
	Security             *SecurityLevel `json:"security,omitempty"`
	Unknowns             map[string]any `json:"-"` // Custom fields
}

// SecurityLevel represents an issue security level.
type SecurityLevel struct {
	Self        string `json:"self,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Progress represents progress information.
type Progress struct {
	Progress int `json:"progress,omitempty"`
	Total    int `json:"total,omitempty"`
	Percent  int `json:"percent,omitempty"`
}

// TimeTracking represents time tracking information.
type TimeTracking struct {
	OriginalEstimate         string `json:"originalEstimate,omitempty"`
	RemainingEstimate        string `json:"remainingEstimate,omitempty"`
	TimeSpent                string `json:"timeSpent,omitempty"`
	OriginalEstimateSeconds  int    `json:"originalEstimateSeconds,omitempty"`
	RemainingEstimateSeconds int    `json:"remainingEstimateSeconds,omitempty"`
	TimeSpentSeconds         int    `json:"timeSpentSeconds,omitempty"`
}

// Comments represents a list of comments on an issue.
type Comments struct {
	StartAt    int        `json:"startAt,omitempty"`
	MaxResults int        `json:"maxResults,omitempty"`
	Total      int        `json:"total,omitempty"`
	Comments   []*Comment `json:"comments,omitempty"`
}

// Worklogs represents a list of worklogs on an issue.
type Worklogs struct {
	StartAt    int        `json:"startAt,omitempty"`
	MaxResults int        `json:"maxResults,omitempty"`
	Total      int        `json:"total,omitempty"`
	Worklogs   []*Worklog `json:"worklogs,omitempty"`
}

// Changelog represents the changelog of an issue.
type Changelog struct {
	StartAt    int              `json:"startAt,omitempty"`
	MaxResults int              `json:"maxResults,omitempty"`
	Total      int              `json:"total,omitempty"`
	Histories  []*ChangeHistory `json:"histories,omitempty"`
}

// ChangeHistory represents a single change in the changelog.
type ChangeHistory struct {
	ID      string        `json:"id,omitempty"`
	Author  *User         `json:"author,omitempty"`
	Created *Time         `json:"created,omitempty"`
	Items   []*ChangeItem `json:"items,omitempty"`
}

// ChangeItem represents a changed field in a change history.
type ChangeItem struct {
	Field      string `json:"field,omitempty"`
	FieldType  string `json:"fieldtype,omitempty"`
	FieldID    string `json:"fieldId,omitempty"`
	From       string `json:"from,omitempty"`
	FromString string `json:"fromString,omitempty"`
	To         string `json:"to,omitempty"`
	ToString   string `json:"toString,omitempty"`
}

// Operations represents available operations on an issue.
type Operations struct {
	LinkGroups []*LinkGroup `json:"linkGroups,omitempty"`
}

// LinkGroup represents a group of operation links.
type LinkGroup struct {
	ID         string        `json:"id,omitempty"`
	StyleClass string        `json:"styleClass,omitempty"`
	Header     *SimpleLink   `json:"header,omitempty"`
	Weight     int           `json:"weight,omitempty"`
	Links      []*SimpleLink `json:"links,omitempty"`
	Groups     []*LinkGroup  `json:"groups,omitempty"`
}

// SimpleLink represents a simple link.
type SimpleLink struct {
	ID         string `json:"id,omitempty"`
	StyleClass string `json:"styleClass,omitempty"`
	IconClass  string `json:"iconClass,omitempty"`
	Label      string `json:"label,omitempty"`
	Title      string `json:"title,omitempty"`
	Href       string `json:"href,omitempty"`
	Weight     int    `json:"weight,omitempty"`
}

// EditMeta represents edit metadata for an issue.
type EditMeta struct {
	Fields map[string]*FieldMeta `json:"fields,omitempty"`
}

// IssueGetOptions specifies optional parameters for Get.
type IssueGetOptions struct {
	// Fields to return for the issue.
	Fields []string `url:"fields,omitempty"`

	// Additional information to include.
	Expand []string `url:"expand,omitempty"`

	// Properties to return for the issue.
	Properties []string `url:"properties,omitempty"`

	// Whether fields should be returned in the response.
	FieldsByKeys bool `url:"fieldsByKeys,omitempty"`

	// Whether to update the issue history.
	UpdateHistory bool `url:"updateHistory,omitempty"`
}

// Get returns a single issue.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-get
func (s *IssuesService) Get(ctx context.Context, issueIDOrKey string, opts *IssueGetOptions) (*Issue, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s", issueIDOrKey)

	if opts != nil {
		query := url.Values{}
		if len(opts.Fields) > 0 {
			query.Set("fields", strings.Join(opts.Fields, ","))
		}
		if len(opts.Expand) > 0 {
			query.Set("expand", strings.Join(opts.Expand, ","))
		}
		if len(opts.Properties) > 0 {
			query.Set("properties", strings.Join(opts.Properties, ","))
		}
		if opts.FieldsByKeys {
			query.Set("fieldsByKeys", "true")
		}
		if opts.UpdateHistory {
			query.Set("updateHistory", "true")
		}
		if len(query) > 0 {
			u += "?" + query.Encode()
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	issue := new(Issue)
	resp, err := s.client.Do(req, issue)
	if err != nil {
		return nil, resp, err
	}

	return issue, resp, nil
}

// IssueCreateRequest represents a request to create an issue.
type IssueCreateRequest struct {
	Fields          map[string]any    `json:"fields,omitempty"`
	Update          map[string]any    `json:"update,omitempty"`
	Transition      *TransitionInput  `json:"transition,omitempty"`
	HistoryMetadata *HistoryMetadata  `json:"historyMetadata,omitempty"`
	Properties      []*EntityProperty `json:"properties,omitempty"`
}

// TransitionInput represents a transition in a create/update request.
type TransitionInput struct {
	ID     string `json:"id,omitempty"`
	Looped bool   `json:"looped,omitempty"`
}

// HistoryMetadata represents history metadata.
type HistoryMetadata struct {
	Type                   string                      `json:"type,omitempty"`
	Description            string                      `json:"description,omitempty"`
	DescriptionKey         string                      `json:"descriptionKey,omitempty"`
	ActivityDescription    string                      `json:"activityDescription,omitempty"`
	ActivityDescriptionKey string                      `json:"activityDescriptionKey,omitempty"`
	EmailDescription       string                      `json:"emailDescription,omitempty"`
	EmailDescriptionKey    string                      `json:"emailDescriptionKey,omitempty"`
	Actor                  *HistoryMetadataParticipant `json:"actor,omitempty"`
	Generator              *HistoryMetadataParticipant `json:"generator,omitempty"`
	Cause                  *HistoryMetadataParticipant `json:"cause,omitempty"`
	ExtraData              map[string]string           `json:"extraData,omitempty"`
}

// HistoryMetadataParticipant represents a participant in history metadata.
type HistoryMetadataParticipant struct {
	ID             string `json:"id,omitempty"`
	DisplayName    string `json:"displayName,omitempty"`
	DisplayNameKey string `json:"displayNameKey,omitempty"`
	Type           string `json:"type,omitempty"`
	AvatarURL      string `json:"avatarUrl,omitempty"`
	URL            string `json:"url,omitempty"`
}

// IssueCreateResponse represents the response from creating an issue.
type IssueCreateResponse struct {
	ID         string            `json:"id,omitempty"`
	Key        string            `json:"key,omitempty"`
	Self       string            `json:"self,omitempty"`
	Transition *TransitionResult `json:"transition,omitempty"`
}

// TransitionResult represents a transition result.
type TransitionResult struct {
	Status          int              `json:"status,omitempty"`
	ErrorCollection *ErrorCollection `json:"errorCollection,omitempty"`
}

// ErrorCollection represents a collection of errors.
type ErrorCollection struct {
	ErrorMessages []string          `json:"errorMessages,omitempty"`
	Errors        map[string]string `json:"errors,omitempty"`
	Status        int               `json:"status,omitempty"`
}

// Create creates a new issue.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-post
func (s *IssuesService) Create(ctx context.Context, issue *IssueCreateRequest) (*IssueCreateResponse, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/issue", issue)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssueCreateResponse)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateBulk creates multiple issues in one request.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-bulk-post
func (s *IssuesService) CreateBulk(ctx context.Context, issues []*IssueCreateRequest) (*IssuesBulkResponse, *Response, error) {
	body := map[string]any{
		"issueUpdates": issues,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/issue/bulk", body)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssuesBulkResponse)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// IssuesBulkResponse represents the response from bulk issue creation.
type IssuesBulkResponse struct {
	Issues []*IssueCreateResponse `json:"issues,omitempty"`
	Errors []*BulkOperationError  `json:"errors,omitempty"`
}

// BulkOperationError represents an error from a bulk operation.
type BulkOperationError struct {
	Status              int              `json:"status,omitempty"`
	ElementErrors       *ErrorCollection `json:"elementErrors,omitempty"`
	FailedElementNumber int              `json:"failedElementNumber,omitempty"`
}

// IssueUpdateRequest represents a request to update an issue.
type IssueUpdateRequest struct {
	Fields          map[string]any    `json:"fields,omitempty"`
	Update          map[string]any    `json:"update,omitempty"`
	Transition      *TransitionInput  `json:"transition,omitempty"`
	HistoryMetadata *HistoryMetadata  `json:"historyMetadata,omitempty"`
	Properties      []*EntityProperty `json:"properties,omitempty"`
}

// Update updates an issue.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-put
func (s *IssuesService) Update(ctx context.Context, issueIDOrKey string, issue *IssueUpdateRequest, opts *IssueUpdateOptions) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s", issueIDOrKey)

	if opts != nil {
		query := url.Values{}
		if opts.NotifyUsers != nil {
			query.Set("notifyUsers", fmt.Sprintf("%t", *opts.NotifyUsers))
		}
		if opts.OverrideScreenSecurity {
			query.Set("overrideScreenSecurity", "true")
		}
		if opts.OverrideEditableFlag {
			query.Set("overrideEditableFlag", "true")
		}
		if opts.ReturnIssue {
			query.Set("returnIssue", "true")
		}
		if len(opts.Expand) > 0 {
			query.Set("expand", strings.Join(opts.Expand, ","))
		}
		if len(query) > 0 {
			u += "?" + query.Encode()
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, issue)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// IssueUpdateOptions specifies optional parameters for Update.
type IssueUpdateOptions struct {
	NotifyUsers            *bool    `url:"notifyUsers,omitempty"`
	OverrideScreenSecurity bool     `url:"overrideScreenSecurity,omitempty"`
	OverrideEditableFlag   bool     `url:"overrideEditableFlag,omitempty"`
	ReturnIssue            bool     `url:"returnIssue,omitempty"`
	Expand                 []string `url:"expand,omitempty"`
}

// Delete deletes an issue.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-delete
func (s *IssuesService) Delete(ctx context.Context, issueIDOrKey string, deleteSubtasks bool) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s", issueIDOrKey)
	if deleteSubtasks {
		u += "?deleteSubtasks=true"
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Assign assigns an issue to a user.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-assignee-put
func (s *IssuesService) Assign(ctx context.Context, issueIDOrKey, accountID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/assignee", issueIDOrKey)

	body := map[string]any{}
	if accountID != "" {
		body["accountId"] = accountID
	}
	// Empty body unassigns

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetTransitions returns the available transitions for an issue.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-transitions-get
func (s *IssuesService) GetTransitions(ctx context.Context, issueIDOrKey string, opts *GetTransitionsOptions) ([]*Transition, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueIDOrKey)

	if opts != nil {
		query := url.Values{}
		if opts.TransitionID != "" {
			query.Set("transitionId", opts.TransitionID)
		}
		if opts.SkipRemoteOnlyCondition {
			query.Set("skipRemoteOnlyCondition", "true")
		}
		if opts.IncludeUnavailableTransitions {
			query.Set("includeUnavailableTransitions", "true")
		}
		if opts.SortByOpsBarAndStatus {
			query.Set("sortByOpsBarAndStatus", "true")
		}
		if len(opts.Expand) > 0 {
			query.Set("expand", strings.Join(opts.Expand, ","))
		}
		if len(query) > 0 {
			u += "?" + query.Encode()
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Transitions []*Transition `json:"transitions"`
	}
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Transitions, resp, nil
}

// GetTransitionsOptions specifies optional parameters for GetTransitions.
type GetTransitionsOptions struct {
	TransitionID                  string   `url:"transitionId,omitempty"`
	SkipRemoteOnlyCondition       bool     `url:"skipRemoteOnlyCondition,omitempty"`
	IncludeUnavailableTransitions bool     `url:"includeUnavailableTransitions,omitempty"`
	SortByOpsBarAndStatus         bool     `url:"sortByOpsBarAndStatus,omitempty"`
	Expand                        []string `url:"expand,omitempty"`
}

// DoTransition performs a transition on an issue.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-transitions-post
func (s *IssuesService) DoTransition(ctx context.Context, issueIDOrKey string, transition *IssueTransitionRequest) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, transition)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// IssueTransitionRequest represents a request to transition an issue.
type IssueTransitionRequest struct {
	Transition      *TransitionInput  `json:"transition,omitempty"`
	Fields          map[string]any    `json:"fields,omitempty"`
	Update          map[string]any    `json:"update,omitempty"`
	HistoryMetadata *HistoryMetadata  `json:"historyMetadata,omitempty"`
	Properties      []*EntityProperty `json:"properties,omitempty"`
}

// GetChangelog returns the changelog for an issue.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-changelog-get
func (s *IssuesService) GetChangelog(ctx context.Context, issueIDOrKey string, opts *ChangelogOptions) (*Changelog, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/changelog", issueIDOrKey)

	if opts != nil {
		query := url.Values{}
		if opts.StartAt > 0 {
			query.Set("startAt", fmt.Sprintf("%d", opts.StartAt))
		}
		if opts.MaxResults > 0 {
			query.Set("maxResults", fmt.Sprintf("%d", opts.MaxResults))
		}
		if len(query) > 0 {
			u += "?" + query.Encode()
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	changelog := new(Changelog)
	resp, err := s.client.Do(req, changelog)
	if err != nil {
		return nil, resp, err
	}

	return changelog, resp, nil
}

// ChangelogOptions specifies optional parameters for GetChangelog.
type ChangelogOptions struct {
	StartAt    int `url:"startAt,omitempty"`
	MaxResults int `url:"maxResults,omitempty"`
}

// Notify sends a notification about an issue.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-notify-post
func (s *IssuesService) Notify(ctx context.Context, issueIDOrKey string, notification *Notification) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/notify", issueIDOrKey)

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, notification)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Notification represents a notification to send.
type Notification struct {
	Subject  string                  `json:"subject,omitempty"`
	TextBody string                  `json:"textBody,omitempty"`
	HTMLBody string                  `json:"htmlBody,omitempty"`
	To       *NotificationRecipients `json:"to,omitempty"`
	Restrict *NotificationRestrict   `json:"restrict,omitempty"`
}

// NotificationRecipients represents the recipients of a notification.
type NotificationRecipients struct {
	Reporter bool     `json:"reporter,omitempty"`
	Assignee bool     `json:"assignee,omitempty"`
	Watchers bool     `json:"watchers,omitempty"`
	Voters   bool     `json:"voters,omitempty"`
	Users    []*User  `json:"users,omitempty"`
	Groups   []*Group `json:"groups,omitempty"`
}

// NotificationRestrict represents restrictions on notifications.
type NotificationRestrict struct {
	Groups      []*Group                `json:"groups,omitempty"`
	Permissions []*RestrictedPermission `json:"permissions,omitempty"`
}

// RestrictedPermission represents a restricted permission.
type RestrictedPermission struct {
	ID  string `json:"id,omitempty"`
	Key string `json:"key,omitempty"`
}

// GetEditMeta returns the edit metadata for an issue.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-editmeta-get
func (s *IssuesService) GetEditMeta(ctx context.Context, issueIDOrKey string, opts *EditMetaOptions) (*EditMeta, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/editmeta", issueIDOrKey)

	if opts != nil {
		query := url.Values{}
		if opts.OverrideScreenSecurity {
			query.Set("overrideScreenSecurity", "true")
		}
		if opts.OverrideEditableFlag {
			query.Set("overrideEditableFlag", "true")
		}
		if len(query) > 0 {
			u += "?" + query.Encode()
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	meta := new(EditMeta)
	resp, err := s.client.Do(req, meta)
	if err != nil {
		return nil, resp, err
	}

	return meta, resp, nil
}

// EditMetaOptions specifies optional parameters for GetEditMeta.
type EditMetaOptions struct {
	OverrideScreenSecurity bool `url:"overrideScreenSecurity,omitempty"`
	OverrideEditableFlag   bool `url:"overrideEditableFlag,omitempty"`
}

// GetCreateMeta returns metadata for creating issues.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-createmeta-get
func (s *IssuesService) GetCreateMeta(ctx context.Context, opts *CreateMetaOptions) (*CreateMeta, *Response, error) {
	u := "/rest/api/3/issue/createmeta"

	if opts != nil {
		query := url.Values{}
		if len(opts.ProjectIDs) > 0 {
			query.Set("projectIds", strings.Join(opts.ProjectIDs, ","))
		}
		if len(opts.ProjectKeys) > 0 {
			query.Set("projectKeys", strings.Join(opts.ProjectKeys, ","))
		}
		if len(opts.IssueTypeIDs) > 0 {
			query.Set("issuetypeIds", strings.Join(opts.IssueTypeIDs, ","))
		}
		if len(opts.IssueTypeNames) > 0 {
			query.Set("issuetypeNames", strings.Join(opts.IssueTypeNames, ","))
		}
		if len(opts.Expand) > 0 {
			query.Set("expand", strings.Join(opts.Expand, ","))
		}
		if len(query) > 0 {
			u += "?" + query.Encode()
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	meta := new(CreateMeta)
	resp, err := s.client.Do(req, meta)
	if err != nil {
		return nil, resp, err
	}

	return meta, resp, nil
}

// CreateMetaOptions specifies optional parameters for GetCreateMeta.
type CreateMetaOptions struct {
	ProjectIDs     []string `url:"projectIds,omitempty"`
	ProjectKeys    []string `url:"projectKeys,omitempty"`
	IssueTypeIDs   []string `url:"issuetypeIds,omitempty"`
	IssueTypeNames []string `url:"issuetypeNames,omitempty"`
	Expand         []string `url:"expand,omitempty"`
}

// CreateMeta represents issue creation metadata.
type CreateMeta struct {
	Expand   string               `json:"expand,omitempty"`
	Projects []*CreateMetaProject `json:"projects,omitempty"`
}

// CreateMetaProject represents a project in create metadata.
type CreateMetaProject struct {
	Self       string                 `json:"self,omitempty"`
	ID         string                 `json:"id,omitempty"`
	Key        string                 `json:"key,omitempty"`
	Name       string                 `json:"name,omitempty"`
	AvatarURLs map[string]string      `json:"avatarUrls,omitempty"`
	IssueTypes []*CreateMetaIssueType `json:"issuetypes,omitempty"`
}

// CreateMetaIssueType represents an issue type in create metadata.
type CreateMetaIssueType struct {
	Self        string                `json:"self,omitempty"`
	ID          string                `json:"id,omitempty"`
	Description string                `json:"description,omitempty"`
	IconURL     string                `json:"iconUrl,omitempty"`
	Name        string                `json:"name,omitempty"`
	Subtask     bool                  `json:"subtask,omitempty"`
	AvatarID    int                   `json:"avatarId,omitempty"`
	Fields      map[string]*FieldMeta `json:"fields,omitempty"`
}

// Archive archives issues.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-archive-put
func (s *IssuesService) Archive(ctx context.Context, issueIDsOrKeys []string) (*Response, error) {
	body := map[string]any{
		"issueIdsOrKeys": issueIDsOrKeys,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, "/rest/api/3/issue/archive", body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Unarchive unarchives issues.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-unarchive-put
func (s *IssuesService) Unarchive(ctx context.Context, issueIDsOrKeys []string) (*Response, error) {
	body := map[string]any{
		"issueIdsOrKeys": issueIDsOrKeys,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, "/rest/api/3/issue/unarchive", body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
