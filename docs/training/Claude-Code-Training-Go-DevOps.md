# Claude Code Training Program: Go Backend / DevOps Engineers

> **Duration:** 10 weeks (2 hours per week, 20 hours total)
> **Audience:** Go Backend and DevOps Engineers
> **Prerequisites:** Go 1.21+, Docker, Terraform basics, Kubernetes basics, GitHub account
> **Format:** 20 hands-on exercises with progressive difficulty

---

## Table of Contents

1. [Program Overview](#program-overview)
2. [Phase 1 — Foundation (Weeks 1–3)](#phase-1--foundation-weeks-13)
3. [Phase 2 — Intermediate (Weeks 4–6)](#phase-2--intermediate-weeks-46)
4. [Phase 3 — Advanced (Weeks 7–9)](#phase-3--advanced-weeks-79)
5. [Phase 4 — Expert (Week 10)](#phase-4--expert-week-10)
6. [Appendix A — Quick Reference Card](#appendix-a--quick-reference-card)
7. [Appendix B — Quiz and Self-Assessment](#appendix-b--quiz-and-self-assessment)

---

## Program Overview

### Learning Path

```
Week  1-3   Foundation    ██████░░░░░░░░░░░░░░  Exercises  1-6
Week  4-6   Intermediate  ████████████░░░░░░░░  Exercises  7-12
Week  7-9   Advanced      ██████████████████░░  Exercises 13-18
Week 10     Expert        ████████████████████  Exercises 19-20
```

### Environment Setup

```bash
# Install Claude Code
npm install -g @anthropic-ai/claude-code

# Verify installation
claude --version

# Authenticate
claude auth login

# Start a session in your project
cd your-go-project
claude
```

---

## Phase 1 — Foundation (Weeks 1–3)

### Exercise 1: First Contact — Navigating a Go Codebase

**Goal:** Learn basic Claude Code interaction and codebase exploration.

**Duration:** 30 minutes

**Setup:**
```bash
mkdir -p ~/training/ex01 && cd ~/training/ex01
go mod init github.com/training/ex01
```

Create a starter file:
```go
// main.go
package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/api/users", usersHandler)

    fmt.Printf("Server starting on :%s\n", port)
    http.ListenAndServe(":"+port, nil)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"users":[]}`))
}
```

**Tasks:**
1. Open Claude Code and ask it to explain the project structure.
2. Ask Claude to identify potential issues (hint: `ListenAndServe` error is ignored).
3. Have Claude fix the issues following Go best practices.

**Key Commands to Learn:**
```
claude                           # Start interactive session
/help                            # Show available commands
/clear                           # Clear conversation context
```

**Expected Outcome:**
Claude should identify the ignored error from `ListenAndServe` and the missing `Content-Type` headers, then fix both:
```go
func main() {
    // ...
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}
```

---

### Exercise 2: CLAUDE.md — Teaching Claude Your Standards

**Goal:** Create and iterate on a CLAUDE.md that encodes your team's Go and DevOps conventions.

**Duration:** 30 minutes

**Background:**
`CLAUDE.md` files are placed at the root of your repository (or in subdirectories for scoped context). Claude Code reads them automatically at session start. They define project conventions, coding standards, and architectural decisions that Claude should follow.

**Tasks:**

1. Create a `CLAUDE.md` in the project root:

```markdown
# Project Context

## Architecture Decisions

- We use structured logging (logrus for Go)
- All pipeline configs are in YAML under /configs
- We prefer explicit error handling over panics
- All errors must be wrapped with context using fmt.Errorf with %w

## Coding Standards

- Go: Follow Effective Go guidelines
- Use errgroup for concurrent operations
- Max function length: 50 lines
- All public functions must have godoc comments
- Use contexts for timeout handling and cancellation

## Infrastructure Patterns

- All Terraform code uses remote state in S3
- K8s manifests use Kustomize, not Helm
- CI/CD runs on GitHub Actions
- Docker images must be multi-stage builds

## Common Gotchas

- Always validate YAML before applying to avoid pipeline breaks
- Use context.Context in all functions that do I/O
- Never use time.Sleep() for retries - use exponential backoff
- Always add database indexes when adding foreign keys
```

2. Start a new Claude session and ask it to add a new HTTP endpoint. Observe how it follows the conventions from `CLAUDE.md` (structured logging, error wrapping, godoc comments).

3. Deliberately ask for something that contradicts `CLAUDE.md` (for example: "use `panic` for error handling"). Claude should push back or note the conflict.

**Verification:**
- Ask Claude: "What are the coding standards for this project?"
- Claude should reference your `CLAUDE.md` conventions.

**Tip:** You can also place `CLAUDE.md` files in subdirectories. A `terraform/CLAUDE.md` only applies when Claude works on files in that subtree.

---

### Exercise 3: Plan Mode — Designing a Go Microservice

**Goal:** Use Plan Mode to architect before writing code.

**Duration:** 30 minutes

**Background:**
Plan Mode (`Shift+Tab` to toggle) tells Claude to research and plan without making changes. Use it when you need to think through architecture before committing to an implementation.

**Scenario:** You need a Go microservice that:
- Reads messages from a Redis queue
- Validates and transforms the data
- Writes results to PostgreSQL
- Exposes health and metrics endpoints

**Tasks:**

1. Start Claude in Plan Mode:
   - Press `Shift+Tab` to switch to Plan Mode (you will see "Plan" indicator)
   - Describe the microservice requirements above
   - Let Claude produce a design document

2. Review the plan. Challenge assumptions:
   - "What if Redis is temporarily unavailable?"
   - "How do we handle schema migrations?"
   - "What about graceful shutdown?"

3. Switch back to Act Mode (`Shift+Tab` again) and ask Claude to implement the plan.

**Expected Plan Structure:**
```
1. Project layout
   cmd/worker/main.go
   internal/queue/redis.go
   internal/store/postgres.go
   internal/transform/pipeline.go
   internal/health/server.go

2. Dependencies
   github.com/redis/go-redis/v9
   github.com/jackc/pgx/v5
   github.com/sirupsen/logrus

3. Error handling strategy
   - Exponential backoff for Redis reconnects
   - Circuit breaker for PostgreSQL
   - Dead-letter queue for poison messages

4. Graceful shutdown
   - Signal handling (SIGTERM/SIGINT)
   - Drain in-flight messages
   - Close connections in order
```

---

### Exercise 4: Git Workflows — Branching and PRs with Claude

**Goal:** Use Claude Code for branch management, commits, and pull requests.

**Duration:** 30 minutes

**Tasks:**

1. Create a feature branch and implement a change:
```
> Create a new branch called feature/add-metrics and add a Prometheus
> metrics endpoint to our service. Include request count, latency
> histogram, and error rate counters.
```

2. Review the diff before committing:
```
> Show me the git diff and explain each change
```

3. Create a commit with a conventional message:
```
> Commit these changes with a descriptive message following
> conventional commits format
```

4. Create a pull request:
```
> Create a PR for this branch. Include a summary of changes,
> testing instructions, and any migration notes.
```

**Key Git Commands in Claude:**
```
/commit           # Stage and commit with a generated message
```

You can also ask Claude directly:
```
> Create a PR with title "Add Prometheus metrics" and a detailed body
```

**Best Practice:** Always review the diff Claude produces before approving the commit. Claude shows you what it plans to do — read it.

---

### Exercise 5: Error Handling Patterns in Go

**Goal:** Build a robust error handling package guided by Claude.

**Duration:** 30 minutes

**Setup:**
```bash
mkdir -p ~/training/ex05 && cd ~/training/ex05
go mod init github.com/training/ex05
```

**Tasks:**

1. Ask Claude to create an `errors` package with:
   - Custom error types for domain errors (NotFound, Validation, Conflict)
   - Error wrapping helpers
   - Sentinel errors for common cases
   - HTTP status code mapping

```
> Create an internal/errors package with custom error types for a REST API.
> Include NotFoundError, ValidationError, and ConflictError. Each should
> implement the error interface, support wrapping with %w, and map to the
> correct HTTP status code.
```

2. Expected output:

```go
// internal/errors/errors.go
package errors

import "fmt"

// NotFoundError indicates a requested resource does not exist.
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with id %s not found", e.Resource, e.ID)
}

// ValidationError indicates invalid input.
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

// StatusCode returns the HTTP status code for a given error.
func StatusCode(err error) int {
    var notFound *NotFoundError
    var validation *ValidationError
    switch {
    case errors.As(err, &notFound):
        return 404
    case errors.As(err, &validation):
        return 400
    default:
        return 500
    }
}
```

3. Ask Claude to write table-driven tests for the error package.

---

### Exercise 6: Multi-File Refactoring

**Goal:** Refactor a monolithic Go file into a clean package structure.

**Duration:** 30 minutes

**Setup:** Create a single monolithic file:

```go
// main.go — intentionally messy, 200+ lines
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    _ "github.com/lib/pq"
)

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    http.HandleFunc("/users", handleUsers)
    http.HandleFunc("/users/create", handleCreateUser)
    http.HandleFunc("/health", handleHealth)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, email, created_at FROM users")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer rows.Close()
    var users []User
    for rows.Next() {
        var u User
        rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
        users = append(users, u)
    }
    json.NewEncoder(w).Encode(users)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
    var u User
    json.NewDecoder(r.Body).Decode(&u)
    _, err := db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", u.Name, u.Email)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.WriteHeader(201)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"status":"ok"}`))
}
```

**Tasks:**

1. Ask Claude to refactor into a clean structure:
```
> Refactor this monolithic main.go into a proper Go project structure
> with separate packages for handlers, models, and database. Follow
> our CLAUDE.md standards.
```

2. The result should produce:
```
cmd/server/main.go
internal/handler/user.go
internal/handler/health.go
internal/model/user.go
internal/store/postgres.go
```

3. Ask Claude to add missing error handling and context propagation throughout.

---

## Phase 2 — Intermediate (Weeks 4–6)

### Exercise 7: Slash Commands and Workflow Efficiency

**Goal:** Master built-in commands and workflow shortcuts.

**Duration:** 30 minutes

**Key Commands:**

| Command | Purpose |
|---------|---------|
| `/help` | Show all available commands |
| `/clear` | Clear conversation context |
| `/compact` | Summarize and compress conversation |
| `/commit` | Stage and commit changes |
| `/review` | Review code changes |
| `/init` | Initialize CLAUDE.md for a project |
| `/memory` | Edit CLAUDE.md memory files |
| `/cost` | Show token usage and costs |

**Tasks:**

1. Start a new session. Use `/compact` after a long conversation to see how Claude summarizes context.

2. Use `/commit` after making changes and observe the generated commit message.

3. Practice context management:
```
> [make several changes across files]
> /compact
> [continue working — Claude retains the summary]
```

4. Use `/review` to review recent changes:
```
> /review
```

**Tip:** `/compact` is essential for long sessions. It compresses your conversation while preserving important context, letting you work longer without losing track.

---

### Exercise 8: Custom Slash Commands

**Goal:** Create project-specific slash commands that automate repetitive workflows.

**Duration:** 30 minutes

**Background:**
Custom slash commands are Markdown files stored in `.claude/commands/` (project-scoped) or `~/.claude/commands/` (global). They appear when you type `/` and can include the `$ARGUMENTS` placeholder for user input.

**Tasks:**

1. Create a project-scoped command for generating Go services:

```bash
mkdir -p .claude/commands
```

```markdown
<!-- .claude/commands/new-service.md -->
Create a new Go microservice with the following specifications:

Service name: $ARGUMENTS

Follow these requirements:
1. Create cmd/$ARGUMENTS/main.go with graceful shutdown
2. Create internal/$ARGUMENTS/handler.go with health endpoint
3. Create internal/$ARGUMENTS/service.go with business logic interface
4. Create Dockerfile with multi-stage build
5. Create configs/$ARGUMENTS.yaml with default configuration
6. Follow all conventions in CLAUDE.md

Use structured logging with logrus, context.Context for all I/O,
and proper error wrapping with fmt.Errorf %w.
```

2. Create a command for Terraform modules:

```markdown
<!-- .claude/commands/tf-module.md -->
Create a new Terraform module for: $ARGUMENTS

Requirements:
- Place in terraform/modules/$ARGUMENTS/
- Include main.tf, variables.tf, outputs.tf
- Use remote state data sources where needed
- Add input validation blocks
- Include a README.md with usage example
- Follow our infrastructure patterns from CLAUDE.md
```

3. Test your commands:
```
/project:new-service payment-processor
/project:tf-module rds-cluster
```

---

### Exercise 9: MCP Integration — Connecting External Tools

**Goal:** Configure Model Context Protocol (MCP) servers to extend Claude's capabilities.

**Duration:** 30 minutes

**Background:**
MCP (Model Context Protocol) lets you connect Claude to external data sources and tools. You configure MCP servers in `.claude/settings.json` or `~/.claude/settings.json`. Claude can then call these tools during your session.

**Tasks:**

1. Configure a filesystem MCP server for documentation access:

```json
// .claude/settings.json
{
  "mcpServers": {
    "docs": {
      "command": "npx",
      "args": [
        "-y",
        "@anthropic-ai/mcp-server-filesystem",
        "/path/to/your/docs"
      ]
    }
  }
}
```

2. Configure a PostgreSQL MCP server for database access:

```json
{
  "mcpServers": {
    "postgres": {
      "command": "npx",
      "args": [
        "-y",
        "@anthropic-ai/mcp-server-postgres",
        "postgresql://user:pass@localhost:5432/mydb"
      ]
    }
  }
}
```

3. Use Claude with the connected tools:
```
> List the tables in our database and show me the schema for the users table
> Generate Go struct definitions that match the database schema
```

**Practical DevOps example — GitHub MCP:**
```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@anthropic-ai/mcp-server-github"],
      "env": {
        "GITHUB_TOKEN": "ghp_your_token_here"
      }
    }
  }
}
```

Then ask:
```
> Show me all open PRs that modify Terraform files
> List issues labeled "bug" in our infrastructure repo
```

---

### Exercise 10: Extended Thinking for Complex Architecture

**Goal:** Use Extended Thinking mode for deep architectural analysis.

**Duration:** 30 minutes

**Background:**
Extended Thinking gives Claude more reasoning budget for complex problems. It is enabled by default in Claude Code — Claude uses it automatically when facing difficult tasks. You can see the thinking process in the output, which helps you understand Claude's reasoning.

**Scenario:** You have a monolithic Go application that needs to be decomposed into microservices.

**Tasks:**

1. Provide the monolith context and ask for a decomposition plan:

```
> I have a monolithic Go application with the following packages:
> - internal/auth (JWT, OAuth, sessions)
> - internal/billing (Stripe integration, invoices)
> - internal/notifications (email, SMS, push)
> - internal/inventory (stock tracking, warehouse)
> - internal/orders (cart, checkout, fulfillment)
>
> Design a microservice decomposition plan. Consider:
> - Service boundaries and data ownership
> - Inter-service communication (sync vs async)
> - Shared libraries vs duplication
> - Migration strategy (strangler fig pattern)
> - Data consistency across services
```

2. Observe the extended thinking — Claude will reason through trade-offs before presenting the plan.

3. Challenge the plan:
```
> What if the orders service needs real-time inventory data?
> How do we handle distributed transactions for checkout?
> What's the testing strategy for inter-service contracts?
```

4. Ask Claude to generate the Terraform infrastructure for the proposed architecture:
```
> Generate Terraform code for deploying these microservices on EKS,
> including service mesh configuration, RDS instances for each
> service, and ElastiCache for shared caching.
```

---

### Exercise 11: Docker Optimization with Claude

**Goal:** Build production-grade, optimized Docker images for Go services.

**Duration:** 30 minutes

**Tasks:**

1. Start with a naive Dockerfile and ask Claude to optimize it:

```dockerfile
# Naive Dockerfile
FROM golang:1.24
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server
EXPOSE 8080
CMD ["./server"]
```

2. Ask Claude to optimize:
```
> Optimize this Dockerfile for production. Minimize image size, use
> multi-stage builds, and follow security best practices. Our
> CLAUDE.md requires multi-stage builds.
```

3. Expected optimized output:

```dockerfile
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /build/server \
    ./cmd/server

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
```

4. Ask Claude to add health check and non-root user:
```
> Add a HEALTHCHECK instruction and run as a non-root user.
> Also add labels for OCI image spec compliance.
```

5. Ask Claude to create a `docker-compose.yaml` for local development with Redis, PostgreSQL, and the Go service:
```
> Create a docker-compose.yaml for local development that includes
> our Go service, PostgreSQL, and Redis. Include volume mounts for
> live reloading with air.
```

---

### Exercise 12: GitHub Actions CI/CD Pipeline

**Goal:** Build a complete CI/CD pipeline with Claude's help.

**Duration:** 30 minutes

**Tasks:**

1. Ask Claude to create a GitHub Actions workflow:

```
> Create a GitHub Actions CI/CD pipeline for our Go project with:
> - Lint (golangci-lint)
> - Test with coverage
> - Build Docker image
> - Push to GHCR
> - Deploy to staging on PR merge
> - Deploy to production on tag
```

2. Expected workflow structure:

```yaml
# .github/workflows/ci.yaml
name: CI/CD

on:
  push:
    branches: [main]
    tags: ['v*']
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: testdb
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7
        ports:
          - 6379:6379
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: go test -race -coverprofile=coverage.out ./...
      - uses: codecov/codecov-action@v4
        with:
          file: coverage.out

  build:
    needs: [lint, test]
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v6
        with:
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:${{ github.sha }}
            ghcr.io/${{ github.repository }}:latest

  deploy-staging:
    needs: [build]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: |
          echo "Deploying to staging..."
          # kubectl set image deployment/app \
          #   app=ghcr.io/${{ github.repository }}:${{ github.sha }}

  deploy-production:
    needs: [build]
    if: startsWith(github.ref, 'refs/tags/v')
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v4
      - run: |
          echo "Deploying to production..."
          # kubectl set image deployment/app \
          #   app=ghcr.io/${{ github.repository }}:${{ github.sha }}
```

3. Ask Claude to add a reusable workflow for deploying to multiple environments.

---

## Phase 3 — Advanced (Weeks 7–9)

### Exercise 13: Subagent Architecture — Parallel Code Generation

**Goal:** Understand how Claude uses subagents for concurrent tasks.

**Duration:** 30 minutes

**Background:**
When you give Claude a complex task that spans multiple independent files or concerns, Claude can internally use subagents — parallel workers that each handle a piece of the task. This happens transparently. You can encourage this by structuring your requests to highlight independent work streams.

**Tasks:**

1. Give Claude a task with clearly independent components:

```
> Create a complete observability stack for our Go service:
>
> 1. internal/metrics/prometheus.go — Prometheus metrics collectors
>    for HTTP requests, database queries, and Redis operations
>
> 2. internal/tracing/opentelemetry.go — OpenTelemetry tracing setup
>    with span creation helpers for each layer
>
> 3. internal/logging/logrus.go — Structured logging middleware that
>    correlates logs with trace IDs
>
> 4. internal/health/checker.go — Deep health checks for all
>    dependencies (DB, Redis, external APIs)
>
> These are independent packages. Generate all four simultaneously.
```

2. Observe how Claude handles the parallel workload. It may create files concurrently.

3. Ask Claude to wire everything together:
```
> Now create cmd/server/main.go that initializes all four observability
> packages and attaches them to our HTTP server as middleware.
```

---

### Exercise 14: Hooks — Automating Quality Gates

**Goal:** Configure Claude Code hooks for automated validation.

**Duration:** 30 minutes

**Background:**
Hooks are shell commands that run automatically in response to Claude Code events. They can run before or after tool calls, and can even modify Claude's behavior. Configure them in `.claude/settings.json`.

**Hook types:**
- `PreToolUse` — Runs before a tool call. Can block the action.
- `PostToolUse` — Runs after a tool call. Can validate the result.
- `Notification` — Triggered when Claude wants to notify you.

**Tasks:**

1. Create a hook that runs `go vet` after every file edit:

```json
// .claude/settings.json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "if echo $CLAUDE_TOOL_INPUT | grep -q '\\.go\"'; then cd /path/to/project && go vet ./... 2>&1 || echo 'HOOK_ERROR: go vet failed'; fi"
          }
        ]
      }
    ]
  }
}
```

2. Create a hook that prevents committing without tests:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "if echo $CLAUDE_TOOL_INPUT | grep -q 'git commit'; then STAGED=$(git diff --cached --name-only | grep '\\.go$' | grep -v '_test.go$'); for f in $STAGED; do TEST=${f%.go}_test.go; if [ ! -f \"$TEST\" ]; then echo \"HOOK_BLOCK: Missing test file for $f\" && exit 1; fi; done; fi"
          }
        ]
      }
    ]
  }
}
```

3. Create a hook that runs `terraform validate` after Terraform file changes:

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "if echo $CLAUDE_TOOL_INPUT | grep -q '\\.tf\"'; then cd terraform && terraform validate 2>&1 || echo 'HOOK_ERROR: terraform validate failed'; fi"
          }
        ]
      }
    ]
  }
}
```

4. Test the hooks by making changes that should trigger them.

---

### Exercise 15: Terraform Infrastructure as Code

**Goal:** Use Claude to generate and manage Terraform infrastructure.

**Duration:** 30 minutes

**Tasks:**

1. Ask Claude to create a complete EKS cluster module:

```
> Create a Terraform module for an EKS cluster with:
> - Managed node groups (spot and on-demand)
> - VPC with public and private subnets
> - IAM roles with least-privilege policies
> - Cluster autoscaler configuration
> - Remote state in S3 with DynamoDB locking
>
> Follow our CLAUDE.md infrastructure patterns.
```

2. Expected structure:
```
terraform/
  modules/
    eks/
      main.tf
      variables.tf
      outputs.tf
      iam.tf
      node_groups.tf
  environments/
    staging/
      main.tf
      terraform.tfvars
    production/
      main.tf
      terraform.tfvars
  backend.tf
```

3. Ask Claude to add an RDS module:

```
> Create a Terraform module for RDS PostgreSQL with:
> - Multi-AZ for production, single-AZ for staging
> - Automated backups with 7-day retention
> - Parameter group optimized for Go connection pooling
> - Security group allowing access only from EKS nodes
> - Outputs that generate Kubernetes secrets
```

4. Example Terraform code Claude should produce:

```hcl
# terraform/modules/rds/main.tf
resource "aws_db_instance" "main" {
  identifier     = "${var.project}-${var.environment}"
  engine         = "postgres"
  engine_version = var.engine_version
  instance_class = var.instance_class

  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.max_allocated_storage

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  multi_az            = var.environment == "production"
  deletion_protection = var.environment == "production"

  backup_retention_period = 7
  backup_window           = "03:00-04:00"
  maintenance_window      = "sun:04:00-sun:05:00"

  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name
  parameter_group_name   = aws_db_parameter_group.main.name

  tags = var.tags
}

resource "aws_db_parameter_group" "main" {
  family = "postgres16"
  name   = "${var.project}-${var.environment}"

  parameter {
    name  = "max_connections"
    value = "200"
  }

  parameter {
    name  = "idle_in_transaction_session_timeout"
    value = "30000"
  }
}
```

---

### Exercise 16: Kubernetes Manifests with Kustomize

**Goal:** Generate Kubernetes deployment manifests using Kustomize (not Helm, per our CLAUDE.md).

**Duration:** 30 minutes

**Tasks:**

1. Ask Claude to create a Kustomize-based deployment:

```
> Create Kubernetes manifests for our Go microservice using Kustomize.
> Include base manifests and overlays for staging and production.
> Requirements:
> - Deployment with rolling update strategy
> - HPA for autoscaling
> - PodDisruptionBudget
> - ConfigMap and Secret references
> - Resource limits and requests
> - Liveness and readiness probes
> - ServiceAccount with IRSA annotation for AWS
```

2. Expected structure:
```
k8s/
  base/
    deployment.yaml
    service.yaml
    hpa.yaml
    pdb.yaml
    serviceaccount.yaml
    kustomization.yaml
  overlays/
    staging/
      kustomization.yaml
      patches/
        replicas.yaml
        resources.yaml
    production/
      kustomization.yaml
      patches/
        replicas.yaml
        resources.yaml
        hpa.yaml
```

3. Example base deployment Claude should produce:

```yaml
# k8s/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: worker
  labels:
    app: worker
spec:
  selector:
    matchLabels:
      app: worker
  template:
    metadata:
      labels:
        app: worker
    spec:
      serviceAccountName: worker
      terminationGracePeriodSeconds: 30
      containers:
        - name: worker
          image: ghcr.io/org/worker:latest
          ports:
            - containerPort: 8080
              name: http
            - containerPort: 9090
              name: metrics
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: worker-secrets
                  key: database-url
            - name: REDIS_ADDR
              valueFrom:
                configMapKeyRef:
                  name: worker-config
                  key: redis-addr
            - name: LOG_LEVEL
              valueFrom:
                configMapKeyRef:
                  name: worker-config
                  key: log-level
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
          livenessProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 5
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
```

4. Ask Claude to generate the Kustomize overlays that patch resource limits for production.

---

### Exercise 17: Security Automation

**Goal:** Use Claude to implement security scanning and hardening.

**Duration:** 30 minutes

**Tasks:**

1. Ask Claude to create a security scanning pipeline:

```
> Create a GitHub Actions workflow for security scanning:
> - gosec for Go static analysis
> - Trivy for container image scanning
> - Checkov for Terraform scanning
> - SBOM generation with syft
> - Run on every PR and weekly on main
```

2. Expected workflow:

```yaml
# .github/workflows/security.yaml
name: Security Scan

on:
  pull_request:
  schedule:
    - cron: '0 6 * * 1'  # Weekly Monday 6am

jobs:
  gosec:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: securego/gosec@master
        with:
          args: ./...

  trivy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          severity: 'CRITICAL,HIGH'
          exit-code: '1'

  trivy-image:
    needs: [gosec]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/build-push-action@v6
        with:
          push: false
          tags: scan-target:latest
          load: true
      - uses: aquasecurity/trivy-action@master
        with:
          image-ref: 'scan-target:latest'
          severity: 'CRITICAL,HIGH'
          exit-code: '1'

  checkov:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: bridgecrewio/checkov-action@master
        with:
          directory: terraform/
          framework: terraform

  sbom:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: anchore/sbom-action@v0
        with:
          artifact-name: sbom.spdx.json
```

3. Ask Claude to add security middleware to the Go service:

```
> Add security middleware to our HTTP server:
> - Rate limiting per IP
> - Request ID injection
> - Security headers (HSTS, CSP, X-Frame-Options)
> - Request body size limiting
> - Panic recovery
```

4. Example Go security middleware:

```go
// internal/middleware/security.go
package middleware

import (
    "context"
    "fmt"
    "net/http"

    "github.com/google/uuid"
    "golang.org/x/time/rate"
)

// RequestID injects a unique request ID into the context and response headers.
func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := uuid.New().String()
        ctx := context.WithValue(r.Context(), requestIDKey, id)
        w.Header().Set("X-Request-ID", id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// SecureHeaders adds security headers to every response.
func SecureHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("Strict-Transport-Security",
            "max-age=63072000; includeSubDomains")
        next.ServeHTTP(w, r)
    })
}

// RateLimiter limits requests per client IP.
func RateLimiter(rps float64, burst int) func(http.Handler) http.Handler {
    limiter := rate.NewLimiter(rate.Limit(rps), burst)
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "rate limit exceeded",
                    http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// Recovery catches panics and returns a 500 response.
func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                http.Error(w,
                    fmt.Sprintf("internal error: %v", err),
                    http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

---

### Exercise 18: Testing Strategies with Claude

**Goal:** Generate comprehensive test suites including unit, integration, and table-driven tests.

**Duration:** 30 minutes

**Tasks:**

1. Ask Claude to generate tests for the worker package:

```
> Write comprehensive tests for our worker.go. Include:
> - Table-driven unit tests for ProcessJob
> - Integration tests with a real Redis instance (use testcontainers)
> - Benchmark tests for the hot path
> - Test helpers for common setup/teardown
```

2. Example table-driven test:

```go
// worker_test.go
package main

import (
    "context"
    "errors"
    "testing"
    "time"
)

func TestProcessJob(t *testing.T) {
    tests := []struct {
        name    string
        job     *Job
        wantErr error
    }{
        {
            name: "valid job",
            job: &Job{
                ID:      "test-1",
                Payload: `{"action":"process","data":"test"}`,
            },
            wantErr: nil,
        },
        {
            name: "empty payload",
            job: &Job{
                ID:      "test-2",
                Payload: "",
            },
            wantErr: ErrEmptyPayload,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx, cancel := context.WithTimeout(
                context.Background(), 5*time.Second)
            defer cancel()

            err := ProcessJob(ctx, tt.job)

            if tt.wantErr != nil {
                if !errors.Is(err, tt.wantErr) {
                    t.Errorf("ProcessJob() error = %v, want %v",
                        err, tt.wantErr)
                }
                return
            }
            if err != nil {
                t.Errorf("ProcessJob() unexpected error: %v", err)
            }
        })
    }
}

func TestProcessJob_CancelledContext(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel() // Cancel immediately

    job := &Job{ID: "cancel-test", Payload: "data"}
    err := ProcessJob(ctx, job)

    if err == nil {
        t.Error("expected error for cancelled context")
    }
    if !errors.Is(err, context.Canceled) {
        t.Errorf("expected context.Canceled, got: %v", err)
    }
}

func BenchmarkProcessJob(b *testing.B) {
    ctx := context.Background()
    job := &Job{
        ID:      "bench",
        Payload: `{"action":"process","data":"benchmark"}`,
    }

    b.ResetTimer()
    for range b.N {
        ProcessJob(ctx, job)
    }
}
```

3. Ask Claude to add testcontainers-based integration tests:

```
> Add integration tests that spin up a real Redis container
> using testcontainers-go. Test the full enqueue-dequeue-process cycle.
```

4. Ask Claude to generate a test coverage report and identify untested paths:

```
> Run the tests with coverage and identify any code paths that are
> not covered. Then write tests to cover those paths.
```

---

## Phase 4 — Expert (Week 10)

### Exercise 19: Multi-Agent Orchestration — Boris's Patterns

**Goal:** Learn patterns for running multiple Claude Code instances in parallel on a single codebase.

**Duration:** 60 minutes

**Background:**
Boris (a well-known Claude Code power user) runs 5-10 parallel Claude instances on the same repository, each working on independent tasks. This dramatically increases throughput for large projects. The key is task independence — each instance must work on files that don't conflict.

**Patterns for Parallel Instances:**

**Pattern 1: Feature-Per-Instance**

Open 5 terminal windows, each with Claude working on a separate feature:

```bash
# Terminal 1 — API endpoints
cd ~/project && claude
> Add CRUD endpoints for the orders resource in internal/handler/orders.go

# Terminal 2 — Database layer
cd ~/project && claude
> Create the database migration and store implementation in internal/store/orders.go

# Terminal 3 — Terraform
cd ~/project && claude
> Add an SQS queue and Lambda trigger to terraform/modules/events/

# Terminal 4 — Kubernetes
cd ~/project && claude
> Create Kustomize manifests for the orders service in k8s/orders/

# Terminal 5 — CI/CD
cd ~/project && claude
> Add a GitHub Actions workflow for the orders service in .github/workflows/orders.yaml
```

**Pattern 2: Headless Mode for Batch Operations**

Use `claude -p` (print mode) for non-interactive batch tasks:

```bash
# Run multiple Claude instances in parallel
claude -p "Add godoc comments to all public functions in internal/store/" &
claude -p "Generate table-driven tests for internal/handler/users.go" &
claude -p "Add input validation to all handler functions" &
wait
```

**Pattern 3: Session Continuation**

Resume previous sessions to continue complex work:

```bash
# Start a session
claude

# Later, resume it
claude --resume     # Resume most recent session
claude --continue   # Continue the last conversation
```

**Tasks:**

1. Set up a project with clear boundaries between 5 independent work areas:
   - `internal/auth/` — Authentication
   - `internal/billing/` — Billing
   - `terraform/` — Infrastructure
   - `k8s/` — Kubernetes manifests
   - `.github/workflows/` — CI/CD

2. Launch 3 parallel Claude instances, each working on a different area. Observe how they don't conflict.

3. Practice headless mode:
```bash
# Generate boilerplate in parallel
claude -p "Create internal/auth/middleware.go with JWT validation" > /dev/null &
claude -p "Create internal/billing/stripe.go with payment intent creation" > /dev/null &
wait
echo "Both files created"
```

**Conflict Management:**
- If two instances edit the same file, you will get a git conflict.
- Prevention: Assign each instance to distinct directories or files.
- Resolution: Use `git merge` or ask Claude to resolve the conflict.

**Cost Awareness:**
- Each parallel instance consumes tokens independently.
- Use `/cost` to monitor spending per session.
- Use `/compact` to reduce context in long-running sessions.

---

### Exercise 20: Production Deployment and Team Collaboration

**Goal:** Bring together all skills for a production-ready deployment pipeline and team setup.

**Duration:** 60 minutes

**Scenario:** Your team of 4 engineers is preparing to deploy a new Go microservice to production. You need to set up everything for safe, collaborative development.

**Tasks:**

**Part A — Team CLAUDE.md Setup**

1. Create a team-wide `~/.claude/CLAUDE.md` with shared conventions:

```markdown
# Team Standards

## Git Workflow
- Branch naming: feature/, fix/, infra/
- All PRs require 1 approval
- Squash merge to main
- Conventional commit messages

## Deployment
- Staging auto-deploys from main
- Production deploys require manual approval
- All deployments need rollback plan documented in PR

## On-Call
- Check Grafana dashboards before and after deploy
- Alert thresholds: p99 latency >500ms, error rate >1%
```

2. Create project-level `.claude/settings.json` for the team:

```json
{
  "permissions": {
    "allow": [
      "Read",
      "Glob",
      "Grep",
      "Write",
      "Edit"
    ],
    "deny": [
      "Bash(rm -rf *)",
      "Bash(git push --force)"
    ]
  }
}
```

**Part B — Production Readiness Checklist**

Ask Claude to run through a production readiness checklist:

```
> Review our entire project for production readiness. Check:
> 1. All errors are wrapped with context
> 2. All I/O functions accept context.Context
> 3. Graceful shutdown is implemented
> 4. Health endpoints exist and check dependencies
> 5. Structured logging is used everywhere
> 6. No hardcoded secrets or credentials
> 7. Docker image follows best practices
> 8. Kubernetes manifests have resource limits
> 9. HPA is configured with appropriate thresholds
> 10. CI/CD pipeline includes security scanning
```

**Part C — Deployment Runbook**

Ask Claude to generate a deployment runbook:

```
> Generate a deployment runbook for our Go service that covers:
> - Pre-deployment checks
> - Deployment steps for staging and production
> - Health check verification
> - Rollback procedure
> - Post-deployment monitoring
> Format it as a step-by-step checklist.
```

**Part D — End-to-End Exercise**

Put it all together. Ask Claude to:

1. Review the codebase for any issues.
2. Fix any problems found.
3. Generate missing tests.
4. Build and verify the Docker image.
5. Create the PR with full description.

```
> Review this project end-to-end. Fix any issues you find, add
> any missing tests, verify the Docker build works, and create
> a PR with a comprehensive description.
```

---

## Appendix A — Quick Reference Card

### Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Enter` | Send message |
| `Shift+Tab` | Toggle Plan Mode / Act Mode |
| `Escape` | Cancel current generation |
| `Ctrl+C` | Interrupt / Exit |
| `/` | Show available commands |

### Essential Commands

| Command | Purpose |
|---------|---------|
| `claude` | Start interactive session |
| `claude -p "prompt"` | Headless single-prompt mode |
| `claude --resume` | Resume last session |
| `/help` | Show help |
| `/clear` | Clear context |
| `/compact` | Compress conversation |
| `/commit` | Commit staged changes |
| `/review` | Review changes |
| `/cost` | Show token usage |
| `/init` | Create CLAUDE.md |
| `/memory` | Edit CLAUDE.md |

### CLAUDE.md File Locations

| Path | Scope |
|------|-------|
| `PROJECT_ROOT/CLAUDE.md` | Whole project (checked in) |
| `PROJECT_ROOT/.claude/CLAUDE.md` | Whole project (gitignored, personal) |
| `subdirectory/CLAUDE.md` | That subtree only |
| `~/.claude/CLAUDE.md` | All your projects (global, personal) |

### MCP Configuration Locations

| Path | Scope |
|------|-------|
| `.claude/settings.json` | Project-level (shared) |
| `~/.claude/settings.json` | User-level (personal) |

---

## Appendix B — Quiz and Self-Assessment

### Foundation (Exercises 1–6)

1. **What is the purpose of CLAUDE.md?**
   - [ ] A) It's a README for humans
   - [ ] B) It provides project context and conventions that Claude reads automatically
   - [ ] C) It's required for Claude to work
   - [ ] D) It replaces Go documentation

   **Answer:** B — CLAUDE.md is read by Claude at session start to understand your project's conventions.

2. **When should you use Plan Mode?**
   - [ ] A) Always
   - [ ] B) When you want Claude to make changes immediately
   - [ ] C) When you need to think through architecture before writing code
   - [ ] D) Only for Terraform

   **Answer:** C — Plan Mode is for research and design phases where you don't want code changes yet.

3. **What's the correct way to handle errors in Go according to our CLAUDE.md?**
   - [ ] A) Use panic/recover
   - [ ] B) Ignore errors with `_`
   - [ ] C) Wrap errors with `fmt.Errorf` and `%w`
   - [ ] D) Log and continue

   **Answer:** C — Our standards require explicit error handling with `fmt.Errorf` wrapping.

4. **What command compresses a long conversation while preserving context?**
   - [ ] A) `/clear`
   - [ ] B) `/compact`
   - [ ] C) `/reset`
   - [ ] D) `/context`

   **Answer:** B — `/compact` summarizes the conversation while retaining important information.

### Intermediate (Exercises 7–12)

5. **Where do custom slash commands live for project scope?**
   - [ ] A) `~/.claude/commands/`
   - [ ] B) `.claude/commands/`
   - [ ] C) `CLAUDE.md`
   - [ ] D) `.github/commands/`

   **Answer:** B — Project-scoped commands go in `.claude/commands/` at the project root.

6. **What is MCP in the context of Claude Code?**
   - [ ] A) A testing framework
   - [ ] B) Model Context Protocol — it connects Claude to external data sources
   - [ ] C) A Git merge strategy
   - [ ] D) A CI/CD platform

   **Answer:** B — MCP servers extend Claude's capabilities by connecting it to databases, APIs, and other tools.

7. **What does our CLAUDE.md say about Docker images?**
   - [ ] A) Use single-stage builds for simplicity
   - [ ] B) Always use `latest` tag
   - [ ] C) Must be multi-stage builds
   - [ ] D) Use Alpine only

   **Answer:** C — Our infrastructure patterns require multi-stage Docker builds.

8. **In a GitHub Actions workflow, how should you run Go tests with race detection?**
   - [ ] A) `go test ./...`
   - [ ] B) `go test -race ./...`
   - [ ] C) `go test -v ./...`
   - [ ] D) `go build && go test`

   **Answer:** B — The `-race` flag enables the race detector, essential for concurrent Go code.

### Advanced (Exercises 13–18)

9. **What hook type runs before a tool call and can block it?**
   - [ ] A) `PostToolUse`
   - [ ] B) `PreToolUse`
   - [ ] C) `Notification`
   - [ ] D) `BeforeAction`

   **Answer:** B — `PreToolUse` hooks can inspect and block tool calls before they execute.

10. **According to CLAUDE.md, should we use Helm or Kustomize?**
    - [ ] A) Helm
    - [ ] B) Kustomize
    - [ ] C) Raw YAML
    - [ ] D) Either

    **Answer:** B — Our infrastructure patterns explicitly specify Kustomize over Helm.

11. **What's the correct retry strategy per our CLAUDE.md?**
    - [ ] A) `time.Sleep()` with fixed delay
    - [ ] B) Exponential backoff
    - [ ] C) Immediate retry
    - [ ] D) No retries

    **Answer:** B — Never use `time.Sleep()` for retries; use exponential backoff.

12. **Which security scanner is used for Go static analysis?**
    - [ ] A) Trivy
    - [ ] B) Checkov
    - [ ] C) gosec
    - [ ] D) Snyk

    **Answer:** C — gosec is the Go security static analysis tool. Trivy handles containers, Checkov handles Terraform.

### Expert (Exercises 19–20)

13. **In Boris's parallel pattern, what's the key to avoiding conflicts?**
    - [ ] A) Using git locks
    - [ ] B) Working on the same files simultaneously
    - [ ] C) Assigning each instance to independent files/directories
    - [ ] D) Using a single branch

    **Answer:** C — Each parallel Claude instance must work on distinct, non-overlapping files.

14. **What command runs Claude in headless (non-interactive) mode?**
    - [ ] A) `claude --headless`
    - [ ] B) `claude -p "prompt"`
    - [ ] C) `claude --batch`
    - [ ] D) `claude run`

    **Answer:** B — `claude -p` sends a single prompt and prints the result without an interactive session.

15. **Where should team-wide personal conventions go?**
    - [ ] A) `CLAUDE.md` in the project root
    - [ ] B) `~/.claude/CLAUDE.md`
    - [ ] C) `.claude/settings.json`
    - [ ] D) `~/.bashrc`

    **Answer:** B — `~/.claude/CLAUDE.md` applies to all your projects and is personal to you.

### Self-Assessment Scoring

| Score | Level | Recommendation |
|-------|-------|----------------|
| 0–5 | Beginner | Revisit Phase 1 exercises |
| 6–9 | Developing | Review Phase 2 and practice daily |
| 10–12 | Proficient | Move to advanced exercises, start parallel workflows |
| 13–15 | Expert | Mentor others, optimize team workflows |

---

## Program Completion Checklist

After completing all 20 exercises, you should be able to:

- [ ] Navigate and understand any Go codebase with Claude
- [ ] Write effective CLAUDE.md files that encode team standards
- [ ] Use Plan Mode for architecture design before implementation
- [ ] Create custom slash commands for repeated workflows
- [ ] Connect MCP servers for database and API access
- [ ] Generate production-grade Dockerfiles with multi-stage builds
- [ ] Build CI/CD pipelines with GitHub Actions
- [ ] Create Terraform modules following infrastructure standards
- [ ] Write Kustomize-based Kubernetes manifests
- [ ] Implement security scanning automation
- [ ] Run 5-10 parallel Claude instances without conflicts
- [ ] Set up Claude Code for team collaboration
- [ ] Execute end-to-end production deployments
- [ ] Write comprehensive test suites with table-driven tests
- [ ] Configure hooks for automated quality gates
