package transition

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

func TestListTransitions(t *testing.T) {
	mockResponse := map[string]any{
		"transitions": []any{
			map[string]any{
				"id":   "11",
				"name": "To Do",
				"to": map[string]any{
					"id":   "10001",
					"name": "To Do",
				},
			},
			map[string]any{
				"id":   "21",
				"name": "In Progress",
				"to": map[string]any{
					"id":   "10002",
					"name": "In Progress",
				},
			},
			map[string]any{
				"id":   "31",
				"name": "Done",
				"to": map[string]any{
					"id":   "10003",
					"name": "Done",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/rest/api/2/issue/TEST-1/transitions", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockResponse))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := listTransitions(client, "TEST-1")
	require.NoError(t, err)
	assert.Len(t, result, 3)

	// Verify transition structure
	first := result[0].(map[string]any)
	assert.Equal(t, "11", first["id"])
	assert.Equal(t, "To Do", first["name"])
	to := first["to"].(map[string]any)
	assert.Equal(t, "To Do", to["name"])
}

func TestListTransitions_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
			"transitions": []any{},
		}))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := listTransitions(client, "TEST-1")
	require.NoError(t, err)
	assert.Len(t, result, 0)
}
