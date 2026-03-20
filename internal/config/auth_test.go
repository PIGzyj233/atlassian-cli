package config

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicAuth_Apply(t *testing.T) {
	auth := &BasicAuth{Username: "user@test.com", Token: "api-token"}
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	err := auth.Apply(req)
	require.NoError(t, err)

	username, password, ok := req.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, "user@test.com", username)
	assert.Equal(t, "api-token", password)
}

func TestPATAuth_Apply(t *testing.T) {
	auth := &PATAuth{Token: "personal-token-123"}
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	err := auth.Apply(req)
	require.NoError(t, err)

	assert.Equal(t, "Bearer personal-token-123", req.Header.Get("Authorization"))
}

func TestNewAuthenticator_Basic(t *testing.T) {
	host := &HostConfig{
		Auth:     "basic",
		Username: "user@test.com",
		Token:    "token",
	}
	auth, err := NewAuthenticator(host)
	require.NoError(t, err)

	_, ok := auth.(*BasicAuth)
	assert.True(t, ok)
}

func TestNewAuthenticator_PAT(t *testing.T) {
	host := &HostConfig{
		Auth:  "pat",
		Token: "personal-token",
	}
	auth, err := NewAuthenticator(host)
	require.NoError(t, err)

	_, ok := auth.(*PATAuth)
	assert.True(t, ok)
}

func TestNewAuthenticator_MissingCredentials(t *testing.T) {
	host := &HostConfig{Auth: "basic"}
	_, err := NewAuthenticator(host)
	assert.Error(t, err)
}
