package cmdutil

import (
	"testing"

	"github.com/PigZyj2333/atlassian-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactoryJiraClient_UsesHostBasePathAndAuth(t *testing.T) {
	cfg := &config.Config{
		Hosts: map[string]*config.HostConfig{
			"jira.example.com": {
				Type:  "server",
				Auth:  "pat",
				Token: "secret",
				Jira: config.ServiceConfig{
					Enabled:  true,
					BasePath: "/jira",
				},
			},
		},
		Defaults: config.Defaults{Host: "jira.example.com"},
	}

	f := &Factory{Config: cfg}

	client, err := f.JiraClient("")
	require.NoError(t, err)

	assert.Equal(t, "https://jira.example.com/jira", client.BaseURL)
	assert.IsType(t, &config.PATAuth{}, client.Auth)
}

func TestFactoryConfluenceClient_DefaultsCloudBasePathToWiki(t *testing.T) {
	cfg := &config.Config{
		Hosts: map[string]*config.HostConfig{
			"company.atlassian.net": {
				Type:     "cloud",
				Auth:     "basic",
				Token:    "secret",
				Username: "user@example.com",
				Confluence: config.ServiceConfig{
					Enabled: true,
				},
			},
		},
		Defaults: config.Defaults{Host: "company.atlassian.net"},
	}

	f := &Factory{Config: cfg}

	client, err := f.ConfluenceClient("")
	require.NoError(t, err)

	assert.Equal(t, "https://company.atlassian.net/wiki", client.BaseURL)
	assert.IsType(t, &config.BasicAuth{}, client.Auth)
}

func TestGetHost_UsesInheritedPersistentFlag(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	root.PersistentFlags().String("host", "", "host")
	require.NoError(t, root.PersistentFlags().Set("host", "example.atlassian.net"))

	child := &cobra.Command{Use: "child"}
	root.AddCommand(child)

	assert.Equal(t, "example.atlassian.net", GetHost(child))
}

func TestGetOutput_DefaultsToJSON(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}

	assert.Equal(t, "json", GetOutput(cmd))
}
