# Boris's Patterns: Parallel Claude Code Workflows

> Boris is a well-known Claude Code power user who runs 5–10 parallel instances
> on a single codebase. This document captures his specific techniques.

---

## Core Philosophy

> "Each Claude instance gets one job, one directory, zero overlap."
> — Boris's Rule #1

The key insight: Claude Code instances are cheap to start but expensive when they conflict. The entire strategy is about **task independence**.

---

## Pattern 1: Feature-Per-Terminal

Open N terminals, each Claude instance owns a distinct slice of the codebase.

```
Terminal 1 (API)          → internal/handler/ + internal/model/
Terminal 2 (Data layer)   → internal/store/ + migrations/
Terminal 3 (Infra)        → terraform/
Terminal 4 (K8s)          → k8s/
Terminal 5 (CI/CD)        → .github/workflows/
```

### Setup Script

```bash
#!/bin/bash
# boris-parallel.sh — Launch N Claude instances in split terminals

PROJECT_DIR=$(pwd)
TASKS=(
  "Add CRUD handlers for the orders resource in internal/handler/orders.go and model in internal/model/order.go"
  "Create PostgreSQL store for orders in internal/store/orders.go with migration in migrations/"
  "Add Terraform module for SQS queue in terraform/modules/events/"
  "Create Kustomize manifests for orders service in k8s/services/orders/"
  ".github/workflows/orders.yaml — Add CI pipeline for the orders service"
)

for i in "${!TASKS[@]}"; do
  echo "=== Instance $((i+1)): ${TASKS[$i]:0:60}..."
  claude -p "${TASKS[$i]}" &
  PIDS+=($!)
done

echo "Waiting for ${#PIDS[@]} instances..."
for pid in "${PIDS[@]}"; do
  wait "$pid"
done
echo "All instances complete."
```

### Why It Works
- Each instance writes to files the others never touch
- No merge conflicts
- Linear speedup for independent features

---

## Pattern 2: Headless Batch Mode

Use `claude -p` for fire-and-forget tasks that don't need interaction.

### Boilerplate Generation

```bash
# Generate CRUD for 5 resources in parallel
for resource in users orders products invoices shipments; do
  claude -p "Create internal/handler/${resource}.go with full CRUD handlers, \
    internal/store/${resource}.go with PostgreSQL implementation, and \
    internal/model/${resource}.go with validation. Follow CLAUDE.md standards." &
done
wait
```

### Documentation Generation

```bash
# Generate godoc in parallel across packages
for pkg in handler store model middleware; do
  claude -p "Add godoc comments to all public types and functions in internal/${pkg}/" &
done
wait
```

### Test Generation

```bash
# Generate tests in parallel
for file in $(find internal/ -name "*.go" ! -name "*_test.go"); do
  test_file="${file%.go}_test.go"
  if [ ! -f "$test_file" ]; then
    claude -p "Generate table-driven tests for ${file}. Include edge cases and error paths." &
  fi
done
wait
```

---

## Pattern 3: Pipeline Stages

Run Claude instances sequentially where output of one feeds the next, but parallelize within each stage.

```
Stage 1 (parallel): Generate interfaces and models
  ├── claude -p "Define Store interface in internal/store/store.go"
  ├── claude -p "Define all model types in internal/model/"
  └── claude -p "Define service interfaces in internal/service/"

Stage 2 (parallel, depends on Stage 1): Implement
  ├── claude -p "Implement PostgreSQL store matching the Store interface"
  ├── claude -p "Implement HTTP handlers using the service interfaces"
  └── claude -p "Implement business logic in service layer"

Stage 3 (parallel, depends on Stage 2): Test + Infra
  ├── claude -p "Generate integration tests for the store layer"
  ├── claude -p "Create Dockerfile and K8s manifests"
  └── claude -p "Create CI/CD pipeline"
```

### Implementation

```bash
#!/bin/bash
set -e

echo "=== Stage 1: Interfaces ==="
claude -p "Define Store interface in internal/store/store.go" &
claude -p "Define model types in internal/model/" &
claude -p "Define service interfaces in internal/service/" &
wait

echo "=== Stage 2: Implementation ==="
claude -p "Implement PostgreSQL store matching internal/store/store.go interface" &
claude -p "Implement HTTP handlers using service interfaces from internal/service/" &
claude -p "Implement business logic in internal/service/ using store and model packages" &
wait

echo "=== Stage 3: Test and Infra ==="
claude -p "Generate integration tests for internal/store/" &
claude -p "Create Dockerfile and k8s/ manifests for this service" &
claude -p "Create .github/workflows/ci.yaml for this service" &
wait

echo "=== Done ==="
```

---

## Pattern 4: Review Swarm

Run parallel Claude instances to review different aspects of the same codebase:

```bash
claude -p "Review all Go code for error handling issues. List every place where an error is ignored or not wrapped." &
claude -p "Review all Terraform code for security issues: overly permissive IAM, public buckets, missing encryption." &
claude -p "Review all Dockerfiles for image size and security: multi-stage builds, non-root user, pinned versions." &
claude -p "Review all K8s manifests for production readiness: resource limits, probes, PDB, HPA." &
wait
```

---

## Conflict Prevention Rules

### Rule 1: One Instance Per Directory

Never assign two instances to the same package or directory.

```
GOOD:
  Instance A → internal/handler/
  Instance B → internal/store/

BAD:
  Instance A → internal/handler/users.go
  Instance B → internal/handler/orders.go   # Same directory = risk
```

### Rule 2: Shared Files Are Owned by One Instance

Files touched by multiple features (like `main.go`, `go.mod`, `router.go`) should be updated by a single final instance after all others complete.

```bash
# After all parallel instances finish
claude -p "Wire all new handlers and stores into cmd/server/main.go. \
  Update the router, dependency injection, and graceful shutdown."
```

### Rule 3: Use Git Worktrees for Safety

For maximum isolation, use git worktrees so each instance has its own working directory:

```bash
git worktree add ../project-api feature/api
git worktree add ../project-infra feature/infra
git worktree add ../project-k8s feature/k8s

# Each instance works in its own worktree
(cd ../project-api && claude -p "Add API handlers") &
(cd ../project-infra && claude -p "Add Terraform modules") &
(cd ../project-k8s && claude -p "Add K8s manifests") &
wait

# Merge all branches
git merge feature/api feature/infra feature/k8s
```

---

## Cost Management

### Monitor Spending

```bash
# Check cost in each session
/cost

# Set a budget alert (in your workflow script)
MAX_COST=5.00  # dollars per run
```

### Reduce Token Usage

- Use `/compact` in long interactive sessions
- For batch mode, keep prompts focused on a single task
- Reference CLAUDE.md conventions instead of repeating them in every prompt
- Use short, precise prompts — Claude reads the code, you don't need to paste it

### Cost Per Pattern (Rough Estimates)

| Pattern | Instances | Typical Cost |
|---------|-----------|-------------|
| Single interactive session | 1 | $0.50–$2.00/hr |
| Feature-per-terminal (5 instances) | 5 | $2.00–$5.00/run |
| Headless batch (10 tasks) | 10 | $1.00–$3.00/run |
| Pipeline stages (3 stages x 3) | 9 | $3.00–$6.00/run |
| Review swarm (4 reviewers) | 4 | $1.00–$2.00/run |

---

## Checklist Before Launching Parallel Instances

- [ ] Map every file/directory each instance will touch
- [ ] Verify zero overlap between instances
- [ ] Identify shared files — assign a single "wiring" instance for the end
- [ ] Ensure CLAUDE.md is up to date (all instances read it)
- [ ] Run `git status` — start from a clean working tree
- [ ] Set up separate branches or worktrees if changing shared areas
