---
name: atlassian-cli
description: |
  Use this skill whenever the user wants to interact with Atlassian Jira or Confluence — including creating, updating, searching, or querying issues, tickets, sprints, boards, epics, projects, pages, spaces, comments, labels, attachments, users, or worklogs. Trigger on any mention of Jira, Confluence, JIRA, JQL, CQL, sprint planning, kanban board, service desk, scrum board, Atlassian, issue tracking, or project management with Jira/Confluence context. Also trigger when user asks to "create a ticket", "search issues", "update story", "move to sprint", "create a page", "publish to confluence", or any similar action.
---

# Atlassian CLI Skill

Two CLI tools for interacting with Atlassian Jira and Confluence.

## Setup: Resolve Binary Path

Binaries are downloaded on first use from GitHub Releases — only the current platform (Linux/macOS/Windows, amd64/arm64).

Before running any command, source the resolver which auto-detects OS/arch and sets `$JIRA` / `$CONFLUENCE`:

```bash
source "${CLAUDE_SKILL_DIR}/scripts/resolve-bin.sh"
```

If binaries are missing, it automatically runs `${CLAUDE_SKILL_DIR}/scripts/install.sh` to download them. Cached locally — only downloads once.

To install a specific version or re-install:
```bash
bash "${CLAUDE_SKILL_DIR}/scripts/install.sh"          # latest
bash "${CLAUDE_SKILL_DIR}/scripts/install.sh" v1.2.3   # specific version
```

Both tools output **JSON by default**, ideal for parsing. Use `--output table` or `--output text` for human-readable output.

---

## Authentication

Credentials persist on the local machine once logged in — no need to re-authenticate each session. Only check auth status if:
- This is the first Atlassian command in the current conversation
- A command returns an authentication error (401/403)

```bash
$JIRA auth status
```

**If already authenticated, proceed directly with the user's request. Do not ask about login.**

### First-time login (only when not authenticated)

Ask the user which deployment type they use:
- **Cloud** (`xxx.atlassian.net`): email + API Token
- **Data Center / Server** (self-hosted): username + password/PAT

**Cloud:**
```bash
$JIRA auth login --host your-company.atlassian.net
# Provide: email + API Token (generate at https://id.atlassian.com/manage-profile/security/api-tokens)
```

**Data Center / Server:**
```bash
$JIRA auth login --host jira.your-company.com
# Provide: username + password or Personal Access Token (PAT)
```

> **Security note**: Credentials are stored only on the local machine (Windows Credential Manager / macOS Keychain / Linux keyring). They are never uploaded to any server. To remove credentials, run `$JIRA auth logout --host <host>`.

## Global Flags (apply to all commands)

```
--host string     Override default Atlassian instance
--output string   json (default) | text | table
--verbose         Debug logging to stderr
--no-color        Disable ANSI colors
```

---

## Jira Commands

### Issues

```bash
$JIRA issue get PROJ-123
$JIRA issue create --project PROJ --type Story --summary "Title" --description "Details"
$JIRA issue update PROJ-123 --fields-json '{"summary":"New title","priority":{"name":"High"}}'
$JIRA issue update PROJ-123 --attachments ./file.pdf --components "Backend,API"
$JIRA issue delete PROJ-123
$JIRA issue transition PROJ-123 --status "In Progress"
$JIRA issue batch-create --file issues.json
```

### Search (JQL)

```bash
$JIRA search "project = PROJ AND status = 'In Progress'"
$JIRA search "assignee = currentUser() AND sprint in openSprints()" \
  --fields summary,status,assignee,priority --limit 20
$JIRA search "status = Open" --projects-filter PROJ1,PROJ2
```

Common JQL patterns:
- `project = PROJ` — filter by project
- `status = "To Do"` — filter by status
- `assignee = currentUser()` — assigned to me
- `sprint in openSprints()` — current sprint
- `issuetype = Story` — issue type filter
- `priority = High` — priority filter
- `created >= -7d` — created in last 7 days

### Sprints & Boards

```bash
$JIRA board list --project PROJ
$JIRA sprint list --board 42
$JIRA sprint create --board 42 --name "Sprint 5" --start-date 2026-03-20 --end-date 2026-04-03
$JIRA sprint issues 123
$JIRA sprint add-issues 123 --issues PROJ-1,PROJ-2,PROJ-3
$JIRA sprint update 123 --state active
```

### Projects

```bash
$JIRA project list
$JIRA project get PROJ
```

### Comments & Worklogs

```bash
$JIRA comment list PROJ-123
$JIRA comment add PROJ-123 --body "This is a comment"
$JIRA worklog list PROJ-123
$JIRA worklog add PROJ-123 --time-spent "2h" --comment "Worked on implementation"
```

### Attachments, Links, Watchers

```bash
$JIRA attachment upload PROJ-123 --file ./report.pdf
$JIRA attachment list PROJ-123
$JIRA link add PROJ-123 --link-type "blocks" --target PROJ-456
$JIRA watcher add PROJ-123 --user john.doe
$JIRA watcher list PROJ-123
```

### Transitions, Fields, Users

```bash
$JIRA transition list PROJ-123
$JIRA field list --project PROJ
$JIRA user search --query "john"
$JIRA user get john.doe@company.com
```

### Epics

```bash
$JIRA epic link PROJ-123 --epic PROJ-100
```

### Service Desk

```bash
$JIRA servicedesk list
$JIRA sla list PROJ-123
```

---

## Confluence Commands

### Pages

```bash
$CONFLUENCE page get 12345678
$CONFLUENCE page create --space MYSPACE --title "My Page" --body "Content here"
$CONFLUENCE page update 12345678 --title "Updated Title" --body "New content"
$CONFLUENCE page delete 12345678
$CONFLUENCE page children 12345678
$CONFLUENCE page history 12345678
$CONFLUENCE page tree --space MYSPACE
$CONFLUENCE page move 12345678 --parent 87654321
$CONFLUENCE page diff 12345678 --from 1 --to 2
```

### Search (CQL)

```bash
$CONFLUENCE search "space = MYSPACE AND type = page"
$CONFLUENCE search "title ~ 'Meeting Notes'" --limit 10
```

Common CQL patterns:
- `space = MYSPACE` — filter by space
- `type = page` — only pages
- `title ~ "keyword"` — title contains keyword
- `ancestor = 12345678` — children of page
- `creator = currentUser()` — created by me
- `lastModified >= now("-7d")` — recently modified

### Comments, Labels, Attachments

```bash
$CONFLUENCE comment list 12345678
$CONFLUENCE comment add 12345678 --body "Great page!"
$CONFLUENCE label list 12345678
$CONFLUENCE label add 12345678 --label meeting-notes
$CONFLUENCE attachment upload 12345678 --file ./diagram.png
$CONFLUENCE attachment list 12345678
```

### Users & Analytics

```bash
$CONFLUENCE user search --query "alice"
$CONFLUENCE analytics page 12345678    # Cloud only
```

---

## Tips

1. **JSON output** (`--output json`) is best when you need to parse results and extract IDs.
2. **Chain commands**: Get an issue to find its sprint ID, then list sprint issues.
3. **JQL is powerful**: Combine conditions for complex Jira searches.
4. **Page IDs**: Confluence operations need numeric page IDs, not titles. Use `search` to find them.
5. **Multiple instances**: Use `--host` to switch instances, or `$JIRA auth switch`.

## Reference

```bash
$JIRA <command> --help
$JIRA <command> <subcommand> --help
```
