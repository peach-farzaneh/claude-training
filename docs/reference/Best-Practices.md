# Best Practices: Claude Code for Go/DevOps

> Condensed dos and don'ts. Print this out, pin it to your monitor.

---

## General

### Do

- **Start every project with a CLAUDE.md** — it's your team's contract with Claude
- **Use Plan Mode for architecture** — think first, code second (`Shift+Tab`)
- **Review every diff** before committing — Claude is good, not infallible
- **Use `/compact`** in long sessions to preserve context
- **Be specific** — "Add a `/health` endpoint to `internal/handler/health.go`" beats "add a health check"
- **Reference existing code** — "Follow the same pattern as `users.go`"
- **Break large tasks into steps** — 5 small prompts beat 1 huge prompt

### Don't

- **Don't paste secrets** into conversations (API keys, tokens, passwords)
- **Don't skip reading the diff** — always verify before committing
- **Don't fight CLAUDE.md** — if you need an exception, explain why explicitly
- **Don't use Claude for one-liners** you could type faster yourself
- **Don't keep stale context** — use `/clear` when switching topics

---

## Go

### Do

- **Wrap all errors with context:** `fmt.Errorf("doing X: %w", err)`
- **Pass `context.Context`** as the first parameter to all I/O functions
- **Use `errgroup`** for concurrent goroutine management
- **Write table-driven tests** for all business logic
- **Add godoc comments** to every public function and type
- **Use structured logging** with logrus and field-based entries
- **Keep functions under 50 lines** — if it's longer, split it

### Don't

- **Don't use `panic`** for error handling — return errors
- **Don't use `time.Sleep`** for retries — use exponential backoff
- **Don't ignore errors:** never `_ = someFunc()` for functions that return errors
- **Don't use global mutable state** — pass dependencies explicitly
- **Don't use `init()` functions** unless absolutely necessary

### Example: Error Handling

```go
// GOOD
func (s *Store) GetUser(ctx context.Context, id string) (*User, error) {
    row := s.db.QueryRowContext(ctx, "SELECT ... WHERE id = $1", id)
    var u User
    if err := row.Scan(&u.ID, &u.Name); err != nil {
        return nil, fmt.Errorf("querying user %s: %w", id, err)
    }
    return &u, nil
}

// BAD
func (s *Store) GetUser(id string) *User {
    row := s.db.QueryRow("SELECT ... WHERE id = $1", id)
    var u User
    row.Scan(&u.ID, &u.Name) // error ignored, no context, no wrapping
    return &u
}
```

---

## Terraform

### Do

- **Use remote state** in S3 with DynamoDB locking
- **Use modules** for reusable infrastructure components
- **Pin provider versions** in `required_providers`
- **Add `description`** to all variables and outputs
- **Use `validation` blocks** on input variables
- **Tag all resources** with project, environment, and owner

### Don't

- **Don't hardcode values** — use variables for everything environment-specific
- **Don't use `terraform apply` without review** — always plan first
- **Don't store secrets in `.tfvars`** — use AWS Secrets Manager or SSM
- **Don't create resources without lifecycle rules** — especially for databases
- **Don't use `count`** when `for_each` with a map is clearer

### Example: Module Structure

```
terraform/modules/rds/
  main.tf          # Resource definitions
  variables.tf     # Input variables with descriptions and validation
  outputs.tf       # Output values
  versions.tf      # Required providers and versions
```

---

## Kubernetes

### Do

- **Use Kustomize** for manifest management (not Helm)
- **Set resource requests AND limits** on every container
- **Configure liveness and readiness probes** for every service
- **Use PodDisruptionBudgets** for production workloads
- **Use HPA** for auto-scaling based on CPU/memory or custom metrics
- **Use ServiceAccounts** with IRSA for AWS access
- **Set `terminationGracePeriodSeconds`** to allow clean shutdown

### Don't

- **Don't use `latest` tag** in production — pin image versions
- **Don't run as root** — set `securityContext.runAsNonRoot: true`
- **Don't skip resource limits** — unbounded pods get OOM-killed unpredictably
- **Don't use NodePort** in production — use LoadBalancer or Ingress
- **Don't store config in images** — use ConfigMaps and Secrets

### Example: Minimal Production Pod Spec

```yaml
spec:
  serviceAccountName: my-service
  terminationGracePeriodSeconds: 30
  securityContext:
    runAsNonRoot: true
  containers:
    - name: app
      image: ghcr.io/org/app:v1.2.3
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
          port: 8080
        initialDelaySeconds: 5
      readinessProbe:
        httpGet:
          path: /ready
          port: 8080
        initialDelaySeconds: 5
```

---

## Docker

### Do

- **Use multi-stage builds** — build in one stage, run in another
- **Use `scratch` or `distroless`** as the final base image for Go
- **Copy `go.mod`/`go.sum` first** — then `go mod download` — then copy source (maximizes layer cache)
- **Set `CGO_ENABLED=0`** for fully static Go binaries
- **Use `-ldflags="-s -w"`** to strip debug info and reduce binary size
- **Run as non-root** with `USER` directive
- **Add a `HEALTHCHECK`** instruction

### Don't

- **Don't use `golang:latest`** as the final image — it's 800MB+
- **Don't copy the entire repo first** — it invalidates the layer cache on every change
- **Don't install dev tools** in the final image
- **Don't leave build artifacts** in the final image
- **Don't use `ADD`** when `COPY` works (ADD has implicit unpack and URL fetch)

### Example: Optimized Dockerfile

```dockerfile
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache ca-certificates
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/app ./cmd/server

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app /app
USER 65534
EXPOSE 8080
HEALTHCHECK --interval=30s CMD ["/app", "-healthcheck"]
ENTRYPOINT ["/app"]
```

---

## CI/CD (GitHub Actions)

### Do

- **Run lint, test, and build as separate jobs** — they can run in parallel
- **Use `actions/setup-go`** with a pinned Go version
- **Enable race detection:** `go test -race ./...`
- **Upload coverage** to Codecov or similar
- **Use service containers** for integration tests (Postgres, Redis)
- **Use GitHub environment protection** for production deploys
- **Cache Go modules** via `actions/setup-go` (automatic with v5)

### Don't

- **Don't skip linting** — use `golangci-lint-action`
- **Don't deploy without tests passing** — enforce branch protection
- **Don't hardcode secrets in workflows** — use GitHub Secrets
- **Don't use `latest`** for action versions — pin to SHA or tag
- **Don't combine lint + test + build** into one job — parallel is faster

---

## Security

### Do

- **Run `gosec`** on every PR
- **Run `trivy`** on container images
- **Run `checkov`** on Terraform code
- **Generate SBOMs** for container images
- **Scan for hardcoded secrets** with `gitleaks` or similar
- **Run `govulncheck`** for Go dependency vulnerabilities

### Don't

- **Don't commit `.env` files** with real credentials
- **Don't use overly permissive IAM policies** — follow least privilege
- **Don't disable security scanning** to make CI pass — fix the issue
- **Don't skip dependency updates** — automate with Dependabot or Renovate

---

## Claude Code Workflow

### Do

- **One task per prompt** — clear, focused instructions get better results
- **Use Plan Mode for design** — Act Mode for implementation
- **Use `/commit`** for consistent commit messages
- **Use custom commands** for repetitive workflows
- **Use parallel instances** for independent work (see Boris's Patterns)
- **Monitor costs** with `/cost`

### Don't

- **Don't give vague instructions** — "make it better" gets vague results
- **Don't fight Claude** on CLAUDE.md — update CLAUDE.md if rules change
- **Don't run 10 instances on shared files** — split by directory
- **Don't forget `/compact`** in long sessions — context limits are real
