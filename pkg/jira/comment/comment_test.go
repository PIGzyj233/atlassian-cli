package comment

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

func TestListComments(t *testing.T) {
	mockResponse := map[string]any{
		"startAt":    0,
		"maxResults": 50,
		"total":      2,
		"comments": []any{
			map[string]any{
				"id":   "10001",
				"body": "First comment",
				"author": map[string]any{
					"displayName": "Alice",
				},
			},
			map[string]any{
				"id":   "10002",
				"body": "Second comment",
				"author": map[string]any{
					"displayName": "Bob",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/rest/api/2/issue/TEST-1/comment", r.URL.Path)
		assert.Equal(t, "50", r.URL.Query().Get("maxResults"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := listComments(client, "TEST-1", 50)
	require.NoError(t, err)
	comments := result["comments"].([]any)
	assert.Len(t, comments, 2)
}

func TestAddComment(t *testing.T) {
	mockResponse := map[string]any{
		"id":      "10003",
		"body":    "New comment",
		"created": "2026-03-20T10:00:00.000+0000",
		"author": map[string]any{
			"displayName": "Alice",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/rest/api/2/issue/TEST-1/comment", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))
		assert.Equal(t, "New comment", payload["body"])
		// Check visibility is set
		vis, ok := payload["visibility"].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "role", vis["type"])
		assert.Equal(t, "Developers", vis["value"])

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := addComment(client, "TEST-1", "New comment", "Developers")
	require.NoError(t, err)
	assert.Equal(t, "10003", result["id"])
}

func TestAddComment_NoVisibility(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))
		assert.Equal(t, "Simple comment", payload["body"])
		_, hasVisibility := payload["visibility"]
		assert.False(t, hasVisibility)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"id": "10004"}))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := addComment(client, "TEST-1", "Simple comment", "")
	require.NoError(t, err)
	assert.Equal(t, "10004", result["id"])
}

func TestEditComment(t *testing.T) {
	mockResponse := map[string]any{
		"id":      "10001",
		"body":    "Updated comment",
		"updated": "2026-03-20T11:00:00.000+0000",
		"author": map[string]any{
			"displayName": "Alice",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/rest/api/2/issue/TEST-1/comment/10001", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))
		assert.Equal(t, "Updated comment", payload["body"])

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := editComment(client, "TEST-1", "10001", "Updated comment", "")
	require.NoError(t, err)
	assert.Equal(t, "10001", result["id"])
	assert.Equal(t, "Updated comment", result["body"])
}

func TestAddServiceDeskComment(t *testing.T) {
	mockResponse := map[string]any{
		"id":     "10005",
		"body":   "Customer response",
		"public": true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/rest/servicedeskapi/request/TEST-1/comment", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))
		assert.Equal(t, "Customer response", payload["body"])
		assert.Equal(t, true, payload["public"])

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := addServiceDeskComment(client, "TEST-1", "Customer response", true)
	require.NoError(t, err)
	assert.Equal(t, "10005", result["id"])
	assert.Equal(t, true, result["public"])
}
