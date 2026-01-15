package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/rest/api/3/project/search" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/project/search")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ProjectListResult{
			Values: []*Project{
				{Key: "TEST", Name: "Test Project"},
				{Key: "DEMO", Name: "Demo Project"},
			},
			Total: 2,
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	result, _, err := client.Projects.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(result.Values) != 2 {
		t.Errorf("len(Values) = %v, want %v", len(result.Values), 2)
	}
	if result.Values[0].Key != "TEST" {
		t.Errorf("Values[0].Key = %v, want %v", result.Values[0].Key, "TEST")
	}
}

func TestProjectsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/project/TEST" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/project/TEST")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Project{
			Key:         "TEST",
			Name:        "Test Project",
			Description: "A test project",
			Lead: &User{
				AccountID:   "123",
				DisplayName: "Lead User",
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	project, _, err := client.Projects.Get(context.Background(), "TEST", nil)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if project.Key != "TEST" {
		t.Errorf("Key = %v, want %v", project.Key, "TEST")
	}
	if project.Lead == nil {
		t.Fatal("Lead is nil")
	}
	if project.Lead.DisplayName != "Lead User" {
		t.Errorf("Lead.DisplayName = %v, want %v", project.Lead.DisplayName, "Lead User")
	}
}

func TestProjectsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/rest/api/3/project" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/project")
		}

		var req ProjectCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Key != "NEW" {
			t.Errorf("Key = %v, want %v", req.Key, "NEW")
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ProjectCreateResponse{
			ID:  10003,
			Key: "NEW",
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	project, _, err := client.Projects.Create(context.Background(), &ProjectCreateRequest{
		Key:            "NEW",
		Name:           "New Project",
		ProjectTypeKey: "software",
		LeadAccountID:  "123",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if project.Key != "NEW" {
		t.Errorf("Key = %v, want %v", project.Key, "NEW")
	}
}

func TestProjectsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodPut)
		}
		if r.URL.Path != "/rest/api/3/project/TEST" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/project/TEST")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Project{
			Key:  "TEST",
			Name: "Updated Project",
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	project, _, err := client.Projects.Update(context.Background(), "TEST", &ProjectUpdateRequest{
		Name: "Updated Project",
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if project.Name != "Updated Project" {
		t.Errorf("Name = %v, want %v", project.Name, "Updated Project")
	}
}

func TestProjectsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodDelete)
		}
		if r.URL.Path != "/rest/api/3/project/TEST" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/project/TEST")
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	_, err := client.Projects.Delete(context.Background(), "TEST", false)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestProjectsService_GetStatuses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/project/TEST/statuses" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/project/TEST/statuses")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*IssueTypeWithStatuses{
			{
				ID:   "10001",
				Name: "Bug",
				Statuses: []*Status{
					{ID: "1", Name: "To Do"},
					{ID: "2", Name: "Done"},
				},
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	statuses, _, err := client.Projects.GetStatuses(context.Background(), "TEST")
	if err != nil {
		t.Fatalf("GetStatuses() error = %v", err)
	}
	if len(statuses) != 1 {
		t.Errorf("len(statuses) = %v, want %v", len(statuses), 1)
	}
}

func TestProjectsService_ListRecent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/project/recent" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/project/recent")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*Project{
			{Key: "TEST", Name: "Test Project"},
			{Key: "DEMO", Name: "Demo Project"},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	projects, _, err := client.Projects.ListRecent(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListRecent() error = %v", err)
	}
	if len(projects) != 2 {
		t.Errorf("len(projects) = %v, want %v", len(projects), 2)
	}
}
