package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration structure.
type Config struct {
	Hosts    map[string]*HostConfig `yaml:"hosts"`
	Defaults Defaults               `yaml:"defaults"`
}

// HostConfig represents a single Atlassian instance.
type HostConfig struct {
	Type       string        `yaml:"type"`       // cloud | server | datacenter
	Auth       string        `yaml:"auth"`       // basic | pat
	Username   string        `yaml:"username"`
	Token      string        `yaml:"token"`
	SSLVerify  *bool         `yaml:"ssl_verify"` // pointer to distinguish unset from false
	Jira       ServiceConfig `yaml:"jira"`
	Confluence ServiceConfig `yaml:"confluence"`
}

// ServiceConfig holds per-service settings for a host.
type ServiceConfig struct {
	Enabled  bool   `yaml:"enabled"`
	BasePath string `yaml:"base_path"`
}

// Defaults holds global default settings.
type Defaults struct {
	Host       string          `yaml:"host"`
	Output     string          `yaml:"output"`
	Jira       ServiceDefaults `yaml:"jira"`
	Confluence ServiceDefaults `yaml:"confluence"`
}

// ServiceDefaults holds default values for a service.
type ServiceDefaults struct {
	Project string `yaml:"project"` // jira default project
	Space   string `yaml:"space"`   // confluence default space
}

// ConfigDir returns the platform-specific config directory.
func ConfigDir() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "atlassian-cli")
		}
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "atlassian-cli")
}

// DefaultConfigPath returns the full path to the config file.
func DefaultConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yml")
}

// LoadFromPath reads and parses a config file from the given path.
func LoadFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	cfg := &Config{
		Hosts: make(map[string]*HostConfig),
	}
	if len(data) == 0 {
		return cfg, nil
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}
	if cfg.Hosts == nil {
		cfg.Hosts = make(map[string]*HostConfig)
	}
	return cfg, nil
}

// Load reads config from the default path.
func Load() (*Config, error) {
	return LoadFromPath(DefaultConfigPath())
}

// Save writes config to the given path with 0600 permissions.
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// DefaultHost returns the active host config, resolving defaults.
func (c *Config) DefaultHost() (*HostConfig, string, error) {
	if c.Defaults.Host == "" {
		if len(c.Hosts) == 1 {
			for name, host := range c.Hosts {
				return host, name, nil
			}
		}
		return nil, "", fmt.Errorf("no default host configured; run `auth login` or set defaults.host in config")
	}
	host, ok := c.Hosts[c.Defaults.Host]
	if !ok {
		return nil, "", fmt.Errorf("default host %q not found in config", c.Defaults.Host)
	}
	return host, c.Defaults.Host, nil
}
