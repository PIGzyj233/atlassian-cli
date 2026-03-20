package config

import (
	"fmt"
	"net/http"
)

// Authenticator applies authentication to an HTTP request.
type Authenticator interface {
	Apply(req *http.Request) error
}

// BasicAuth uses HTTP Basic Authentication (email + API token).
type BasicAuth struct {
	Username string
	Token    string
}

// Apply sets the Basic Auth header.
func (a *BasicAuth) Apply(req *http.Request) error {
	req.SetBasicAuth(a.Username, a.Token)
	return nil
}

// PATAuth uses a Personal Access Token (Bearer token).
type PATAuth struct {
	Token string
}

// Apply sets the Bearer token header.
func (a *PATAuth) Apply(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+a.Token)
	return nil
}

// NewAuthenticator creates the appropriate Authenticator from a HostConfig.
func NewAuthenticator(host *HostConfig) (Authenticator, error) {
	switch host.Auth {
	case "basic":
		if host.Username == "" || host.Token == "" {
			return nil, fmt.Errorf("basic auth requires both username and token")
		}
		return &BasicAuth{Username: host.Username, Token: host.Token}, nil
	case "pat":
		if host.Token == "" {
			return nil, fmt.Errorf("PAT auth requires a token")
		}
		return &PATAuth{Token: host.Token}, nil
	default:
		return nil, fmt.Errorf("unsupported auth type: %q", host.Auth)
	}
}
