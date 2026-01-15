package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchService_Do(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodGet)
		}

		jql := r.URL.Query().Get("jql")
		if jql != "project = TEST" {
			t.Errorf("JQL = %v, want %v", jql, "project = TEST")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResult{
			Total: 2,
			Issues: []*Issue{
				{Key: "TEST-1", Fields: &IssueFields{Summary: "Issue 1"}},
				{Key: "TEST-2", Fields: &IssueFields{Summary: "Issue 2"}},
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	result, _, err := client.Search.Do(context.Background(), "project = TEST", nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if result.Total != 2 {
		t.Errorf("Total = %v, want %v", result.Total, 2)
	}
	if len(result.Issues) != 2 {
		t.Errorf("len(Issues) = %v, want %v", len(result.Issues), 2)
	}
}

func TestSearchService_Do_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		maxResults := r.URL.Query().Get("maxResults")
		if maxResults != "50" {
			t.Errorf("maxResults = %v, want %v", maxResults, "50")
		}

		// Fields are passed as separate query params
		fields := r.URL.Query()["fields"]
		if len(fields) != 2 {
			t.Errorf("len(fields) = %v, want %v", len(fields), 2)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResult{Issues: []*Issue{}})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	_, _, err := client.Search.Do(context.Background(), "project = TEST", &SearchOptions{
		MaxResults: 50,
		Fields:     []string{"summary", "status"},
	})
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
}

func TestSearchService_DoPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodPost)
		}

		var req SearchRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.JQL != "project = TEST" {
			t.Errorf("JQL = %v, want %v", req.JQL, "project = TEST")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResult{
			Total:  1,
			Issues: []*Issue{{Key: "TEST-1"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	result, _, err := client.Search.DoPost(context.Background(), &SearchRequest{
		JQL: "project = TEST",
	})
	if err != nil {
		t.Fatalf("DoPost() error = %v", err)
	}
	if result.Total != 1 {
		t.Errorf("Total = %v, want %v", result.Total, 1)
	}
}

func TestSearchService_Picker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/picker" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/issue/picker")
		}

		query := r.URL.Query().Get("query")
		if query != "test" {
			t.Errorf("query = %v, want %v", query, "test")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PickerSuggestions{
			Sections: []*PickerSection{
				{
					ID:    "cs",
					Label: "Current Search",
					Issues: []*PickerIssue{
						{Key: "TEST-1", Summary: "Test issue"},
					},
				},
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	result, _, err := client.Search.Picker(context.Background(), &PickerOptions{
		Query: "test",
	})
	if err != nil {
		t.Fatalf("Picker() error = %v", err)
	}
	if len(result.Sections) != 1 {
		t.Errorf("len(Sections) = %v, want %v", len(result.Sections), 1)
	}
}

func TestSearchResult_JSON(t *testing.T) {
	input := `{
		"expand": "names,schema",
		"startAt": 0,
		"maxResults": 50,
		"total": 125,
		"issues": [
			{
				"key": "TEST-1",
				"id": "10001",
				"fields": {
					"summary": "Test issue 1",
					"status": {"name": "To Do"}
				}
			},
			{
				"key": "TEST-2",
				"id": "10002",
				"fields": {
					"summary": "Test issue 2",
					"status": {"name": "In Progress"}
				}
			}
		]
	}`

	var result SearchResult
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if result.Total != 125 {
		t.Errorf("Total = %v, want %v", result.Total, 125)
	}
	if result.MaxResults != 50 {
		t.Errorf("MaxResults = %v, want %v", result.MaxResults, 50)
	}
	if len(result.Issues) != 2 {
		t.Errorf("len(Issues) = %v, want %v", len(result.Issues), 2)
	}
	if result.Issues[0].Key != "TEST-1" {
		t.Errorf("Issues[0].Key = %v, want %v", result.Issues[0].Key, "TEST-1")
	}
}
