Analyze and optimize Dockerfiles for image size, build speed, and security.

Target: $ARGUMENTS (defaults to all Dockerfiles in the project if empty)

## Optimization Checklist

Run through every item and report PASS, FAIL, or N/A for each.

### Build Strategy (CLAUDE.md: "Docker images must be multi-stage builds")
- [ ] Multi-stage build used (separate builder and runtime stages)
- [ ] Final stage uses minimal base (`scratch`, `distroless`, or `alpine`)
- [ ] Builder stage uses an appropriate SDK image
- [ ] Build arguments don't contain secrets

### Layer Caching
- [ ] `go.mod` and `go.sum` copied before source code
- [ ] `go mod download` runs as a separate layer (before `COPY . .`)
- [ ] Infrequently changing layers come before frequently changing ones
- [ ] `.dockerignore` exists and excludes: `.git`, `*.md`, `docs/`, test fixtures, IDE files

### Go-Specific
- [ ] `CGO_ENABLED=0` set for static binary (required for `scratch`/`distroless`)
- [ ] `GOOS=linux` explicitly set
- [ ] `-ldflags="-s -w"` used to strip debug info and symbol tables
- [ ] Build target specifies the exact main package (`./cmd/server`)

### Security
- [ ] Final image runs as non-root user (`USER` directive or numeric UID)
- [ ] No secrets in environment variables or build args
- [ ] CA certificates copied if the binary makes HTTPS calls
- [ ] Base image tags are pinned (e.g., `golang:1.24-alpine`, not `golang:latest`)
- [ ] No unnecessary packages installed in the final stage
- [ ] No `COPY --from=builder /` (only copy what's needed)

### Health and Metadata
- [ ] `HEALTHCHECK` instruction present
- [ ] `EXPOSE` declares the application port
- [ ] OCI labels present (maintainer, version, description)

### Size Optimization
- [ ] `apk add --no-cache` used (no leftover cache in Alpine)
- [ ] Temporary files cleaned up in the same RUN layer they're created
- [ ] No development tools (git, make, gcc) in the final image
- [ ] Final image size estimated — flag if > 50MB for a Go service

## Output Format

For each Dockerfile reviewed:
1. File path
2. Current estimated image size (if buildable) or size assessment
3. Checklist results (PASS/FAIL/N/A per item)
4. Optimization opportunities ranked by impact (HIGH / MEDIUM / LOW)
5. Rewritten Dockerfile with all optimizations applied

Provide a before/after comparison with estimated size savings.

## Optimized Dockerfile Template

If the current Dockerfile needs a complete rewrite, use this as a starting point:

```dockerfile
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache ca-certificates
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/app ./cmd/server

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app /app
USER 65534
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD ["/app", "-healthcheck"]
ENTRYPOINT ["/app"]
```
