# Atlassian CLI

English | [简体中文](README_CN.md)

Atlassian CLI gives your coding agents — and you — direct access to Jira and Confluence from the command line. Two standalone binaries, `jira` and `confluence`, that work with **Cloud**, **Server**, and **Data Center** instances.

## How it works

You tell your coding agent "create a Jira ticket for the login bug" or "find the architecture doc in Confluence" — and it just does it. No browser, no copy-pasting, no context switching.

Under the hood, the agent uses the `jira` and `confluence` CLI binaries to talk to the Atlassian REST API. Output is JSON by default, perfect for agents to parse and act on. The binaries auto-detect whether your instance is Cloud or Server/Data Center, handle authentication via OS keychain, and retry on rate limits.

Because it's a skill, your agent picks it up automatically. You don't need to teach it anything — it already knows the commands.

And because it's also a plain CLI tool, you can use it yourself in the terminal, in scripts, in CI/CD pipelines. `--output table` for humans, `--output json` for machines.

## Acknowledgments

This project owes a great deal to **[mcp-atlassian](https://github.com/sooperset/mcp-atlassian)** by [@sooperset](https://github.com/sooperset) and its contributors. Their pioneering work on building an MCP server for Atlassian products was a major source of inspiration and reference. The API interaction patterns, the approach to supporting Cloud/Server/Data Center deployments, and the overall vision of bridging AI assistants with Atlassian tools — much of our design thinking was informed by studying their excellent implementation. If you need an MCP-native solution for Atlassian integration, we highly recommend checking out their project.

## Installation

**Note:** Installation differs by platform. Claude Code has a built-in skill system. Codex and other agents require manual setup.

### Claude Code (Marketplace)

Install the skill directly:

```bash
claude skill add --from https://github.com/PigZyj2333/atlassian-cli --subpath skill
```

### Claude Code (Manual)

Add to your project's `.claude/settings.json`:

```json
{
  "skills": [
    {
      "source": "https://github.com/PigZyj2333/atlassian-cli",
      "subpath": "skill"
    }
  ]
}
```

### Codex

Tell Codex:

```
Fetch and follow instructions from https://raw.githubusercontent.com/PigZyj2333/atlassian-cli/refs/heads/master/.codex/install.md
```

### Other AI Agents

Any agent that can run shell commands can use this CLI. Tell it:

```
Fetch and follow instructions from https://raw.githubusercontent.com/PigZyj2333/atlassian-cli/refs/heads/master/.codex/install.md
```

### Standalone (no agent)

If you just want the CLI tools for yourself:

**From GitHub Releases:**
```bash
# Linux / macOS
curl -fsSL https://github.com/PigZyj2333/atlassian-cli/releases/latest/download/atlassian-cli_*_$(uname -s | tr A-Z a-z)_$(uname -m).tar.gz | tar xz
sudo mv jira confluence /usr/local/bin/
```

**Go install:**
```bash
go install github.com/PigZyj2333/atlassian-cli/cmd/jira@latest
go install github.com/PigZyj2333/atlassian-cli/cmd/confluence@latest
```

### Verify Installation

Start a new session and try something like "search my open Jira tickets" or "find the design doc in Confluence". The agent should automatically invoke the skill and run the right commands.

For standalone usage:

```bash
jira --help
confluence --help
```

## The Basic Workflow

1. **Authenticate** — Run `jira auth login --host your-company.atlassian.net` once. Credentials persist in your OS keychain. Cloud uses email + API token; Server/DC uses username + password or PAT.

2. **Search** — Find issues with JQL (`jira search "assignee = currentUser()"`), or pages with CQL (`confluence search "space = ENG AND title ~ 'API'"`).

3. **Create** — Create issues (`jira issue create --project PROJ --type Bug --summary "..."`) or pages (`confluence page create --space ENG --title "..." --body "..."`).

4. **Update** — Transition issues (`jira issue transition PROJ-123 --status "In Progress"`), update fields, add comments, upload attachments.

5. **Manage** — Work with sprints, boards, epics, worklogs, watchers, links, labels, and more.

**The agent handles all of this automatically.** You describe what you want in natural language, it figures out the commands.

## What's Inside

### Jira Commands

- **auth** — Login, logout, check status, switch instances
- **issue** — Get, create, update, delete, transition, batch-create
- **search** — JQL-powered search with field selection
- **board** — List agile boards by project
- **sprint** — List, create, update sprints; manage sprint backlog
- **project** — List and inspect projects
- **comment** — List and add issue comments
- **worklog** — Time tracking (list, add)
- **attachment** — Upload and list file attachments
- **link** — Link issues together (blocks, relates to, etc.)
- **watcher** — Add and list issue watchers
- **transition** — List available status transitions
- **field** — List custom and system fields
- **user** — Search and get user details
- **epic** — Link issues to epics
- **servicedesk** — List service desk queues
- **sla** — SLA information for service desk tickets
- **changelog** — Issue change history
- **dev** — Development info (branches, commits, PRs)
- **form** — Proforma/form data

### Confluence Commands

- **auth** — Login, logout, check status
- **page** — Get, create, update, delete, children, tree, move, diff, history
- **search** — CQL-powered search
- **comment** — List and add page comments
- **label** — List and add page labels
- **attachment** — Upload and list file attachments
- **user** — Search users
- **analytics** — Page view analytics (Cloud only)

### Global Flags

```
--host string     Target Atlassian instance (overrides default)
--output string   Output format: json (default) | text | table
--verbose         Enable debug logging to stderr
--no-color        Disable colored output
```

## Architecture

Two independent binaries share a common infrastructure layer:

```
cmd/{jira,confluence}/main.go       Entry points
pkg/cmdutil/factory.go              Dependency injection (config → client)
pkg/{jira,confluence}/<feature>/    One package per feature
internal/api/                       HTTP client with retry, pagination, error handling
internal/config/                    YAML config, auth (BasicAuth/PAT), host resolution
internal/output/                    JSON / Table / Text formatters
internal/convert/                   ADF ↔ storage format (Confluence)
```

Key design decisions:
- **Factory pattern** — every command gets a `*cmdutil.Factory` for config and authenticated clients
- **Cloud vs Server auto-detection** — `*.atlassian.net` uses REST API v3, everything else uses v2
- **Retry with backoff** — 3x exponential retry on 429 and 5xx
- **Static binaries** — `CGO_ENABLED=0`, zero runtime dependencies

## Build from Source

```bash
git clone https://github.com/PigZyj2333/atlassian-cli.git
cd atlassian-cli

go build -ldflags "-s -w -X main.version=dev" -o bin/jira ./cmd/jira
go build -ldflags "-s -w -X main.version=dev" -o bin/confluence ./cmd/confluence

go test ./...
```

## Updating

For Claude Code:

```bash
claude skill update atlassian-cli
```

For standalone installations, download the latest release from [GitHub Releases](https://github.com/PigZyj2333/atlassian-cli/releases).

## Contributing

1. Fork the repository
2. Create a branch for your feature
3. Submit a PR

See [CLAUDE.md](CLAUDE.md) for development conventions and architecture details.

## License

MIT License — see [LICENSE](LICENSE) for details.

## Support

- **Issues**: https://github.com/PigZyj2333/atlassian-cli/issues
- **Inspired by**: [mcp-atlassian](https://github.com/sooperset/mcp-atlassian) — the MCP server for Atlassian
