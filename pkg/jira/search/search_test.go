package search

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_ServerDC(t *testing.T) {
	mockResponse := map[string]any{
		"startAt":    0,
		"maxResults": 50,
		"total":      2,
		"issues": []any{
			map[string]any{"key": "TEST-1", "fields": map[string]any{"summary": "First"}},
			map[string]any{"key": "TEST-2", "fields": map[string]any{"summary": "Second"}},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/rest/api/2/search", r.URL.Path)
		assert.Equal(t, "project = TEST", r.URL.Query().Get("jql"))
		assert.Equal(t, "summary,status", r.URL.Query().Get("fields"))
		assert.Equal(t, "50", r.URL.Query().Get("maxResults"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := searchIssues(client, SearchOpts{JQL: "project = TEST", Fields: "summary,status", Limit: 50})
	require.NoError(t, err)
	issues := result["issues"].([]any)
	assert.Len(t, issues, 2)
}

func TestSearch_Cloud(t *testing.T) {
	mockResponse := map[string]any{
		"maxResults": 10,
		"issues": []any{
			map[string]any{"key": "CLOUD-1"},
		},
		"nextPageToken": "next-token",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/rest/api/3/search/jql"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))
		assert.Equal(t, "project = CLOUD", payload["jql"])
		assert.Equal(t, float64(10), payload["maxResults"])
		assert.Equal(t, "summary", payload["fields"])
		assert.Equal(t, "renderedFields", payload["expand"])

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "u", Token: "t"}
	client := api.NewClient(server.URL+"/company.atlassian.net", auth)

	result, err := searchIssues(client, SearchOpts{JQL: "project = CLOUD", Fields: "summary", Limit: 10, Expand: "renderedFields"})
	require.NoError(t, err)
	issues := result["issues"].([]any)
	assert.Len(t, issues, 1)
	assert.Equal(t, "next-token", result["nextPageToken"])
}
