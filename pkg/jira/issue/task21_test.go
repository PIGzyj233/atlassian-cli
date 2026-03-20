package issue

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/rest/api/2/issue/TEST-1", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))

		fields := payload["fields"].(map[string]any)
		assert.Equal(t, "Updated summary", fields["summary"])
		components := fields["components"].([]any)
		require.Len(t, components, 1)
		assert.Equal(t, "API", components[0].(map[string]any)["name"])

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "u", Token: "t"}
	client := api.NewClient(server.URL, auth)

	err := updateIssue(client, "TEST-1", map[string]any{
		"summary":    "Updated summary",
		"components": []map[string]any{{"name": "API"}},
	})
	require.NoError(t, err)
}

func TestDeleteIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/rest/api/2/issue/TEST-9", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	err := deleteIssue(client, "TEST-9")
	require.NoError(t, err)
}

func TestTransitionIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/rest/api/2/issue/TEST-7/transitions", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))

		assert.Equal(t, "31", payload["transition"].(map[string]any)["id"])
		assert.Equal(t, "Done", payload["fields"].(map[string]any)["resolution"])

		update := payload["update"].(map[string]any)
		comments := update["comment"].([]any)
		require.Len(t, comments, 1)
		assert.Equal(t, "Ship it", comments[0].(map[string]any)["add"].(map[string]any)["body"])

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "u", Token: "t"}
	client := api.NewClient(server.URL, auth)

	err := transitionIssue(client, "TEST-7", "31", map[string]any{"resolution": "Done"}, "Ship it")
	require.NoError(t, err)
}

func TestBatchCreateIssues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/rest/api/2/issue/bulk", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))

		issueUpdates := payload["issueUpdates"].([]any)
		require.Len(t, issueUpdates, 2)

		firstFields := issueUpdates[0].(map[string]any)["fields"].(map[string]any)
		assert.Equal(t, "TEST", firstFields["project"].(map[string]any)["key"])
		assert.Equal(t, "Issue 1", firstFields["summary"])
		assert.Equal(t, "Task", firstFields["issuetype"].(map[string]any)["name"])

		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
			"issues": []map[string]any{{"key": "TEST-1"}, {"key": "TEST-2"}},
			"errors": []any{},
		}))
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "u", Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := batchCreateIssues(client, []map[string]any{
		{"project_key": "TEST", "summary": "Issue 1", "issue_type": "Task"},
		{"project_key": "TEST", "summary": "Issue 2", "issue_type": "Bug", "components": []string{"Frontend"}},
	})
	require.NoError(t, err)
	issues := result["issues"].([]any)
	assert.Len(t, issues, 2)
}
