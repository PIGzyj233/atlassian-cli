package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_ValidYAML(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yml")

	yaml := `hosts:
  company.atlassian.net:
    type: cloud
    auth: basic
    username: user@company.com
    token: "ATATT3xFake"
    jira:
      enabled: true
      base_path: ""
    confluence:
      enabled: true
      base_path: /wiki
defaults:
  host: company.atlassian.net
  output: json
`
	err := os.WriteFile(configPath, []byte(yaml), 0600)
	require.NoError(t, err)

	cfg, err := LoadFromPath(configPath)
	require.NoError(t, err)

	assert.Equal(t, "company.atlassian.net", cfg.Defaults.Host)
	assert.Equal(t, "json", cfg.Defaults.Output)

	host, ok := cfg.Hosts["company.atlassian.net"]
	require.True(t, ok)
	assert.Equal(t, "cloud", host.Type)
	assert.Equal(t, "basic", host.Auth)
	assert.Equal(t, "user@company.com", host.Username)
	assert.Equal(t, "ATATT3xFake", host.Token)
	assert.True(t, host.Jira.Enabled)
	assert.True(t, host.Confluence.Enabled)
	assert.Equal(t, "/wiki", host.Confluence.BasePath)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadFromPath("/nonexistent/config.yml")
	assert.Error(t, err)
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yml")
	err := os.WriteFile(configPath, []byte(""), 0600)
	require.NoError(t, err)

	cfg, err := LoadFromPath(configPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Empty(t, cfg.Hosts)
}

func TestConfig_Save(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "subdir", "config.yml")

	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"test.atlassian.net": {
				Type:     "cloud",
				Auth:     "basic",
				Username: "test@test.com",
				Token:    "token123",
			},
		},
		Defaults: Defaults{
			Host:   "test.atlassian.net",
			Output: "json",
		},
	}

	err := cfg.Save(configPath)
	require.NoError(t, err)

	info, err := os.Stat(configPath)
	require.NoError(t, err)
	assert.NotZero(t, info.Size())

	loaded, err := LoadFromPath(configPath)
	require.NoError(t, err)
	assert.Equal(t, "test.atlassian.net", loaded.Defaults.Host)
	assert.Equal(t, "token123", loaded.Hosts["test.atlassian.net"].Token)
}

func TestConfig_DefaultHost(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"a.atlassian.net": {Type: "cloud"},
		},
		Defaults: Defaults{Host: "a.atlassian.net"},
	}

	host, name, err := cfg.DefaultHost()
	require.NoError(t, err)
	assert.Equal(t, "a.atlassian.net", name)
	assert.Equal(t, "cloud", host.Type)
}

func TestConfig_DefaultHost_SingleHost(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"only.atlassian.net": {Type: "cloud"},
		},
	}

	host, name, err := cfg.DefaultHost()
	require.NoError(t, err)
	assert.Equal(t, "only.atlassian.net", name)
	assert.Equal(t, "cloud", host.Type)
}

func TestConfig_DefaultHost_NoDefault(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"a.atlassian.net": {Type: "cloud"},
			"b.atlassian.net": {Type: "server"},
		},
	}

	_, _, err := cfg.DefaultHost()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default host")
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("ATLASSIAN_HOST", "env.atlassian.net")
	t.Setenv("ATLASSIAN_USERNAME", "envuser@test.com")
	t.Setenv("ATLASSIAN_TOKEN", "env-token-123")
	t.Setenv("ATLASSIAN_AUTH_TYPE", "basic")

	host, hostname := HostFromEnv()
	require.NotNil(t, host)
	assert.Equal(t, "env.atlassian.net", hostname)
	assert.Equal(t, "envuser@test.com", host.Username)
	assert.Equal(t, "env-token-123", host.Token)
	assert.Equal(t, "basic", host.Auth)
}

func TestConfigFromEnv_NotSet(t *testing.T) {
	// Ensure env vars are clean (t.Setenv restores after test)
	host, hostname := HostFromEnv()
	assert.Nil(t, host)
	assert.Empty(t, hostname)
}
