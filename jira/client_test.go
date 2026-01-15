package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("https://example.atlassian.net")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
	if client.baseURL.String() != "https://example.atlassian.net" {
		t.Errorf("baseURL = %v, want %v", client.baseURL.String(), "https://example.atlassian.net")
	}
}

func TestNewClient_InvalidURL(t *testing.T) {
	_, err := NewClient("://invalid")
	if err == nil {
		t.Error("NewClient() expected error for invalid URL")
	}
}

func TestClient_WithBasicAuth(t *testing.T) {
	client, _ := NewClient("https://example.atlassian.net", WithBasicAuth("user@example.com", "token"))
	if client.auth == nil {
		t.Fatal("WithBasicAuth() did not set auth")
	}
}

func TestClient_WithBearerToken(t *testing.T) {
	client, _ := NewClient("https://example.atlassian.net", WithBearerToken("my-token"))
	if client.auth == nil {
		t.Fatal("WithBearerToken() did not set auth")
	}
}

func TestClient_WithHTTPClient(t *testing.T) {
	customClient := &http.Client{}
	client, _ := NewClient("https://example.atlassian.net", WithHTTPClient(customClient))
	if client.client != customClient {
		t.Error("WithHTTPClient() did not set custom client")
	}
}

func TestClient_WithUserAgent(t *testing.T) {
	client, _ := NewClient("https://example.atlassian.net", WithUserAgent("my-custom-agent"))
	if client.UserAgent != "my-custom-agent" {
		t.Errorf("UserAgent = %v, want %v", client.UserAgent, "my-custom-agent")
	}
}

func TestClient_NewRequest(t *testing.T) {
	client, _ := NewClient("https://example.atlassian.net")
	req, err := client.NewRequest(context.Background(), http.MethodGet, "/rest/api/3/issue/TEST-1", nil)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	if req.URL.String() != "https://example.atlassian.net/rest/api/3/issue/TEST-1" {
		t.Errorf("URL = %v, want %v", req.URL.String(), "https://example.atlassian.net/rest/api/3/issue/TEST-1")
	}
	// Content-Type is only set when there's a body
	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("Accept = %v, want %v", req.Header.Get("Accept"), "application/json")
	}
}

func TestClient_NewRequest_WithBody(t *testing.T) {
	client, _ := NewClient("https://example.atlassian.net")
	body := map[string]string{"summary": "Test issue"}
	req, err := client.NewRequest(context.Background(), http.MethodPost, "/rest/api/3/issue", body)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	if req.Body == nil {
		t.Error("NewRequest() body is nil")
	}
}

func TestClient_Do(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/serverInfo" {
			t.Errorf("URL path = %v, want %v", r.URL.Path, "/rest/api/3/serverInfo")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ServerInfo{
			Version:        "9.0.0",
			DeploymentType: "Cloud",
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	req, _ := client.NewRequest(context.Background(), http.MethodGet, "/rest/api/3/serverInfo", nil)

	var info ServerInfo
	resp, err := client.Do(req, &info)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}
	if info.Version != "9.0.0" {
		t.Errorf("Version = %v, want %v", info.Version, "9.0.0")
	}
}

func TestClient_Do_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"errorMessages": []string{"Issue does not exist"},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	req, _ := client.NewRequest(context.Background(), http.MethodGet, "/rest/api/3/issue/NOTFOUND-1", nil)

	var issue Issue
	_, err := client.Do(req, &issue)
	if err == nil {
		t.Error("Do() expected error for 404 response")
	}
}

func TestBasicAuth_Apply(t *testing.T) {
	auth := &BasicAuth{
		Email:    "user@example.com",
		APIToken: "token",
	}
	req, _ := http.NewRequest(http.MethodGet, "https://example.atlassian.net", nil)
	auth.Apply(req)

	if req.Header.Get("Authorization") == "" {
		t.Error("Apply() did not set Authorization header")
	}
}

func TestBearerAuth_Apply(t *testing.T) {
	auth := &BearerAuth{
		Token: "my-token",
	}
	req, _ := http.NewRequest(http.MethodGet, "https://example.atlassian.net", nil)
	auth.Apply(req)

	expected := "Bearer my-token"
	if req.Header.Get("Authorization") != expected {
		t.Errorf("Authorization = %v, want %v", req.Header.Get("Authorization"), expected)
	}
}

// setupTestServer creates a test server and client for testing.
func setupTestServer(handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)
	client, _ := NewClient(server.URL)
	return server, client
}
