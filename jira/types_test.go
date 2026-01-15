package jira

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*Time) bool
	}{
		{
			name:    "full datetime with offset",
			input:   `"2024-01-15T10:30:45.123-0500"`,
			wantErr: false,
			check: func(jt *Time) bool {
				return jt.Year() == 2024 && jt.Month() == time.January && jt.Day() == 15
			},
		},
		{
			name:    "datetime with Z suffix",
			input:   `"2024-01-15T10:30:45.123Z"`,
			wantErr: false,
			check: func(jt *Time) bool {
				return jt.Year() == 2024 && jt.Hour() == 10
			},
		},
		{
			name:    "datetime without milliseconds",
			input:   `"2024-01-15T10:30:45Z"`,
			wantErr: false,
			check: func(jt *Time) bool {
				return jt.Minute() == 30 && jt.Second() == 45
			},
		},
		{
			name:    "date only",
			input:   `"2024-01-15"`,
			wantErr: false,
			check: func(jt *Time) bool {
				return jt.Year() == 2024 && jt.Month() == time.January && jt.Day() == 15
			},
		},
		{
			name:    "empty string",
			input:   `""`,
			wantErr: false,
			check: func(jt *Time) bool {
				return jt.IsZero()
			},
		},
		{
			name:    "invalid format",
			input:   `"not-a-date"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jt Time
			err := json.Unmarshal([]byte(tt.input), &jt)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !tt.check(&jt) {
				t.Errorf("UnmarshalJSON() time check failed for %s", tt.input)
			}
		})
	}
}

func TestTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    Time
		expected string
	}{
		{
			name:     "zero time",
			input:    Time{},
			expected: "null",
		},
		{
			name:     "valid time",
			input:    Time{time.Date(2024, time.January, 15, 10, 30, 45, 123000000, time.UTC)},
			expected: `"2024-01-15T10:30:45.123+0000"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.expected {
				t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.expected)
			}
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*Date) bool
	}{
		{
			name:    "valid date",
			input:   `"2024-01-15"`,
			wantErr: false,
			check: func(d *Date) bool {
				return d.Year() == 2024 && d.Month() == time.January && d.Day() == 15
			},
		},
		{
			name:    "empty string",
			input:   `""`,
			wantErr: false,
			check: func(d *Date) bool {
				return d.IsZero()
			},
		},
		{
			name:    "invalid format",
			input:   `"2024/01/15"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Date
			err := json.Unmarshal([]byte(tt.input), &d)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !tt.check(&d) {
				t.Errorf("UnmarshalJSON() date check failed for %s", tt.input)
			}
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    Date
		expected string
	}{
		{
			name:     "zero date",
			input:    Date{},
			expected: "null",
		},
		{
			name:     "valid date",
			input:    Date{time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)},
			expected: `"2024-01-15"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.expected {
				t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.expected)
			}
		})
	}
}

func TestUser_JSON(t *testing.T) {
	input := `{
		"self": "https://example.atlassian.net/rest/api/3/user?accountId=123",
		"accountId": "123",
		"emailAddress": "user@example.com",
		"displayName": "Test User",
		"active": true,
		"timeZone": "America/New_York"
	}`

	var user User
	if err := json.Unmarshal([]byte(input), &user); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if user.AccountID != "123" {
		t.Errorf("AccountID = %v, want %v", user.AccountID, "123")
	}
	if user.DisplayName != "Test User" {
		t.Errorf("DisplayName = %v, want %v", user.DisplayName, "Test User")
	}
	if !user.Active {
		t.Error("Active = false, want true")
	}
}

func TestProject_JSON(t *testing.T) {
	input := `{
		"self": "https://example.atlassian.net/rest/api/3/project/10000",
		"id": "10000",
		"key": "TEST",
		"name": "Test Project",
		"projectTypeKey": "software",
		"simplified": false
	}`

	var project Project
	if err := json.Unmarshal([]byte(input), &project); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if project.Key != "TEST" {
		t.Errorf("Key = %v, want %v", project.Key, "TEST")
	}
	if project.Name != "Test Project" {
		t.Errorf("Name = %v, want %v", project.Name, "Test Project")
	}
}

func TestIssueType_JSON(t *testing.T) {
	input := `{
		"self": "https://example.atlassian.net/rest/api/3/issuetype/10001",
		"id": "10001",
		"name": "Bug",
		"subtask": false,
		"hierarchyLevel": 0
	}`

	var issueType IssueType
	if err := json.Unmarshal([]byte(input), &issueType); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if issueType.Name != "Bug" {
		t.Errorf("Name = %v, want %v", issueType.Name, "Bug")
	}
	if issueType.Subtask {
		t.Error("Subtask = true, want false")
	}
}

func TestStatus_JSON(t *testing.T) {
	input := `{
		"self": "https://example.atlassian.net/rest/api/3/status/10001",
		"id": "10001",
		"name": "In Progress",
		"statusCategory": {
			"self": "https://example.atlassian.net/rest/api/3/statuscategory/4",
			"id": 4,
			"key": "indeterminate",
			"colorName": "yellow",
			"name": "In Progress"
		}
	}`

	var status Status
	if err := json.Unmarshal([]byte(input), &status); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if status.Name != "In Progress" {
		t.Errorf("Name = %v, want %v", status.Name, "In Progress")
	}
	if status.StatusCategory == nil {
		t.Fatal("StatusCategory is nil")
	}
	if status.StatusCategory.ColorName != "yellow" {
		t.Errorf("StatusCategory.ColorName = %v, want %v", status.StatusCategory.ColorName, "yellow")
	}
}

func TestTransition_JSON(t *testing.T) {
	input := `{
		"id": "11",
		"name": "Start Progress",
		"hasScreen": false,
		"isGlobal": false,
		"isInitial": false,
		"isAvailable": true,
		"isConditional": false,
		"to": {
			"id": "3",
			"name": "In Progress"
		}
	}`

	var transition Transition
	if err := json.Unmarshal([]byte(input), &transition); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if transition.ID != "11" {
		t.Errorf("ID = %v, want %v", transition.ID, "11")
	}
	if transition.Name != "Start Progress" {
		t.Errorf("Name = %v, want %v", transition.Name, "Start Progress")
	}
	if !transition.IsAvailable {
		t.Error("IsAvailable = false, want true")
	}
}
