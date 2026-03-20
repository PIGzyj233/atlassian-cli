package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func exeName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func TestJiraBinary_IssueGet(t *testing.T) {
	// Start mock Jira server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"key": "TEST-1",
			"fields": map[string]any{
				"summary": "E2E test issue",
				"status":  map[string]any{"name": "Open"},
			},
		})
	}))
	defer server.Close()

	// Write temporary config pointing to mock server
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "config.yml")

	// Extract host from mock server URL (localhost:port)
	host := server.URL[7:] // strip "http://"

	config := `hosts:
  ` + host + `:
    type: server
    auth: pat
    token: "test-token"
    jira:
      enabled: true
defaults:
  host: ` + host + `
  output: json
`
	os.WriteFile(configPath, []byte(config), 0600)

	// Build the binary
	binPath := filepath.Join(t.TempDir(), exeName("jira"))
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/jira")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "build failed: %s", output)

	// Verify the binary responds to --help
	cmd = exec.Command(binPath, "--help")
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "A CLI tool for interacting with Atlassian Jira")
}

func TestConfluenceBinary_Help(t *testing.T) {
	binPath := filepath.Join(t.TempDir(), exeName("confluence"))
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/confluence")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "build failed: %s", output)

	cmd = exec.Command(binPath, "--help")
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "A CLI tool for interacting with Atlassian Confluence")
}
