package jira

import (
	"context"
	"net/http"
)

// ServerInfoService handles server info operations for the Jira API.
type ServerInfoService struct {
	client *Client
}

// ServerInfo represents Jira server information.
type ServerInfo struct {
	BaseURL        string         `json:"baseUrl,omitempty"`
	Version        string         `json:"version,omitempty"`
	VersionNumbers []int          `json:"versionNumbers,omitempty"`
	DeploymentType string         `json:"deploymentType,omitempty"`
	BuildNumber    int            `json:"buildNumber,omitempty"`
	BuildDate      string         `json:"buildDate,omitempty"`
	ServerTime     string         `json:"serverTime,omitempty"`
	ScmInfo        string         `json:"scmInfo,omitempty"`
	ServerTitle    string         `json:"serverTitle,omitempty"`
	HealthChecks   []*HealthCheck `json:"healthChecks,omitempty"`
}

// HealthCheck represents a server health check.
type HealthCheck struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Passed      bool   `json:"passed,omitempty"`
}

// Get returns server information.
func (s *ServerInfoService) Get(ctx context.Context) (*ServerInfo, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/serverInfo", nil)
	if err != nil {
		return nil, nil, err
	}

	info := new(ServerInfo)
	resp, err := s.client.Do(req, info)
	if err != nil {
		return nil, resp, err
	}

	return info, resp, nil
}
