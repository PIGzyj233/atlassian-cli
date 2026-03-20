# Installing Atlassian CLI

Install the `jira` and `confluence` CLI binaries for your platform. These tools give you direct access to the Atlassian REST API from the command line.

## Prerequisites

- Bash shell (Linux, macOS, or Git Bash / WSL on Windows)
- `curl` and `tar` (or `unzip` on Windows)

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/PigZyj2333/atlassian-cli.git ~/.atlassian-cli
   ```

2. **Run the installer:**
   ```bash
   bash ~/.atlassian-cli/skill/scripts/install.sh
   ```

   This auto-detects your OS and architecture, downloads the correct binaries from GitHub Releases, and caches them locally.

3. **Source the resolver** to set up `$JIRA` and `$CONFLUENCE` environment variables:
   ```bash
   source ~/.atlassian-cli/skill/scripts/resolve-bin.sh
   ```

4. **Verify:**
   ```bash
   $JIRA --help
   $CONFLUENCE --help
   ```

## Authentication

Before running any commands, authenticate with your Atlassian instance:

**Cloud** (`*.atlassian.net`):
```bash
$JIRA auth login --host your-company.atlassian.net
# Provide: email + API Token
```

**Server / Data Center** (self-hosted):
```bash
$JIRA auth login --host jira.your-company.com
# Provide: username + password or PAT
```

Credentials persist locally (OS keychain). The `confluence` tool shares the same credential store.

## Usage

Output is JSON by default. Use `--output table` for human-readable output.

```bash
# Search Jira issues
$JIRA search "assignee = currentUser() AND sprint in openSprints()"

# Create an issue
$JIRA issue create --project PROJ --type Bug --summary "Title" --description "Details"

# Search Confluence pages
$CONFLUENCE search "space = ENG AND title ~ 'Design'"

# Get a page
$CONFLUENCE page get 12345678
```

Run `$JIRA --help` or `$CONFLUENCE --help` for full command reference.

## Updating

```bash
cd ~/.atlassian-cli && git pull
bash ~/.atlassian-cli/skill/scripts/install.sh
```

## Uninstalling

```bash
rm -rf ~/.atlassian-cli
```
