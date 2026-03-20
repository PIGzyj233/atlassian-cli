package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveHost_FlagOverride(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"a.atlassian.net": {Type: "cloud", Auth: "basic", Username: "u", Token: "t"},
			"b.atlassian.net": {Type: "server", Auth: "pat", Token: "p"},
		},
		Defaults: Defaults{Host: "a.atlassian.net"},
	}

	host, name, err := ResolveHost(cfg, "b.atlassian.net")
	require.NoError(t, err)
	assert.Equal(t, "b.atlassian.net", name)
	assert.Equal(t, "server", host.Type)
}

func TestResolveHost_EnvOverride(t *testing.T) {
	t.Setenv("ATLASSIAN_HOST", "env.atlassian.net")
	t.Setenv("ATLASSIAN_TOKEN", "env-tok")
	t.Setenv("ATLASSIAN_AUTH_TYPE", "pat")

	cfg := &Config{
		Hosts:    map[string]*HostConfig{},
		Defaults: Defaults{Host: "other.atlassian.net"},
	}

	host, name, err := ResolveHost(cfg, "")
	require.NoError(t, err)
	assert.Equal(t, "env.atlassian.net", name)
	assert.Equal(t, "env-tok", host.Token)
}

func TestResolveHost_ConfigDefault(t *testing.T) {
	cfg := &Config{
		Hosts: map[string]*HostConfig{
			"default.atlassian.net": {Type: "cloud", Auth: "basic", Username: "u", Token: "t"},
		},
		Defaults: Defaults{Host: "default.atlassian.net"},
	}

	host, name, err := ResolveHost(cfg, "")
	require.NoError(t, err)
	assert.Equal(t, "default.atlassian.net", name)
	assert.Equal(t, "cloud", host.Type)
}
