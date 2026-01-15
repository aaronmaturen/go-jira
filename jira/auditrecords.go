package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// AuditRecordsService handles audit record operations for the Jira API.
type AuditRecordsService struct {
	client *Client
}

// AuditRecord represents an audit record.
type AuditRecord struct {
	ID              int64                `json:"id,omitempty"`
	Summary         string               `json:"summary,omitempty"`
	RemoteAddress   string               `json:"remoteAddress,omitempty"`
	AuthorKey       string               `json:"authorKey,omitempty"`
	AuthorAccountID string               `json:"authorAccountId,omitempty"`
	Created         string               `json:"created,omitempty"`
	Category        string               `json:"category,omitempty"`
	EventSource     string               `json:"eventSource,omitempty"`
	Description     string               `json:"description,omitempty"`
	ObjectItem      *AssociatedItem      `json:"objectItem,omitempty"`
	ChangedValues   []*ChangedValue      `json:"changedValues,omitempty"`
	AssociatedItems []*AssociatedItem    `json:"associatedItems,omitempty"`
}

// AssociatedItem represents an item associated with an audit record.
type AssociatedItem struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	TypeName   string `json:"typeName,omitempty"`
	ParentID   string `json:"parentId,omitempty"`
	ParentName string `json:"parentName,omitempty"`
}

// ChangedValue represents a changed value in an audit record.
type ChangedValue struct {
	FieldName   string `json:"fieldName,omitempty"`
	ChangedFrom string `json:"changedFrom,omitempty"`
	ChangedTo   string `json:"changedTo,omitempty"`
}

// AuditRecordsResult represents a paginated list of audit records.
type AuditRecordsResult struct {
	Offset  int            `json:"offset,omitempty"`
	Limit   int            `json:"limit,omitempty"`
	Total   int            `json:"total,omitempty"`
	Records []*AuditRecord `json:"records,omitempty"`
}

// ListOptions specifies options for listing audit records.
type AuditRecordsListOptions struct {
	Offset int    `url:"offset,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Filter string `url:"filter,omitempty"`
	From   string `url:"from,omitempty"` // Date format: yyyy-MM-dd
	To     string `url:"to,omitempty"`   // Date format: yyyy-MM-dd
}

// List returns audit records.
func (s *AuditRecordsService) List(ctx context.Context, opts *AuditRecordsListOptions) (*AuditRecordsResult, *Response, error) {
	u := "/rest/api/3/auditing/record"

	if opts != nil {
		params := url.Values{}
		if opts.Offset > 0 {
			params.Set("offset", strconv.Itoa(opts.Offset))
		}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Filter != "" {
			params.Set("filter", opts.Filter)
		}
		if opts.From != "" {
			params.Set("from", opts.From)
		}
		if opts.To != "" {
			params.Set("to", opts.To)
		}
		if len(params) > 0 {
			u = fmt.Sprintf("%s?%s", u, params.Encode())
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(AuditRecordsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
