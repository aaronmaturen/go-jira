package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// WorkflowsService handles workflow operations for the Jira API.
type WorkflowsService struct {
	client *Client
}

// Workflow represents a Jira workflow.
type Workflow struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	IsDefault   bool              `json:"isDefault,omitempty"`
	Scope       *Scope            `json:"scope,omitempty"`
	Transitions []*WorkflowTransition `json:"transitions,omitempty"`
	Statuses    []*WorkflowStatus `json:"statuses,omitempty"`
}

// WorkflowTransition represents a transition in a workflow.
type WorkflowTransition struct {
	ID          string             `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	From        []string           `json:"from,omitempty"`
	To          string             `json:"to,omitempty"`
	Type        string             `json:"type,omitempty"`
	Screen      *TransitionScreen  `json:"screen,omitempty"`
	Rules       *TransitionRules   `json:"rules,omitempty"`
	Properties  map[string]string  `json:"properties,omitempty"`
}

// TransitionScreen represents a screen for a transition.
type TransitionScreen struct {
	ID string `json:"id,omitempty"`
}

// TransitionRules represents rules for a transition.
type TransitionRules struct {
	Conditions      []*WorkflowCondition `json:"conditions,omitempty"`
	Validators      []*WorkflowValidator `json:"validators,omitempty"`
	PostFunctions   []*WorkflowFunction  `json:"postFunctions,omitempty"`
	ConditionGroups []*ConditionGroup    `json:"conditionGroups,omitempty"`
}

// WorkflowCondition represents a condition in a workflow.
type WorkflowCondition struct {
	Type          string            `json:"type,omitempty"`
	Configuration map[string]string `json:"configuration,omitempty"`
}

// WorkflowValidator represents a validator in a workflow.
type WorkflowValidator struct {
	Type          string            `json:"type,omitempty"`
	Configuration map[string]string `json:"configuration,omitempty"`
}

// WorkflowFunction represents a post function in a workflow.
type WorkflowFunction struct {
	Type          string            `json:"type,omitempty"`
	Configuration map[string]string `json:"configuration,omitempty"`
}

// ConditionGroup represents a group of conditions.
type ConditionGroup struct {
	Operation  string               `json:"operation,omitempty"`
	Conditions []*WorkflowCondition `json:"conditions,omitempty"`
	Groups     []*ConditionGroup    `json:"conditionGroups,omitempty"`
}

// WorkflowStatus represents a status in a workflow.
type WorkflowStatus struct {
	ID         string              `json:"id,omitempty"`
	Name       string              `json:"name,omitempty"`
	Properties map[string]string   `json:"properties,omitempty"`
}

// WorkflowListResult represents a paginated list of workflows.
type WorkflowListResult struct {
	Self       string      `json:"self,omitempty"`
	NextPage   string      `json:"nextPage,omitempty"`
	MaxResults int         `json:"maxResults,omitempty"`
	StartAt    int         `json:"startAt,omitempty"`
	Total      int         `json:"total,omitempty"`
	IsLast     bool        `json:"isLast,omitempty"`
	Values     []*Workflow `json:"values,omitempty"`
}

// List returns all workflows.
func (s *WorkflowsService) List(ctx context.Context, startAt, maxResults int, workflowName []string, expand, queryString, orderBy string, isActive bool) (*WorkflowListResult, *Response, error) {
	u := "/rest/api/3/workflow/search"

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	for _, name := range workflowName {
		params.Add("workflowName", name)
	}
	if expand != "" {
		params.Set("expand", expand)
	}
	if queryString != "" {
		params.Set("queryString", queryString)
	}
	if orderBy != "" {
		params.Set("orderBy", orderBy)
	}
	if isActive {
		params.Set("isActive", "true")
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(WorkflowListResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Get returns a workflow by ID.
func (s *WorkflowsService) Get(ctx context.Context, workflowID string, expand string) (*Workflow, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflow/%s", workflowID)

	if expand != "" {
		u = fmt.Sprintf("%s?expand=%s", u, expand)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	workflow := new(Workflow)
	resp, err := s.client.Do(req, workflow)
	if err != nil {
		return nil, resp, err
	}

	return workflow, resp, nil
}

// WorkflowCreateRequest represents a request to create a workflow.
type WorkflowCreateRequest struct {
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	Transitions []*WorkflowTransitionCreate `json:"transitions,omitempty"`
	Statuses    []*WorkflowStatusCreate `json:"statuses,omitempty"`
}

// WorkflowTransitionCreate represents a transition for workflow creation.
type WorkflowTransitionCreate struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	From        []string          `json:"from,omitempty"`
	To          string            `json:"to"`
	Type        string            `json:"type,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
}

// WorkflowStatusCreate represents a status for workflow creation.
type WorkflowStatusCreate struct {
	ID         string            `json:"id,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// Create creates a workflow.
func (s *WorkflowsService) Create(ctx context.Context, workflow *WorkflowCreateRequest) (*Workflow, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/workflow", workflow)
	if err != nil {
		return nil, nil, err
	}

	result := new(Workflow)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete removes a workflow.
func (s *WorkflowsService) Delete(ctx context.Context, workflowID string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflow/%s", workflowID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetTransitionProperties returns properties for a workflow transition.
func (s *WorkflowsService) GetTransitionProperties(ctx context.Context, transitionID int64, workflowName string, includeReservedKeys bool, key string, workflowMode string) (map[string]string, *Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflow/transitions/%d/properties", transitionID)

	params := url.Values{}
	params.Set("workflowName", workflowName)
	if includeReservedKeys {
		params.Set("includeReservedKeys", "true")
	}
	if key != "" {
		params.Set("key", key)
	}
	if workflowMode != "" {
		params.Set("workflowMode", workflowMode)
	}

	u = fmt.Sprintf("%s?%s", u, params.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var props map[string]string
	resp, err := s.client.Do(req, &props)
	if err != nil {
		return nil, resp, err
	}

	return props, resp, nil
}

// SetTransitionProperty sets a property for a workflow transition.
func (s *WorkflowsService) SetTransitionProperty(ctx context.Context, transitionID int64, key, workflowName string, value string, workflowMode string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflow/transitions/%d/properties", transitionID)

	params := url.Values{}
	params.Set("key", key)
	params.Set("workflowName", workflowName)
	if workflowMode != "" {
		params.Set("workflowMode", workflowMode)
	}

	u = fmt.Sprintf("%s?%s", u, params.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, map[string]string{"value": value})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteTransitionProperty removes a property from a workflow transition.
func (s *WorkflowsService) DeleteTransitionProperty(ctx context.Context, transitionID int64, key, workflowName string, workflowMode string) (*Response, error) {
	u := fmt.Sprintf("/rest/api/3/workflow/transitions/%d/properties", transitionID)

	params := url.Values{}
	params.Set("key", key)
	params.Set("workflowName", workflowName)
	if workflowMode != "" {
		params.Set("workflowMode", workflowMode)
	}

	u = fmt.Sprintf("%s?%s", u, params.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// TransitionRule represents a transition rule.
type TransitionRule struct {
	ID            string            `json:"id,omitempty"`
	Type          string            `json:"type,omitempty"`
	Configuration map[string]string `json:"configuration,omitempty"`
}

// TransitionRulesResult represents rules for a workflow transition.
type TransitionRulesResult struct {
	WorkflowID       string           `json:"workflowId,omitempty"`
	Rules            []*TransitionRule `json:"rules,omitempty"`
}

// GetTransitionRuleConfigurations returns rule configurations for transitions.
func (s *WorkflowsService) GetTransitionRuleConfigurations(ctx context.Context, startAt, maxResults int, types []string, keys []string, workflowNames []string, withTags []string, draft bool, expand []string) ([]*TransitionRulesResult, *Response, error) {
	u := "/rest/api/3/workflow/rule/config"

	params := url.Values{}
	if startAt > 0 {
		params.Set("startAt", strconv.Itoa(startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	for _, t := range types {
		params.Add("types", t)
	}
	for _, k := range keys {
		params.Add("keys", k)
	}
	for _, wn := range workflowNames {
		params.Add("workflowNames", wn)
	}
	for _, t := range withTags {
		params.Add("withTags", t)
	}
	if draft {
		params.Set("draft", "true")
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

	var results []*TransitionRulesResult
	resp, err := s.client.Do(req, &results)
	if err != nil {
		return nil, resp, err
	}

	return results, resp, nil
}
