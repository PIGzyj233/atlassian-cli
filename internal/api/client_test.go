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

func TestClient_Get_Basic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header is present
		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "user@test.com", username)
		assert.Equal(t, "api-token", password)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "user@test.com", Token: "api-token"}
	client := NewClient(server.URL, auth)

	var result map[string]string
	resp, err := client.Get("/test/endpoint", &result)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "ok", result["status"])
}

func TestClient_Get_PAT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer my-pat", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "my-pat"}
	client := NewClient(server.URL, auth)

	var result map[string]string
	_, err := client.Get("/test", &result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestClient_RetryOn429(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "tok"}
	client := NewClient(server.URL, auth)
	client.MaxRetries = 3

	var result map[string]string
	_, err := client.Get("/test", &result)
	require.NoError(t, err)
	assert.Equal(t, 3, attempts)
	assert.Equal(t, "ok", result["status"])
}

func TestClient_RetryExhausted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "tok"}
	client := NewClient(server.URL, auth)
	client.MaxRetries = 2

	_, err := client.Get("/test", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "429")
}

func TestClient_IsCloud(t *testing.T) {
	client := NewClient("https://company.atlassian.net", nil)
	assert.True(t, client.IsCloud())

	client2 := NewClient("https://jira.internal.corp", nil)
	assert.False(t, client2.IsCloud())
}

func TestClient_JiraAPIPath(t *testing.T) {
	cloud := NewClient("https://x.atlassian.net", nil)
	assert.Equal(t, "/rest/api/3/issue", cloud.JiraAPIPath("issue"))

	server := NewClient("https://jira.corp.com", nil)
	assert.Equal(t, "/rest/api/2/issue", server.JiraAPIPath("issue"))
}

func TestClient_Get_NilDest_ClosesBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "tok"}
	client := NewClient(server.URL, auth)

	resp, err := client.Get("/test", nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	// If resp.Body were not closed, this would leak; the test validates no error path.
}

func TestClient_Delete_NilDest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "tok"}
	client := NewClient(server.URL, auth)

	resp, err := client.Delete("/test/123", nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
