# Reference Materials: Advanced Techniques

> Lookup guide for advanced Claude Code techniques in Go/DevOps workflows.
> For basics, see [ONBOARDING.md](../training/ONBOARDING.md).

---

## Table of Contents

1. [Context Management](#context-management)
2. [Advanced Prompting Patterns](#advanced-prompting-patterns)
3. [MCP Server Configuration](#mcp-server-configuration)
4. [Hooks Deep Dive](#hooks-deep-dive)
5. [Headless and CI Integration](#headless-and-ci-integration)
6. [Multi-File Refactoring Strategies](#multi-file-refactoring-strategies)
7. [Infrastructure as Code Patterns](#infrastructure-as-code-patterns)
8. [Security Automation](#security-automation)

---

## Context Management

### When to Use /compact

Use `/compact` when:
- Your conversation exceeds ~50 exchanges
- Claude starts losing track of earlier decisions
- You are switching focus areas within the same session

```
> /compact
```

Claude summarizes the conversation while preserving key decisions, file changes, and architectural context.

### When to Use /clear

Use `/clear` when:
- You are starting a completely new task
- Previous context is no longer relevant
- Claude is confused by earlier conversation

```
> /clear
```

This resets context entirely. Claude re-reads `CLAUDE.md` but forgets everything else.

### Session Continuation

```bash
# Resume the most recent session
claude --resume

# Continue the last conversation
claude --continue
```

Use `--resume` for long-running tasks that span multiple sittings.

---

## Advanced Prompting Patterns

### Constraint-Based Prompts

Be specific about what you want AND what you don't want:

```
> Create a Redis connection pool manager.
> DO: use context.Context for all operations
> DO: implement exponential backoff for reconnects
> DO: use logrus structured logging
> DON'T: use time.Sleep for retries
> DON'T: panic on connection failure
> DON'T: use global variables
```

### Iterative Refinement

Build complex features in layers:

```
Step 1: > Create the data model for Order with basic fields
Step 2: > Add validation methods to the Order type
Step 3: > Create the PostgreSQL store with CRUD operations
Step 4: > Add HTTP handlers that use the store
Step 5: > Wire everything together in main.go
```

### Reference-Based Prompts

Point Claude at existing code as a template:

```
> Create a billing handler following the same patterns as
> internal/handler/users.go — same error handling, same
> middleware chain, same response format.
```

### Diff-Review Prompts

Ask Claude to review before committing:

```
> Show me the git diff of everything you changed.
> For each change, explain why it was necessary.
```

---

## MCP Server Configuration

### Project-Level Configuration

```json
// .claude/settings.json
{
  "mcpServers": {
    "postgres": {
      "command": "npx",
      "args": ["-y", "@anthropic-ai/mcp-server-postgres",
               "postgresql://user:pass@localhost:5432/appdb"]
    },
    "github": {
      "command": "npx",
      "args": ["-y", "@anthropic-ai/mcp-server-github"],
      "env": { "GITHUB_TOKEN": "ghp_..." }
    },
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@anthropic-ai/mcp-server-filesystem", "/opt/docs"]
    }
  }
}
```

### Common MCP Use Cases

| MCP Server | Use Case |
|-----------|----------|
| postgres | Generate Go structs from live schema, write queries |
| github | Query PRs, issues, check runs |
| filesystem | Access documentation outside the repo |

### Usage Example

After configuring the PostgreSQL MCP:

```
> List all tables in the database
> Generate Go struct definitions matching the users table
> Write a migration to add an "updated_at" column to orders
```

---

## Hooks Deep Dive

### Hook Types

| Type | When It Runs | Can Block? |
|------|-------------|------------|
| `PreToolUse` | Before a tool call | Yes |
| `PostToolUse` | After a tool call | No (but can flag errors) |
| `Notification` | When Claude sends a notification | No |

### Hook Environment Variables

Inside a hook command, these variables are available:

| Variable | Contents |
|----------|----------|
| `$CLAUDE_TOOL_NAME` | Name of the tool being called |
| `$CLAUDE_TOOL_INPUT` | JSON string of tool input parameters |
| `$CLAUDE_TOOL_OUTPUT` | (PostToolUse only) JSON string of tool output |

### Production Hook Examples

**Run go vet after every Go file edit:**
```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [{
          "type": "command",
          "command": "if echo $CLAUDE_TOOL_INPUT | grep -q '\\.go\"'; then go vet ./... 2>&1 || echo 'HOOK_ERROR: go vet failed'; fi"
        }]
      }
    ]
  }
}
```

**Run terraform validate after Terraform edits:**
```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [{
          "type": "command",
          "command": "if echo $CLAUDE_TOOL_INPUT | grep -q '\\.tf\"'; then cd terraform && terraform validate 2>&1 || echo 'HOOK_ERROR: terraform validate failed'; fi"
        }]
      }
    ]
  }
}
```

**Block force-push:**
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [{
          "type": "command",
          "command": "if echo $CLAUDE_TOOL_INPUT | grep -q 'push.*--force'; then echo 'HOOK_BLOCK: Force push is not allowed' && exit 1; fi"
        }]
      }
    ]
  }
}
```

---

## Headless and CI Integration

### Single-Prompt Mode

```bash
# Generate code without an interactive session
claude -p "Add input validation to internal/handler/users.go"

# Pipe output to a file
claude -p "Generate a Terraform module for an S3 bucket" > main.tf
```

### Batch Operations

```bash
# Run multiple independent tasks in parallel
claude -p "Add godoc to all public functions in internal/store/" &
claude -p "Generate tests for internal/handler/orders.go" &
claude -p "Create a Dockerfile for cmd/worker/" &
wait
```

### CI Pipeline Integration

```yaml
# .github/workflows/claude-review.yaml
name: Claude Code Review
on: [pull_request]
jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: npm install -g @anthropic-ai/claude-code
      - run: |
          claude -p "Review the changes in this PR for:
          1. Error handling following our CLAUDE.md standards
          2. Missing tests
          3. Security issues
          Provide feedback as GitHub-flavored markdown."
```

---

## Multi-File Refactoring Strategies

### Strategy 1: Plan First, Execute Second

```
[Plan Mode]
> I need to move all database logic from internal/handler/ into
> internal/store/. List every file that needs to change and what
> changes are needed in each.

[Act Mode]
> Execute the refactoring plan. Start with the store package,
> then update the handlers, then update main.go wiring.
```

### Strategy 2: Interface-Driven

```
> Define an interface in internal/store/store.go for all database
> operations. Then refactor each handler to depend on the interface
> instead of the concrete implementation.
```

### Strategy 3: Parallel Refactoring (Boris's Pattern)

Split the work across Claude instances:
```bash
# Terminal 1: Refactor store layer
claude -p "Extract all SQL queries from internal/handler/users.go into internal/store/users.go"

# Terminal 2: Refactor a different handler (no file overlap)
claude -p "Extract all SQL queries from internal/handler/orders.go into internal/store/orders.go"
```

---

## Infrastructure as Code Patterns

### Terraform Module Composition

```
> Create a Terraform module that composes our EKS, RDS, and
> ElastiCache modules into a single "service-stack" module.
> Each service team should be able to deploy their own stack
> by providing just a service name and environment.
```

### Kustomize Overlay Management

```
> Create a new Kustomize overlay for the "canary" environment
> that deploys only 1 replica with debug logging enabled.
> Base it on the staging overlay.
```

### State Management

```
> Generate the Terraform backend configuration for a new
> service called "payments". Use our standard S3 backend
> with DynamoDB locking. The bucket is "company-tf-state"
> and the region is us-east-1.
```

---

## Security Automation

### Secret Scanning

```
> Scan this repository for any hardcoded secrets, API keys,
> passwords, or credentials. Check all file types including
> YAML, Go, Terraform, and Docker files.
```

### Dependency Auditing

```bash
# Check Go dependencies for known vulnerabilities
claude -p "Run govulncheck on this project and explain any findings"
```

### Terraform Security

```
> Review our Terraform code for security issues:
> - Overly permissive IAM policies
> - Public S3 buckets
> - Missing encryption at rest
> - Security groups with 0.0.0.0/0 ingress
```

### Container Security

```
> Review our Dockerfile for security best practices:
> - Running as non-root
> - Minimal base image
> - No secrets in build args
> - Pinned dependency versions
```
