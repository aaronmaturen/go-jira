# go-jira

A Go client library for the [Jira Cloud REST API v3](https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/).

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronmaturen/go-jira.svg)](https://pkg.go.dev/github.com/aaronmaturen/go-jira)
[![CI](https://github.com/aaronmaturen/go-jira/actions/workflows/ci.yml/badge.svg)](https://github.com/aaronmaturen/go-jira/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aaronmaturen/go-jira)](https://goreportcard.com/report/github.com/aaronmaturen/go-jira)

## Installation

```bash
go get github.com/aaronmaturen/go-jira
```

Requires Go 1.22 or later.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/aaronmaturen/go-jira/jira"
)

func main() {
    // Create a client with basic authentication
    client, err := jira.NewClient(
        "https://yourinstance.atlassian.net",
        jira.WithBasicAuth("your-email@example.com", "your-api-token"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Search for issues
    result, _, err := client.Search.Do(ctx, "project = MYPROJECT AND status = 'In Progress'", nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, issue := range result.Issues {
        fmt.Printf("%s: %s\n", issue.Key, issue.Fields.Summary)
    }
}
```

## Authentication

### Basic Auth (API Token)

For Jira Cloud, use your email and an [API token](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/):

```go
client, err := jira.NewClient(
    "https://yourinstance.atlassian.net",
    jira.WithBasicAuth("email@example.com", "your-api-token"),
)
```

### Bearer Token (OAuth 2.0)

For OAuth 2.0 access tokens:

```go
client, err := jira.NewClient(
    "https://yourinstance.atlassian.net",
    jira.WithBearerToken("your-oauth-token"),
)
```

## Usage Examples

### Get an Issue

```go
issue, _, err := client.Issues.Get(ctx, "PROJ-123", &jira.IssueGetOptions{
    Expand: []string{"changelog", "transitions"},
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Summary: %s\n", issue.Fields.Summary)
fmt.Printf("Status: %s\n", issue.Fields.Status.Name)
fmt.Printf("Assignee: %s\n", issue.Fields.Assignee.DisplayName)
```

### Create an Issue

```go
issue, _, err := client.Issues.Create(ctx, &jira.IssueCreateRequest{
    Fields: map[string]any{
        "project":   map[string]any{"key": "PROJ"},
        "issuetype": map[string]any{"name": "Bug"},
        "summary":   "Something is broken",
        "description": map[string]any{
            "type":    "doc",
            "version": 1,
            "content": []any{
                map[string]any{
                    "type": "paragraph",
                    "content": []any{
                        map[string]any{
                            "type": "text",
                            "text": "Detailed description here",
                        },
                    },
                },
            },
        },
        "priority": map[string]any{"name": "High"},
    },
})
```

### Update an Issue

```go
_, err := client.Issues.Update(ctx, "PROJ-123", &jira.IssueUpdateRequest{
    Fields: map[string]any{
        "summary": "Updated summary",
    },
}, nil)
```

### Transition an Issue

```go
// Get available transitions
transitions, _, err := client.Issues.GetTransitions(ctx, "PROJ-123", nil)

// Perform a transition
_, err = client.Issues.DoTransition(ctx, "PROJ-123", &jira.IssueTransitionRequest{
    Transition: &jira.TransitionInput{ID: "31"}, // "Done" transition ID
})
```

### Search with JQL

```go
// Using GET (simpler queries)
result, _, err := client.Search.Do(ctx, "project = PROJ ORDER BY created DESC", &jira.SearchOptions{
    MaxResults: 50,
    Fields:     []string{"summary", "status", "assignee"},
})

// Using POST (complex queries, more options)
result, _, err := client.Search.DoPost(ctx, &jira.SearchRequest{
    JQL:        "project = PROJ AND status IN ('To Do', 'In Progress')",
    MaxResults: 100,
    Fields:     []string{"summary", "status", "priority"},
    Expand:     []string{"changelog"},
})
```

### Work with Projects

```go
// List all projects
projects, _, err := client.Projects.List(ctx, nil)

// Get a specific project
project, _, err := client.Projects.Get(ctx, "PROJ", nil)

// Get project statuses
statuses, _, err := client.Projects.GetStatuses(ctx, "PROJ")
```

### Manage Comments

```go
// Add a comment
comment, _, err := client.Comments.Create(ctx, "PROJ-123", &jira.CommentCreateRequest{
    Body: "This is a comment",
})

// Get all comments
comments, _, err := client.Comments.List(ctx, "PROJ-123", nil)
```

### Work with Users

```go
// Search for users
users, _, err := client.Users.Search(ctx, &jira.UserSearchOptions{
    Query: "john",
})

// Get current user
myself, _, err := client.Myself.Get(ctx, nil)

// Find assignable users for a project
assignable, _, err := client.Users.FindAssignableUsers(ctx, &jira.FindAssignableOptions{
    Project: "PROJ",
})
```

## Available Services

The client provides access to the following Jira API services:

| Service | Description |
|---------|-------------|
| `Issues` | Create, read, update, delete issues; transitions; assignments |
| `Search` | JQL search, issue picker |
| `Projects` | Project CRUD, versions, components, statuses |
| `Users` | User search, bulk operations, assignable users |
| `Groups` | Group management, membership |
| `Comments` | Issue comments |
| `Attachments` | File attachments |
| `Worklogs` | Time tracking |
| `IssueLinks` | Issue relationships |
| `IssueLinkTypes` | Link type definitions |
| `Watchers` | Issue watchers |
| `Votes` | Issue voting |
| `Fields` | Custom and system fields |
| `Screens` | Screen configurations |
| `Workflows` | Workflow definitions |
| `WorkflowSchemes` | Workflow scheme mappings |
| `Filters` | Saved JQL filters |
| `Dashboards` | Dashboard management |
| `Priorities` | Priority levels |
| `Statuses` | Status definitions |
| `Resolutions` | Resolution types |
| `IssueTypes` | Issue type definitions |
| `Components` | Project components |
| `Versions` | Project versions |
| `Labels` | Issue labels |
| `Permissions` | Permission schemes |
| `ProjectRoles` | Project role management |
| `ServerInfo` | Server information |
| `Myself` | Current user operations |
| `ApplicationRoles` | Application role management |
| `AuditRecords` | Audit log access |
| `Avatars` | Avatar management |
| `JQL` | JQL autocomplete and validation |

## Configuration Options

```go
client, err := jira.NewClient(
    "https://yourinstance.atlassian.net",
    jira.WithBasicAuth("email", "token"),
    jira.WithHTTPClient(customHTTPClient),
    jira.WithUserAgent("my-app/1.0"),
)
```

## Error Handling

The library returns detailed error information from the Jira API:

```go
issue, resp, err := client.Issues.Get(ctx, "INVALID-123", nil)
if err != nil {
    // Check HTTP status code
    if resp != nil && resp.StatusCode == 404 {
        fmt.Println("Issue not found")
    }
    log.Fatal(err)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Running Tests

```bash
go test -v ./...
```

### Running Linter

```bash
golangci-lint run
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Atlassian Jira Cloud REST API](https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/)
- Inspired by [go-github](https://github.com/google/go-github)
