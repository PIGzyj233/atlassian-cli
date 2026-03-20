package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PigZyj2333/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIError_JSON(t *testing.T) {
	apiErr := &APIError{
		StatusCode: 404,
		Message:    "Issue NOT-123 not found",
		Hint:       "Check that the issue key is correct",
	}

	data, err := json.Marshal(apiErr)
	require.NoError(t, err)

	var result map[string]any
	json.Unmarshal(data, &result)
	assert.Equal(t, float64(404), result["status"])
	assert.Equal(t, "Issue NOT-123 not found", result["error"])
	assert.Equal(t, "Check that the issue key is correct", result["hint"])
}

func TestAPIError_ErrorString(t *testing.T) {
	apiErr := &APIError{StatusCode: 401, Message: "Unauthorized"}
	assert.Contains(t, apiErr.Error(), "401")
	assert.Contains(t, apiErr.Error(), "Unauthorized")
}

func TestParseAPIError_AtlassianFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]any{
			"errorMessages": []string{"Issue does not exist or you do not have permission"},
		})
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := NewClient(server.URL, auth)

	_, err := client.Get("/test", nil)
	require.Error(t, err)

	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 404, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "Issue does not exist")
}
