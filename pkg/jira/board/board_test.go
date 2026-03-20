package board

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

func TestListBoards(t *testing.T) {
	mockResponse := map[string]any{
		"maxResults": 50,
		"startAt":    0,
		"total":      2,
		"values": []any{
			map[string]any{"id": float64(1), "name": "Scrum Board", "type": "scrum"},
			map[string]any{"id": float64(2), "name": "Kanban Board", "type": "kanban"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/rest/agile/1.0/board", r.URL.Path)
		assert.Equal(t, "scrum", r.URL.Query().Get("type"))
		assert.Equal(t, "50", r.URL.Query().Get("maxResults"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := listBoards(client, "", "", "scrum", 50, 0)
	require.NoError(t, err)
	values := result["values"].([]any)
	assert.Len(t, values, 2)
}

func TestListBoards_WithName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Sprint Board", r.URL.Query().Get("name"))
		assert.Equal(t, "PROJ", r.URL.Query().Get("projectKeyOrId"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
			"values": []any{},
			"total":  0,
		}))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := listBoards(client, "Sprint Board", "PROJ", "", 50, 0)
	require.NoError(t, err)
	values := result["values"].([]any)
	assert.Len(t, values, 0)
}

func TestListBoardIssues(t *testing.T) {
	mockResponse := map[string]any{
		"startAt":    0,
		"maxResults": 50,
		"total":      1,
		"issues": []any{
			map[string]any{
				"key": "PROJ-1",
				"fields": map[string]any{
					"summary": "Test issue",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/rest/agile/1.0/board/42/issue", r.URL.Path)
		assert.Equal(t, "sprint = 100", r.URL.Query().Get("jql"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := listBoardIssues(client, "42", "sprint = 100", "", 50, 0)
	require.NoError(t, err)
	issues := result["issues"].([]any)
	assert.Len(t, issues, 1)
	first := issues[0].(map[string]any)
	assert.Equal(t, "PROJ-1", first["key"])
}
