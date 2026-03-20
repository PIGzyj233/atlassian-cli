# Atlassian CLI Design Spec

> Two standalone Go CLI tools (`jira`, `confluence`) that expose all current MCP-Atlassian capabilities as shell commands, designed for both human operators and LLM-driven automation.

---

## 1. Goals and Non-Goals

### Goals

- Replace MCP tool interface with two native CLI binaries (`jira`, `confluence`)
- Persistent local credential storage so users authenticate once
- Full feature parity with the existing 49 Jira tools and 24 Confluence tools
- JSON-first output suitable for LLM consumption; switchable to table/text for humans
- Cross-platform support: Linux, macOS, Windows (amd64 + arm64)
- Inspired by GitHub CLI (`gh`) in UX patterns

### Non-Goals

- OAuth 2.0 support (Basic Auth and PAT cover the primary use cases; OAuth can be added later)
- MCP protocol compatibility (this is a clean break from MCP)
- GUI or TUI interactive mode
- Plugin/extension system

---

## 2. Technology Stack

| Component | Choice | Rationale |
|---|---|---|
| Language | Go 1.22+ | Fast compilation, easy cross-compilation, gh CLI precedent |
| CLI framework | cobra | Industry standard, same as gh/kubectl |
| Config management | viper | Integrates with cobra, supports YAML/env/flags |
| YAML | go-yaml v3 | Config file format |
| HTTP client | net/http | Standard library, sufficient for REST APIs |
| Table output | tablewriter or lipgloss | Terminal-friendly formatting |
| Release | goreleaser | Cross-platform builds and GitHub Releases |
| Testing | testing + testify | Standard Go testing with assertion helpers |

---

## 3. Repository Structure

```
atlassian-cli/
├── cmd/
│   ├── jira/
│   │   └── main.go                 # jira binary entry point
│   └── confluence/
│       └── main.go                 # confluence binary entry point
├── internal/
│   ├── config/
│   │   ├── config.go               # Config struct, Load/Save, XDG paths
│   │   ├── auth.go                 # Authenticator interface, Basic, PAT
│   │   └── hosts.go                # Multi-instance host management
│   ├── api/
│   │   ├── client.go               # HTTP client with auth, retry, timeout
│   │   ├── request.go              # Request builder (headers, query params)
│   │   ├── errors.go               # API error parsing and formatting
│   │   └── pagination.go           # Auto-pagination helper
│   ├── output/
│   │   ├── formatter.go            # Formatter interface
│   │   ├── json.go                 # JSON output
│   │   ├── table.go                # Table output
│   │   └── text.go                 # Plain text output
│   └── convert/
│       ├── adf.go                  # ADF (Atlassian Document Format) ↔ Markdown
│       └── storage.go              # Storage Format ↔ Markdown
├── pkg/
│   ├── jira/
│   │   ├── root.go                 # jira root command
│   │   ├── auth/                   # auth login/logout/status/list/switch
│   │   ├── issue/                  # issue get/create/batch-create/update/delete/transition
│   │   ├── search/                 # search (JQL)
│   │   ├── comment/                # comment list/add/edit
│   │   ├── field/                  # field search/options
│   │   ├── transition/             # transition list
│   │   ├── board/                  # board list/issues
│   │   ├── sprint/                 # sprint list/issues/create/update/add-issues
│   │   ├── project/                # project list/versions/components/version create
│   │   ├── worklog/                # worklog list/add
│   │   ├── watcher/                # watcher list/add/remove
│   │   ├── link/                   # link types/create/create-remote/remove
│   │   ├── attachment/             # attachment download/images
│   │   ├── sla/                    # sla get/dates
│   │   ├── dev/                    # dev info/batch-info
│   │   ├── changelog/              # changelog batch (Cloud only)
│   │   ├── form/                   # form list/get/update (ProForma)
│   │   ├── servicedesk/            # servicedesk get/queues/queue-issues (Server/DC)
│   │   ├── user/                   # user get
│   │   └── epic/                   # epic link
│   └── confluence/
│       ├── root.go                 # confluence root command
│       ├── auth/                   # auth (shared with jira)
│       ├── page/                   # page get/create/update/delete/move/children/tree/history/diff
│       ├── search/                 # search
│       ├── comment/                # comment list/add/reply
│       ├── label/                  # label list/add
│       ├── attachment/             # attachment list/upload/upload-batch/download/download-all/delete/images
│       ├── user/                   # user search
│       └── analytics/              # analytics views (Cloud only)
├── go.mod
├── go.sum
├── Makefile
├── .goreleaser.yml
└── README.md
```

---

## 4. Configuration and Authentication

### 4.1 Config File

Location follows XDG Base Directory spec:
- Linux/macOS: `~/.config/atlassian-cli/config.yml`
- Windows: `%APPDATA%\atlassian-cli\config.yml`

```yaml
hosts:
  company.atlassian.net:
    type: cloud                     # cloud | server | datacenter
    auth: basic                     # basic | pat
    username: user@company.com
    token: "ATATT3x..."
    jira:
      enabled: true
      base_path: ""
    confluence:
      enabled: true
      base_path: /wiki

  jira.internal.corp:
    type: datacenter
    auth: pat
    token: "NjM2..."
    ssl_verify: false
    jira:
      enabled: true
    confluence:
      enabled: false

defaults:
  host: company.atlassian.net
  output: json                      # json | text | table
  jira:
    project: ""                     # optional default project
  confluence:
    space: ""                       # optional default space
```

File permissions: 0600 (owner read/write only).

### 4.2 Authentication Methods

| Method | Scope | Credentials |
|---|---|---|
| Basic Auth | Cloud + Server/DC | Email/username + API token |
| PAT | Server/DC only | Personal Access Token |

### 4.3 Auth Commands

```
jira auth login           # Interactive: URL → type → method → credentials → verify → save
jira auth logout [--host] # Remove credentials for a host
jira auth status          # Show current auth info
jira auth list            # List all configured hosts
jira auth switch <host>   # Change default host
```

The `confluence auth` commands share the same config file and behave identically.

### 4.4 Environment Variable Overrides

For CI/CD or non-interactive use:

```bash
ATLASSIAN_HOST=company.atlassian.net
ATLASSIAN_USERNAME=user@company.com
ATLASSIAN_TOKEN=xxx
ATLASSIAN_AUTH_TYPE=basic   # basic | pat
```

Environment variables override config file values when set.

---

## 5. Command Reference

### 5.1 Global Flags

All commands accept:

| Flag | Description | Default |
|---|---|---|
| `--host <hostname>` | Target Atlassian instance | Config default |
| `--output json\|text\|table` | Output format | `json` |
| `--verbose` | Enable debug logging to stderr | `false` |
| `--no-color` | Disable ANSI colors | `false` |
| `--help` / `-h` | Display command-specific usage | — |
| `--version` | Display installed version (root command only) | — |

### 5.1.1 Pagination

Commands that return lists support pagination via:

| Flag | Description | Default |
|---|---|---|
| `--limit <N>` | Maximum number of results to return | Command-specific (typically 50) |
| `--start-at <offset>` | Starting index for offset-based pagination | 0 |

By default, commands return a single page of results up to `--limit`. The JSON output includes a `pagination` object with `startAt`, `maxResults`, `total`, and `isLast` fields so callers can implement their own pagination loop.

Batch operations (e.g., `jira issue batch-create`, `jira dev batch-info`) process items serially. Parallelization is not planned for v1.

### 5.2 Jira Commands

#### Auth (5 commands)

| Command | MCP Equivalent | Description |
|---|---|---|
| `jira auth login` | N/A | Interactive authentication setup |
| `jira auth logout` | N/A | Remove stored credentials |
| `jira auth status` | N/A | Show current authentication state |
| `jira auth list` | N/A | List all configured instances |
| `jira auth switch` | N/A | Change default instance |

#### Issue (6 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira issue get <key>` | `jira_get_issue` | `--fields`, `--expand`, `--comment-limit` |
| `jira issue create` | `jira_create_issue` | `--project`, `--summary`, `--type`, `--assignee`, `--description`, `--components`, `--fields-json` |
| `jira issue batch-create` | `jira_batch_create_issues` | `--file <json>`, `--validate-only` |
| `jira issue update <key>` | `jira_update_issue` | `--fields-json`, `--components`, `--attachments` |
| `jira issue delete <key>` | `jira_delete_issue` | `--confirm` |
| `jira issue transition <key>` | `jira_transition_issue` | `--transition-id`, `--fields-json`, `--comment` |

#### Search (1 command)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira search <jql>` | `jira_search` | `--fields`, `--limit`, `--start-at`, `--projects-filter`, `--expand` |

#### Comment (3 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira comment list <issue-key>` | (part of jira_get_issue) | — |
| `jira comment add <issue-key>` | `jira_add_comment` | `--body`, `--visibility`, `--public` |
| `jira comment edit <issue-key> <comment-id>` | `jira_edit_comment` | `--body`, `--visibility` |

> Note: Comment deletion is not supported (not available in the current Atlassian REST API for all deployment types).

#### Field (2 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira field search` | `jira_search_fields` | `--keyword`, `--limit`, `--refresh` |
| `jira field options <field-id>` | `jira_get_field_options` | `--context-id`, `--project`, `--issue-type`, `--contains` |

#### Transition (1 command)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira transition list <issue-key>` | `jira_get_transitions` | — |

#### Board (2 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira board list` | `jira_get_agile_boards` | `--name`, `--project`, `--type`, `--limit` |
| `jira board issues <board-id>` | `jira_get_board_issues` | `--jql`, `--fields`, `--limit` |

#### Sprint (5 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira sprint list <board-id>` | `jira_get_sprints_from_board` | `--state`, `--limit` |
| `jira sprint issues <sprint-id>` | `jira_get_sprint_issues` | `--fields`, `--limit` |
| `jira sprint create <board-id>` | `jira_create_sprint` | `--name`, `--start-date`, `--end-date`, `--goal` |
| `jira sprint update <sprint-id>` | `jira_update_sprint` | `--name`, `--state`, `--start-date`, `--end-date`, `--goal` |
| `jira sprint add-issues <sprint-id>` | `jira_add_issues_to_sprint` | `--keys` |

#### Project (4 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira project list` | `jira_get_all_projects` | `--include-archived` |
| `jira project versions <key>` | `jira_get_project_versions` | — |
| `jira project components <key>` | `jira_get_project_components` | — |
| `jira project version create <key>` | `jira_create_version` / `jira_batch_create_versions` | `--name`, `--start-date`, `--release-date`, `--file` |

> Version create modes: If `--file <json>` is provided, batch creation mode is used and `--name`/`--start-date`/`--release-date`/`--description` flags are ignored. Otherwise, `--name` is required for single-version creation.

#### Worklog (2 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira worklog list <issue-key>` | `jira_get_worklog` | — |
| `jira worklog add <issue-key>` | `jira_add_worklog` | `--time-spent`, `--comment`, `--started` |

> Note: Worklog deletion/update is not supported (not available in the current MCP implementation).

#### Watcher (3 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira watcher list <issue-key>` | `jira_get_issue_watchers` | — |
| `jira watcher add <issue-key>` | `jira_add_watcher` | `--user` |
| `jira watcher remove <issue-key>` | `jira_remove_watcher` | `--username`, `--account-id` |

#### Link (4 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira link types` | `jira_get_link_types` | `--filter` |
| `jira link create` | `jira_create_issue_link` | `--type`, `--inward`, `--outward`, `--comment` |
| `jira link create-remote <issue-key>` | `jira_create_remote_issue_link` | `--url`, `--title`, `--summary` |
| `jira link remove <link-id>` | `jira_remove_issue_link` | — |

#### Attachment (2 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira attachment download <issue-key>` | `jira_download_attachments` | — |
| `jira attachment images <issue-key>` | `jira_get_issue_images` | — |

#### SLA (2 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira sla get <issue-key>` | `jira_get_issue_sla` | `--metrics`, `--working-hours-only`, `--include-raw-dates` |
| `jira sla dates <issue-key>` | `jira_get_issue_dates` | `--include-status-changes`, `--include-status-summary` |

#### Dev (2 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira dev info <issue-key>` | `jira_get_issue_development_info` | `--app-type`, `--data-type` |
| `jira dev batch-info` | `jira_get_issues_development_info` | `--keys`, `--app-type`, `--data-type` |

#### Changelog (1 command, Cloud only)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira changelog batch` | `jira_batch_get_changelogs` | `--keys`, `--fields`, `--limit` |

#### Form / ProForma (3 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira form list <issue-key>` | `jira_get_issue_proforma_forms` | — |
| `jira form get <issue-key> <form-id>` | `jira_get_proforma_form_details` | — |
| `jira form update <issue-key> <form-id>` | `jira_update_proforma_form_answers` | `--answers-json` |

#### Service Desk (3 commands, Server/DC only)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira servicedesk get <project-key>` | `jira_get_service_desk_for_project` | — |
| `jira servicedesk queues <desk-id>` | `jira_get_service_desk_queues` | `--limit` |
| `jira servicedesk queue-issues <desk-id> <queue-id>` | `jira_get_queue_issues` | `--limit` |

#### User (1 command)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira user get <identifier>` | `jira_get_user_profile` | — |

#### Epic (1 command)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `jira epic link <issue-key>` | `jira_link_to_epic` | `--epic` |

### 5.3 Confluence Commands

#### Page (9 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `confluence page get` | `confluence_get_page` | `--id`, `--title`, `--space`, `--include-metadata`, `--raw` |
| `confluence page create` | `confluence_create_page` | `--space`, `--title`, `--content`, `--parent-id`, `--format`, `--emoji` |
| `confluence page update <page-id>` | `confluence_update_page` | `--title`, `--content`, `--minor-edit`, `--comment`, `--format`, `--emoji` |
| `confluence page delete <page-id>` | `confluence_delete_page` | `--confirm` |
| `confluence page move <page-id>` | `confluence_move_page` | `--target-parent-id`, `--target-space`, `--position` |
| `confluence page children <parent-id>` | `confluence_get_page_children` | `--limit`, `--include-content`, `--include-folders` |
| `confluence page tree <space-key>` | `confluence_get_space_page_tree` | `--limit` |
| `confluence page history <page-id>` | `confluence_get_page_history` | `--version`, `--raw` |
| `confluence page diff <page-id>` | `confluence_get_page_diff` | `--from <version>`, `--to <version>` |

> `confluence page diff` parameters: `--from` and `--to` are both **required** and accept integer version numbers (e.g., `--from 3 --to 5`). Output is a unified diff of the Markdown-converted content between the two versions.

#### Search (1 command)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `confluence search <query>` | `confluence_search` | `--limit`, `--spaces-filter` |

#### Comment (3 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `confluence comment list <page-id>` | `confluence_get_comments` | — |
| `confluence comment add <page-id>` | `confluence_add_comment` | `--body` |
| `confluence comment reply <comment-id>` | `confluence_reply_to_comment` | `--body` |

#### Label (2 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `confluence label list <page-id>` | `confluence_get_labels` | — |
| `confluence label add <page-id>` | `confluence_add_label` | `--name` |

> Note: Label deletion is not supported (not available in the current MCP implementation).

#### Attachment (7 commands)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `confluence attachment list <content-id>` | `confluence_get_attachments` | `--limit`, `--filename`, `--media-type` |
| `confluence attachment upload <content-id>` | `confluence_upload_attachment` | `--file`, `--comment`, `--minor-edit` |
| `confluence attachment upload-batch <content-id>` | `confluence_upload_attachments` | `--files`, `--comment`, `--minor-edit` |
| `confluence attachment download <attachment-id>` | `confluence_download_attachment` | — |
| `confluence attachment download-all <content-id>` | `confluence_download_content_attachments` | — |
| `confluence attachment delete <attachment-id>` | `confluence_delete_attachment` | `--confirm` |
| `confluence attachment images <content-id>` | `confluence_get_page_images` | — |

#### User (1 command)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `confluence user search <query>` | `confluence_search_user` | `--limit`, `--group` |

#### Analytics (1 command, Cloud only)

| Command | MCP Equivalent | Key Parameters |
|---|---|---|
| `confluence analytics views <page-id>` | `confluence_get_page_views` | `--include-title` |

---

## 6. Architecture

### 6.1 HTTP Client

```go
// internal/api/client.go
type Client struct {
    HTTP       *http.Client
    BaseURL    string
    Auth       Authenticator
    Headers    map[string]string
    MaxRetries int
}

type Authenticator interface {
    Apply(req *http.Request) error
}

type BasicAuth struct {
    Username string
    Token    string
}

type PATAuth struct {
    Token string
}
```

The client handles:
- Auth header injection via the `Authenticator` interface
- Retry with exponential backoff on 429 (rate limit) and 5xx errors
- Configurable timeout (default 75s, matching current MCP behavior)
- SSL verification toggle for self-signed certificates

### 6.2 Cloud vs Server/DC

```go
func (c *Client) IsCloud() bool {
    return strings.Contains(c.BaseURL, "atlassian.net")
}

func (c *Client) JiraAPIPath(path string) string {
    if c.IsCloud() {
        return "/rest/api/3/" + path
    }
    return "/rest/api/2/" + path
}
```

Commands that differ between Cloud and Server/DC (e.g., `jira servicedesk` is Server/DC only, `jira changelog batch` is Cloud only) validate the instance type before executing and return a clear error if unsupported.

### 6.3 Content Conversion

The current MCP project converts between ADF (Cloud), Storage Format (Server/DC), and Markdown. The CLI replicates this:

- **Input**: User writes Markdown → CLI converts to ADF (Cloud) or Storage Format (Server/DC) before sending to API
- **Output**: API returns ADF/Storage → CLI converts to Markdown for display (default). Use `--raw` for original format.
- **Implementation**: `internal/convert/adf.go` and `internal/convert/storage.go`

### 6.4 Output Formatting

```go
// internal/output/formatter.go
type Formatter interface {
    Format(data any) (string, error)
}

// Usage in commands:
func runIssueGet(cmd *cobra.Command, args []string) error {
    issue, err := client.GetIssue(args[0], opts)
    if err != nil {
        return err
    }
    return output.Print(cmd, issue)  // respects --output flag
}
```

- **JSON** (default): `encoding/json` with indent. Stable field names matching API response.
- **Table**: Predefined column layouts per resource type. Uses tablewriter.
- **Text**: Human-readable prose format with key fields highlighted.

All output goes to stdout. Logs and progress go to stderr.

---

## 7. Error Handling

### Exit Codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | User/input error (missing flags, invalid args) |
| 2 | API error (4xx/5xx from Atlassian) |
| 3 | Network/connectivity error |

### Error Format

Errors are written to stderr. In JSON output mode, errors are also written as JSON:

```json
{"error": "Issue INVALID-999 not found", "status": 404, "hint": "Check that the issue key is correct"}
```

The `hint` field is optional and included only when a standard remediation is known (e.g., 404 → "Check that the issue key is correct", 401 → "Run `jira auth login` to authenticate"). Hints are hardcoded for common error patterns, not parsed from API responses.

In text mode:
```
Error: Issue INVALID-999 not found (HTTP 404)
```

### Common Error Patterns

- Missing required flags → cobra validation, exit 1
- Auth not configured → suggest `jira auth login`, exit 1
- API 4xx → parse Atlassian error body, display message + hint, exit 2
- API 429 → retry with backoff, then fail with rate limit message, exit 2
- Network timeout/DNS → clear connectivity error, exit 3
- Unsupported operation (e.g., Cloud-only command on Server) → clear message, exit 1

---

## 8. Testing Strategy

### Unit Tests

- Config parsing and validation
- Output formatters (JSON, table, text)
- Markdown ↔ ADF/Storage conversion
- Command flag parsing and validation
- Error formatting

### Integration Tests

- `httptest.Server` mocking Atlassian REST API responses
- End-to-end command execution with mock server
- Auth flow testing (login, status, switch)
- Pagination handling

### E2E Tests (optional, requires live instance)

- Real API calls against a test Atlassian instance
- Gated by build tag (similar to current `@pytest.mark.integration`)

### Framework

- Standard `testing` package
- `testify/assert` and `testify/require` for assertions
- Table-driven tests for parameter variations

---

## 9. Build and Release

### goreleaser Configuration

```yaml
builds:
  - id: jira
    main: ./cmd/jira
    binary: jira
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w -X main.version={{.Version}}
  - id: confluence
    main: ./cmd/confluence
    binary: confluence
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w -X main.version={{.Version}}
```

### Distribution Channels

1. **GitHub Releases**: Automated via goreleaser on tag push
2. **Homebrew**: Tap repository for `brew install jira` / `brew install confluence`
3. **Scoop** (Windows): Bucket for `scoop install jira` / `scoop install confluence`

### Makefile Targets

```makefile
build:      Build both binaries
test:       Run all tests
lint:       Run golangci-lint
fmt:        Format code
release:    Run goreleaser
install:    Install binaries to $GOPATH/bin
```

---

## 10. LLM Integration Design

The CLI is designed as a stepping stone toward LLM skill integration:

- **JSON default**: Structured, parseable output without extra flags
- **Stable schema**: Field names match Atlassian API conventions
- **Exit codes**: Semantic (0/1/2/3) for programmatic error handling
- **No ANSI in stdout**: Colors only on stderr; `--no-color` disables globally
- **stderr/stdout separation**: stdout = data only, stderr = logs/progress
- **Self-documenting**: `--help` output is structured enough for an LLM to derive tool definitions
- **Idempotent reads**: All read commands are safe to retry
- **Deterministic output**: Same input → same output structure (pagination tokens aside)
