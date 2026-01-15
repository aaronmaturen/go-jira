package jira

import (
	"encoding/json"
	"time"
)

// Time is a wrapper around time.Time that handles Jira's date formats.
type Time struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for Time.
func (t *Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}

	// Try different formats
	formats := []string{
		"2006-01-02T15:04:05.000-0700",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	var parseErr error
	for _, format := range formats {
		parsed, err := time.Parse(format, s)
		if err == nil {
			t.Time = parsed
			return nil
		}
		parseErr = err
	}
	return parseErr
}

// MarshalJSON implements json.Marshaler for Time.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Format("2006-01-02T15:04:05.000-0700"))
}

// Date is a wrapper for date-only values.
type Date struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for Date.
func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

// MarshalJSON implements json.Marshaler for Date.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Format("2006-01-02"))
}

// User represents a Jira user.
type User struct {
	Self         string            `json:"self,omitempty"`
	AccountID    string            `json:"accountId,omitempty"`
	AccountType  string            `json:"accountType,omitempty"`
	EmailAddress string            `json:"emailAddress,omitempty"`
	AvatarURLs   map[string]string `json:"avatarUrls,omitempty"`
	DisplayName  string            `json:"displayName,omitempty"`
	Active       bool              `json:"active,omitempty"`
	TimeZone     string            `json:"timeZone,omitempty"`
	Locale       string            `json:"locale,omitempty"`
}

// Project represents a Jira project.
type Project struct {
	Self              string            `json:"self,omitempty"`
	ID                string            `json:"id,omitempty"`
	Key               string            `json:"key,omitempty"`
	Name              string            `json:"name,omitempty"`
	Description       string            `json:"description,omitempty"`
	Lead              *User             `json:"lead,omitempty"`
	Components        []*Component      `json:"components,omitempty"`
	IssueTypes        []*IssueType      `json:"issueTypes,omitempty"`
	URL               string            `json:"url,omitempty"`
	Email             string            `json:"email,omitempty"`
	AssigneeType      string            `json:"assigneeType,omitempty"`
	Versions          []*Version        `json:"versions,omitempty"`
	Roles             map[string]string `json:"roles,omitempty"`
	AvatarURLs        map[string]string `json:"avatarUrls,omitempty"`
	ProjectCategory   *ProjectCategory  `json:"projectCategory,omitempty"`
	ProjectTypeKey    string            `json:"projectTypeKey,omitempty"`
	Simplified        bool              `json:"simplified,omitempty"`
	Style             string            `json:"style,omitempty"`
	Favourite         bool              `json:"favourite,omitempty"`
	IsPrivate         bool              `json:"isPrivate,omitempty"`
	Properties        map[string]any    `json:"properties,omitempty"`
	UUID              string            `json:"uuid,omitempty"`
	Insight           *ProjectInsight   `json:"insight,omitempty"`
	Deleted           bool              `json:"deleted,omitempty"`
	RetentionTillDate string            `json:"retentionTillDate,omitempty"`
	DeletedDate       *Time             `json:"deletedDate,omitempty"`
	DeletedBy         *User             `json:"deletedBy,omitempty"`
	Archived          bool              `json:"archived,omitempty"`
	ArchivedDate      *Time             `json:"archivedDate,omitempty"`
	ArchivedBy        *User             `json:"archivedBy,omitempty"`
}

// ProjectInsight represents project insight data.
type ProjectInsight struct {
	TotalIssueCount     int   `json:"totalIssueCount,omitempty"`
	LastIssueUpdateTime *Time `json:"lastIssueUpdateTime,omitempty"`
}

// ProjectCategory represents a project category.
type ProjectCategory struct {
	Self        string `json:"self,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Component represents a project component.
type Component struct {
	Self                string `json:"self,omitempty"`
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	Lead                *User  `json:"lead,omitempty"`
	LeadAccountID       string `json:"leadAccountId,omitempty"`
	AssigneeType        string `json:"assigneeType,omitempty"`
	Assignee            *User  `json:"assignee,omitempty"`
	RealAssigneeType    string `json:"realAssigneeType,omitempty"`
	RealAssignee        *User  `json:"realAssignee,omitempty"`
	IsAssigneeTypeValid bool   `json:"isAssigneeTypeValid,omitempty"`
	Project             string `json:"project,omitempty"`
	ProjectID           int    `json:"projectId,omitempty"`
}

// IssueType represents an issue type.
type IssueType struct {
	Self           string `json:"self,omitempty"`
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	IconURL        string `json:"iconUrl,omitempty"`
	Subtask        bool   `json:"subtask,omitempty"`
	AvatarID       int    `json:"avatarId,omitempty"`
	EntityID       string `json:"entityId,omitempty"`
	HierarchyLevel int    `json:"hierarchyLevel,omitempty"`
	Scope          *Scope `json:"scope,omitempty"`
}

// Priority represents an issue priority.
type Priority struct {
	Self        string `json:"self,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"iconUrl,omitempty"`
	StatusColor string `json:"statusColor,omitempty"`
}

// Resolution represents an issue resolution.
type Resolution struct {
	Self        string `json:"self,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Status represents an issue status.
type Status struct {
	Self           string          `json:"self,omitempty"`
	ID             string          `json:"id,omitempty"`
	Name           string          `json:"name,omitempty"`
	Description    string          `json:"description,omitempty"`
	IconURL        string          `json:"iconUrl,omitempty"`
	StatusCategory *StatusCategory `json:"statusCategory,omitempty"`
	Scope          *Scope          `json:"scope,omitempty"`
}

// Transition represents a workflow transition.
type Transition struct {
	ID            string                `json:"id,omitempty"`
	Name          string                `json:"name,omitempty"`
	To            *Status               `json:"to,omitempty"`
	HasScreen     bool                  `json:"hasScreen,omitempty"`
	IsGlobal      bool                  `json:"isGlobal,omitempty"`
	IsInitial     bool                  `json:"isInitial,omitempty"`
	IsAvailable   bool                  `json:"isAvailable,omitempty"`
	IsConditional bool                  `json:"isConditional,omitempty"`
	Fields        map[string]*FieldMeta `json:"fields,omitempty"`
	IsLooped      bool                  `json:"isLooped,omitempty"`
}

// FieldMeta represents metadata about a field in a transition.
type FieldMeta struct {
	Required        bool     `json:"required,omitempty"`
	Schema          *Schema  `json:"schema,omitempty"`
	Name            string   `json:"name,omitempty"`
	Key             string   `json:"key,omitempty"`
	AutoCompleteURL string   `json:"autoCompleteUrl,omitempty"`
	HasDefaultValue bool     `json:"hasDefaultValue,omitempty"`
	Operations      []string `json:"operations,omitempty"`
	AllowedValues   []any    `json:"allowedValues,omitempty"`
	DefaultValue    any      `json:"defaultValue,omitempty"`
}

// Schema represents a field schema.
type Schema struct {
	Type     string `json:"type,omitempty"`
	Items    string `json:"items,omitempty"`
	System   string `json:"system,omitempty"`
	Custom   string `json:"custom,omitempty"`
	CustomID int    `json:"customId,omitempty"`
}

// AvatarURLs represents avatar URLs in different sizes.
type AvatarURLs struct {
	Size16x16 string `json:"16x16,omitempty"`
	Size24x24 string `json:"24x24,omitempty"`
	Size32x32 string `json:"32x32,omitempty"`
	Size48x48 string `json:"48x48,omitempty"`
}

// Comment represents an issue comment.
type Comment struct {
	Self         string            `json:"self,omitempty"`
	ID           string            `json:"id,omitempty"`
	Author       *User             `json:"author,omitempty"`
	Body         any               `json:"body,omitempty"` // Can be string or ADF
	RenderedBody string            `json:"renderedBody,omitempty"`
	UpdateAuthor *User             `json:"updateAuthor,omitempty"`
	Created      *Time             `json:"created,omitempty"`
	Updated      *Time             `json:"updated,omitempty"`
	Visibility   *Visibility       `json:"visibility,omitempty"`
	JsdPublic    bool              `json:"jsdPublic,omitempty"`
	Properties   []*EntityProperty `json:"properties,omitempty"`
}

// RemoteLink represents a remote link on an issue.
type RemoteLink struct {
	Self         string            `json:"self,omitempty"`
	ID           int               `json:"id,omitempty"`
	GlobalID     string            `json:"globalId,omitempty"`
	Application  *Application      `json:"application,omitempty"`
	Relationship string            `json:"relationship,omitempty"`
	Object       *RemoteLinkObject `json:"object,omitempty"`
}

// Application represents a remote application.
type Application struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

// RemoteLinkObject represents the object of a remote link.
type RemoteLinkObject struct {
	URL     string            `json:"url,omitempty"`
	Title   string            `json:"title,omitempty"`
	Summary string            `json:"summary,omitempty"`
	Icon    *Icon             `json:"icon,omitempty"`
	Status  *RemoteLinkStatus `json:"status,omitempty"`
}

// Icon represents an icon.
type Icon struct {
	URL16x16 string `json:"url16x16,omitempty"`
	Title    string `json:"title,omitempty"`
	Link     string `json:"link,omitempty"`
}

// RemoteLinkStatus represents the status of a remote link.
type RemoteLinkStatus struct {
	Resolved bool  `json:"resolved,omitempty"`
	Icon     *Icon `json:"icon,omitempty"`
}

// Watches represents watchers on an issue.
type Watches struct {
	Self       string  `json:"self,omitempty"`
	WatchCount int     `json:"watchCount,omitempty"`
	IsWatching bool    `json:"isWatching,omitempty"`
	Watchers   []*User `json:"watchers,omitempty"`
}

// Actor represents an actor in a project role.
type Actor struct {
	ID          int    `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	ActorUser   *User  `json:"actorUser,omitempty"`
	ActorGroup  *Group `json:"actorGroup,omitempty"`
}

// Subscription represents a filter subscription.
type Subscription struct {
	ID    int    `json:"id,omitempty"`
	User  *User  `json:"user,omitempty"`
	Group *Group `json:"group,omitempty"`
}

// Label represents a label.
type Label struct {
	Label string `json:"label,omitempty"`
}

// WorkflowOperations represents available operations on a workflow.
type WorkflowOperations struct {
	CanEdit   bool `json:"canEdit,omitempty"`
	CanDelete bool `json:"canDelete,omitempty"`
}
