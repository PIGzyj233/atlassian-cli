package sprint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListSprints(t *testing.T) {
	mockResponse := map[string]any{
		"maxResults": 50,
		"startAt":    0,
		"values": []any{
			map[string]any{
				"id":    float64(1),
				"name":  "Sprint 1",
				"state": "active",
			},
			map[string]any{
				"id":    float64(2),
				"name":  "Sprint 2",
				"state": "future",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/rest/agile/1.0/board/10/sprint", r.URL.Path)
		assert.Equal(t, "active", r.URL.Query().Get("state"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := listSprints(client, "10", "active", 50)
	require.NoError(t, err)
	values := result["values"].([]any)
	assert.Len(t, values, 2)
}

func TestCreateSprint(t *testing.T) {
	mockResponse := map[string]any{
		"id":            float64(100),
		"name":          "New Sprint",
		"originBoardId": float64(10),
		"state":         "future",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/rest/agile/1.0/sprint", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))
		assert.Equal(t, "New Sprint", payload["name"])
		assert.Equal(t, "10", payload["originBoardId"])
		assert.Equal(t, "2027-01-01", payload["startDate"])
		assert.Equal(t, "2027-01-15", payload["endDate"])
		assert.Equal(t, "Deliver feature X", payload["goal"])

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := createSprint(client, "10", "New Sprint", "2027-01-01", "2027-01-15", "Deliver feature X")
	require.NoError(t, err)
	assert.Equal(t, float64(100), result["id"])
	assert.Equal(t, "New Sprint", result["name"])
}

func TestUpdateSprint(t *testing.T) {
	mockResponse := map[string]any{
		"id":    float64(100),
		"name":  "Updated Sprint",
		"state": "active",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/rest/agile/1.0/sprint/100", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))
		assert.Equal(t, "Updated Sprint", payload["name"])
		assert.Equal(t, "active", payload["state"])
		// goal should not be sent if empty
		_, hasGoal := payload["goal"]
		assert.False(t, hasGoal)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := updateSprint(client, "100", "Updated Sprint", "active", "", "", "")
	require.NoError(t, err)
	assert.Equal(t, "Updated Sprint", result["name"])
	assert.Equal(t, "active", result["state"])
}

func TestAddIssuesToSprint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/rest/agile/1.0/sprint/100/issue", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))
		issues := payload["issues"].([]any)
		assert.Len(t, issues, 2)
		assert.Equal(t, "PROJ-1", issues[0])
		assert.Equal(t, "PROJ-2", issues[1])

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	// Test the POST body construction directly
	path := fmt.Sprintf("/rest/agile/1.0/sprint/%s/issue", "100")
	keys := parseKeys("PROJ-1, PROJ-2")
	assert.Equal(t, []string{"PROJ-1", "PROJ-2"}, keys)

	body := map[string]any{
		"issues": keys,
	}
	_, err := client.Post(path, body, nil)
	require.NoError(t, err)
}

func TestListSprintIssues(t *testing.T) {
	mockResponse := map[string]any{
		"startAt":    0,
		"maxResults": 50,
		"total":      1,
		"issues": []any{
			map[string]any{
				"key":    "PROJ-1",
				"fields": map[string]any{"summary": "Sprint task"},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/rest/agile/1.0/sprint/50/issue", r.URL.Path)
		assert.Equal(t, "summary,status", r.URL.Query().Get("fields"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := listSprintIssues(client, "50", "summary,status", 50)
	require.NoError(t, err)
	issues := result["issues"].([]any)
	assert.Len(t, issues, 1)
}

func TestParseKeys(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"PROJ-1,PROJ-2", []string{"PROJ-1", "PROJ-2"}},
		{"PROJ-1, PROJ-2, PROJ-3", []string{"PROJ-1", "PROJ-2", "PROJ-3"}},
		{" PROJ-1 ", []string{"PROJ-1"}},
		{"", nil},
		{",,,", nil},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, parseKeys(tt.input))
		})
	}
}

func TestValidSprintStates(t *testing.T) {
	assert.True(t, validSprintStates["future"])
	assert.True(t, validSprintStates["active"])
	assert.True(t, validSprintStates["closed"])
	assert.False(t, validSprintStates["canceled"])
	assert.False(t, validSprintStates[""])
}
