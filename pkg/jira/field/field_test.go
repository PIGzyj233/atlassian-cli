package field

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchFields_NoKeyword(t *testing.T) {
	mockFields := []any{
		map[string]any{"id": "summary", "name": "Summary"},
		map[string]any{"id": "status", "name": "Status"},
		map[string]any{"id": "priority", "name": "Priority"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/rest/api/2/field", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockFields))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := searchFields(client, "", 10)
	require.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestSearchFields_WithKeyword(t *testing.T) {
	mockFields := []any{
		map[string]any{"id": "summary", "name": "Summary"},
		map[string]any{"id": "status", "name": "Status"},
		map[string]any{"id": "customfield_10001", "name": "Epic Link", "clauseNames": []any{"cf[10001]", "Epic Link"}},
		map[string]any{"id": "priority", "name": "Priority"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockFields))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	// Test keyword matching on name
	result, err := searchFields(client, "epic", 10)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	firstField := result[0].(map[string]any)
	assert.Equal(t, "Epic Link", firstField["name"])
}

func TestSearchFields_WithLimit(t *testing.T) {
	mockFields := []any{
		map[string]any{"id": "f1", "name": "Status One"},
		map[string]any{"id": "f2", "name": "Status Two"},
		map[string]any{"id": "f3", "name": "Status Three"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockFields))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := searchFields(client, "status", 2)
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestSearchFields_MatchByClauseName(t *testing.T) {
	mockFields := []any{
		map[string]any{"id": "customfield_10001", "name": "Custom Field", "clauseNames": []any{"cf[10001]", "story points"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(mockFields))
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := searchFields(client, "story", 10)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestFieldOptionsCloud(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		callCount++

		switch {
		case strings.Contains(r.URL.Path, "/field/customfield_10001/context/10100/option"):
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"values": []any{
					map[string]any{"id": "1", "value": "Option A"},
					map[string]any{"id": "2", "value": "Option B"},
					map[string]any{"id": "3", "value": "Option C"},
				},
				"total": 3,
			}))
		case strings.Contains(r.URL.Path, "/field/customfield_10001/context"):
			// Return contexts; first is global
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"values": []any{
					map[string]any{"id": "10100", "name": "Default", "isGlobalContext": true},
				},
			}))
		default:
			http.Error(w, fmt.Sprintf("unexpected path: %s", r.URL.Path), http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Use .atlassian.net URL so IsCloud() returns true
	auth := &config.BasicAuth{Username: "u", Token: "t"}
	client := api.NewClient(server.URL+"/company.atlassian.net", auth)

	result, err := getFieldOptionsCloud(client, "customfield_10001", "", "")
	require.NoError(t, err)
	options := result.([]any)
	assert.Len(t, options, 3)

	// Verify context was auto-resolved
	assert.GreaterOrEqual(t, callCount, 2)
}

func TestFieldOptionsCloud_WithContains(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "/field/customfield_10001/context") && !strings.Contains(r.URL.Path, "/option") {
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"values": []any{
					map[string]any{"id": "10100", "isGlobalContext": true},
				},
			}))
		} else {
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"values": []any{
					map[string]any{"id": "1", "value": "High Priority"},
					map[string]any{"id": "2", "value": "Low Priority"},
					map[string]any{"id": "3", "value": "Medium"},
				},
				"total": 3,
			}))
		}
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "u", Token: "t"}
	client := api.NewClient(server.URL+"/company.atlassian.net", auth)

	result, err := getFieldOptionsCloud(client, "customfield_10001", "", "priority")
	require.NoError(t, err)
	options := result.([]any)
	assert.Len(t, options, 2)
}

func TestMatchesKeyword(t *testing.T) {
	tests := []struct {
		name    string
		field   map[string]any
		keyword string
		want    bool
	}{
		{"match by id", map[string]any{"id": "summary"}, "summ", true},
		{"match by name", map[string]any{"id": "f1", "name": "Epic Link"}, "epic", true},
		{"match by clauseName", map[string]any{"id": "f1", "clauseNames": []any{"story points"}}, "story", true},
		{"no match", map[string]any{"id": "f1", "name": "Summary"}, "epic", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, matchesKeyword(tt.field, tt.keyword))
		})
	}
}
