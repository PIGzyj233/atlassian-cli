package config

import "fmt"

// ResolveHost determines the active host using priority:
// 1. --host flag (if non-empty)
// 2. Environment variables (ATLASSIAN_HOST)
// 3. Config file defaults
func ResolveHost(cfg *Config, flagHost string) (*HostConfig, string, error) {
	// Priority 1: explicit --host flag
	if flagHost != "" {
		host, ok := cfg.Hosts[flagHost]
		if !ok {
			return nil, "", fmt.Errorf("host %q not found in config", flagHost)
		}
		return host, flagHost, nil
	}

	// Priority 2: environment variables
	if envHost, envName := HostFromEnv(); envHost != nil {
		return envHost, envName, nil
	}

	// Priority 3: config default
	return cfg.DefaultHost()
}
