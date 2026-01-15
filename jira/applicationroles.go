package jira

import (
	"context"
	"fmt"
	"net/http"
)

// ApplicationRolesService handles application role operations for the Jira API.
type ApplicationRolesService struct {
	client *Client
}

// ApplicationRole represents an application role.
type ApplicationRole struct {
	Key                  string   `json:"key,omitempty"`
	Groups               []string `json:"groups,omitempty"`
	GroupDetails         []*Group `json:"groupDetails,omitempty"`
	Name                 string   `json:"name,omitempty"`
	DefaultGroups        []string `json:"defaultGroups,omitempty"`
	DefaultGroupsDetails []*Group `json:"defaultGroupsDetails,omitempty"`
	SelectedByDefault    bool     `json:"selectedByDefault,omitempty"`
	Defined              bool     `json:"defined,omitempty"`
	NumberOfSeats        int      `json:"numberOfSeats,omitempty"`
	RemainingSeats       int      `json:"remainingSeats,omitempty"`
	UserCount            int      `json:"userCount,omitempty"`
	UserCountDescription string   `json:"userCountDescription,omitempty"`
	HasUnlimitedSeats    bool     `json:"hasUnlimitedSeats,omitempty"`
	Platform             bool     `json:"platform,omitempty"`
}

// ListAll returns all application roles.
func (s *ApplicationRolesService) ListAll(ctx context.Context) ([]*ApplicationRole, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/applicationrole", nil)
	if err != nil {
		return nil, nil, err
	}

	var roles []*ApplicationRole
	resp, err := s.client.Do(req, &roles)
	if err != nil {
		return nil, resp, err
	}

	return roles, resp, nil
}

// Get returns an application role by key.
func (s *ApplicationRolesService) Get(ctx context.Context, key string) (*ApplicationRole, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/applicationrole/%s", key)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	role := new(ApplicationRole)
	resp, err := s.client.Do(req, role)
	if err != nil {
		return nil, resp, err
	}

	return role, resp, nil
}
