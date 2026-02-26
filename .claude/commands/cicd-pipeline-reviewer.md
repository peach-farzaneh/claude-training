Review CI/CD pipeline configurations for best practices, security, and reliability.

Target: $ARGUMENTS (defaults to all YAML files in .github/workflows/ if empty)

## Review Checklist

Run through every item and report PASS, FAIL, or N/A for each.

### Workflow Structure (CLAUDE.md: "CI/CD runs on GitHub Actions")
- [ ] Jobs separated by concern (lint, test, build, deploy are distinct jobs)
- [ ] Independent jobs run in parallel (lint and test don't depend on each other)
- [ ] Dependent jobs use `needs:` correctly (build needs lint+test)
- [ ] Workflow triggers are appropriate (`push`, `pull_request`, `tags`)
- [ ] Branch filters limit runs to relevant branches

### Actions and Versions
- [ ] Actions pinned to specific versions or SHA (not `@main` or `@latest`)
- [ ] `actions/checkout@v4` used (current version)
- [ ] `actions/setup-go@v5` with pinned Go version
- [ ] Third-party actions reviewed for trustworthiness

### Testing
- [ ] Tests run with `-race` flag for race detection
- [ ] Code coverage collected and uploaded
- [ ] Service containers used for integration tests (Postgres, Redis)
- [ ] Service containers have health checks configured
- [ ] Test environment variables set via `env:` (not hardcoded in commands)

### Linting
- [ ] `golangci-lint` runs on every PR
- [ ] Lint version pinned
- [ ] Lint does not block on warnings (only errors)

### Build and Artifacts
- [ ] Docker build uses `docker/build-push-action`
- [ ] Image tags include SHA, branch, and semver where applicable
- [ ] Images pushed only on main/tag (not on PRs)
- [ ] Build cache utilized (GitHub Actions cache or registry cache)

### Security
- [ ] `permissions:` block set at workflow or job level (least privilege)
- [ ] Secrets accessed via `${{ secrets.* }}` (never hardcoded)
- [ ] No secrets printed to logs (no `echo ${{ secrets.* }}`)
- [ ] `pull_request` target used for PRs (not `pull_request_target` without review)
- [ ] OIDC used for cloud provider auth where possible (not long-lived keys)
- [ ] Security scanning jobs present (gosec, trivy, checkov)
- [ ] Dependency review or audit step included

### Deployment
- [ ] Staging deploys automatically on main branch merge
- [ ] Production deploys require manual approval (`environment: production`)
- [ ] Deployment uses immutable image tags (SHA-based, not `latest`)
- [ ] Rollback strategy documented or automated
- [ ] Post-deploy health check or smoke test exists

### Reliability
- [ ] `timeout-minutes` set on jobs to prevent runaway workflows
- [ ] Retry logic for flaky external calls (npm install, docker push)
- [ ] Caching used for dependencies (`actions/setup-go` caches by default)
- [ ] Matrix strategy used if testing multiple Go versions or platforms
- [ ] `concurrency` group set to cancel superseded runs on the same branch

## Output Format

For each workflow file reviewed:
1. File path
2. Workflow trigger summary
3. Job dependency graph (visual)
4. Checklist results (PASS/FAIL/N/A per item)
5. Issues with severity (CRITICAL / WARNING / INFO)
6. Suggested fixes with corrected YAML snippets

Summarize with:
- Security score: X/Y security checks passed
- Reliability score: X/Y reliability checks passed
- Overall: List CRITICAL issues that must be fixed before merging.
