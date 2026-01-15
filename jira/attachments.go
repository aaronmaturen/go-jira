package jira

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// AttachmentsService handles attachment operations for the Jira API.
type AttachmentsService struct {
	client *Client
}

// Attachment represents a Jira attachment.
type Attachment struct {
	Self      string `json:"self,omitempty"`
	ID        string `json:"id,omitempty"`
	Filename  string `json:"filename,omitempty"`
	Author    *User  `json:"author,omitempty"`
	Created   *Time  `json:"created,omitempty"`
	Size      int64  `json:"size,omitempty"`
	MimeType  string `json:"mimeType,omitempty"`
	Content   string `json:"content,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

// AttachmentMeta represents attachment metadata.
type AttachmentMeta struct {
	Enabled     bool `json:"enabled,omitempty"`
	UploadLimit int  `json:"uploadLimit,omitempty"`
}

// GetMeta returns attachment metadata.
func (s *AttachmentsService) GetMeta(ctx context.Context) (*AttachmentMeta, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/attachment/meta", nil)
	if err != nil {
		return nil, nil, err
	}

	meta := new(AttachmentMeta)
	resp, err := s.client.Do(req, meta)
	if err != nil {
		return nil, resp, err
	}

	return meta, resp, nil
}

// Get returns an attachment by ID.
func (s *AttachmentsService) Get(ctx context.Context, attachmentID string) (*Attachment, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/attachment/%s", attachmentID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	attachment := new(Attachment)
	resp, err := s.client.Do(req, attachment)
	if err != nil {
		return nil, resp, err
	}

	return attachment, resp, nil
}

// Delete removes an attachment.
func (s *AttachmentsService) Delete(ctx context.Context, attachmentID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/attachment/%s", attachmentID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ExpandedContent represents expanded attachment content.
type ExpandedContent struct {
	ID              string             `json:"id,omitempty"`
	MediaType       string             `json:"mediaType,omitempty"`
	Name            string             `json:"name,omitempty"`
	Self            string             `json:"self,omitempty"`
	Entries         []*AttachmentEntry `json:"entries,omitempty"`
	TotalEntryCount int                `json:"totalEntryCount,omitempty"`
}

// AttachmentEntry represents an entry in an expanded attachment.
type AttachmentEntry struct {
	EntryIndex int    `json:"entryIndex,omitempty"`
	Name       string `json:"name,omitempty"`
	Size       int64  `json:"size,omitempty"`
	MediaType  string `json:"mediaType,omitempty"`
}

// Expand returns the contents of an attachment (for zip/tar files).
func (s *AttachmentsService) Expand(ctx context.Context, attachmentID string) (*ExpandedContent, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/attachment/%s/expand/human", attachmentID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	content := new(ExpandedContent)
	resp, err := s.client.Do(req, content)
	if err != nil {
		return nil, resp, err
	}

	return content, resp, nil
}

// ExpandRaw returns the raw contents of an attachment.
func (s *AttachmentsService) ExpandRaw(ctx context.Context, attachmentID string) (*ExpandedContent, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/attachment/%s/expand/raw", attachmentID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	content := new(ExpandedContent)
	resp, err := s.client.Do(req, content)
	if err != nil {
		return nil, resp, err
	}

	return content, resp, nil
}

// Download downloads an attachment.
func (s *AttachmentsService) Download(ctx context.Context, attachmentID string) (io.ReadCloser, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/attachment/content/%s", attachmentID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return resp.Body, newResponse(resp), nil
}

// AddToIssue adds attachments to an issue.
func (s *AttachmentsService) AddToIssue(ctx context.Context, issueIDOrKey string, files map[string]io.Reader) ([]*Attachment, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/issue/%s/attachments", issueIDOrKey)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for filename, reader := range files {
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			return nil, nil, fmt.Errorf("create form file: %w", err)
		}
		if _, err := io.Copy(part, reader); err != nil {
			return nil, nil, fmt.Errorf("copy file content: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, nil, fmt.Errorf("close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.client.baseURL.String()+u, &body)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Atlassian-Token", "no-check")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", s.client.UserAgent)

	if s.client.auth != nil {
		s.client.auth.Apply(req)
	}

	var attachments []*Attachment
	resp, err := s.client.Do(req, &attachments)
	if err != nil {
		return nil, resp, err
	}

	return attachments, resp, nil
}

// AddToIssueFromBytes adds a single attachment from bytes.
func (s *AttachmentsService) AddToIssueFromBytes(ctx context.Context, issueIDOrKey, filename string, data []byte) ([]*Attachment, *Response, error) {
	return s.AddToIssue(ctx, issueIDOrKey, map[string]io.Reader{
		filename: bytes.NewReader(data),
	})
}

// AttachmentSettings represents global attachment settings.
type AttachmentSettings struct {
	Enabled     bool `json:"enabled,omitempty"`
	UploadLimit int  `json:"uploadLimit,omitempty"`
}

// GetSettings returns global attachment settings.
func (s *AttachmentsService) GetSettings(ctx context.Context) (*AttachmentSettings, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/attachment/meta", nil)
	if err != nil {
		return nil, nil, err
	}

	settings := new(AttachmentSettings)
	resp, err := s.client.Do(req, settings)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, nil
}

// GetThumbnail returns the thumbnail for an attachment.
func (s *AttachmentsService) GetThumbnail(ctx context.Context, attachmentID string, width, height int, fallbackToDefault bool) (io.ReadCloser, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/attachment/thumbnail/%s", attachmentID)

	params := make(map[string]string)
	if width > 0 {
		params["width"] = fmt.Sprintf("%d", width)
	}
	if height > 0 {
		params["height"] = fmt.Sprintf("%d", height)
	}
	if fallbackToDefault {
		params["fallbackToDefault"] = "true"
	}

	if len(params) > 0 {
		first := true
		for k, v := range params {
			if first {
				u += "?"
				first = false
			} else {
				u += "&"
			}
			u += k + "=" + v
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return resp.Body, newResponse(resp), nil
}
