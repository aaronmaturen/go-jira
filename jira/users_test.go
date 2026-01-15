package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodGet)
		}

		accountID := r.URL.Query().Get("accountId")
		if accountID != "123abc" {
			t.Errorf("accountId = %v, want %v", accountID, "123abc")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{
			AccountID:    "123abc",
			DisplayName:  "Test User",
			EmailAddress: "test@example.com",
			Active:       true,
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	user, _, err := client.Users.Get(context.Background(), "123abc", nil)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if user.AccountID != "123abc" {
		t.Errorf("AccountID = %v, want %v", user.AccountID, "123abc")
	}
	if user.DisplayName != "Test User" {
		t.Errorf("DisplayName = %v, want %v", user.DisplayName, "Test User")
	}
}

func TestUsersService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/user/search" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/user/search")
		}

		query := r.URL.Query().Get("query")
		if query != "test" {
			t.Errorf("query = %v, want %v", query, "test")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*User{
			{AccountID: "123", DisplayName: "Test User 1"},
			{AccountID: "456", DisplayName: "Test User 2"},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	users, _, err := client.Users.Search(context.Background(), &UserSearchOptions{
		Query: "test",
	})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(users) != 2 {
		t.Errorf("len(users) = %v, want %v", len(users), 2)
	}
}

func TestUsersService_FindAssignableUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/user/assignable/search" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/user/assignable/search")
		}

		project := r.URL.Query().Get("project")
		if project != "TEST" {
			t.Errorf("project = %v, want %v", project, "TEST")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*User{
			{AccountID: "123", DisplayName: "Assignable User"},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	users, _, err := client.Users.FindAssignableUsers(context.Background(), &FindAssignableOptions{
		Project: "TEST",
	})
	if err != nil {
		t.Fatalf("FindAssignableUsers() error = %v", err)
	}
	if len(users) != 1 {
		t.Errorf("len(users) = %v, want %v", len(users), 1)
	}
}

func TestUsersService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/rest/api/3/user" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/user")
		}

		var req UserCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.EmailAddress != "new@example.com" {
			t.Errorf("EmailAddress = %v, want %v", req.EmailAddress, "new@example.com")
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{
			AccountID:    "newuser123",
			EmailAddress: "new@example.com",
			DisplayName:  "New User",
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	user, _, err := client.Users.Create(context.Background(), &UserCreateRequest{
		EmailAddress: "new@example.com",
		DisplayName:  "New User",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if user.AccountID != "newuser123" {
		t.Errorf("AccountID = %v, want %v", user.AccountID, "newuser123")
	}
}

func TestUsersService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %v, want %v", r.Method, http.MethodDelete)
		}

		accountID := r.URL.Query().Get("accountId")
		if accountID != "123abc" {
			t.Errorf("accountId = %v, want %v", accountID, "123abc")
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	_, err := client.Users.Delete(context.Background(), "123abc")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestUsersService_GetAllUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/users/search" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/users/search")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*User{
			{AccountID: "123", DisplayName: "User 1"},
			{AccountID: "456", DisplayName: "User 2"},
			{AccountID: "789", DisplayName: "User 3"},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	users, _, err := client.Users.GetAllUsers(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetAllUsers() error = %v", err)
	}
	if len(users) != 3 {
		t.Errorf("len(users) = %v, want %v", len(users), 3)
	}
}

func TestUsersService_BulkGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/user/bulk" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/user/bulk")
		}

		accountIDs := r.URL.Query()["accountId"]
		if len(accountIDs) != 2 {
			t.Errorf("len(accountId) = %v, want %v", len(accountIDs), 2)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BulkGetResult{
			Values: []*User{
				{AccountID: "123", DisplayName: "User 1"},
				{AccountID: "456", DisplayName: "User 2"},
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	result, _, err := client.Users.BulkGet(context.Background(), &BulkGetOptions{
		AccountIDs: []string{"123", "456"},
	})
	if err != nil {
		t.Fatalf("BulkGet() error = %v", err)
	}
	if len(result.Values) != 2 {
		t.Errorf("len(Values) = %v, want %v", len(result.Values), 2)
	}
}
