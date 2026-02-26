# Troubleshooting Guide

> Common issues and solutions when using Claude Code for Go/DevOps workflows.

---

## Table of Contents

1. [Installation and Authentication](#installation-and-authentication)
2. [Session Issues](#session-issues)
3. [Go-Specific Issues](#go-specific-issues)
4. [Terraform Issues](#terraform-issues)
5. [Kubernetes Issues](#kubernetes-issues)
6. [Docker Issues](#docker-issues)
7. [Git and GitHub Issues](#git-and-github-issues)
8. [MCP Issues](#mcp-issues)
9. [Hooks Issues](#hooks-issues)
10. [Performance Issues](#performance-issues)

---

## Installation and Authentication

### Claude Code won't install

**Symptom:** `npm install -g @anthropic-ai/claude-code` fails.

**Solutions:**
```bash
# Check Node.js version (requires 18+)
node --version

# Try with explicit registry
npm install -g @anthropic-ai/claude-code --registry https://registry.npmjs.org

# Permission issues on Linux/macOS
sudo npm install -g @anthropic-ai/claude-code

# Or fix npm permissions (preferred)
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
export PATH=~/.npm-global/bin:$PATH
npm install -g @anthropic-ai/claude-code
```

### Authentication fails

**Symptom:** `claude auth login` hangs or errors.

**Solutions:**
```bash
# Clear cached credentials
rm -rf ~/.claude/auth*

# Re-authenticate
claude auth login

# Check API key if using direct key auth
echo $ANTHROPIC_API_KEY
```

---

## Session Issues

### Claude loses context mid-conversation

**Symptom:** Claude forgets earlier decisions or files it already read.

**Solutions:**
- Use `/compact` to compress context before hitting limits
- Break large tasks into smaller, focused sessions
- Reference specific files by path: "In `internal/handler/users.go`, the function on line 45..."
- Use CLAUDE.md for conventions so they don't consume conversation context

### Claude ignores CLAUDE.md conventions

**Symptom:** Claude generates code that violates CLAUDE.md rules (e.g., uses panic instead of error returns).

**Solutions:**
- Verify CLAUDE.md is at the project root
- Start a fresh session (`/clear` or new terminal)
- Be explicit: "Follow the error handling convention from CLAUDE.md"
- Check for conflicting instructions in subdirectory CLAUDE.md files

### Claude modifies the wrong file

**Symptom:** Claude edits a file you didn't intend.

**Solutions:**
- Always specify file paths explicitly: "Edit `internal/handler/users.go`"
- Review the diff before approving: "Show me the diff first"
- Use Plan Mode (`Shift+Tab`) to confirm the approach before changes
- If a bad edit was made, undo with git: `git checkout -- path/to/file`

---

## Go-Specific Issues

### Generated code doesn't compile

**Symptom:** Claude writes code with syntax errors or unresolved imports.

**Solutions:**
```
> Run go build ./... and fix any errors
```
Or ask Claude to self-check:
```
> Verify this compiles by checking all imports and type signatures
```

### Wrong import paths

**Symptom:** Claude uses incorrect module paths.

**Solutions:**
- Ensure `go.mod` is at the project root with the correct module path
- Tell Claude explicitly: "The module path is github.com/org/project"
- Add the module path to CLAUDE.md

### go mod tidy fails

**Symptom:** Dependency resolution errors after Claude adds new imports.

**Solutions:**
```bash
# Clear module cache and retry
go clean -modcache
go mod tidy

# If a specific version is needed
go get github.com/some/package@v1.2.3
go mod tidy
```

### Tests use wrong package name

**Symptom:** Claude generates `package main_test` instead of `package main`.

**Solutions:**
- Specify: "Use black-box testing with `package main_test`" or "Use white-box testing with `package main`"
- Add testing conventions to CLAUDE.md

---

## Terraform Issues

### Claude generates HCL syntax errors

**Symptom:** `terraform validate` fails on Claude-generated code.

**Solutions:**
- Use a PostToolUse hook to auto-validate (see [Hooks](#hooks-issues))
- Ask Claude: "Run terraform validate and fix any issues"
- Specify the Terraform version: "We use Terraform 1.7+"

### State file conflicts

**Symptom:** Claude suggests changes that conflict with existing state.

**Solutions:**
- Always tell Claude the current state: "We have an existing RDS instance managed by Terraform"
- Use Plan Mode to review infrastructure changes before applying
- Never let Claude run `terraform apply` without your review

### Provider version mismatches

**Symptom:** Claude pins different provider versions than your project uses.

**Solutions:**
Add to CLAUDE.md:
```markdown
## Terraform Versions
- Terraform: >= 1.7
- AWS Provider: ~> 5.0
- Kubernetes Provider: ~> 2.25
```

---

## Kubernetes Issues

### Claude generates Helm charts instead of Kustomize

**Symptom:** Claude creates `Chart.yaml` and templates despite CLAUDE.md saying Kustomize.

**Solutions:**
- Be explicit: "Use Kustomize, not Helm. Create base/ and overlays/ directories."
- Verify CLAUDE.md contains: "K8s manifests use Kustomize, not Helm"
- Start a fresh session if Claude was previously discussing Helm

### Invalid YAML in manifests

**Symptom:** `kubectl apply` fails with parsing errors.

**Solutions:**
```
> Validate the YAML syntax of all files in k8s/
> Run: kubectl apply --dry-run=client -k k8s/overlays/staging/
```

### Missing resource limits

**Symptom:** Claude doesn't set resource requests/limits.

**Solutions:**
Add to CLAUDE.md:
```markdown
## Kubernetes Standards
- All containers must have resource requests AND limits
- CPU requests: 100m minimum
- Memory requests: 128Mi minimum
```

---

## Docker Issues

### Large image sizes

**Symptom:** Docker image is hundreds of MB for a Go binary.

**Solutions:**
Ensure CLAUDE.md has: "Docker images must be multi-stage builds"

Ask Claude:
```
> Optimize the Dockerfile. Use multi-stage build with scratch
> or distroless as the final stage. The Go binary should be
> statically compiled with CGO_ENABLED=0.
```

### Build cache not utilized

**Symptom:** Docker builds are slow because dependencies re-download every time.

**Solutions:**
```
> Reorder the Dockerfile to copy go.mod and go.sum first,
> then run go mod download, then copy the rest of the source.
> This maximizes layer caching.
```

---

## Git and GitHub Issues

### Claude creates commits with wrong author

**Symptom:** Commits show the wrong author information.

**Solutions:**
```bash
# Set git config for the repo
git config user.name "Your Name"
git config user.email "you@example.com"
```

### PR creation fails

**Symptom:** Claude can't create a pull request.

**Solutions:**
```bash
# Check gh CLI auth
gh auth status

# Re-authenticate if needed
gh auth login

# Check you have push access
git push -u origin your-branch
```

### Merge conflicts after parallel instances

**Symptom:** Multiple Claude instances edited overlapping files.

**Solutions:**
- See [Boris's Patterns](Boris-Patterns.md) for conflict prevention
- Use git worktrees for full isolation
- Assign a single "wiring" instance to handle shared files last
- Ask Claude to resolve: "Resolve the merge conflict in internal/handler/router.go"

---

## MCP Issues

### MCP server won't start

**Symptom:** Claude can't connect to configured MCP server.

**Solutions:**
```bash
# Test the MCP server manually
npx -y @anthropic-ai/mcp-server-postgres "postgresql://..."

# Check .claude/settings.json syntax
cat .claude/settings.json | python3 -m json.tool

# Ensure npx is available
which npx
```

### MCP queries return errors

**Symptom:** Database queries through MCP fail.

**Solutions:**
- Verify connection string in `.claude/settings.json`
- Check database is running and accessible
- Verify credentials have the right permissions
- Test connection outside Claude: `psql "postgresql://..."`

---

## Hooks Issues

### Hook never fires

**Symptom:** PostToolUse hook doesn't run after file edits.

**Solutions:**
- Check the `matcher` regex matches the tool name exactly
- Verify `.claude/settings.json` is valid JSON
- Test the hook command manually in a terminal
- Check that the hook command has execute permissions

### Hook blocks everything

**Symptom:** PreToolUse hook blocks unintended operations.

**Solutions:**
- Make the matcher more specific (e.g., match only certain file patterns)
- Add logging to the hook to see what's being matched
- Temporarily disable: rename `.claude/settings.json` to test

### Hook command errors

**Symptom:** Hook runs but produces errors.

**Solutions:**
- Test the command standalone: copy it and run in a terminal
- Ensure paths are absolute (hooks may run from a different working directory)
- Check that required tools (`go`, `terraform`, `kubectl`) are in PATH

---

## Performance Issues

### Claude is slow to respond

**Causes and solutions:**
- **Large codebase:** Use `/compact` more frequently
- **Too many files open:** Focus Claude on specific directories
- **Complex task:** Break into smaller steps
- **Long conversation:** Start a new session for unrelated tasks

### High token usage

**Solutions:**
- Use `/cost` to monitor spending
- Use `/compact` to reduce context size
- Write focused prompts — one task per message
- Use headless mode (`claude -p`) for simple tasks
- Reference CLAUDE.md instead of repeating conventions

### Parallel instances overwhelming the machine

**Solutions:**
- Limit to 5 concurrent instances on 8GB RAM machines
- Use `claude -p` (headless) instead of interactive for batch work
- Stagger instance launches by a few seconds
- Monitor with `htop` or `top`
