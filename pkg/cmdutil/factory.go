package cmdutil

import "github.com/PigZyj2333/atlassian-cli/internal/config"

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
		// Config may not exist yet (first run); use empty config
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
