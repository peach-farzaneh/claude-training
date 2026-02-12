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
