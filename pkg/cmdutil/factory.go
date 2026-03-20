package cmdutil

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/config"
)

// Factory holds shared dependencies for all commands.
type Factory struct {
	Config     *config.Config
	ConfigPath string
	Version    string
}

// NewFactory creates a Factory, loading config from the default path.
func NewFactory(version string) *Factory {
	configPath := config.DefaultConfigPath()
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{
			Hosts: make(map[string]*config.HostConfig),
		}
	}
	return &Factory{
		Config:     cfg,
		ConfigPath: configPath,
		Version:    version,
	}
}

// JiraClient creates an API client for Jira, resolving the active host.
func (f *Factory) JiraClient(flagHost string) (*api.Client, error) {
	host, hostname, err := config.ResolveHost(f.Config, flagHost)
	if err != nil {
		return nil, err
	}

	auth, err := config.NewAuthenticator(host)
	if err != nil {
		return nil, fmt.Errorf("auth for %s: %w", hostname, err)
	}

	baseURL := "https://" + hostname + host.Jira.BasePath
	return api.NewClient(baseURL, auth), nil
}

// ConfluenceClient creates an API client for Confluence, resolving the active host.
func (f *Factory) ConfluenceClient(flagHost string) (*api.Client, error) {
	host, hostname, err := config.ResolveHost(f.Config, flagHost)
	if err != nil {
		return nil, err
	}

	auth, err := config.NewAuthenticator(host)
	if err != nil {
		return nil, fmt.Errorf("auth for %s: %w", hostname, err)
	}

	basePath := host.Confluence.BasePath
	if basePath == "" && config.IsCloudHost(hostname) {
		basePath = "/wiki"
	}

	baseURL := "https://" + hostname + basePath
	return api.NewClient(baseURL, auth), nil
}
