package jira

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// MyselfService handles current user operations for the Jira API.
type MyselfService struct {
	client *Client
}

// Get returns the current user.
func (s *MyselfService) Get(ctx context.Context, expand []string) (*User, *Response, error) {
	u := "/rest/api/3/myself"

	if len(expand) > 0 {
		u = fmt.Sprintf("%s?expand=%s", u, strings.Join(expand, ","))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

// CurrentUserPreferences represents preferences for the current user.
type CurrentUserPreferences struct {
	Locale string `json:"locale,omitempty"`
}

// GetPreference returns a preference for the current user.
func (s *MyselfService) GetPreference(ctx context.Context, key string) (string, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/mypreferences?key=%s", key)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", nil, err
	}

	var result string
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return "", resp, err
	}

	return result, resp, nil
}

// SetPreference sets a preference for the current user.
func (s *MyselfService) SetPreference(ctx context.Context, key, value string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/mypreferences?key=%s", key)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, value)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeletePreference removes a preference for the current user.
func (s *MyselfService) DeletePreference(ctx context.Context, key string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/mypreferences?key=%s", key)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetLocale returns the locale for the current user.
func (s *MyselfService) GetLocale(ctx context.Context) (*CurrentUserPreferences, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/mypreferences/locale", nil)
	if err != nil {
		return nil, nil, err
	}

	prefs := new(CurrentUserPreferences)
	resp, err := s.client.Do(req, prefs)
	if err != nil {
		return nil, resp, err
	}

	return prefs, resp, nil
}

// SetLocale sets the locale for the current user.
func (s *MyselfService) SetLocale(ctx context.Context, locale string) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, "/rest/api/3/mypreferences/locale", map[string]string{"locale": locale})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteLocale removes the locale preference for the current user.
func (s *MyselfService) DeleteLocale(ctx context.Context) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, "/rest/api/3/mypreferences/locale", nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
