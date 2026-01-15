package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIssuesService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/issue/TEST-1")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Issue{
			Key: "TEST-1",
			ID:  "10001",
			Fields: &IssueFields{
				Summary: "Test issue",
				Status: &Status{
					Name: "To Do",
				},
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	issue, _, err := client.Issues.Get(context.Background(), "TEST-1", nil)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if issue.Key != "TEST-1" {
		t.Errorf("Key = %v, want %v", issue.Key, "TEST-1")
	}
	if issue.Fields.Summary != "Test issue" {
		t.Errorf("Summary = %v, want %v", issue.Fields.Summary, "Test issue")
	}
}

func TestIssuesService_Get_WithExpand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expand := r.URL.Query().Get("expand")
		if expand != "changelog,transitions" {
			t.Errorf("expand = %v, want %v", expand, "changelog,transitions")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Issue{Key: "TEST-1"})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	opts := &IssueGetOptions{Expand: []string{"changelog", "transitions"}}
	_, _, err := client.Issues.Get(context.Background(), "TEST-1", opts)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
}

func TestIssuesService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/rest/api/3/issue" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/issue")
		}

		var req IssueCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if summary, ok := req.Fields["summary"].(string); !ok || summary != "New issue" {
			t.Errorf("Summary = %v, want %v", req.Fields["summary"], "New issue")
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Issue{
			Key: "TEST-2",
			ID:  "10002",
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	issue, _, err := client.Issues.Create(context.Background(), &IssueCreateRequest{
		Fields: map[string]any{
			"summary":   "New issue",
			"project":   map[string]any{"key": "TEST"},
			"issuetype": map[string]any{"name": "Bug"},
		},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if issue.Key != "TEST-2" {
		t.Errorf("Key = %v, want %v", issue.Key, "TEST-2")
	}
}

func TestIssuesService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodPut)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/issue/TEST-1")
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	_, err := client.Issues.Update(context.Background(), "TEST-1", &IssueUpdateRequest{
		Fields: map[string]any{
			"summary": "Updated summary",
		},
	}, nil)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
}

func TestIssuesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodDelete)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/issue/TEST-1")
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	_, err := client.Issues.Delete(context.Background(), "TEST-1", false)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestIssuesService_GetTransitions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/TEST-1/transitions" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/issue/TEST-1/transitions")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"transitions": []*Transition{
				{ID: "11", Name: "Start Progress"},
				{ID: "21", Name: "Done"},
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	transitions, _, err := client.Issues.GetTransitions(context.Background(), "TEST-1", nil)
	if err != nil {
		t.Fatalf("GetTransitions() error = %v", err)
	}
	if len(transitions) != 2 {
		t.Errorf("len(transitions) = %v, want %v", len(transitions), 2)
	}
}

func TestIssuesService_DoTransition(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1/transitions" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/issue/TEST-1/transitions")
		}

		var req IssueTransitionRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Transition == nil || req.Transition.ID != "11" {
			t.Errorf("Transition.ID = %v, want %v", req.Transition, "11")
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	_, err := client.Issues.DoTransition(context.Background(), "TEST-1", &IssueTransitionRequest{
		Transition: &TransitionInput{ID: "11"},
	})
	if err != nil {
		t.Fatalf("DoTransition() error = %v", err)
	}
}

func TestIssuesService_Assign(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodPut)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1/assignee" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/issue/TEST-1/assignee")
		}

		var req map[string]string
		json.NewDecoder(r.Body).Decode(&req)
		if req["accountId"] != "123abc" {
			t.Errorf("accountId = %v, want %v", req["accountId"], "123abc")
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	_, err := client.Issues.Assign(context.Background(), "TEST-1", "123abc")
	if err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
}

func TestIssueFields_JSON(t *testing.T) {
	input := `{
		"summary": "Test issue",
		"description": "Test description",
		"issuetype": {
			"id": "10001",
			"name": "Bug"
		},
		"project": {
			"key": "TEST"
		},
		"status": {
			"name": "To Do"
		},
		"priority": {
			"id": "3",
			"name": "Medium"
		},
		"assignee": {
			"accountId": "123",
			"displayName": "Test User"
		},
		"labels": ["bug", "urgent"]
	}`

	var fields IssueFields
	if err := json.Unmarshal([]byte(input), &fields); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if fields.Summary != "Test issue" {
		t.Errorf("Summary = %v, want %v", fields.Summary, "Test issue")
	}
	if fields.Type == nil || fields.Type.Name != "Bug" {
		t.Error("Type is nil or incorrect")
	}
	if fields.Status == nil || fields.Status.Name != "To Do" {
		t.Error("Status is nil or incorrect")
	}
	if len(fields.Labels) != 2 {
		t.Errorf("len(Labels) = %v, want %v", len(fields.Labels), 2)
	}
}
