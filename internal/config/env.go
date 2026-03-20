package config

import "os"

// HostFromEnv creates a HostConfig from environment variables.
// Returns nil if ATLASSIAN_HOST is not set.
func HostFromEnv() (*HostConfig, string) {
	hostname := os.Getenv("ATLASSIAN_HOST")
	if hostname == "" {
		return nil, ""
	}

	authType := os.Getenv("ATLASSIAN_AUTH_TYPE")
	if authType == "" {
		authType = "basic"
	}

	host := &HostConfig{
		Auth:       authType,
		Username:   os.Getenv("ATLASSIAN_USERNAME"),
		Token:      os.Getenv("ATLASSIAN_TOKEN"),
		Jira:       ServiceConfig{Enabled: true},
		Confluence: ServiceConfig{Enabled: true},
	}

	// Detect instance type from hostname
	if IsCloudHost(hostname) {
		host.Type = "cloud"
		host.Confluence.BasePath = "/wiki"
	} else {
		host.Type = "server"
	}

	return host, hostname
}

// IsCloudHost checks if a hostname is an Atlassian Cloud instance.
func IsCloudHost(hostname string) bool {
	// Match: *.atlassian.net
	return len(hostname) > 14 && hostname[len(hostname)-14:] == ".atlassian.net"
}
