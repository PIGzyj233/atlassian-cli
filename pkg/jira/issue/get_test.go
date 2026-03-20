package issue

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunIssueGet(t *testing.T) {
	mockIssue := map[string]any{
		"key": "TEST-1",
		"fields": map[string]any{
			"summary": "Test issue summary",
			"status": map[string]any{
				"name": "Open",
			},
			"issuetype": map[string]any{
				"name": "Task",
			},
			"assignee": map[string]any{
				"displayName": "John Doe",
			},
			"description": nil,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/rest/api/2/issue/TEST-1")
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockIssue))
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "u", Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := fetchIssue(client, "TEST-1", "", "")
	require.NoError(t, err)

	assert.Equal(t, "TEST-1", result["key"])
	fields := result["fields"].(map[string]any)
	assert.Equal(t, "Test issue summary", fields["summary"])
}
