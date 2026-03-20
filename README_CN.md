# Atlassian CLI

[English](README.md) | 简体中文

Atlassian CLI 让你的编程 Agent 和你自己都能直接从命令行访问 Jira 和 Confluence。两个独立的二进制文件 `jira` 和 `confluence`，支持 **Cloud**、**Server** 和 **Data Center** 实例。

## 工作原理

你对编程 Agent 说"给登录 bug 创建一个 Jira 工单"或"在 Confluence 里找一下架构文档"——它就直接做了。不需要打开浏览器，不需要复制粘贴，不需要切换上下文。

底层原理是 Agent 调用 `jira` 和 `confluence` CLI 二进制文件与 Atlassian REST API 通信。默认输出 JSON，非常适合 Agent 解析和处理。二进制文件会自动检测你的实例是 Cloud 还是 Server/Data Center，通过操作系统密钥链管理认证信息，并在遇到限流时自动重试。

因为它是一个 Skill，你的 Agent 会自动加载它。你不需要教它任何东西——它已经知道所有命令。

同时它也是一个普通的 CLI 工具，你可以在终端、脚本、CI/CD 流水线中直接使用。`--output table` 给人看，`--output json` 给机器用。

## 致谢

本项目深受 **[mcp-atlassian](https://github.com/sooperset/mcp-atlassian)** 项目的启发，由 [@sooperset](https://github.com/sooperset) 及其贡献者构建。他们在构建 Atlassian MCP 服务器方面的开创性工作，为本项目提供了重要的灵感和参考。API 交互模式、对 Cloud/Server/Data Center 多部署方式的支持方案、以及将 AI 助手与 Atlassian 工具桥接的整体愿景——我们的设计思路很大程度上受益于对他们优秀实现的研究。如果你需要一个原生 MCP 方案来集成 Atlassian，我们强烈推荐他们的项目。

## 安装

**注意：** 不同平台的安装方式不同。Claude Code 有内置的 Skill 系统，Codex 和其他 Agent 需要手动配置。

### Claude Code（应用市场）

直接安装 Skill：

```bash
claude skill add --from https://github.com/PigZyj2333/atlassian-cli --subpath skill
```

### Claude Code（手动配置）

添加到项目的 `.claude/settings.json`：

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

告诉 Codex：

```
Fetch and follow instructions from https://raw.githubusercontent.com/PigZyj2333/atlassian-cli/refs/heads/master/.codex/install.md
```

### 其他 AI Agent

任何能执行 Shell 命令的 Agent 都可以使用。告诉它：

```
Fetch and follow instructions from https://raw.githubusercontent.com/PigZyj2333/atlassian-cli/refs/heads/master/.codex/install.md
```

### 独立安装（不用 Agent）

如果你只想自己用这个 CLI 工具：

**从 GitHub Releases 下载：**
```bash
# Linux / macOS
curl -fsSL https://github.com/PigZyj2333/atlassian-cli/releases/latest/download/atlassian-cli_*_$(uname -s | tr A-Z a-z)_$(uname -m).tar.gz | tar xz
sudo mv jira confluence /usr/local/bin/
```

**Go install：**
```bash
go install github.com/PigZyj2333/atlassian-cli/cmd/jira@latest
go install github.com/PigZyj2333/atlassian-cli/cmd/confluence@latest
```

### 验证安装

启动新会话，试试说"搜索我的 Jira 工单"或"在 Confluence 里找设计文档"。Agent 会自动调用 Skill 并执行正确的命令。

独立使用时：

```bash
jira --help
confluence --help
```

## 基本工作流程

1. **认证** — 执行 `jira auth login --host your-company.atlassian.net`，只需一次。凭证持久化存储在操作系统密钥链中。Cloud 使用邮箱 + API Token；Server/DC 使用用户名 + 密码或 PAT。

2. **搜索** — 用 JQL 搜索 Issue（`jira search "assignee = currentUser()"`），用 CQL 搜索页面（`confluence search "space = ENG AND title ~ 'API'"`）。

3. **创建** — 创建 Issue（`jira issue create --project PROJ --type Bug --summary "..."`）或页面（`confluence page create --space ENG --title "..." --body "..."`）。

4. **更新** — 变更 Issue 状态（`jira issue transition PROJ-123 --status "In Progress"`），更新字段、添加评论、上传附件。

5. **管理** — Sprint、看板、Epic、工时记录、关注者、链接、标签等。

**Agent 会自动完成以上所有操作。** 你用自然语言描述需求，它来确定命令。

## 功能一览

### Jira 命令

- **auth** — 登录、登出、查看状态、切换实例
- **issue** — 获取、创建、更新、删除、状态变更、批量创建
- **search** — 基于 JQL 的搜索，支持字段选择
- **board** — 按项目列出敏捷看板
- **sprint** — 列出、创建、更新 Sprint；管理 Sprint 待办
- **project** — 列出和查看项目
- **comment** — 列出和添加 Issue 评论
- **worklog** — 工时记录（列出、添加）
- **attachment** — 上传和列出附件
- **link** — 关联 Issue（阻塞、关联等）
- **watcher** — 添加和列出关注者
- **transition** — 列出可用的状态变更
- **field** — 列出自定义和系统字段
- **user** — 搜索和查看用户
- **epic** — 关联 Issue 到 Epic
- **servicedesk** — 列出服务台队列
- **sla** — 服务台工单的 SLA 信息
- **changelog** — Issue 变更历史
- **dev** — 开发信息（分支、提交、PR）
- **form** — Proforma 表单数据

### Confluence 命令

- **auth** — 登录、登出、查看状态
- **page** — 获取、创建、更新、删除、子页面、页面树、移动、对比、历史
- **search** — 基于 CQL 的搜索
- **comment** — 列出和添加页面评论
- **label** — 列出和添加页面标签
- **attachment** — 上传和列出附件
- **user** — 搜索用户
- **analytics** — 页面浏览分析（仅 Cloud）

### 全局标志

```
--host string     目标 Atlassian 实例（覆盖默认值）
--output string   输出格式：json（默认）| text | table
--verbose         输出调试日志到 stderr
--no-color        禁用彩色输出
```

## 架构

两个独立二进制共享同一基础设施层：

```
cmd/{jira,confluence}/main.go       入口点
pkg/cmdutil/factory.go              依赖注入（配置 → 客户端）
pkg/{jira,confluence}/<feature>/    每个功能一个包
internal/api/                       HTTP 客户端（重试、分页、错误处理）
internal/config/                    YAML 配置、认证（BasicAuth/PAT）、主机解析
internal/output/                    JSON / Table / Text 格式化器
internal/convert/                   ADF ↔ storage 格式转换（Confluence）
```

核心设计决策：
- **工厂模式** — 每个命令通过 `*cmdutil.Factory` 获取配置和认证客户端
- **Cloud vs Server 自动检测** — `*.atlassian.net` 使用 REST API v3，其他使用 v2
- **指数退避重试** — 遇到 429 和 5xx 自动重试 3 次
- **静态二进制** — `CGO_ENABLED=0`，零运行时依赖

## 从源码构建

```bash
git clone https://github.com/PigZyj2333/atlassian-cli.git
cd atlassian-cli

go build -ldflags "-s -w -X main.version=dev" -o bin/jira ./cmd/jira
go build -ldflags "-s -w -X main.version=dev" -o bin/confluence ./cmd/confluence

go test ./...
```

## 更新

Claude Code：

```bash
claude skill update atlassian-cli
```

独立安装请从 [GitHub Releases](https://github.com/PigZyj2333/atlassian-cli/releases) 下载最新版本。

## 贡献

1. Fork 本仓库
2. 创建功能分支
3. 提交 PR

开发规范和架构细节请参阅 [CLAUDE.md](CLAUDE.md)。

## 许可证

MIT License — 详见 [LICENSE](LICENSE)。

## 支持

- **Issues**：https://github.com/PigZyj2333/atlassian-cli/issues
- **灵感来源**：[mcp-atlassian](https://github.com/sooperset/mcp-atlassian) — Atlassian 的 MCP 服务器
