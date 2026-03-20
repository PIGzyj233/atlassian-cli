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

func TestCreateIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/rest/api/2/issue", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))

		fields := payload["fields"].(map[string]any)
		project := fields["project"].(map[string]any)
		assert.Equal(t, "TEST", project["key"])
		assert.Equal(t, "Bug title", fields["summary"])
		assert.Equal(t, "Bug", fields["issuetype"].(map[string]any)["name"])

		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
			"id":   "10001",
			"key":  "TEST-42",
			"self": "https://example.com/rest/api/3/issue/10001",
		}))
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "u", Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := createIssue(client, CreateIssueOpts{
		Project:   "TEST",
		Summary:   "Bug title",
		IssueType: "Bug",
	})
	require.NoError(t, err)
	assert.Equal(t, "TEST-42", result["key"])
}
