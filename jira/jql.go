package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// JQLService handles JQL operations for the Jira API.
type JQLService struct {
	client *Client
}

// AutocompleteData represents JQL autocomplete suggestions.
type AutocompleteData struct {
	VisibleFieldNames    []*FieldRef    `json:"visibleFieldNames,omitempty"`
	VisibleFunctionNames []*FunctionRef `json:"visibleFunctionNames,omitempty"`
	JQLReservedWords     []string       `json:"jqlReservedWords,omitempty"`
}

// FieldRef represents a field reference for autocomplete.
type FieldRef struct {
	Value       string   `json:"value,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	Auto        string   `json:"auto,omitempty"`
	Orderable   string   `json:"orderable,omitempty"`
	Searchable  string   `json:"searchable,omitempty"`
	CFID        string   `json:"cfid,omitempty"`
	Operators   []string `json:"operators,omitempty"`
	Types       []string `json:"types,omitempty"`
}

// FunctionRef represents a function reference for autocomplete.
type FunctionRef struct {
	Value       string   `json:"value,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	IsList      string   `json:"isList,omitempty"`
	Types       []string `json:"types,omitempty"`
}

// GetAutocompleteData returns JQL autocomplete data.
func (s *JQLService) GetAutocompleteData(ctx context.Context) (*AutocompleteData, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/jql/autocompletedata", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(AutocompleteData)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// AutocompleteSuggestion represents a JQL autocomplete suggestion.
type AutocompleteSuggestion struct {
	Value       string `json:"value,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// AutocompleteSuggestionsResult represents autocomplete suggestions.
type AutocompleteSuggestionsResult struct {
	Results []*AutocompleteSuggestion `json:"results,omitempty"`
}

// GetAutocompleteSuggestions returns JQL autocomplete suggestions.
func (s *JQLService) GetAutocompleteSuggestions(ctx context.Context, fieldName, fieldValue, predicateName, predicateValue string) (*AutocompleteSuggestionsResult, *Response, error) {
	u := "/rest/api/3/jql/autocompletedata/suggestions"

	params := url.Values{}
	if fieldName != "" {
		params.Set("fieldName", fieldName)
	}
	if fieldValue != "" {
		params.Set("fieldValue", fieldValue)
	}
	if predicateName != "" {
		params.Set("predicateName", predicateName)
	}
	if predicateValue != "" {
		params.Set("predicateValue", predicateValue)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(AutocompleteSuggestionsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// FieldReferenceData represents field reference data for JQL.
type FieldReferenceData struct {
	Value       string   `json:"value,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	Orderable   string   `json:"orderable,omitempty"`
	Searchable  string   `json:"searchable,omitempty"`
	Auto        string   `json:"auto,omitempty"`
	CFID        string   `json:"cfid,omitempty"`
	Operators   []string `json:"operators,omitempty"`
	Types       []string `json:"types,omitempty"`
}

// GetFieldReferenceData returns field reference data for JQL.
func (s *JQLService) GetFieldReferenceData(ctx context.Context) ([]*FieldReferenceData, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/rest/api/3/jql/autocompletedata/fields", nil)
	if err != nil {
		return nil, nil, err
	}

	var result []*FieldReferenceData
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ParsedJQL represents a parsed JQL query.
type ParsedJQL struct {
	JQL       string        `json:"jql,omitempty"`
	Structure *JQLStructure `json:"structure,omitempty"`
	Errors    []string      `json:"errors,omitempty"`
}

// JQLStructure represents the structure of a parsed JQL query.
type JQLStructure struct {
	Where   *JQLClause    `json:"where,omitempty"`
	OrderBy []*JQLOrderBy `json:"orderBy,omitempty"`
}

// JQLClause represents a clause in a JQL query.
type JQLClause struct {
	Type      string       `json:"type,omitempty"` // field, compound, not
	Field     *JQLField    `json:"field,omitempty"`
	Operator  string       `json:"operator,omitempty"`
	Operand   *JQLOperand  `json:"operand,omitempty"`
	Clauses   []*JQLClause `json:"clauses,omitempty"`
	Predicate *JQLClause   `json:"predicate,omitempty"`
}

// JQLField represents a field in a JQL clause.
type JQLField struct {
	Name     string      `json:"name,omitempty"`
	Property []*JQLField `json:"property,omitempty"`
}

// JQLOperand represents an operand in a JQL clause.
type JQLOperand struct {
	Type         string        `json:"type,omitempty"` // value, list, function, keyword
	Value        string        `json:"value,omitempty"`
	Values       []*JQLOperand `json:"values,omitempty"`
	Function     string        `json:"function,omitempty"`
	Arguments    []string      `json:"arguments,omitempty"`
	Keyword      string        `json:"keyword,omitempty"`
	EncodedValue string        `json:"encodedValue,omitempty"`
}

// JQLOrderBy represents an order by clause.
type JQLOrderBy struct {
	Field     *JQLField `json:"field,omitempty"`
	Direction string    `json:"direction,omitempty"`
}

// ParseJQLRequest represents a request to parse JQL queries.
type ParseJQLRequest struct {
	Queries []string `json:"queries"`
}

// ParseJQLResult represents the result of parsing JQL queries.
type ParseJQLResult struct {
	Queries []*ParsedJQL `json:"queries,omitempty"`
}

// Parse parses JQL queries.
func (s *JQLService) Parse(ctx context.Context, queries []string, validation string) (*ParseJQLResult, *Response, error) {
	u := "/rest/api/3/jql/parse"

	if validation != "" {
		u = fmt.Sprintf("%s?validation=%s", u, validation)
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, &ParseJQLRequest{Queries: queries})
	if err != nil {
		return nil, nil, err
	}

	result := new(ParseJQLResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ConvertedJQL represents a converted JQL query.
type ConvertedJQL struct {
	JQL   string `json:"jql,omitempty"`
	Error string `json:"error,omitempty"`
}

// ConvertJQLResult represents the result of converting JQL queries.
type ConvertJQLResult struct {
	QueryStrings []*ConvertedJQL `json:"queryStrings,omitempty"`
}

// ConvertToIDs converts JQL queries to use IDs instead of keys/names.
func (s *JQLService) ConvertToIDs(ctx context.Context, queries []string) (*ConvertJQLResult, *Response, error) {
	u := "/rest/api/3/jql/pdcleaner"

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, map[string][]string{"queryStrings": queries})
	if err != nil {
		return nil, nil, err
	}

	result := new(ConvertJQLResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SanitizeJQLRequest represents a request to sanitize JQL queries.
type SanitizeJQLRequest struct {
	Queries []*SanitizeJQLInput `json:"queries"`
}

// SanitizeJQLInput represents input for sanitizing a JQL query.
type SanitizeJQLInput struct {
	AccountID string `json:"accountId,omitempty"`
	Query     string `json:"query,omitempty"`
}

// SanitizedJQL represents a sanitized JQL query.
type SanitizedJQL struct {
	InitialQuery   string          `json:"initialQuery,omitempty"`
	SanitizedQuery string          `json:"sanitizedQuery,omitempty"`
	Errors         *SanitizeErrors `json:"errors,omitempty"`
	AccountID      string          `json:"accountId,omitempty"`
}

// SanitizeErrors represents errors from sanitization.
type SanitizeErrors struct {
	ErrorMessages []string          `json:"errorMessages,omitempty"`
	Errors        map[string]string `json:"errors,omitempty"`
}

// SanitizeJQLResult represents the result of sanitizing JQL queries.
type SanitizeJQLResult struct {
	Queries []*SanitizedJQL `json:"queries,omitempty"`
}

// Sanitize sanitizes JQL queries for a user.
func (s *JQLService) Sanitize(ctx context.Context, queries []*SanitizeJQLInput) (*SanitizeJQLResult, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/jql/sanitize", &SanitizeJQLRequest{Queries: queries})
	if err != nil {
		return nil, nil, err
	}

	result := new(SanitizeJQLResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// FunctionPrecomputation represents a precomputed JQL function.
type FunctionPrecomputation struct {
	Arguments    []string `json:"arguments,omitempty"`
	Created      string   `json:"created,omitempty"`
	Error        string   `json:"error,omitempty"`
	Field        string   `json:"field,omitempty"`
	FunctionKey  string   `json:"functionKey,omitempty"`
	FunctionName string   `json:"functionName,omitempty"`
	ID           string   `json:"id,omitempty"`
	Operator     string   `json:"operator,omitempty"`
	Updated      string   `json:"updated,omitempty"`
	Value        string   `json:"value,omitempty"`
}

// FunctionPrecomputationsResult represents function precomputations.
type FunctionPrecomputationsResult struct {
	NextPageToken string                    `json:"nextPageToken,omitempty"`
	Self          string                    `json:"self,omitempty"`
	Values        []*FunctionPrecomputation `json:"values,omitempty"`
}

// GetFunctionPrecomputations returns JQL function precomputations.
func (s *JQLService) GetFunctionPrecomputations(ctx context.Context, functionKey []string, startAt, maxResults int, orderBy string, filter string) (*FunctionPrecomputationsResult, *Response, error) {
	u := "/rest/api/3/jql/function/computation"

	params := url.Values{}
	for _, fk := range functionKey {
		params.Add("functionKey", fk)
	}
	if startAt > 0 {
		params.Set("startAt", fmt.Sprintf("%d", startAt))
	}
	if maxResults > 0 {
		params.Set("maxResults", fmt.Sprintf("%d", maxResults))
	}
	if orderBy != "" {
		params.Set("orderBy", orderBy)
	}
	if filter != "" {
		params.Set("filter", filter)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(FunctionPrecomputationsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateFunctionPrecomputations updates JQL function precomputations.
func (s *JQLService) UpdateFunctionPrecomputations(ctx context.Context, values []*FunctionPrecomputation) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/jql/function/computation", map[string][]*FunctionPrecomputation{"values": values})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// MatchIssuesRequest represents a request to match issues against JQL.
type MatchIssuesRequest struct {
	IssueIDs []int64  `json:"issueIds"`
	JQLs     []string `json:"jqls"`
}

// MatchIssuesResult represents the result of matching issues against JQL.
type MatchIssuesResult struct {
	Matches []*MatchEntry `json:"matches,omitempty"`
}

// Match matches issues against JQL queries.
func (s *JQLService) Match(ctx context.Context, issueIDs []int64, jqls []string) (*MatchIssuesResult, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/jql/match", &MatchIssuesRequest{
		IssueIDs: issueIDs,
		JQLs:     jqls,
	})
	if err != nil {
		return nil, nil, err
	}

	result := new(MatchIssuesResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// MigrateJQLRequest represents a request to migrate JQL queries.
type MigrateJQLRequest struct {
	QueryStrings []string `json:"queryStrings"`
}

// MigrateJQLResult represents the result of migrating JQL queries.
type MigrateJQLResult struct {
	QueryStrings []*MigratedJQL `json:"queryStrings,omitempty"`
}

// MigratedJQL represents a migrated JQL query.
type MigratedJQL struct {
	OriginalQuery string `json:"originalQuery,omitempty"`
	MigratedQuery string `json:"migratedQuery,omitempty"`
}

// Migrate migrates JQL queries to use account IDs.
func (s *JQLService) Migrate(ctx context.Context, queries []string) (*MigrateJQLResult, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/rest/api/3/jql/pdcleaner", &MigrateJQLRequest{
		QueryStrings: queries,
	})
	if err != nil {
		return nil, nil, err
	}

	result := new(MigrateJQLResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetFieldAutocompleteSuggestions returns autocomplete suggestions for a field.
func (s *JQLService) GetFieldAutocompleteSuggestions(ctx context.Context, fieldName, fieldValue string) (*AutocompleteSuggestionsResult, *Response, error) {
	u := "/rest/api/3/jql/autocompletedata/suggestions"

	params := url.Values{}
	params.Set("fieldName", fieldName)
	if fieldValue != "" {
		params.Set("fieldValue", fieldValue)
	}
	u = fmt.Sprintf("%s?%s", u, params.Encode())

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(AutocompleteSuggestionsResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetVisibleFields returns fields visible in JQL for the user.
func (s *JQLService) GetVisibleFields(ctx context.Context, projectKey string, issueTypeID string) ([]*FieldReferenceData, *Response, error) {
	u := "/rest/api/3/jql/autocompletedata/fields"

	params := url.Values{}
	if projectKey != "" {
		params.Set("projectKey", projectKey)
	}
	if issueTypeID != "" {
		params.Set("issueTypeId", issueTypeID)
	}
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result []*FieldReferenceData
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ValidateJQL validates a JQL query.
func (s *JQLService) ValidateJQL(ctx context.Context, jql string) (bool, []string, *Response, error) {
	result, resp, err := s.Parse(ctx, []string{jql}, "strict")
	if err != nil {
		return false, nil, resp, err
	}

	if len(result.Queries) > 0 && len(result.Queries[0].Errors) > 0 {
		return false, result.Queries[0].Errors, resp, nil
	}

	return true, nil, resp, nil
}
