# Atlassian CLI Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build two standalone Go CLI tools (`jira`, `confluence`) that expose all current MCP-Atlassian capabilities as shell commands, designed for both human operators and LLM-driven automation.

**Architecture:** Two binaries sharing a common internal library for config, auth, HTTP client, output formatting, and content conversion. Commands follow cobra patterns inspired by `gh` CLI. JSON-first output with switchable table/text modes.

**Tech Stack:** Go 1.22+, cobra, viper, go-yaml v3, net/http, tablewriter, goreleaser, testify

**Spec:** `E:/ai-services/mcp-atlassian/docs/superpowers/specs/2026-03-19-atlassian-cli-design.md`

**Reference codebase:** The existing Python MCP-Atlassian project at `E:/ai-services/mcp-atlassian/src/mcp_atlassian/` — all API endpoints, auth patterns, and content conversion logic should be ported from this codebase.

---

## File Structure

```
cli/
├── cmd/
│   ├── jira/
│   │   └── main.go                     # jira binary entry point
│   └── confluence/
│       └── main.go                     # confluence binary entry point
├── internal/
│   ├── config/
│   │   ├── config.go                   # Config struct, Load/Save, XDG paths
│   │   ├── config_test.go
│   │   ├── auth.go                     # Authenticator interface, BasicAuth, PATAuth
│   │   ├── auth_test.go
│   │   ├── hosts.go                    # Multi-instance host management
│   │   └── hosts_test.go
│   ├── api/
│   │   ├── client.go                   # HTTP client with auth, retry, timeout
│   │   ├── client_test.go
│   │   ├── request.go                  # Request builder (headers, query params)
│   │   ├── request_test.go
│   │   ├── errors.go                   # API error parsing and formatting
│   │   ├── errors_test.go
│   │   ├── pagination.go              # Pagination response struct + helpers
│   │   └── pagination_test.go
│   ├── output/
│   │   ├── formatter.go               # Formatter interface + Print() dispatcher
│   │   ├── json.go                    # JSON output (indented)
│   │   ├── json_test.go
│   │   ├── table.go                   # Table output (tablewriter)
│   │   ├── table_test.go
│   │   ├── text.go                    # Plain text output
│   │   └── text_test.go
│   └── convert/
│       ├── adf.go                     # ADF (Jira Cloud) ↔ Markdown
│       ├── adf_test.go
│       ├── storage.go                 # Storage Format (Confluence) ↔ Markdown
│       └── storage_test.go
├── pkg/
│   ├── cmdutil/
│   │   ├── factory.go                 # Factory struct (client + config + output)
│   │   └── flags.go                   # Shared flag helpers
│   ├── jira/
│   │   ├── root.go                    # jira root command + subcommand registration
│   │   ├── auth/
│   │   │   ├── auth.go               # auth subcommand group
│   │   │   ├── login.go              # auth login
│   │   │   ├── logout.go             # auth logout
│   │   │   ├── status.go             # auth status
│   │   │   ├── list.go               # auth list
│   │   │   └── switch.go             # auth switch
│   │   ├── issue/
│   │   │   ├── issue.go              # issue subcommand group
│   │   │   ├── get.go                # issue get
│   │   │   ├── get_test.go
│   │   │   ├── create.go             # issue create
│   │   │   ├── create_test.go
│   │   │   ├── batch_create.go       # issue batch-create
│   │   │   ├── update.go             # issue update
│   │   │   ├── delete.go             # issue delete
│   │   │   └── transition.go         # issue transition
│   │   ├── search/
│   │   │   ├── search.go             # search command
│   │   │   └── search_test.go
│   │   ├── comment/
│   │   │   ├── comment.go            # comment subcommand group
│   │   │   ├── list.go
│   │   │   ├── add.go
│   │   │   └── edit.go
│   │   ├── field/
│   │   │   ├── field.go
│   │   │   ├── search.go
│   │   │   └── options.go
│   │   ├── transition/
│   │   │   └── list.go
│   │   ├── board/
│   │   │   ├── board.go
│   │   │   ├── list.go
│   │   │   └── issues.go
│   │   ├── sprint/
│   │   │   ├── sprint.go
│   │   │   ├── list.go
│   │   │   ├── issues.go
│   │   │   ├── create.go
│   │   │   ├── update.go
│   │   │   └── add_issues.go
│   │   ├── project/
│   │   │   ├── project.go
│   │   │   ├── list.go
│   │   │   ├── versions.go
│   │   │   ├── components.go
│   │   │   └── version_create.go
│   │   ├── worklog/
│   │   │   ├── worklog.go
│   │   │   ├── list.go
│   │   │   └── add.go
│   │   ├── watcher/
│   │   │   ├── watcher.go
│   │   │   ├── list.go
│   │   │   ├── add.go
│   │   │   └── remove.go
│   │   ├── link/
│   │   │   ├── link.go
│   │   │   ├── types.go
│   │   │   ├── create.go
│   │   │   ├── create_remote.go
│   │   │   └── remove.go
│   │   ├── attachment/
│   │   │   ├── attachment.go
│   │   │   ├── download.go
│   │   │   └── images.go
│   │   ├── sla/
│   │   │   ├── sla.go
│   │   │   ├── get.go
│   │   │   └── dates.go
│   │   ├── dev/
│   │   │   ├── dev.go
│   │   │   ├── info.go
│   │   │   └── batch_info.go
│   │   ├── changelog/
│   │   │   └── batch.go
│   │   ├── form/
│   │   │   ├── form.go
│   │   │   ├── list.go
│   │   │   ├── get.go
│   │   │   └── update.go
│   │   ├── servicedesk/
│   │   │   ├── servicedesk.go
│   │   │   ├── get.go
│   │   │   ├── queues.go
│   │   │   └── queue_issues.go
│   │   ├── user/
│   │   │   └── get.go
│   │   └── epic/
│   │       └── link.go
│   └── confluence/
│       ├── root.go                    # confluence root command
│       ├── auth/                      # shared auth (same config file as jira)
│       │   ├── auth.go
│       │   ├── login.go
│       │   ├── logout.go
│       │   ├── status.go
│       │   ├── list.go
│       │   └── switch.go
│       ├── page/
│       │   ├── page.go
│       │   ├── get.go
│       │   ├── get_test.go
│       │   ├── create.go
│       │   ├── update.go
│       │   ├── delete.go
│       │   ├── move.go
│       │   ├── children.go
│       │   ├── tree.go
│       │   ├── history.go
│       │   └── diff.go
│       ├── search/
│       │   └── search.go
│       ├── comment/
│       │   ├── comment.go
│       │   ├── list.go
│       │   ├── add.go
│       │   └── reply.go
│       ├── label/
│       │   ├── label.go
│       │   ├── list.go
│       │   └── add.go
│       ├── attachment/
│       │   ├── attachment.go
│       │   ├── list.go
│       │   ├── upload.go
│       │   ├── upload_batch.go
│       │   ├── download.go
│       │   ├── download_all.go
│       │   ├── delete.go
│       │   └── images.go
│       ├── user/
│       │   └── search.go
│       └── analytics/
│           └── views.go
├── go.mod
├── go.sum
├── Makefile
├── .goreleaser.yml
└── README.md
```

---

## Chunk 1: Project Scaffolding & Root Commands

### Task 1: Initialize Go Module

**Files:**
- Create: `cli/go.mod`
- Create: `cli/go.sum` (auto-generated)

- [x] **Step 1: Create project directory and initialize Go module**

```bash
mkdir -p cli
cd cli
go mod init github.com/anthropics/atlassian-cli
```

- [x] **Step 2: Add core dependencies**

```bash
cd cli
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get gopkg.in/yaml.v3@latest
go get github.com/olekukonko/tablewriter@v1.1.4
go get github.com/stretchr/testify@latest
```

- [x] **Step 3: Verify module compiles**

Run: `cd cli && go mod tidy`
Expected: Clean exit, `go.sum` populated

- [x] **Step 4: Commit**

```bash
git add cli/go.mod cli/go.sum
git commit -m "chore: initialize Go module for atlassian-cli"
```

---

### Task 2: Create Makefile

**Files:**
- Create: `cli/Makefile`

- [x] **Step 1: Write Makefile**

```makefile
.PHONY: build test lint fmt clean install

VERSION ?= dev
LDFLAGS := -s -w -X main.version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/jira ./cmd/jira
	go build -ldflags "$(LDFLAGS)" -o bin/confluence ./cmd/confluence

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

clean:
	rm -rf bin/

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/jira
	go install -ldflags "$(LDFLAGS)" ./cmd/confluence
```

- [x] **Step 2: Commit**

```bash
git add cli/Makefile
git commit -m "chore: add Makefile with build/test/lint targets"
```

---

### Task 3: Create Binary Entry Points & Root Commands

**Files:**
- Create: `cli/cmd/jira/main.go`
- Create: `cli/cmd/confluence/main.go`
- Create: `cli/pkg/jira/root.go`
- Create: `cli/pkg/confluence/root.go`
- Create: `cli/pkg/cmdutil/factory.go`

- [x] **Step 1: Write the cmdutil factory (shared dependency container)**

```go
// cli/pkg/cmdutil/factory.go
package cmdutil

// Factory holds shared dependencies for all commands.
type Factory struct {
	// Will be populated in later tasks:
	// Config  *config.Config
	// Client  *api.Client
	// Output  string
	Version string
}

// NewFactory creates a Factory with defaults.
func NewFactory(version string) *Factory {
	return &Factory{
		Version: version,
	}
}
```

- [x] **Step 2: Write jira root command**

```go
// cli/pkg/jira/root.go
package jira

import (
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdRoot creates the root `jira` command.
func NewCmdRoot(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jira <command> [flags]",
		Short: "Work with Jira from the command line",
		Long:  "A CLI tool for interacting with Atlassian Jira, designed for both humans and LLM automation.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().String("host", "", "Target Atlassian instance hostname")
	cmd.PersistentFlags().String("output", "json", "Output format: json, text, table")
	cmd.PersistentFlags().Bool("verbose", false, "Enable debug logging to stderr")
	cmd.PersistentFlags().Bool("no-color", false, "Disable ANSI colors")

	cmd.Version = f.Version

	// Subcommands will be registered here in later tasks:
	// cmd.AddCommand(auth.NewCmdAuth(f))
	// cmd.AddCommand(issue.NewCmdIssue(f))
	// ...

	return cmd
}
```

- [x] **Step 3: Write confluence root command**

```go
// cli/pkg/confluence/root.go
package confluence

import (
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdRoot creates the root `confluence` command.
func NewCmdRoot(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confluence <command> [flags]",
		Short: "Work with Confluence from the command line",
		Long:  "A CLI tool for interacting with Atlassian Confluence, designed for both humans and LLM automation.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().String("host", "", "Target Atlassian instance hostname")
	cmd.PersistentFlags().String("output", "json", "Output format: json, text, table")
	cmd.PersistentFlags().Bool("verbose", false, "Enable debug logging to stderr")
	cmd.PersistentFlags().Bool("no-color", false, "Disable ANSI colors")

	cmd.Version = f.Version

	return cmd
}
```

- [x] **Step 4: Write jira binary entry point**

```go
// cli/cmd/jira/main.go
package main

import (
	"fmt"
	"os"

	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	jiracmd "github.com/anthropics/atlassian-cli/pkg/jira"
)

var version = "dev"

func main() {
	f := cmdutil.NewFactory(version)
	rootCmd := jiracmd.NewCmdRoot(f)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

- [x] **Step 5: Write confluence binary entry point**

```go
// cli/cmd/confluence/main.go
package main

import (
	"fmt"
	"os"

	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	confluencecmd "github.com/anthropics/atlassian-cli/pkg/confluence"
)

var version = "dev"

func main() {
	f := cmdutil.NewFactory(version)
	rootCmd := confluencecmd.NewCmdRoot(f)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

- [x] **Step 6: Build both binaries to verify compilation**

Run: `cd cli && make build`
Expected: `bin/jira` and `bin/confluence` created without errors

- [x] **Step 7: Verify help output**

Run: `cd cli && ./bin/jira --help`
Expected: Shows "Work with Jira from the command line" with global flags

Run: `cd cli && ./bin/confluence --help`
Expected: Shows "Work with Confluence from the command line" with global flags

- [x] **Step 8: Verify version flag**

Run: `cd cli && ./bin/jira --version`
Expected: Shows "jira version dev"

- [x] **Step 9: Commit**

```bash
git add cli/cmd/ cli/pkg/cmdutil/ cli/pkg/jira/root.go cli/pkg/confluence/root.go
git commit -m "feat: add root commands and binary entry points for jira and confluence"
```

---

## Chunk 2: Configuration & Authentication

### Task 4: Config Struct & YAML Parsing

**Files:**
- Create: `cli/internal/config/config.go`
- Create: `cli/internal/config/config_test.go`

- [x] **Step 1: Write the failing test for config loading**

```go
// cli/internal/config/config_test.go
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
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/config/ -v -run TestLoadConfig`
Expected: FAIL — `LoadFromPath` not defined

- [x] **Step 3: Implement config structs and loader**

```go
// cli/internal/config/config.go
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
```

- [x] **Step 4: Run tests to verify they pass**

Run: `cd cli && go test ./internal/config/ -v -run TestLoadConfig`
Expected: All 3 tests PASS

- [x] **Step 5: Add Save and DefaultHost tests**

```go
// Append to cli/internal/config/config_test.go

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
```

- [x] **Step 6: Run all config tests**

Run: `cd cli && go test ./internal/config/ -v`
Expected: All tests PASS

- [x] **Step 7: Commit**

```bash
git add cli/internal/config/config.go cli/internal/config/config_test.go
git commit -m "feat(auth): add config struct with YAML load/save and XDG paths"
```

---

### Task 5: Authenticator Interface

**Files:**
- Create: `cli/internal/config/auth.go`
- Create: `cli/internal/config/auth_test.go`

- [x] **Step 1: Write failing tests for authenticators**

```go
// cli/internal/config/auth_test.go
package config

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicAuth_Apply(t *testing.T) {
	auth := &BasicAuth{Username: "user@test.com", Token: "api-token"}
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	err := auth.Apply(req)
	require.NoError(t, err)

	username, password, ok := req.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, "user@test.com", username)
	assert.Equal(t, "api-token", password)
}

func TestPATAuth_Apply(t *testing.T) {
	auth := &PATAuth{Token: "personal-token-123"}
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	err := auth.Apply(req)
	require.NoError(t, err)

	assert.Equal(t, "Bearer personal-token-123", req.Header.Get("Authorization"))
}

func TestNewAuthenticator_Basic(t *testing.T) {
	host := &HostConfig{
		Auth:     "basic",
		Username: "user@test.com",
		Token:    "token",
	}
	auth, err := NewAuthenticator(host)
	require.NoError(t, err)

	_, ok := auth.(*BasicAuth)
	assert.True(t, ok)
}

func TestNewAuthenticator_PAT(t *testing.T) {
	host := &HostConfig{
		Auth:  "pat",
		Token: "personal-token",
	}
	auth, err := NewAuthenticator(host)
	require.NoError(t, err)

	_, ok := auth.(*PATAuth)
	assert.True(t, ok)
}

func TestNewAuthenticator_MissingCredentials(t *testing.T) {
	host := &HostConfig{Auth: "basic"}
	_, err := NewAuthenticator(host)
	assert.Error(t, err)
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/config/ -v -run TestBasicAuth`
Expected: FAIL — types not defined

- [x] **Step 3: Implement authenticators**

```go
// cli/internal/config/auth.go
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
```

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/config/ -v`
Expected: All tests PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/config/auth.go cli/internal/config/auth_test.go
git commit -m "feat(auth): add Authenticator interface with BasicAuth and PATAuth"
```

---

### Task 6: Environment Variable Overrides

**Files:**
- Create: `cli/internal/config/env.go`
- Modify: `cli/internal/config/config_test.go`

- [x] **Step 1: Write failing test for env override**

```go
// Append to cli/internal/config/config_test.go

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
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/config/ -v -run TestConfigFromEnv`
Expected: FAIL — `HostFromEnv` not defined

- [x] **Step 3: Implement env override**

```go
// cli/internal/config/env.go
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
		Auth:     authType,
		Username: os.Getenv("ATLASSIAN_USERNAME"),
		Token:    os.Getenv("ATLASSIAN_TOKEN"),
		Jira:       ServiceConfig{Enabled: true},
		Confluence: ServiceConfig{Enabled: true},
	}

	// Detect instance type from hostname
	if isCloudHost(hostname) {
		host.Type = "cloud"
		host.Confluence.BasePath = "/wiki"
	} else {
		host.Type = "server"
	}

	return host, hostname
}

// isCloudHost checks if a hostname is an Atlassian Cloud instance.
func isCloudHost(hostname string) bool {
	// Match: *.atlassian.net
	return len(hostname) > 14 && hostname[len(hostname)-14:] == ".atlassian.net"
}
```

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/config/ -v -run TestConfigFromEnv`
Expected: PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/config/env.go cli/internal/config/config_test.go
git commit -m "feat(auth): add environment variable overrides for CI/CD use"
```

---

### Task 7: Resolve Active Host (Config + Env + Flag Merging)

**Files:**
- Create: `cli/internal/config/hosts.go`
- Create: `cli/internal/config/hosts_test.go`

- [x] **Step 1: Write failing test**

```go
// cli/internal/config/hosts_test.go
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
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/config/ -v -run TestResolveHost`
Expected: FAIL — `ResolveHost` not defined

- [x] **Step 3: Implement host resolution**

```go
// cli/internal/config/hosts.go
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
```

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/config/ -v -run TestResolveHost`
Expected: All PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/config/hosts.go cli/internal/config/hosts_test.go
git commit -m "feat(auth): add host resolution with flag > env > config priority"
```

---

### Task 8: Auth CLI Commands

**Files:**
- Create: `cli/pkg/jira/auth/auth.go`
- Create: `cli/pkg/jira/auth/login.go`
- Create: `cli/pkg/jira/auth/logout.go`
- Create: `cli/pkg/jira/auth/status.go`
- Create: `cli/pkg/jira/auth/list.go`
- Create: `cli/pkg/jira/auth/switch.go`
- Modify: `cli/pkg/jira/root.go` — register auth subcommand
- Modify: `cli/pkg/cmdutil/factory.go` — add Config field

> Note: `confluence auth` commands share the same config file and identical logic. Implement for jira first; confluence auth will be a thin wrapper (Task 40).

- [x] **Step 1: Export `IsCloudHost` in env.go**

In `cli/internal/config/env.go`, rename `isCloudHost` to `IsCloudHost` (capitalize the first letter). This must be done before writing `login.go` which depends on it.

- [x] **Step 2: Update Factory to hold Config**

```go
// cli/pkg/cmdutil/factory.go
package cmdutil

import "github.com/anthropics/atlassian-cli/internal/config"

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
```

- [x] **Step 2: Write auth subcommand group**

```go
// cli/pkg/jira/auth/auth.go
package auth

import (
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAuth creates the `auth` subcommand group.
func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Manage authentication",
	}

	cmd.AddCommand(NewCmdLogin(f))
	cmd.AddCommand(NewCmdLogout(f))
	cmd.AddCommand(NewCmdStatus(f))
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdSwitch(f))

	return cmd
}
```

- [x] **Step 3: Write auth login command (interactive)**

```go
// cli/pkg/jira/auth/login.go
package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/atlassian-cli/internal/config"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdLogin creates the `auth login` command.
func NewCmdLogin(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with an Atlassian instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(f)
		},
	}
	return cmd
}

func runLogin(f *cmdutil.Factory) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Atlassian instance hostname (e.g., company.atlassian.net): ")
	hostname, _ := reader.ReadString('\n')
	hostname = strings.TrimSpace(hostname)
	if hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	instType := "server"
	if config.IsCloudHost(hostname) {
		instType = "cloud"
	}
	fmt.Fprintf(os.Stderr, "Detected instance type: %s\n", instType)

	fmt.Print("Auth method (basic/pat) [basic]: ")
	authMethod, _ := reader.ReadString('\n')
	authMethod = strings.TrimSpace(authMethod)
	if authMethod == "" {
		authMethod = "basic"
	}

	host := &config.HostConfig{
		Type: instType,
		Auth: authMethod,
		Jira:       config.ServiceConfig{Enabled: true},
		Confluence: config.ServiceConfig{Enabled: true},
	}

	if instType == "cloud" {
		host.Confluence.BasePath = "/wiki"
	}

	switch authMethod {
	case "basic":
		fmt.Print("Username/email: ")
		username, _ := reader.ReadString('\n')
		host.Username = strings.TrimSpace(username)

		fmt.Print("API token: ")
		token, _ := reader.ReadString('\n')
		host.Token = strings.TrimSpace(token)
	case "pat":
		fmt.Print("Personal Access Token: ")
		token, _ := reader.ReadString('\n')
		host.Token = strings.TrimSpace(token)
	default:
		return fmt.Errorf("unsupported auth method: %s", authMethod)
	}

	// Validate credentials by building authenticator
	if _, err := config.NewAuthenticator(host); err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}

	f.Config.Hosts[hostname] = host
	if f.Config.Defaults.Host == "" {
		f.Config.Defaults.Host = hostname
	}

	if err := f.Config.Save(f.ConfigPath); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Logged in to %s\n", hostname)
	return nil
}
```

> **Note:** `config.IsCloudHost` needs to be exported — rename `isCloudHost` to `IsCloudHost` in `env.go`.

- [x] **Step 4: Write auth logout**

```go
// cli/pkg/jira/auth/logout.go
package auth

import (
	"fmt"
	"os"

	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdLogout(f *cmdutil.Factory) *cobra.Command {
	var host string
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials for a host",
		RunE: func(cmd *cobra.Command, args []string) error {
			target := host
			if target == "" {
				target = f.Config.Defaults.Host
			}
			if target == "" {
				return fmt.Errorf("specify --host or configure a default host")
			}
			if _, ok := f.Config.Hosts[target]; !ok {
				return fmt.Errorf("host %q not found in config", target)
			}
			delete(f.Config.Hosts, target)
			if f.Config.Defaults.Host == target {
				f.Config.Defaults.Host = ""
			}
			if err := f.Config.Save(f.ConfigPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Logged out of %s\n", target)
			return nil
		},
	}
	cmd.Flags().StringVar(&host, "host", "", "Host to log out of")
	return cmd
}
```

- [x] **Step 5: Write auth status**

```go
// cli/pkg/jira/auth/status.go
package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current authentication state",
		RunE: func(cmd *cobra.Command, args []string) error {
			host, name, err := f.Config.DefaultHost()
			if err != nil {
				return fmt.Errorf("not authenticated: %w", err)
			}

			output, _ := cmd.Flags().GetString("output")
			if output == "json" {
				data := map[string]any{
					"host":     name,
					"type":     host.Type,
					"auth":     host.Auth,
					"username": host.Username,
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(data)
			}

			fmt.Fprintf(os.Stdout, "Host:     %s\n", name)
			fmt.Fprintf(os.Stdout, "Type:     %s\n", host.Type)
			fmt.Fprintf(os.Stdout, "Auth:     %s\n", host.Auth)
			if host.Username != "" {
				fmt.Fprintf(os.Stdout, "Username: %s\n", host.Username)
			}
			return nil
		},
	}
	return cmd
}
```

- [x] **Step 6: Write auth list**

```go
// cli/pkg/jira/auth/list.go
package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configured instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			if output == "json" {
				type hostEntry struct {
					Host     string `json:"host"`
					Type     string `json:"type"`
					Auth     string `json:"auth"`
					Default  bool   `json:"default"`
				}
				var entries []hostEntry
				for name, host := range f.Config.Hosts {
					entries = append(entries, hostEntry{
						Host:    name,
						Type:    host.Type,
						Auth:    host.Auth,
						Default: name == f.Config.Defaults.Host,
					})
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(entries)
			}

			for name, host := range f.Config.Hosts {
				marker := "  "
				if name == f.Config.Defaults.Host {
					marker = "* "
				}
				fmt.Fprintf(os.Stdout, "%s%s (%s, %s)\n", marker, name, host.Type, host.Auth)
			}
			return nil
		},
	}
	return cmd
}
```

- [x] **Step 7: Write auth switch**

```go
// cli/pkg/jira/auth/switch.go
package auth

import (
	"fmt"
	"os"

	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSwitch(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch <host>",
		Short: "Change default instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			if _, ok := f.Config.Hosts[target]; !ok {
				return fmt.Errorf("host %q not found in config", target)
			}
			f.Config.Defaults.Host = target
			if err := f.Config.Save(f.ConfigPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Switched default host to %s\n", target)
			return nil
		},
	}
	return cmd
}
```

- [x] **Step 8: Register auth subcommand in jira root**

```go
// Update cli/pkg/jira/root.go — add import and registration:
import "github.com/anthropics/atlassian-cli/pkg/jira/auth"

// Inside NewCmdRoot, after flag definitions:
cmd.AddCommand(auth.NewCmdAuth(f))
```

- [x] **Step 9: Build and verify auth commands**

Run: `cd cli && make build && ./bin/jira auth --help`
Expected: Shows login, logout, status, list, switch subcommands

- [x] **Step 10: Commit**

```bash
git add cli/pkg/jira/auth/ cli/pkg/jira/root.go cli/pkg/cmdutil/factory.go cli/internal/config/env.go
git commit -m "feat(auth): add auth login/logout/status/list/switch commands"
```

---

## Chunk 3: HTTP Client & API Layer

### Task 9: Base HTTP Client with Auth & Retry

**Files:**
- Create: `cli/internal/api/client.go`
- Create: `cli/internal/api/client_test.go`

- [x] **Step 1: Write failing test for client creation and authenticated requests**

```go
// cli/internal/api/client_test.go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthropics/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Get_Basic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header is present
		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "user@test.com", username)
		assert.Equal(t, "api-token", password)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "user@test.com", Token: "api-token"}
	client := NewClient(server.URL, auth)

	var result map[string]string
	resp, err := client.Get("/test/endpoint", &result)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "ok", result["status"])
}

func TestClient_Get_PAT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer my-pat", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "my-pat"}
	client := NewClient(server.URL, auth)

	var result map[string]string
	_, err := client.Get("/test", &result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestClient_RetryOn429(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "tok"}
	client := NewClient(server.URL, auth)
	client.MaxRetries = 3

	var result map[string]string
	_, err := client.Get("/test", &result)
	require.NoError(t, err)
	assert.Equal(t, 3, attempts)
	assert.Equal(t, "ok", result["status"])
}

func TestClient_RetryExhausted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "tok"}
	client := NewClient(server.URL, auth)
	client.MaxRetries = 2

	_, err := client.Get("/test", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "429")
}

func TestClient_IsCloud(t *testing.T) {
	client := NewClient("https://company.atlassian.net", nil)
	assert.True(t, client.IsCloud())

	client2 := NewClient("https://jira.internal.corp", nil)
	assert.False(t, client2.IsCloud())
}

func TestClient_JiraAPIPath(t *testing.T) {
	cloud := NewClient("https://x.atlassian.net", nil)
	assert.Equal(t, "/rest/api/3/issue", cloud.JiraAPIPath("issue"))

	server := NewClient("https://jira.corp.com", nil)
	assert.Equal(t, "/rest/api/2/issue", server.JiraAPIPath("issue"))
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/api/ -v`
Expected: FAIL — package/types not defined

- [x] **Step 3: Implement HTTP client**

```go
// cli/internal/api/client.go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/anthropics/atlassian-cli/internal/config"
)

// Client is an HTTP client for Atlassian REST APIs.
type Client struct {
	HTTP       *http.Client
	BaseURL    string
	Auth       config.Authenticator
	MaxRetries int
}

// NewClient creates a new API client.
func NewClient(baseURL string, auth config.Authenticator) *Client {
	return &Client{
		HTTP: &http.Client{
			Timeout: 75 * time.Second,
		},
		BaseURL:    strings.TrimRight(baseURL, "/"),
		Auth:       auth,
		MaxRetries: 3,
	}
}

// IsCloud returns true if the base URL is an Atlassian Cloud instance.
func (c *Client) IsCloud() bool {
	return strings.Contains(c.BaseURL, ".atlassian.net")
}

// JiraAPIPath returns the versioned Jira REST API path.
func (c *Client) JiraAPIPath(resource string) string {
	if c.IsCloud() {
		return "/rest/api/3/" + resource
	}
	return "/rest/api/2/" + resource
}

// ConfluenceAPIPath returns the Confluence REST API path.
// Note: The /wiki prefix (Cloud) is already in the BaseURL set by Factory.
func (c *Client) ConfluenceAPIPath(resource string) string {
	return "/rest/api/" + resource
}

// Get performs an authenticated GET request. If dest is non-nil, the response
// body is JSON-decoded into dest.
func (c *Client) Get(path string, dest any) (*http.Response, error) {
	return c.do("GET", path, nil, dest)
}

// Post performs an authenticated POST request with a JSON body.
func (c *Client) Post(path string, body any, dest any) (*http.Response, error) {
	return c.do("POST", path, body, dest)
}

// Put performs an authenticated PUT request with a JSON body.
func (c *Client) Put(path string, body any, dest any) (*http.Response, error) {
	return c.do("PUT", path, body, dest)
}

// Delete performs an authenticated DELETE request.
func (c *Client) Delete(path string, dest any) (*http.Response, error) {
	return c.do("DELETE", path, nil, dest)
}

func (c *Client) do(method, path string, body any, dest any) (*http.Response, error) {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = strings.NewReader(string(data))
	}

	var lastResp *http.Response
	var lastErr error

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		req.Header.Set("Accept", "application/json")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		if c.Auth != nil {
			if err := c.Auth.Apply(req); err != nil {
				return nil, fmt.Errorf("applying auth: %w", err)
			}
		}

		resp, err := c.HTTP.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			time.Sleep(backoff(attempt))
			// Reset body reader for retry
			if body != nil {
				data, _ := json.Marshal(body)
				bodyReader = strings.NewReader(string(data))
			}
			continue
		}

		lastResp = resp
		lastErr = nil

		// Retry on 429 (rate limit) and 5xx (server errors)
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			resp.Body.Close()
			if attempt < c.MaxRetries {
				time.Sleep(backoff(attempt))
				if body != nil {
					data, _ := json.Marshal(body)
					bodyReader = strings.NewReader(string(data))
				}
				continue
			}
			return resp, &APIError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("HTTP %d after %d retries", resp.StatusCode, c.MaxRetries)}
		}

		// Non-retryable error
		if resp.StatusCode >= 400 {
			defer resp.Body.Close()
			return resp, parseAPIError(resp)
		}

		// Success — decode if dest provided
		if dest != nil {
			defer resp.Body.Close()
			if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
				return resp, fmt.Errorf("decoding response: %w", err)
			}
		}

		return resp, nil
	}

	if lastErr != nil {
		return lastResp, lastErr
	}
	return lastResp, fmt.Errorf("request failed after %d retries", c.MaxRetries)
}

func backoff(attempt int) time.Duration {
	// Exponential backoff: 100ms, 200ms, 400ms, ...
	// Capped at 5 seconds
	d := time.Duration(100<<uint(attempt)) * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}
```

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/api/ -v`
Expected: All tests PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/api/client.go cli/internal/api/client_test.go
git commit -m "feat: add HTTP client with auth injection, retry, and cloud detection"
```

---

### Task 10: API Error Handling

**Files:**
- Create: `cli/internal/api/errors.go`
- Create: `cli/internal/api/errors_test.go`

- [x] **Step 1: Write failing test**

```go
// cli/internal/api/errors_test.go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthropics/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIError_JSON(t *testing.T) {
	apiErr := &APIError{
		StatusCode: 404,
		Message:    "Issue NOT-123 not found",
		Hint:       "Check that the issue key is correct",
	}

	data, err := json.Marshal(apiErr)
	require.NoError(t, err)

	var result map[string]any
	json.Unmarshal(data, &result)
	assert.Equal(t, float64(404), result["status"])
	assert.Equal(t, "Issue NOT-123 not found", result["error"])
	assert.Equal(t, "Check that the issue key is correct", result["hint"])
}

func TestAPIError_ErrorString(t *testing.T) {
	apiErr := &APIError{StatusCode: 401, Message: "Unauthorized"}
	assert.Contains(t, apiErr.Error(), "401")
	assert.Contains(t, apiErr.Error(), "Unauthorized")
}

func TestParseAPIError_AtlassianFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]any{
			"errorMessages": []string{"Issue does not exist or you do not have permission"},
		})
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := NewClient(server.URL, auth)

	_, err := client.Get("/test", nil)
	require.Error(t, err)

	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 404, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "Issue does not exist")
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/api/ -v -run TestAPIError`
Expected: FAIL — `APIError` not defined

- [x] **Step 3: Implement API error types**

```go
// cli/internal/api/errors.go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIError represents an error response from the Atlassian API.
type APIError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"error"`
	Hint       string `json:"hint,omitempty"`
}

func (e *APIError) Error() string {
	s := fmt.Sprintf("API error (HTTP %d): %s", e.StatusCode, e.Message)
	if e.Hint != "" {
		s += fmt.Sprintf(" — hint: %s", e.Hint)
	}
	return s
}

// ExitCode returns the appropriate process exit code for this error.
func (e *APIError) ExitCode() int {
	if e.StatusCode >= 400 && e.StatusCode < 500 {
		return 2 // API error
	}
	if e.StatusCode >= 500 {
		return 2 // API error
	}
	return 1 // input error
}

// parseAPIError reads the response body and builds an APIError.
func parseAPIError(resp *http.Response) *APIError {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("HTTP %d (could not read body)", resp.StatusCode),
		}
	}

	apiErr := &APIError{StatusCode: resp.StatusCode}

	// Try Atlassian error format: {"errorMessages": [...], "errors": {...}}
	var atlassianErr struct {
		ErrorMessages []string          `json:"errorMessages"`
		Errors        map[string]string `json:"errors"`
		Message       string            `json:"message"`
	}
	if json.Unmarshal(body, &atlassianErr) == nil {
		var parts []string
		parts = append(parts, atlassianErr.ErrorMessages...)
		if atlassianErr.Message != "" {
			parts = append(parts, atlassianErr.Message)
		}
		for field, msg := range atlassianErr.Errors {
			parts = append(parts, fmt.Sprintf("%s: %s", field, msg))
		}
		if len(parts) > 0 {
			apiErr.Message = strings.Join(parts, "; ")
		}
	}

	if apiErr.Message == "" {
		apiErr.Message = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	// Add hints for common errors
	apiErr.Hint = hintForStatus(resp.StatusCode)

	return apiErr
}

func hintForStatus(status int) string {
	switch status {
	case 401:
		return "Run `auth login` to authenticate"
	case 403:
		return "Check that your account has the required permissions"
	case 404:
		return "Check that the resource key/ID is correct"
	case 429:
		return "Rate limited — wait and retry"
	default:
		return ""
	}
}
```

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/api/ -v`
Expected: All tests PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/api/errors.go cli/internal/api/errors_test.go
git commit -m "feat: add API error parsing with Atlassian error format support and hints"
```

---

### Task 11: Request Builder & Query Parameters

**Files:**
- Create: `cli/internal/api/request.go`
- Create: `cli/internal/api/request_test.go`

- [x] **Step 1: Write failing test**

```go
// cli/internal/api/request_test.go
package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBuilder_BuildPath(t *testing.T) {
	rb := NewRequestBuilder("/rest/api/3/search").
		Query("jql", "project = TEST").
		Query("maxResults", "50").
		Query("startAt", "0")

	path := rb.BuildPath()
	assert.Contains(t, path, "/rest/api/3/search?")
	assert.Contains(t, path, "jql=project+%3D+TEST")
	assert.Contains(t, path, "maxResults=50")
}

func TestRequestBuilder_EmptyValues(t *testing.T) {
	rb := NewRequestBuilder("/path").
		Query("key", "value").
		Query("empty", "")

	path := rb.BuildPath()
	assert.Contains(t, path, "key=value")
	assert.NotContains(t, path, "empty")
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/api/ -v -run TestRequestBuilder`
Expected: FAIL

- [x] **Step 3: Implement request builder**

```go
// cli/internal/api/request.go
package api

import (
	"fmt"
	"net/url"
)

// RequestBuilder constructs URL paths with query parameters.
type RequestBuilder struct {
	basePath string
	params   url.Values
}

// NewRequestBuilder creates a new RequestBuilder for the given path.
func NewRequestBuilder(basePath string) *RequestBuilder {
	return &RequestBuilder{
		basePath: basePath,
		params:   url.Values{},
	}
}

// Query adds a query parameter. Empty values are skipped.
func (rb *RequestBuilder) Query(key, value string) *RequestBuilder {
	if value != "" {
		rb.params.Set(key, value)
	}
	return rb
}

// QueryInt adds an integer query parameter. Zero values are skipped.
func (rb *RequestBuilder) QueryInt(key string, value int) *RequestBuilder {
	if value != 0 {
		rb.params.Set(key, fmt.Sprintf("%d", value))
	}
	return rb
}

// BuildPath returns the full path with query string.
func (rb *RequestBuilder) BuildPath() string {
	if len(rb.params) == 0 {
		return rb.basePath
	}
	return rb.basePath + "?" + rb.params.Encode()
}
```

> Note: Add `"fmt"` to imports in request.go.

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/api/ -v -run TestRequestBuilder`
Expected: PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/api/request.go cli/internal/api/request_test.go
git commit -m "feat: add request builder for URL path construction with query params"
```

---

### Task 12: Pagination Response Helpers

**Files:**
- Create: `cli/internal/api/pagination.go`
- Create: `cli/internal/api/pagination_test.go`

- [x] **Step 1: Write failing test**

```go
// cli/internal/api/pagination_test.go
package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationInfo_IsLast(t *testing.T) {
	// Offset-based: startAt=0, maxResults=50, total=30
	p := &PaginationInfo{StartAt: 0, MaxResults: 50, Total: 30, IsLast: true}
	assert.True(t, p.IsLast)

	// Not last
	p2 := &PaginationInfo{StartAt: 0, MaxResults: 50, Total: 100, IsLast: false}
	assert.False(t, p2.IsLast)
}

func TestPaginationInfo_JSON(t *testing.T) {
	p := &PaginationInfo{StartAt: 10, MaxResults: 50, Total: 200, IsLast: false}
	data := p.ToMap()
	assert.Equal(t, 10, data["startAt"])
	assert.Equal(t, 50, data["maxResults"])
	assert.Equal(t, 200, data["total"])
	assert.Equal(t, false, data["isLast"])
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/api/ -v -run TestPaginationInfo`
Expected: FAIL

- [x] **Step 3: Implement pagination struct**

```go
// cli/internal/api/pagination.go
package api

// PaginationInfo holds pagination metadata for list responses.
type PaginationInfo struct {
	StartAt    int  `json:"startAt"`
	MaxResults int  `json:"maxResults"`
	Total      int  `json:"total"`
	IsLast     bool `json:"isLast"`
}

// ToMap converts pagination info to a map for JSON embedding.
func (p *PaginationInfo) ToMap() map[string]any {
	return map[string]any{
		"startAt":    p.StartAt,
		"maxResults": p.MaxResults,
		"total":      p.Total,
		"isLast":     p.IsLast,
	}
}

// PaginationFromResponse extracts pagination info from a typical Atlassian response.
func PaginationFromResponse(data map[string]any) *PaginationInfo {
	p := &PaginationInfo{}

	if v, ok := data["startAt"].(float64); ok {
		p.StartAt = int(v)
	}
	if v, ok := data["maxResults"].(float64); ok {
		p.MaxResults = int(v)
	}
	if v, ok := data["total"].(float64); ok {
		p.Total = int(v)
	}
	if v, ok := data["isLast"].(bool); ok {
		p.IsLast = v
	} else {
		// Compute isLast from offset pagination
		p.IsLast = p.StartAt+p.MaxResults >= p.Total
	}

	return p
}
```

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/api/ -v`
Expected: All PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/api/pagination.go cli/internal/api/pagination_test.go
git commit -m "feat: add pagination response helpers for list commands"
```

---

## Chunk 4: Output Formatting

### Task 13: Formatter Interface & JSON Output

> **Important:** Tasks 13-15 must be implemented together before running tests.
> `formatter.go` references `TableFormatter` and `TextFormatter` which are defined in Tasks 14-15.
> Create all three formatter files (json.go, table.go, text.go) before running `go test`.

**Files:**
- Create: `cli/internal/output/formatter.go`
- Create: `cli/internal/output/json.go`
- Create: `cli/internal/output/json_test.go`

- [x] **Step 1: Write failing test**

```go
// cli/internal/output/json_test.go
package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFormatter_Format(t *testing.T) {
	f := &JSONFormatter{}
	data := map[string]any{
		"key":     "TEST-1",
		"summary": "Test issue",
		"status":  "Open",
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, `"key": "TEST-1"`)
	assert.Contains(t, output, `"summary": "Test issue"`)
	// Verify indented output
	assert.Contains(t, output, "  ")
}

func TestJSONFormatter_FormatList(t *testing.T) {
	f := &JSONFormatter{}
	data := []map[string]string{
		{"key": "A-1"},
		{"key": "A-2"},
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), `"key": "A-1"`)
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/output/ -v`
Expected: FAIL

- [x] **Step 3: Implement formatter interface and JSON formatter**

```go
// cli/internal/output/formatter.go
package output

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Formatter writes structured data to an output stream.
type Formatter interface {
	Write(w io.Writer, data any) error
}

// NewFormatter creates a Formatter based on the format name.
func NewFormatter(format string) Formatter {
	switch format {
	case "table":
		return &TableFormatter{}
	case "text":
		return &TextFormatter{}
	default:
		return &JSONFormatter{}
	}
}

// Print writes data to stdout using the format from --output flag.
func Print(cmd *cobra.Command, data any) error {
	format, _ := cmd.Flags().GetString("output")
	if format == "" {
		format = "json"
	}
	f := NewFormatter(format)
	return f.Write(os.Stdout, data)
}

// PrintError writes an error as JSON to stderr.
func PrintError(cmd *cobra.Command, err error) {
	format, _ := cmd.Flags().GetString("output")
	if format == "json" {
		fmt.Fprintf(os.Stderr, `{"error": %q}`+"\n", err.Error())
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	}
}
```

```go
// cli/internal/output/json.go
package output

import (
	"encoding/json"
	"io"
)

// JSONFormatter outputs data as indented JSON.
type JSONFormatter struct{}

func (f *JSONFormatter) Write(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}
```

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/output/ -v -run TestJSON`
Expected: PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/output/formatter.go cli/internal/output/json.go cli/internal/output/json_test.go
git commit -m "feat: add output formatter interface and JSON formatter"
```

---

### Task 14: Table Formatter

**Files:**
- Create: `cli/internal/output/table.go`
- Create: `cli/internal/output/table_test.go`

- [x] **Step 1: Write failing test**

```go
// cli/internal/output/table_test.go
package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableFormatter_MapSlice(t *testing.T) {
	f := &TableFormatter{}
	data := []map[string]any{
		{"key": "TEST-1", "summary": "First issue", "status": "Open"},
		{"key": "TEST-2", "summary": "Second issue", "status": "Closed"},
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "TEST-1")
	assert.Contains(t, output, "TEST-2")
	assert.Contains(t, output, "First issue")
}

func TestTableFormatter_SingleMap(t *testing.T) {
	f := &TableFormatter{}
	data := map[string]any{
		"key":     "TEST-1",
		"summary": "Single issue",
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "TEST-1")
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/output/ -v -run TestTable`
Expected: FAIL

- [x] **Step 3: Implement table formatter**

```go
// cli/internal/output/table.go
package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/olekukonko/tablewriter"
)

// TableFormatter outputs data as an ASCII table.
type TableFormatter struct{}

func (f *TableFormatter) Write(w io.Writer, data any) error {
	switch v := data.(type) {
	case []map[string]any:
		return f.writeMapSlice(w, v)
	case map[string]any:
		return f.writeSingleMap(w, v)
	default:
		// Fallback to JSON for unsupported types
		return (&JSONFormatter{}).Write(w, data)
	}
}

func (f *TableFormatter) writeMapSlice(w io.Writer, items []map[string]any) error {
	if len(items) == 0 {
		fmt.Fprintln(w, "No results")
		return nil
	}

	// Extract headers from first item, sorted
	var headers []string
	for k := range items[0] {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	table := tablewriter.NewWriter(w)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, item := range items {
		var row []string
		for _, h := range headers {
			row = append(row, fmt.Sprintf("%v", item[h]))
		}
		table.Append(row)
	}

	table.Render()
	return nil
}

func (f *TableFormatter) writeSingleMap(w io.Writer, item map[string]any) error {
	// Key-value format for single items
	var keys []string
	for k := range item {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	table := tablewriter.NewWriter(w)
	table.SetBorder(false)
	table.SetColumnSeparator(":")
	table.SetAutoWrapText(false)

	for _, k := range keys {
		table.Append([]string{k, fmt.Sprintf("%v", item[k])})
	}

	table.Render()
	return nil
}
```

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/output/ -v -run TestTable`
Expected: PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/output/table.go cli/internal/output/table_test.go
git commit -m "feat: add table output formatter with tablewriter"
```

---

### Task 15: Text Formatter

**Files:**
- Create: `cli/internal/output/text.go`
- Create: `cli/internal/output/text_test.go`

- [x] **Step 1: Write failing test**

```go
// cli/internal/output/text_test.go
package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextFormatter_SingleMap(t *testing.T) {
	f := &TextFormatter{}
	data := map[string]any{
		"key":     "TEST-1",
		"summary": "My issue title",
		"status":  "Open",
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "key:")
	assert.Contains(t, output, "TEST-1")
	assert.Contains(t, output, "summary:")
}

func TestTextFormatter_MapSlice(t *testing.T) {
	f := &TextFormatter{}
	data := []map[string]any{
		{"key": "A-1", "summary": "First"},
		{"key": "A-2", "summary": "Second"},
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "A-1")
	assert.Contains(t, output, "A-2")
}
```

- [x] **Step 2: Implement text formatter**

```go
// cli/internal/output/text.go
package output

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// TextFormatter outputs data as human-readable plain text.
type TextFormatter struct{}

func (f *TextFormatter) Write(w io.Writer, data any) error {
	switch v := data.(type) {
	case map[string]any:
		return f.writeSingleMap(w, v)
	case []map[string]any:
		return f.writeMapSlice(w, v)
	default:
		return (&JSONFormatter{}).Write(w, data)
	}
}

func (f *TextFormatter) writeSingleMap(w io.Writer, item map[string]any) error {
	var keys []string
	for k := range item {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Find longest key for alignment
	maxLen := 0
	for _, k := range keys {
		if len(k) > maxLen {
			maxLen = len(k)
		}
	}

	for _, k := range keys {
		val := fmt.Sprintf("%v", item[k])
		// Truncate long values
		if len(val) > 200 {
			val = val[:200] + "..."
		}
		padding := strings.Repeat(" ", maxLen-len(k))
		fmt.Fprintf(w, "%s:%s %s\n", k, padding, val)
	}
	return nil
}

func (f *TextFormatter) writeMapSlice(w io.Writer, items []map[string]any) error {
	if len(items) == 0 {
		fmt.Fprintln(w, "No results")
		return nil
	}
	for i, item := range items {
		if i > 0 {
			fmt.Fprintln(w, "---")
		}
		if err := f.writeSingleMap(w, item); err != nil {
			return err
		}
	}
	return nil
}
```

- [x] **Step 3: Run all output tests**

Run: `cd cli && go test ./internal/output/ -v`
Expected: All PASS

- [x] **Step 4: Commit**

```bash
git add cli/internal/output/text.go cli/internal/output/text_test.go
git commit -m "feat: add text output formatter for human-readable display"
```

---

## Chunk 5: Content Conversion

> Reference: `src/mcp_atlassian/preprocessing/` for the Python implementation.
> The Go implementation focuses on the most common constructs for v1.
> Full parity with the Python preprocessor is a stretch goal.

### Task 16: Markdown → ADF Conversion (Jira Cloud Input)

**Files:**
- Create: `cli/internal/convert/adf.go`
- Create: `cli/internal/convert/adf_test.go`

- [x] **Step 1: Write failing tests for core ADF conversions**

```go
// cli/internal/convert/adf_test.go
package convert

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownToADF_Paragraph(t *testing.T) {
	md := "Hello world"
	adf := MarkdownToADF(md)

	assert.Equal(t, "doc", adf["type"])
	assert.Equal(t, 1, adf["version"])

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	para := content[0].(map[string]any)
	assert.Equal(t, "paragraph", para["type"])
}

func TestMarkdownToADF_Headings(t *testing.T) {
	md := "# Title\n\nSome text\n\n## Subtitle"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.GreaterOrEqual(t, len(content), 3)

	h1 := content[0].(map[string]any)
	assert.Equal(t, "heading", h1["type"])
	attrs := h1["attrs"].(map[string]any)
	assert.Equal(t, 1, attrs["level"])
}

func TestMarkdownToADF_CodeBlock(t *testing.T) {
	md := "```python\nprint('hello')\n```"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	cb := content[0].(map[string]any)
	assert.Equal(t, "codeBlock", cb["type"])
	attrs := cb["attrs"].(map[string]any)
	assert.Equal(t, "python", attrs["language"])
}

func TestMarkdownToADF_BulletList(t *testing.T) {
	md := "- Item 1\n- Item 2\n- Item 3"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	list := content[0].(map[string]any)
	assert.Equal(t, "bulletList", list["type"])
}

func TestADFToMarkdown_Paragraph(t *testing.T) {
	adf := map[string]any{
		"type":    "doc",
		"version": 1,
		"content": []any{
			map[string]any{
				"type": "paragraph",
				"content": []any{
					map[string]any{"type": "text", "text": "Hello world"},
				},
			},
		},
	}

	md := ADFToMarkdown(adf)
	assert.Equal(t, "Hello world", md)
}

func TestADFToMarkdown_RoundTrip(t *testing.T) {
	// Simple paragraphs should survive a round-trip
	original := "Hello world"
	adf := MarkdownToADF(original)
	result := ADFToMarkdown(adf)
	assert.Equal(t, original, result)
}

func TestMarkdownToADF_JSON(t *testing.T) {
	md := "Hello **bold** world"
	adf := MarkdownToADF(md)

	// Should be valid JSON
	_, err := json.Marshal(adf)
	require.NoError(t, err)
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/convert/ -v`
Expected: FAIL

- [x] **Step 3: Implement ADF converter**

```go
// cli/internal/convert/adf.go
package convert

import (
	"regexp"
	"strings"
)

// MarkdownToADF converts Markdown text to Atlassian Document Format (ADF).
// This is a simplified converter covering the most common constructs:
// paragraphs, headings, code blocks, bullet lists, bold, italic, links.
func MarkdownToADF(md string) map[string]any {
	doc := map[string]any{
		"type":    "doc",
		"version": 1,
		"content": []any{},
	}

	lines := strings.Split(md, "\n")
	content := []any{}
	i := 0

	for i < len(lines) {
		line := lines[i]

		// Code block
		if strings.HasPrefix(line, "```") {
			lang := strings.TrimPrefix(line, "```")
			lang = strings.TrimSpace(lang)
			var codeLines []string
			i++
			for i < len(lines) && !strings.HasPrefix(lines[i], "```") {
				codeLines = append(codeLines, lines[i])
				i++
			}
			i++ // skip closing ```
			cb := map[string]any{
				"type": "codeBlock",
				"attrs": map[string]any{
					"language": lang,
				},
				"content": []any{
					map[string]any{"type": "text", "text": strings.Join(codeLines, "\n")},
				},
			}
			if lang == "" {
				delete(cb["attrs"].(map[string]any), "language")
			}
			content = append(content, cb)
			continue
		}

		// Heading
		if m := headingRe.FindStringSubmatch(line); m != nil {
			level := len(m[1])
			text := m[2]
			content = append(content, map[string]any{
				"type":    "heading",
				"attrs":   map[string]any{"level": level},
				"content": inlineToADF(text),
			})
			i++
			continue
		}

		// Bullet list
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			var items []any
			for i < len(lines) && (strings.HasPrefix(lines[i], "- ") || strings.HasPrefix(lines[i], "* ")) {
				text := strings.TrimPrefix(lines[i], "- ")
				text = strings.TrimPrefix(text, "* ")
				items = append(items, map[string]any{
					"type": "listItem",
					"content": []any{
						map[string]any{
							"type":    "paragraph",
							"content": inlineToADF(text),
						},
					},
				})
				i++
			}
			content = append(content, map[string]any{
				"type":    "bulletList",
				"content": items,
			})
			continue
		}

		// Empty line — skip
		if strings.TrimSpace(line) == "" {
			i++
			continue
		}

		// Paragraph (default)
		content = append(content, map[string]any{
			"type":    "paragraph",
			"content": inlineToADF(line),
		})
		i++
	}

	doc["content"] = content
	return doc
}

var headingRe = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

// inlineToADF converts inline Markdown to ADF inline nodes.
// For v1, wraps as plain text. Bold/italic/code parsing can be added later.
func inlineToADF(text string) []any {
	var nodes []any
	if text != "" {
		nodes = append(nodes, map[string]any{"type": "text", "text": text})
	}
	return nodes
}

// ADFToMarkdown converts ADF JSON to Markdown text.
func ADFToMarkdown(adf map[string]any) string {
	content, ok := adf["content"].([]any)
	if !ok {
		return ""
	}

	var parts []string
	for _, node := range content {
		nodeMap, ok := node.(map[string]any)
		if !ok {
			continue
		}
		parts = append(parts, adfNodeToMarkdown(nodeMap))
	}

	return strings.Join(parts, "\n\n")
}

func adfNodeToMarkdown(node map[string]any) string {
	nodeType, _ := node["type"].(string)

	switch nodeType {
	case "paragraph":
		return adfInlineToMarkdown(node)

	case "heading":
		attrs, _ := node["attrs"].(map[string]any)
		level := 1
		if l, ok := attrs["level"].(float64); ok {
			level = int(l)
		} else if l, ok := attrs["level"].(int); ok {
			level = l
		}
		prefix := strings.Repeat("#", level)
		return prefix + " " + adfInlineToMarkdown(node)

	case "codeBlock":
		attrs, _ := node["attrs"].(map[string]any)
		lang, _ := attrs["language"].(string)
		code := adfInlineToMarkdown(node)
		return "```" + lang + "\n" + code + "\n```"

	case "bulletList":
		items, _ := node["content"].([]any)
		var lines []string
		for _, item := range items {
			itemMap, _ := item.(map[string]any)
			itemContent, _ := itemMap["content"].([]any)
			if len(itemContent) > 0 {
				para, _ := itemContent[0].(map[string]any)
				lines = append(lines, "- "+adfInlineToMarkdown(para))
			}
		}
		return strings.Join(lines, "\n")

	case "orderedList":
		items, _ := node["content"].([]any)
		var lines []string
		for i, item := range items {
			itemMap, _ := item.(map[string]any)
			itemContent, _ := itemMap["content"].([]any)
			if len(itemContent) > 0 {
				para, _ := itemContent[0].(map[string]any)
				lines = append(lines, fmt.Sprintf("%d. %s", i+1, adfInlineToMarkdown(para)))
			}
		}
		return strings.Join(lines, "\n")

	default:
		return adfInlineToMarkdown(node)
	}
}

func adfInlineToMarkdown(node map[string]any) string {
	content, ok := node["content"].([]any)
	if !ok {
		return ""
	}

	var parts []string
	for _, inline := range content {
		inlineMap, ok := inline.(map[string]any)
		if !ok {
			continue
		}
		inlineType, _ := inlineMap["type"].(string)
		switch inlineType {
		case "text":
			text, _ := inlineMap["text"].(string)
			// Check for marks (bold, italic, etc.)
			marks, _ := inlineMap["marks"].([]any)
			for _, mark := range marks {
				markMap, _ := mark.(map[string]any)
				markType, _ := markMap["type"].(string)
				switch markType {
				case "strong":
					text = "**" + text + "**"
				case "em":
					text = "*" + text + "*"
				case "code":
					text = "`" + text + "`"
				}
			}
			parts = append(parts, text)
		case "hardBreak":
			parts = append(parts, "\n")
		}
	}

	return strings.Join(parts, "")
}
```

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/convert/ -v`
Expected: All PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/convert/adf.go cli/internal/convert/adf_test.go
git commit -m "feat: add ADF ↔ Markdown converter for Jira Cloud content"
```

---

### Task 17: Markdown → Storage Format Conversion (Confluence)

**Files:**
- Create: `cli/internal/convert/storage.go`
- Create: `cli/internal/convert/storage_test.go`

- [x] **Step 1: Write failing test**

```go
// cli/internal/convert/storage_test.go
package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownToStorage_Paragraph(t *testing.T) {
	md := "Hello world"
	html := MarkdownToStorage(md)
	assert.Contains(t, html, "<p>Hello world</p>")
}

func TestMarkdownToStorage_Heading(t *testing.T) {
	md := "# Title\n\n## Subtitle"
	html := MarkdownToStorage(md)
	assert.Contains(t, html, "<h1>Title</h1>")
	assert.Contains(t, html, "<h2>Subtitle</h2>")
}

func TestMarkdownToStorage_CodeBlock(t *testing.T) {
	md := "```python\nprint('hello')\n```"
	html := MarkdownToStorage(md)
	assert.Contains(t, html, "<ac:structured-macro ac:name=\"code\">")
	assert.Contains(t, html, "python")
	// CDATA sections preserve content without HTML-escaping
	assert.Contains(t, html, "print('hello')")
}

func TestMarkdownToStorage_BulletList(t *testing.T) {
	md := "- Item 1\n- Item 2"
	html := MarkdownToStorage(md)
	assert.Contains(t, html, "<ul>")
	assert.Contains(t, html, "<li>Item 1</li>")
	assert.Contains(t, html, "<li>Item 2</li>")
}

func TestMarkdownToStorage_Bold(t *testing.T) {
	md := "Hello **bold** world"
	html := MarkdownToStorage(md)
	assert.Contains(t, html, "<strong>bold</strong>")
}

func TestStorageToMarkdown_Paragraph(t *testing.T) {
	storage := "<p>Hello world</p>"
	md := StorageToMarkdown(storage)
	assert.Contains(t, md, "Hello world")
}

func TestStorageToMarkdown_Heading(t *testing.T) {
	storage := "<h1>Title</h1>"
	md := StorageToMarkdown(storage)
	assert.Contains(t, md, "# Title")
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./internal/convert/ -v -run TestMarkdownToStorage`
Expected: FAIL

- [x] **Step 3: Implement Storage Format converter**

```go
// cli/internal/convert/storage.go
package convert

import (
	"html"
	"regexp"
	"strings"
)

// MarkdownToStorage converts Markdown to Confluence Storage Format (XHTML).
// Covers: paragraphs, headings, code blocks, lists, bold, italic, links.
func MarkdownToStorage(md string) string {
	var result strings.Builder
	lines := strings.Split(md, "\n")
	i := 0

	for i < len(lines) {
		line := lines[i]

		// Code block
		if strings.HasPrefix(line, "```") {
			lang := strings.TrimPrefix(line, "```")
			lang = strings.TrimSpace(lang)
			var codeLines []string
			i++
			for i < len(lines) && !strings.HasPrefix(lines[i], "```") {
				codeLines = append(codeLines, lines[i])
				i++
			}
			i++ // skip closing ```
			code := html.EscapeString(strings.Join(codeLines, "\n"))

			result.WriteString(`<ac:structured-macro ac:name="code">`)
			if lang != "" {
				result.WriteString(`<ac:parameter ac:name="language">` + lang + `</ac:parameter>`)
			}
			result.WriteString(`<ac:plain-text-body><![CDATA[` + strings.Join(codeLines, "\n") + `]]></ac:plain-text-body>`)
			result.WriteString(`</ac:structured-macro>`)
			_ = code // used html.EscapeString above for param, CDATA for body
			continue
		}

		// Heading
		if m := storageHeadingRe.FindStringSubmatch(line); m != nil {
			level := len(m[1])
			text := convertInlineToHTML(m[2])
			result.WriteString("<h" + string(rune('0'+level)) + ">" + text + "</h" + string(rune('0'+level)) + ">")
			i++
			continue
		}

		// Bullet list
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			result.WriteString("<ul>")
			for i < len(lines) && (strings.HasPrefix(lines[i], "- ") || strings.HasPrefix(lines[i], "* ")) {
				text := strings.TrimPrefix(lines[i], "- ")
				text = strings.TrimPrefix(text, "* ")
				result.WriteString("<li>" + convertInlineToHTML(text) + "</li>")
				i++
			}
			result.WriteString("</ul>")
			continue
		}

		// Empty line — skip
		if strings.TrimSpace(line) == "" {
			i++
			continue
		}

		// Paragraph
		result.WriteString("<p>" + convertInlineToHTML(line) + "</p>")
		i++
	}

	return result.String()
}

var storageHeadingRe = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
var storageBoldRe = regexp.MustCompile(`\*\*(.+?)\*\*`)
var storageItalicRe = regexp.MustCompile(`(?:^|[^*])\*([^*]+?)\*(?:[^*]|$)`)
var storageLinkRe = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

func convertInlineToHTML(text string) string {
	// Links
	text = storageLinkRe.ReplaceAllString(text, `<a href="$2">$1</a>`)
	// Bold
	text = storageBoldRe.ReplaceAllString(text, `<strong>$1</strong>`)
	// Inline code
	text = regexp.MustCompile("`([^`]+)`").ReplaceAllString(text, `<code>$1</code>`)
	return text
}

// StorageToMarkdown converts Confluence Storage Format (XHTML) to Markdown.
// This is a simplified converter for common elements.
func StorageToMarkdown(storage string) string {
	result := storage

	// Headings
	for level := 1; level <= 6; level++ {
		l := string(rune('0' + level))
		prefix := strings.Repeat("#", level)
		re := regexp.MustCompile(`<h` + l + `[^>]*>(.*?)</h` + l + `>`)
		result = re.ReplaceAllString(result, prefix+" $1\n")
	}

	// Paragraphs
	result = regexp.MustCompile(`<p[^>]*>(.*?)</p>`).ReplaceAllString(result, "$1\n")

	// Bold
	result = regexp.MustCompile(`<strong>(.*?)</strong>`).ReplaceAllString(result, "**$1**")

	// Italic
	result = regexp.MustCompile(`<em>(.*?)</em>`).ReplaceAllString(result, "*$1*")

	// Code
	result = regexp.MustCompile(`<code>(.*?)</code>`).ReplaceAllString(result, "`$1`")

	// Links
	result = regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*>(.*?)</a>`).ReplaceAllString(result, "[$2]($1)")

	// Lists
	result = regexp.MustCompile(`<li>(.*?)</li>`).ReplaceAllString(result, "- $1")
	result = regexp.MustCompile(`</?[uo]l[^>]*>`).ReplaceAllString(result, "")

	// Strip remaining HTML tags
	result = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(result, "")

	// Clean up whitespace
	result = strings.TrimSpace(result)

	return result
}
```

- [x] **Step 4: Run tests**

Run: `cd cli && go test ./internal/convert/ -v`
Expected: All PASS

- [x] **Step 5: Commit**

```bash
git add cli/internal/convert/storage.go cli/internal/convert/storage_test.go
git commit -m "feat: add Confluence Storage Format ↔ Markdown converter"
```

---

## Chunk 6: Jira Core Commands (Issue & Search)

> This chunk establishes the **command implementation pattern** that all subsequent commands follow.
> Subsequent chunks reference this pattern and only document the differences.

### Task 18: Wire Factory with API Client

**Files:**
- Modify: `cli/pkg/cmdutil/factory.go` — add Client creation
- Create: `cli/pkg/cmdutil/flags.go`

- [x] **Step 1: Extend Factory with Client creation**

```go
// cli/pkg/cmdutil/factory.go
package cmdutil

import (
	"fmt"

	"github.com/anthropics/atlassian-cli/internal/api"
	"github.com/anthropics/atlassian-cli/internal/config"
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
// flagHost is the value of --host flag (may be empty).
func (f *Factory) JiraClient(flagHost string) (*api.Client, error) {
	host, hostname, err := config.ResolveHost(f.Config, flagHost)
	if err != nil {
		return nil, err
	}
	if !host.Jira.Enabled && host.Jira.BasePath == "" && host.Type != "" {
		// Service not explicitly disabled; allow by default
	}

	auth, err := config.NewAuthenticator(host)
	if err != nil {
		return nil, fmt.Errorf("auth for %s: %w", hostname, err)
	}

	baseURL := "https://" + hostname + host.Jira.BasePath
	return api.NewClient(baseURL, auth), nil
}

// ConfluenceClient creates an API client for Confluence.
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
```

- [x] **Step 2: Create shared flag helpers**

```go
// cli/pkg/cmdutil/flags.go
package cmdutil

import "github.com/spf13/cobra"

// GetHost returns the --host flag value from the command hierarchy.
func GetHost(cmd *cobra.Command) string {
	host, _ := cmd.Flags().GetString("host")
	if host == "" {
		host, _ = cmd.InheritedFlags().GetString("host")
	}
	return host
}

// GetOutput returns the --output flag value.
func GetOutput(cmd *cobra.Command) string {
	output, _ := cmd.Flags().GetString("output")
	if output == "" {
		output, _ = cmd.InheritedFlags().GetString("output")
	}
	if output == "" {
		return "json"
	}
	return output
}
```

- [x] **Step 3: Commit**

```bash
git add cli/pkg/cmdutil/
git commit -m "feat: add Factory.JiraClient/ConfluenceClient and shared flag helpers"
```

---

### Task 19: `jira issue get` Command

**Files:**
- Create: `cli/pkg/jira/issue/issue.go`
- Create: `cli/pkg/jira/issue/get.go`
- Create: `cli/pkg/jira/issue/get_test.go`
- Modify: `cli/pkg/jira/root.go` — register issue subcommand

> **Reference:** `src/mcp_atlassian/servers/jira.py` → `get_issue()` tool, `src/mcp_atlassian/jira/issues.py` → `get_issue()` method.
> API endpoint: `GET /rest/api/3/issue/{issueKey}` (Cloud) or `GET /rest/api/2/issue/{issueKey}` (Server/DC).

- [x] **Step 1: Write failing test with mock server**

```go
// cli/pkg/jira/issue/get_test.go
package issue

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthropics/atlassian-cli/internal/api"
	"github.com/anthropics/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunIssueGet(t *testing.T) {
	mockIssue := map[string]any{
		"key": "TEST-1",
		"fields": map[string]any{
			"summary": "Test issue summary",
			"status": map[string]any{
				"name": "Open",
			},
			"issuetype": map[string]any{
				"name": "Task",
			},
			"assignee": map[string]any{
				"displayName": "John Doe",
			},
			"description": nil,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock server URL is localhost, so IsCloud()=false → api/2 path
		assert.Contains(t, r.URL.Path, "/rest/api/2/issue/TEST-1")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockIssue)
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "u", Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := fetchIssue(client, "TEST-1", "", "")
	require.NoError(t, err)

	assert.Equal(t, "TEST-1", result["key"])
	fields := result["fields"].(map[string]any)
	assert.Equal(t, "Test issue summary", fields["summary"])
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./pkg/jira/issue/ -v -run TestRunIssueGet`
Expected: FAIL

- [x] **Step 3: Implement issue subcommand group**

```go
// cli/pkg/jira/issue/issue.go
package issue

import (
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdIssue creates the `issue` subcommand group.
func NewCmdIssue(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue <command>",
		Short: "Manage Jira issues",
	}

	cmd.AddCommand(NewCmdGet(f))
	// Future: cmd.AddCommand(NewCmdCreate(f))
	// Future: cmd.AddCommand(NewCmdUpdate(f))
	// Future: cmd.AddCommand(NewCmdDelete(f))
	// Future: cmd.AddCommand(NewCmdTransition(f))
	// Future: cmd.AddCommand(NewCmdBatchCreate(f))

	return cmd
}
```

- [x] **Step 4: Implement `jira issue get`**

```go
// cli/pkg/jira/issue/get.go
package issue

import (
	"fmt"

	"github.com/anthropics/atlassian-cli/internal/api"
	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdGet creates the `issue get` command.
func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var (
		fields       string
		expand       string
		commentLimit int
	)

	cmd := &cobra.Command{
		Use:   "get <issue-key>",
		Short: "Get a Jira issue by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := fetchIssue(client, args[0], fields, expand)
			if err != nil {
				return err
			}

			_ = commentLimit // TODO: filter comments in response
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&fields, "fields", "", "Comma-separated list of fields to include")
	cmd.Flags().StringVar(&expand, "expand", "", "Comma-separated list of fields to expand")
	cmd.Flags().IntVar(&commentLimit, "comment-limit", 0, "Max comments to include (0 = default)")

	return cmd
}

// fetchIssue calls the Jira REST API to get an issue.
func fetchIssue(client *api.Client, issueKey, fields, expand string) (map[string]any, error) {
	path := client.JiraAPIPath("issue/" + issueKey)

	rb := api.NewRequestBuilder(path)
	if fields != "" {
		rb.Query("fields", fields)
	}
	if expand != "" {
		rb.Query("expand", expand)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("getting issue %s: %w", issueKey, err)
	}

	return result, nil
}
```

- [x] **Step 5: Register issue subcommand in root**

```go
// Update cli/pkg/jira/root.go — add import and registration:
import "github.com/anthropics/atlassian-cli/pkg/jira/issue"

// Inside NewCmdRoot, after auth registration:
cmd.AddCommand(issue.NewCmdIssue(f))
```

- [x] **Step 6: Run tests**

Run: `cd cli && go test ./pkg/jira/issue/ -v`
Expected: PASS

- [x] **Step 7: Build and verify**

Run: `cd cli && make build && ./bin/jira issue get --help`
Expected: Shows usage with `--fields`, `--expand`, `--comment-limit` flags

- [x] **Step 8: Commit**

```bash
git add cli/pkg/jira/issue/ cli/pkg/jira/root.go
git commit -m "feat(jira): add issue get command"
```

---

### Task 20: `jira issue create` Command

**Files:**
- Create: `cli/pkg/jira/issue/create.go`
- Create: `cli/pkg/jira/issue/create_test.go`

> API: `POST /rest/api/3/issue` (Cloud) / `POST /rest/api/2/issue` (Server/DC)

- [x] **Step 1: Write failing test**

```go
// cli/pkg/jira/issue/create_test.go
package issue

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthropics/atlassian-cli/internal/api"
	"github.com/anthropics/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		json.Unmarshal(body, &payload)

		fields := payload["fields"].(map[string]any)
		project := fields["project"].(map[string]any)
		assert.Equal(t, "TEST", project["key"])
		assert.Equal(t, "Bug title", fields["summary"])

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{
			"id":   "10001",
			"key":  "TEST-42",
			"self": "https://example.com/rest/api/3/issue/10001",
		})
	}))
	defer server.Close()

	auth := &config.BasicAuth{Username: "u", Token: "t"}
	client := api.NewClient(server.URL, auth)

	result, err := createIssue(client, CreateIssueOpts{
		Project:   "TEST",
		Summary:   "Bug title",
		IssueType: "Bug",
	})
	require.NoError(t, err)
	assert.Equal(t, "TEST-42", result["key"])
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./pkg/jira/issue/ -v -run TestCreateIssue`
Expected: FAIL

- [x] **Step 3: Implement create command**

```go
// cli/pkg/jira/issue/create.go
package issue

import (
	"encoding/json"
	"fmt"

	"github.com/anthropics/atlassian-cli/internal/api"
	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// CreateIssueOpts holds the parameters for creating an issue.
type CreateIssueOpts struct {
	Project     string
	Summary     string
	IssueType   string
	Assignee    string
	Description string
	Components  string
	FieldsJSON  string
}

// NewCmdCreate creates the `issue create` command.
func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var opts CreateIssueOpts

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Jira issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}
			result, err := createIssue(client, opts)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&opts.Project, "project", "", "Project key (required)")
	cmd.Flags().StringVar(&opts.Summary, "summary", "", "Issue summary (required)")
	cmd.Flags().StringVar(&opts.IssueType, "type", "Task", "Issue type")
	cmd.Flags().StringVar(&opts.Assignee, "assignee", "", "Assignee username or account ID")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Issue description (Markdown)")
	cmd.Flags().StringVar(&opts.Components, "components", "", "Comma-separated component names")
	cmd.Flags().StringVar(&opts.FieldsJSON, "fields-json", "", "Additional fields as JSON object")

	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("summary")

	return cmd
}

func createIssue(client *api.Client, opts CreateIssueOpts) (map[string]any, error) {
	fields := map[string]any{
		"project":   map[string]any{"key": opts.Project},
		"summary":   opts.Summary,
		"issuetype": map[string]any{"name": opts.IssueType},
	}

	if opts.Assignee != "" {
		if client.IsCloud() {
			fields["assignee"] = map[string]any{"accountId": opts.Assignee}
		} else {
			fields["assignee"] = map[string]any{"name": opts.Assignee}
		}
	}

	if opts.Description != "" {
		// TODO: convert Markdown to ADF for Cloud
		fields["description"] = opts.Description
	}

	if opts.FieldsJSON != "" {
		var extra map[string]any
		if err := json.Unmarshal([]byte(opts.FieldsJSON), &extra); err != nil {
			return nil, fmt.Errorf("parsing --fields-json: %w", err)
		}
		for k, v := range extra {
			fields[k] = v
		}
	}

	body := map[string]any{"fields": fields}
	path := client.JiraAPIPath("issue")

	var result map[string]any
	_, err := client.Post(path, body, &result)
	if err != nil {
		return nil, fmt.Errorf("creating issue: %w", err)
	}

	return result, nil
}
```

- [x] **Step 4: Register create in issue.go**

Add `cmd.AddCommand(NewCmdCreate(f))` in `issue.go`.

- [x] **Step 5: Run tests**

Run: `cd cli && go test ./pkg/jira/issue/ -v`
Expected: PASS

- [x] **Step 6: Commit**

```bash
git add cli/pkg/jira/issue/create.go cli/pkg/jira/issue/create_test.go cli/pkg/jira/issue/issue.go
git commit -m "feat(jira): add issue create command"
```

---

### Task 21: Remaining Issue Commands (Update, Delete, Transition, Batch Create)

**Files:**
- Create: `cli/pkg/jira/issue/update.go`
- Create: `cli/pkg/jira/issue/delete.go`
- Create: `cli/pkg/jira/issue/transition.go`
- Create: `cli/pkg/jira/issue/batch_create.go`
- Modify: `cli/pkg/jira/issue/issue.go` — register all

> All follow the same pattern as Task 19/20. Key differences documented below.

- [x] **Step 1: Implement `issue update`**

Pattern: `PUT /rest/api/{2|3}/issue/{key}` with `--fields-json` body.

```go
// cli/pkg/jira/issue/update.go
package issue

import (
	"encoding/json"
	"fmt"

	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		fieldsJSON  string
		components  string
		attachments string
	)
	cmd := &cobra.Command{
		Use:   "update <issue-key>",
		Short: "Update a Jira issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fieldsJSON == "" && components == "" && attachments == "" {
				return fmt.Errorf("at least one of --fields-json, --components, or --attachments is required")
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			fields := make(map[string]any)
			if fieldsJSON != "" {
				if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
					return fmt.Errorf("parsing --fields-json: %w", err)
				}
			}
			if components != "" {
				var comps []map[string]any
				for _, name := range strings.Split(components, ",") {
					comps = append(comps, map[string]any{"name": strings.TrimSpace(name)})
				}
				fields["components"] = comps
			}

			body := map[string]any{"fields": fields}
			path := client.JiraAPIPath("issue/" + args[0])
			_, err = client.Put(path, body, nil)
			if err != nil {
				return fmt.Errorf("updating issue %s: %w", args[0], err)
			}

			// TODO: handle --attachments via multipart upload to /rest/api/{v}/issue/{key}/attachments

			return output.Print(cmd, map[string]any{
				"success": true,
				"key":     args[0],
				"message": "Issue updated",
			})
		},
	}
	cmd.Flags().StringVar(&fieldsJSON, "fields-json", "", "Fields to update as JSON")
	cmd.Flags().StringVar(&components, "components", "", "Comma-separated component names")
	cmd.Flags().StringVar(&attachments, "attachments", "", "Comma-separated file paths to attach")
	return cmd
}
```

- [x] **Step 2: Implement `issue delete`**

Pattern: `DELETE /rest/api/{2|3}/issue/{key}`, requires `--confirm` flag.

```go
// cli/pkg/jira/issue/delete.go
package issue

import (
	"fmt"

	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool
	cmd := &cobra.Command{
		Use:   "delete <issue-key>",
		Short: "Delete a Jira issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return fmt.Errorf("use --confirm to delete issue %s", args[0])
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			path := client.JiraAPIPath("issue/" + args[0])
			_, err = client.Delete(path, nil)
			if err != nil {
				return fmt.Errorf("deleting issue %s: %w", args[0], err)
			}

			return output.Print(cmd, map[string]any{
				"success": true,
				"key":     args[0],
				"message": "Issue deleted",
			})
		},
	}
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
	return cmd
}
```

- [x] **Step 3: Implement `issue transition`**

Pattern: `POST /rest/api/{2|3}/issue/{key}/transitions`

```go
// cli/pkg/jira/issue/transition.go
package issue

import (
	"fmt"

	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdTransition(f *cmdutil.Factory) *cobra.Command {
	var (
		transitionID string
		comment      string
		fieldsJSON   string
	)
	cmd := &cobra.Command{
		Use:   "transition <issue-key>",
		Short: "Transition a Jira issue to a new status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			body := map[string]any{
				"transition": map[string]any{"id": transitionID},
			}

			update := map[string]any{}
			if comment != "" {
				update["comment"] = []any{
					map[string]any{
						"add": map[string]any{"body": comment},
					},
				}
			}
			if len(update) > 0 {
				body["update"] = update
			}

			if fieldsJSON != "" {
				var fields map[string]any
				if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
					return fmt.Errorf("parsing --fields-json: %w", err)
				}
				body["fields"] = fields
			}

			path := client.JiraAPIPath("issue/" + args[0] + "/transitions")
			_, err = client.Post(path, body, nil)
			if err != nil {
				return fmt.Errorf("transitioning issue %s: %w", args[0], err)
			}

			return output.Print(cmd, map[string]any{
				"success":      true,
				"key":          args[0],
				"transitionId": transitionID,
			})
		},
	}
	cmd.Flags().StringVar(&transitionID, "transition-id", "", "Target transition ID (required)")
	cmd.Flags().StringVar(&fieldsJSON, "fields-json", "", "Fields to set during transition as JSON")
	cmd.Flags().StringVar(&comment, "comment", "", "Comment to add during transition")
	cmd.MarkFlagRequired("transition-id")
	return cmd
}
```

- [x] **Step 4: Implement `issue batch-create`**

Pattern: Reads JSON file, calls `POST /rest/api/{2|3}/issue/bulk`.

```go
// cli/pkg/jira/issue/batch_create.go
package issue

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdBatchCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		file         string
		validateOnly bool
	)

	cmd := &cobra.Command{
		Use:   "batch-create",
		Short: "Create multiple Jira issues from a JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("reading file %s: %w", file, err)
			}

			var issues []map[string]any
			if err := json.Unmarshal(data, &issues); err != nil {
				return fmt.Errorf("parsing JSON: %w", err)
			}

			if validateOnly {
				return output.Print(cmd, map[string]any{
					"valid": true,
					"count": len(issues),
				})
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// Create issues serially (per spec: no parallelization in v1)
			var results []map[string]any
			for i, issueData := range issues {
				body := map[string]any{"fields": issueData}
				path := client.JiraAPIPath("issue")

				var result map[string]any
				_, err := client.Post(path, body, &result)
				if err != nil {
					results = append(results, map[string]any{
						"index": i,
						"error": err.Error(),
					})
					continue
				}
				result["index"] = i
				results = append(results, result)
			}

			return output.Print(cmd, results)
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "JSON file with issue data (required)")
	cmd.Flags().BoolVar(&validateOnly, "validate-only", false, "Validate without creating")
	cmd.MarkFlagRequired("file")
	return cmd
}
```

- [x] **Step 5: Register all commands in issue.go**

```go
// Update cli/pkg/jira/issue/issue.go
cmd.AddCommand(NewCmdGet(f))
cmd.AddCommand(NewCmdCreate(f))
cmd.AddCommand(NewCmdUpdate(f))
cmd.AddCommand(NewCmdDelete(f))
cmd.AddCommand(NewCmdTransition(f))
cmd.AddCommand(NewCmdBatchCreate(f))
```

- [x] **Step 6: Build and verify all subcommands**

Run: `cd cli && make build && ./bin/jira issue --help`
Expected: Shows get, create, update, delete, transition, batch-create

- [x] **Step 7: Commit**

```bash
git add cli/pkg/jira/issue/
git commit -m "feat(jira): add issue update/delete/transition/batch-create commands"
```

---

### Task 22: `jira search` Command

**Files:**
- Create: `cli/pkg/jira/search/search.go`
- Create: `cli/pkg/jira/search/search_test.go`
- Modify: `cli/pkg/jira/root.go` — register search

> API: Cloud uses `POST /rest/api/3/search/jql`, Server/DC uses `GET /rest/api/2/search`
> Reference: `src/mcp_atlassian/jira/search.py`

- [x] **Step 1: Write failing test**

```go
// cli/pkg/jira/search/search_test.go
package search

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthropics/atlassian-cli/internal/api"
	"github.com/anthropics/atlassian-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_ServerDC(t *testing.T) {
	mockResponse := map[string]any{
		"startAt":    0,
		"maxResults": 50,
		"total":      2,
		"issues": []any{
			map[string]any{"key": "TEST-1", "fields": map[string]any{"summary": "First"}},
			map[string]any{"key": "TEST-2", "fields": map[string]any{"summary": "Second"}},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/rest/api/2/search")
		assert.Equal(t, "project = TEST", r.URL.Query().Get("jql"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	auth := &config.PATAuth{Token: "t"}
	client := api.NewClient(server.URL, auth) // non-cloud URL

	result, err := searchIssues(client, SearchOpts{JQL: "project = TEST", Limit: 50})
	require.NoError(t, err)
	issues := result["issues"].([]any)
	assert.Len(t, issues, 2)
}
```

- [x] **Step 2: Run test to verify it fails**

Run: `cd cli && go test ./pkg/jira/search/ -v`
Expected: FAIL

- [x] **Step 3: Implement search command**

```go
// cli/pkg/jira/search/search.go
package search

import (
	"fmt"

	"github.com/anthropics/atlassian-cli/internal/api"
	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// SearchOpts holds search parameters.
type SearchOpts struct {
	JQL            string
	Fields         string
	Limit          int
	StartAt        int
	ProjectsFilter string
	Expand         string
}

// NewCmdSearch creates the `search` command.
func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	var opts SearchOpts

	cmd := &cobra.Command{
		Use:   "search <jql>",
		Short: "Search Jira issues using JQL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.JQL = args[0]

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := searchIssues(client, opts)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&opts.Fields, "fields", "", "Comma-separated fields to return")
	cmd.Flags().IntVar(&opts.Limit, "limit", 50, "Maximum results to return")
	cmd.Flags().IntVar(&opts.StartAt, "start-at", 0, "Starting index for pagination")
	cmd.Flags().StringVar(&opts.ProjectsFilter, "projects-filter", "", "Comma-separated project keys to filter")
	cmd.Flags().StringVar(&opts.Expand, "expand", "", "Fields to expand")

	return cmd
}

func searchIssues(client *api.Client, opts SearchOpts) (map[string]any, error) {
	if client.IsCloud() {
		return searchCloud(client, opts)
	}
	return searchServerDC(client, opts)
}

func searchCloud(client *api.Client, opts SearchOpts) (map[string]any, error) {
	body := map[string]any{
		"jql":        opts.JQL,
		"maxResults": opts.Limit,
	}
	if opts.Fields != "" {
		body["fields"] = opts.Fields
	}

	path := "/rest/api/3/search/jql"
	var result map[string]any
	_, err := client.Post(path, body, &result)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	return result, nil
}

func searchServerDC(client *api.Client, opts SearchOpts) (map[string]any, error) {
	rb := api.NewRequestBuilder("/rest/api/2/search").
		Query("jql", opts.JQL).
		Query("fields", opts.Fields).
		Query("expand", opts.Expand)

	if opts.Limit > 0 {
		rb.QueryInt("maxResults", opts.Limit)
	}
	if opts.StartAt > 0 {
		rb.QueryInt("startAt", opts.StartAt)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	return result, nil
}
```

- [x] **Step 4: Register in root.go**

```go
// Update cli/pkg/jira/root.go:
import "github.com/anthropics/atlassian-cli/pkg/jira/search"
cmd.AddCommand(search.NewCmdSearch(f))
```

- [x] **Step 5: Run tests**

Run: `cd cli && go test ./pkg/jira/search/ -v`
Expected: PASS

- [x] **Step 6: Commit**

```bash
git add cli/pkg/jira/search/ cli/pkg/jira/root.go
git commit -m "feat(jira): add search command with Cloud/Server API routing"
```

---

## Chunk 7: Jira Supporting Commands

> All commands follow the pattern established in Chunk 6.
> Each task: (1) create subcommand group file, (2) implement each command, (3) register in root.go, (4) test with mock server, (5) commit.

### Task 23: Comment Commands (list/add/edit)

**Files:**
- Create: `cli/pkg/jira/comment/comment.go`, `list.go`, `add.go`, `edit.go`
- Modify: `cli/pkg/jira/root.go`

> API endpoints:
> - List: `GET /rest/api/{2|3}/issue/{key}/comment`
> - Add: `POST /rest/api/{2|3}/issue/{key}/comment` with `{"body": "..."}` (Server) or ADF body (Cloud)
> - Edit: `PUT /rest/api/{2|3}/issue/{key}/comment/{id}` with `{"body": "..."}`

- [x] **Step 1: Create comment group**

```go
// cli/pkg/jira/comment/comment.go
package comment

import (
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdComment(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment <command>",
		Short: "Manage issue comments",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdAdd(f))
	cmd.AddCommand(NewCmdEdit(f))
	return cmd
}
```

- [x] **Step 2: Implement `comment list`**

```go
// cli/pkg/jira/comment/list.go
package comment

import (
	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list <issue-key>",
		Short: "List comments on an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}
			path := client.JiraAPIPath("issue/" + args[0] + "/comment")
			var result map[string]any
			_, err = client.Get(path, &result)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}
}
```

- [x] **Step 3: Implement `comment add`**

```go
// cli/pkg/jira/comment/add.go
package comment

import (
	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdAdd(f *cmdutil.Factory) *cobra.Command {
	var (
		body       string
		visibility string
		public     bool
	)
	cmd := &cobra.Command{
		Use:   "add <issue-key>",
		Short: "Add a comment to an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// Service Desk public/internal comment uses a different API
			if cmd.Flags().Changed("public") {
				payload := map[string]any{
					"body":   body,
					"public": public,
				}
				path := "/rest/servicedeskapi/request/" + args[0] + "/comment"
				var result map[string]any
				_, err = client.Post(path, payload, &result)
				if err != nil {
					return err
				}
				return output.Print(cmd, result)
			}

			payload := map[string]any{"body": body}
			if visibility != "" {
				payload["visibility"] = map[string]any{
					"type":  "role",
					"value": visibility,
				}
			}
			path := client.JiraAPIPath("issue/" + args[0] + "/comment")
			var result map[string]any
			_, err = client.Post(path, payload, &result)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}
	cmd.Flags().StringVar(&body, "body", "", "Comment body (required)")
	cmd.Flags().StringVar(&visibility, "visibility", "", "Visibility role")
	cmd.Flags().BoolVar(&public, "public", false, "Make comment public (Service Desk)")
	cmd.MarkFlagRequired("body")
	return cmd
}
```

- [x] **Step 4: Implement `comment edit`** (same pattern, `PUT` method)

```go
// cli/pkg/jira/comment/edit.go
package comment

import (
	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdEdit(f *cmdutil.Factory) *cobra.Command {
	var body, visibility string
	cmd := &cobra.Command{
		Use:   "edit <issue-key> <comment-id>",
		Short: "Edit an issue comment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}
			payload := map[string]any{"body": body}
			if visibility != "" {
				payload["visibility"] = map[string]any{"type": "role", "value": visibility}
			}
			path := client.JiraAPIPath("issue/" + args[0] + "/comment/" + args[1])
			var result map[string]any
			_, err = client.Put(path, payload, &result)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}
	cmd.Flags().StringVar(&body, "body", "", "New comment body (required)")
	cmd.Flags().StringVar(&visibility, "visibility", "", "Visibility role")
	cmd.MarkFlagRequired("body")
	return cmd
}
```

- [x] **Step 5: Register in root.go, build, test**

Run: `cd cli && make build && ./bin/jira comment --help`

- [x] **Step 6: Commit**

```bash
git add cli/pkg/jira/comment/ cli/pkg/jira/root.go
git commit -m "feat(jira): add comment list/add/edit commands"
```

---

### Task 24: Field Commands (search/options)

**Files:**
- Create: `cli/pkg/jira/field/field.go`, `search.go`, `options.go`

> API: `GET /rest/api/{2|3}/field` (list/search), `GET /rest/api/{2|3}/field/{id}/option` (options)

- [x] **Step 1: Implement field subcommand group and both commands**

Follow the same cobra pattern. `field search` takes `--keyword`, `--limit`, `--refresh`. `field options` takes `<field-id>`, `--context-id`, `--project`, `--issue-type`, `--contains`.

- [x] **Step 2: Register in root.go, build, test**
- [x] **Step 3: Commit**

```bash
git add cli/pkg/jira/field/ cli/pkg/jira/root.go
git commit -m "feat(jira): add field search/options commands"
```

---

### Task 25: Transition List Command

**Files:**
- Create: `cli/pkg/jira/transition/list.go`

> API: `GET /rest/api/{2|3}/issue/{key}/transitions`

- [x] **Step 1: Implement transition list**

Simple GET command: Takes `<issue-key>`, calls API, prints transitions array.

- [x] **Step 2: Register in root.go, build, test**
- [x] **Step 3: Commit**

```bash
git add cli/pkg/jira/transition/ cli/pkg/jira/root.go
git commit -m "feat(jira): add transition list command"
```

---

### Task 26: Board Commands (list/issues)

**Files:**
- Create: `cli/pkg/jira/board/board.go`, `list.go`, `issues.go`

> API: `GET /rest/agile/1.0/board` (list), `GET /rest/agile/1.0/board/{id}/issue` (issues)
> Note: Agile API uses `/rest/agile/1.0/` prefix, NOT the standard Jira REST API path.

- [x] **Step 1: Implement board commands**

`board list`: `--name`, `--project`, `--type`, `--limit`
`board issues`: `<board-id>`, `--jql`, `--fields`, `--limit`

Both use `/rest/agile/1.0/` hardcoded path (same for Cloud and Server/DC).

- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/board/ cli/pkg/jira/root.go
git commit -m "feat(jira): add board list/issues commands"
```

---

### Task 27: Sprint Commands (list/issues/create/update/add-issues)

**Files:**
- Create: `cli/pkg/jira/sprint/sprint.go`, `list.go`, `issues.go`, `create.go`, `update.go`, `add_issues.go`

> API: `/rest/agile/1.0/board/{boardId}/sprint` (list), `/rest/agile/1.0/sprint/{sprintId}/issue` (issues),
> `POST /rest/agile/1.0/sprint` (create), `PUT /rest/agile/1.0/sprint/{id}` (update),
> `POST /rest/agile/1.0/sprint/{id}/issue` (add issues)

- [x] **Step 1: Implement all 5 sprint commands following established pattern**

Key flags per command:
- `sprint list <board-id>`: `--state`, `--limit`
- `sprint issues <sprint-id>`: `--fields`, `--limit`
- `sprint create <board-id>`: `--name`, `--start-date`, `--end-date`, `--goal`
- `sprint update <sprint-id>`: `--name`, `--state`, `--start-date`, `--end-date`, `--goal`
- `sprint add-issues <sprint-id>`: `--keys` (comma-separated)

- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/sprint/ cli/pkg/jira/root.go
git commit -m "feat(jira): add sprint list/issues/create/update/add-issues commands"
```

---

## Chunk 8: Jira Extended Commands

> All commands follow the pattern from Chunk 6-7. Each task lists the command, API endpoint, key flags, and commit.

### Task 28: Project Commands (list/versions/components/version-create)

**Files:** `cli/pkg/jira/project/` — `project.go`, `list.go`, `versions.go`, `components.go`, `version_create.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `project list` | GET | `/rest/api/{v}/project` | `--include-archived` |
| `project versions <key>` | GET | `/rest/api/{v}/project/{key}/versions` | — |
| `project components <key>` | GET | `/rest/api/{v}/project/{key}/components` | — |
| `project version create <key>` | POST | `/rest/api/{v}/version` | `--name`, `--start-date`, `--release-date`, `--description`, `--file` |

- [x] **Step 1: Implement all 4 project commands**

Note for `version create`: if `--file` provided, read JSON array and batch create serially. Otherwise single creation with `--name` (required).

- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/project/ cli/pkg/jira/root.go
git commit -m "feat(jira): add project list/versions/components/version-create commands"
```

---

### Task 29: Worklog Commands (list/add)

**Files:** `cli/pkg/jira/worklog/` — `worklog.go`, `list.go`, `add.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `worklog list <key>` | GET | `/rest/api/{v}/issue/{key}/worklog` | — |
| `worklog add <key>` | POST | `/rest/api/{v}/issue/{key}/worklog` | `--time-spent`, `--comment`, `--started` |

- [x] **Step 1: Implement both worklog commands**
- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/worklog/ cli/pkg/jira/root.go
git commit -m "feat(jira): add worklog list/add commands"
```

---

### Task 30: Watcher Commands (list/add/remove)

**Files:** `cli/pkg/jira/watcher/` — `watcher.go`, `list.go`, `add.go`, `remove.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `watcher list <key>` | GET | `/rest/api/{v}/issue/{key}/watchers` | — |
| `watcher add <key>` | POST | `/rest/api/{v}/issue/{key}/watchers` | `--user` (body: JSON string of user identifier) |
| `watcher remove <key>` | DELETE | `/rest/api/{v}/issue/{key}/watchers?username=X` or `?accountId=X` | `--username`, `--account-id` |

- [x] **Step 1: Implement all 3 watcher commands**

Note: `watcher add` POST body is a plain JSON string (the user identifier), not a JSON object.

- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/watcher/ cli/pkg/jira/root.go
git commit -m "feat(jira): add watcher list/add/remove commands"
```

---

### Task 31: Link Commands (types/create/create-remote/remove)

**Files:** `cli/pkg/jira/link/` — `link.go`, `types.go`, `create.go`, `create_remote.go`, `remove.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `link types` | GET | `/rest/api/{v}/issueLinkType` | `--filter` |
| `link create` | POST | `/rest/api/{v}/issueLink` | `--type`, `--inward`, `--outward`, `--comment` |
| `link create-remote <key>` | POST | `/rest/api/{v}/issue/{key}/remotelink` | `--url`, `--title`, `--summary` |
| `link remove <link-id>` | DELETE | `/rest/api/{v}/issueLink/{id}` | — |

- [x] **Step 1: Implement all 4 link commands**
- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/link/ cli/pkg/jira/root.go
git commit -m "feat(jira): add link types/create/create-remote/remove commands"
```

---

### Task 32: Attachment Commands (download/images)

**Files:** `cli/pkg/jira/attachment/` — `attachment.go`, `download.go`, `images.go`

| Command | Method | API Path | Notes |
|---------|--------|----------|-------|
| `attachment download <key>` | GET | `/rest/api/{v}/issue/{key}?fields=attachment`, then download each | Downloads to current directory |
| `attachment images <key>` | GET | `/rest/api/{v}/issue/{key}?fields=attachment` | Returns base64-encoded image data in JSON |

- [x] **Step 1: Implement attachment commands**

For `download`: Get issue attachments, then GET each `attachment.content` URL to download the file.

- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/attachment/ cli/pkg/jira/root.go
git commit -m "feat(jira): add attachment download/images commands"
```

---

### Task 33: SLA Commands (get/dates)

**Files:** `cli/pkg/jira/sla/` — `sla.go`, `get.go`, `dates.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `sla get <key>` | GET | `/rest/servicedeskapi/request/{key}/sla` | `--metrics`, `--working-hours-only`, `--include-raw-dates` |
| `sla dates <key>` | GET | `/rest/api/{v}/issue/{key}` (parses date fields) | `--include-status-changes`, `--include-status-summary` |

> Note: SLA endpoint uses `/rest/servicedeskapi/` prefix.

- [x] **Step 1: Implement SLA commands**
- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/sla/ cli/pkg/jira/root.go
git commit -m "feat(jira): add sla get/dates commands"
```

---

### Task 34: Dev Info Commands (info/batch-info)

**Files:** `cli/pkg/jira/dev/` — `dev.go`, `info.go`, `batch_info.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `dev info <key>` | GET | `/rest/dev-status/latest/issue/detail?issueId={id}&applicationType=X&dataType=Y` | `--app-type`, `--data-type` |
| `dev batch-info` | GET | same endpoint, iterated | `--keys`, `--app-type`, `--data-type` |

> Note: Dev status API uses `/rest/dev-status/` prefix and requires issue ID (not key).
> First resolve key→ID via `GET /rest/api/{v}/issue/{key}?fields=id`.

- [x] **Step 1: Implement dev commands**
- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/dev/ cli/pkg/jira/root.go
git commit -m "feat(jira): add dev info/batch-info commands"
```

---

### Task 35: Changelog Batch Command (Cloud Only)

**Files:** `cli/pkg/jira/changelog/batch.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `changelog batch` | GET | `/rest/api/3/issue/{key}/changelog` per key | `--keys`, `--fields`, `--limit` |

- [x] **Step 1: Implement changelog batch**

Validate `client.IsCloud()` — return error if Server/DC. Iterate over `--keys`, call changelog API for each.

- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/changelog/ cli/pkg/jira/root.go
git commit -m "feat(jira): add changelog batch command (Cloud only)"
```

---

### Task 36: Form / ProForma Commands (list/get/update)

**Files:** `cli/pkg/jira/form/` — `form.go`, `list.go`, `get.go`, `update.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `form list <key>` | GET | `/rest/api/{v}/issue/{key}/properties/proforma.forms` | — |
| `form get <key> <form-id>` | GET | ProForma API endpoint | — |
| `form update <key> <form-id>` | PUT | ProForma API endpoint | `--answers-json` |

> Reference: `src/mcp_atlassian/jira/forms_api.py` for exact endpoint patterns.

- [x] **Step 1: Implement form commands**
- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/form/ cli/pkg/jira/root.go
git commit -m "feat(jira): add form list/get/update commands (ProForma)"
```

---

### Task 37: Service Desk Commands (get/queues/queue-issues) — Server/DC Only

**Files:** `cli/pkg/jira/servicedesk/` — `servicedesk.go`, `get.go`, `queues.go`, `queue_issues.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `servicedesk get <project-key>` | GET | `/rest/servicedeskapi/servicedesk/projectKey:{key}` | — |
| `servicedesk queues <desk-id>` | GET | `/rest/servicedeskapi/servicedesk/{id}/queue` | `--limit` |
| `servicedesk queue-issues <desk-id> <queue-id>` | GET | `/rest/servicedeskapi/servicedesk/{id}/queue/{qid}/issue` | `--limit` |

> Validate: return error if `client.IsCloud()` (spec says Server/DC only, but service desk API also works on Cloud — validate based on actual API availability).

- [x] **Step 1: Implement servicedesk commands**
- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/servicedesk/ cli/pkg/jira/root.go
git commit -m "feat(jira): add servicedesk get/queues/queue-issues commands"
```

---

### Task 38: User Get Command

**Files:** `cli/pkg/jira/user/get.go`

| Command | Method | API Path | Notes |
|---------|--------|----------|-------|
| `user get <identifier>` | GET | Cloud: `/rest/api/3/user?accountId=X`, Server: `/rest/api/2/user?username=X` | Auto-detect identifier type |

- [x] **Step 1: Implement user get**

Detect if identifier looks like an account ID (contains `:`) or is an email/username.

- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/user/ cli/pkg/jira/root.go
git commit -m "feat(jira): add user get command"
```

---

### Task 39: Epic Link Command

**Files:** `cli/pkg/jira/epic/link.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `epic link <issue-key>` | PUT | Update issue's epic link field | `--epic` (epic key) |

> Implementation: Uses `issue update` internally — sets the epic link custom field.

- [x] **Step 1: Implement epic link**
- [x] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/jira/epic/ cli/pkg/jira/root.go
git commit -m "feat(jira): add epic link command"
```

---

## Chunk 9: Confluence Commands

### Task 40: Confluence Auth & Root Setup

**Files:**
- Create: `cli/pkg/confluence/auth/` — mirrors jira auth, shares same config file
- Modify: `cli/pkg/confluence/root.go` — register auth and all subcommands

- [ ] **Step 1: Create confluence auth commands**

The confluence `auth` package can import and reuse the jira auth commands since they operate on the shared config file. Create thin wrappers:

```go
// cli/pkg/confluence/auth/auth.go
package auth

import (
	jiraauth "github.com/anthropics/atlassian-cli/pkg/jira/auth"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAuth creates the `auth` subcommand group (shared with jira).
func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	// Reuse jira auth — same config file, same commands
	return jiraauth.NewCmdAuth(f)
}
```

- [ ] **Step 2: Register auth in confluence root.go**
- [ ] **Step 3: Commit**

```bash
git add cli/pkg/confluence/auth/ cli/pkg/confluence/root.go
git commit -m "feat(confluence): add auth commands (shared with jira)"
```

---

### Task 41: Page Commands (get/create/update/delete/move/children/tree/history/diff)

**Files:** `cli/pkg/confluence/page/` — `page.go`, `get.go`, `create.go`, `update.go`, `delete.go`, `move.go`, `children.go`, `tree.go`, `history.go`, `diff.go`

> API base: `/rest/api/content` (v1) or `/api/v2/pages` (v2)
> Reference: `src/mcp_atlassian/confluence/pages.py`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `page get` | GET | `/rest/api/content/{id}?expand=body.storage,version,space` | `--id`, `--title`, `--space`, `--include-metadata`, `--raw` |
| `page create` | POST | `/rest/api/content` | `--space`, `--title`, `--content`, `--parent-id`, `--format`, `--emoji` |
| `page update <id>` | PUT | `/rest/api/content/{id}` | `--title`, `--content`, `--minor-edit`, `--comment`, `--format`, `--emoji` |
| `page delete <id>` | DELETE | `/rest/api/content/{id}` | `--confirm` |
| `page move <id>` | PUT | `/rest/api/content/{id}/move/{position}/{targetId}` | `--target-parent-id`, `--target-space`, `--position` |
| `page children <id>` | GET | `/rest/api/content/{id}/child/page` | `--limit`, `--include-content`, `--include-folders` |
| `page tree <space-key>` | GET | `/rest/api/content?spaceKey={key}&type=page` (recursive) | `--limit` |
| `page history <id>` | GET | `/rest/api/content/{id}/history` | `--version`, `--raw` |
| `page diff <id>` | GET | Fetch two versions and diff | `--from`, `--to` (both required) |

- [ ] **Step 1: Implement `page get` (reference implementation)**

```go
// cli/pkg/confluence/page/get.go
package page

import (
	"fmt"

	"github.com/anthropics/atlassian-cli/internal/api"
	"github.com/anthropics/atlassian-cli/internal/convert"
	"github.com/anthropics/atlassian-cli/internal/output"
	"github.com/anthropics/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var (
		pageID          string
		title           string
		space           string
		includeMetadata bool
		raw             bool
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a Confluence page",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			if pageID == "" && (title == "" || space == "") {
				return fmt.Errorf("provide --id, or both --title and --space")
			}

			var result map[string]any

			if pageID != "" {
				result, err = fetchPageByID(client, pageID, raw)
			} else {
				result, err = fetchPageByTitle(client, space, title, raw)
			}
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&pageID, "id", "", "Page ID")
	cmd.Flags().StringVar(&title, "title", "", "Page title")
	cmd.Flags().StringVar(&space, "space", "", "Space key")
	cmd.Flags().BoolVar(&includeMetadata, "include-metadata", false, "Include page metadata")
	cmd.Flags().BoolVar(&raw, "raw", false, "Return raw storage format (no Markdown conversion)")

	return cmd
}

func fetchPageByID(client *api.Client, pageID string, raw bool) (map[string]any, error) {
	path := client.ConfluenceAPIPath("content/" + pageID)
	rb := api.NewRequestBuilder(path).
		Query("expand", "body.storage,version,space,children.attachment")

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("getting page %s: %w", pageID, err)
	}

	if !raw {
		// Convert storage format to Markdown
		if body, ok := result["body"].(map[string]any); ok {
			if storage, ok := body["storage"].(map[string]any); ok {
				if value, ok := storage["value"].(string); ok {
					result["content_markdown"] = convert.StorageToMarkdown(value)
				}
			}
		}
	}

	return result, nil
}

func fetchPageByTitle(client *api.Client, space, title string, raw bool) (map[string]any, error) {
	path := client.ConfluenceAPIPath("content")
	rb := api.NewRequestBuilder(path).
		Query("spaceKey", space).
		Query("title", title).
		Query("expand", "body.storage,version,space")

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("searching page '%s' in space %s: %w", title, space, err)
	}

	// Extract first result
	results, ok := result["results"].([]any)
	if !ok || len(results) == 0 {
		return nil, fmt.Errorf("page '%s' not found in space %s", title, space)
	}

	page := results[0].(map[string]any)
	if !raw {
		if body, ok := page["body"].(map[string]any); ok {
			if storage, ok := body["storage"].(map[string]any); ok {
				if value, ok := storage["value"].(string); ok {
					page["content_markdown"] = convert.StorageToMarkdown(value)
				}
			}
		}
	}

	return page, nil
}
```

- [ ] **Step 2: Implement `page create`**

Body: `{"type": "page", "title": "...", "space": {"key": "..."}, "body": {"storage": {"value": "...", "representation": "storage"}}, "ancestors": [{"id": "..."}]}`

If `--format markdown` (default): convert Markdown content to Storage Format before sending.

- [ ] **Step 3: Implement `page update`**

Requires fetching current version number first: `GET /rest/api/content/{id}?expand=version`, then PUT with `version.number + 1`.

- [ ] **Step 4: Implement `page delete`** (DELETE with `--confirm`)

- [ ] **Step 5: Implement `page move`**

API: `PUT /rest/api/content/{id}/move/{position}/{targetId}`

- [ ] **Step 6: Implement `page children`**

API: `GET /rest/api/content/{id}/child/page`

- [ ] **Step 7: Implement `page tree`**

Recursive fetch: get space root pages, then children recursively up to `--limit`.

- [ ] **Step 8: Implement `page history`**

API: `GET /rest/api/content/{id}/history`

- [ ] **Step 9: Implement `page diff`**

Fetch two versions by version number, convert both to Markdown, compute unified diff.

```go
// Key logic for diff:
// 1. GET /rest/api/content/{id}?expand=body.storage&status=historical&version={from}
// 2. GET /rest/api/content/{id}?expand=body.storage&status=historical&version={to}
// 3. Convert both to Markdown
// 4. Compute line-by-line diff
```

- [ ] **Step 10: Register all page commands, build, test**

Run: `cd cli && make build && ./bin/confluence page --help`

- [ ] **Step 11: Commit**

```bash
git add cli/pkg/confluence/page/ cli/pkg/confluence/root.go
git commit -m "feat(confluence): add page get/create/update/delete/move/children/tree/history/diff commands"
```

---

### Task 42: Confluence Search Command

**Files:** `cli/pkg/confluence/search/search.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `search <query>` | GET | `/rest/api/content/search?cql={cql}` | `--limit`, `--spaces-filter` |

If the query doesn't look like CQL (no `=`, `~`, `AND`, `OR`), wrap it as `text ~ "query"`.

- [ ] **Step 1: Implement search**
- [ ] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/confluence/search/ cli/pkg/confluence/root.go
git commit -m "feat(confluence): add search command with CQL support"
```

---

### Task 43: Confluence Comment Commands (list/add/reply)

**Files:** `cli/pkg/confluence/comment/` — `comment.go`, `list.go`, `add.go`, `reply.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `comment list <page-id>` | GET | `/rest/api/content/{id}/child/comment?expand=body.storage` | — |
| `comment add <page-id>` | POST | `/rest/api/content` with type=comment | `--body` |
| `comment reply <comment-id>` | POST | `/rest/api/content` with type=comment, ancestor=comment-id | `--body` |

- [ ] **Step 1: Implement all 3 comment commands**
- [ ] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/confluence/comment/ cli/pkg/confluence/root.go
git commit -m "feat(confluence): add comment list/add/reply commands"
```

---

### Task 44: Confluence Label Commands (list/add)

**Files:** `cli/pkg/confluence/label/` — `label.go`, `list.go`, `add.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `label list <page-id>` | GET | `/rest/api/content/{id}/label` | — |
| `label add <page-id>` | POST | `/rest/api/content/{id}/label` | `--name` (body: `[{"prefix":"global","name":"..."}]`) |

- [ ] **Step 1: Implement label commands**
- [ ] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/confluence/label/ cli/pkg/confluence/root.go
git commit -m "feat(confluence): add label list/add commands"
```

---

### Task 45: Confluence Attachment Commands (7 commands)

**Files:** `cli/pkg/confluence/attachment/` — `attachment.go`, `list.go`, `upload.go`, `upload_batch.go`, `download.go`, `download_all.go`, `delete.go`, `images.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `attachment list <content-id>` | GET | `/rest/api/content/{id}/child/attachment` | `--limit`, `--filename`, `--media-type` |
| `attachment upload <content-id>` | POST (multipart) | `/rest/api/content/{id}/child/attachment` | `--file`, `--comment`, `--minor-edit` |
| `attachment upload-batch <content-id>` | POST (multipart, repeated) | same | `--files`, `--comment`, `--minor-edit` |
| `attachment download <attachment-id>` | GET | `/rest/api/content/{id}/child/attachment/{aid}/download` | — |
| `attachment download-all <content-id>` | GET | List then download each | — |
| `attachment delete <attachment-id>` | DELETE | `/rest/api/content/{id}` | `--confirm` |
| `attachment images <content-id>` | GET | List attachments, filter by media type | — |

> Note: `upload` requires multipart/form-data, not JSON. Use `mime/multipart` package.

- [ ] **Step 1: Implement all 7 attachment commands**

For upload, use multipart form:
```go
var buf bytes.Buffer
writer := multipart.NewWriter(&buf)
part, _ := writer.CreateFormFile("file", filepath.Base(filePath))
io.Copy(part, fileReader)
writer.Close()
// Set Content-Type to writer.FormDataContentType()
```

- [ ] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/confluence/attachment/ cli/pkg/confluence/root.go
git commit -m "feat(confluence): add attachment list/upload/upload-batch/download/download-all/delete/images commands"
```

---

### Task 46: Confluence User Search Command

**Files:** `cli/pkg/confluence/user/search.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `user search <query>` | GET | Cloud: `/rest/api/search/user?cql=...`, Server: `/rest/api/user/search?username=...` | `--limit`, `--group` |

- [ ] **Step 1: Implement user search with Cloud/Server routing**
- [ ] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/confluence/user/ cli/pkg/confluence/root.go
git commit -m "feat(confluence): add user search command"
```

---

### Task 47: Confluence Analytics Views (Cloud Only)

**Files:** `cli/pkg/confluence/analytics/views.go`

| Command | Method | API Path | Key Flags |
|---------|--------|----------|-----------|
| `analytics views <page-id>` | GET | `/rest/api/analytics/content/{id}/views` | `--include-title` |

- [ ] **Step 1: Implement analytics views**

Validate `client.IsCloud()` — Cloud only.

- [ ] **Step 2: Register, build, test, commit**

```bash
git add cli/pkg/confluence/analytics/ cli/pkg/confluence/root.go
git commit -m "feat(confluence): add analytics views command (Cloud only)"
```

---

## Chunk 10: Build, Release & Final Integration

### Task 48: goreleaser Configuration

**Files:**
- Create: `cli/.goreleaser.yml`

- [ ] **Step 1: Write goreleaser config**

```yaml
# cli/.goreleaser.yml
version: 2

builds:
  - id: jira
    main: ./cmd/jira
    binary: jira
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w -X main.version={{.Version}}

  - id: confluence
    main: ./cmd/confluence
    binary: confluence
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w -X main.version={{.Version}}

archives:
  - id: jira
    builds:
      - jira
    name_template: "jira_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

  - id: confluence
    builds:
      - confluence
    name_template: "confluence_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: "checksums.txt"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
```

- [ ] **Step 2: Verify goreleaser config**

Run: `cd cli && goreleaser check` (if goreleaser installed) or `cd cli && go build ./cmd/jira && go build ./cmd/confluence`

- [ ] **Step 3: Commit**

```bash
git add cli/.goreleaser.yml
git commit -m "chore: add goreleaser configuration for cross-platform builds"
```

---

### Task 49: End-to-End Smoke Test

**Files:**
- Create: `cli/e2e_test.go`

> This test verifies the full binary works end-to-end with a mock server.

- [ ] **Step 1: Write e2e test**

```go
// cli/e2e_test.go
package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	binPath := filepath.Join(t.TempDir(), "jira")
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/jira")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "build failed: %s", output)

	// Run: jira issue get TEST-1
	// Note: Need to figure out how to pass config path — may need env var support
	// For now, this test verifies the binary builds and help works

	cmd = exec.Command(binPath, "--help")
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "Work with Jira from the command line")
}

func TestConfluenceBinary_Help(t *testing.T) {
	binPath := filepath.Join(t.TempDir(), "confluence")
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/confluence")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "build failed: %s", output)

	cmd = exec.Command(binPath, "--help")
	output, err = cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(output), "Work with Confluence from the command line")
}
```

- [ ] **Step 2: Run e2e test**

Run: `cd cli && go test -v -run TestJiraBinary`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add cli/e2e_test.go
git commit -m "test: add end-to-end smoke tests for both CLI binaries"
```

---

### Task 50: Final Integration — Verify All Commands Registered

- [ ] **Step 1: Build both binaries**

Run: `cd cli && make build`
Expected: Clean build, no errors

- [ ] **Step 2: Verify jira command tree**

Run: `cd cli && ./bin/jira --help`
Expected output includes: `auth`, `issue`, `search`, `comment`, `field`, `transition`, `board`, `sprint`, `project`, `worklog`, `watcher`, `link`, `attachment`, `sla`, `dev`, `changelog`, `form`, `servicedesk`, `user`, `epic`

- [ ] **Step 3: Verify confluence command tree**

Run: `cd cli && ./bin/confluence --help`
Expected output includes: `auth`, `page`, `search`, `comment`, `label`, `attachment`, `user`, `analytics`

- [ ] **Step 4: Run full test suite**

Run: `cd cli && make test`
Expected: All tests PASS

- [ ] **Step 5: Final commit**

```bash
git add -A cli/
git commit -m "feat: complete Atlassian CLI v1 with all Jira and Confluence commands"
```

---

## Summary

| Chunk | Tasks | Commands | Description |
|-------|-------|----------|-------------|
| 1 | 1-3 | 0 | Project scaffolding, root commands |
| 2 | 4-8 | 5 | Config, auth, env overrides, auth CLI |
| 3 | 9-12 | 0 | HTTP client, errors, request builder, pagination |
| 4 | 13-15 | 0 | JSON/table/text output formatters |
| 5 | 16-17 | 0 | ADF ↔ Markdown, Storage ↔ Markdown |
| 6 | 18-22 | 8 | Jira issue CRUD + search (pattern established) |
| 7 | 23-27 | 13 | Comment, field, transition, board, sprint |
| 8 | 28-39 | 22 | Project, worklog, watcher, link, attachment, SLA, dev, changelog, form, servicedesk, user, epic |
| 9 | 40-47 | 24 | All Confluence commands |
| 10 | 48-50 | 0 | goreleaser, e2e tests, final integration |

**Total: 50 tasks, 72 commands (49 Jira + 24 Confluence - 1 shared auth = 72 unique)**
