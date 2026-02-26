Review Kubernetes manifests for security, production readiness, and adherence to project standards.

Target: $ARGUMENTS (defaults to all YAML files in k8s/ if empty)

## Validation Checklist

Run through every item and report PASS, FAIL, or N/A for each.

### Manifest Structure (CLAUDE.md: "K8s manifests use Kustomize, not Helm")
- [ ] Kustomize used for environment management (not Helm)
- [ ] Base manifests in `base/` with overlays per environment
- [ ] `kustomization.yaml` present in base and each overlay
- [ ] No hardcoded environment-specific values in base

### Security
- [ ] Containers run as non-root (`runAsNonRoot: true`)
- [ ] Read-only root filesystem where possible (`readOnlyRootFilesystem: true`)
- [ ] No privileged containers (`privileged: false` or omitted)
- [ ] No privilege escalation (`allowPrivilegeEscalation: false`)
- [ ] SecurityContext set at pod and/or container level
- [ ] ServiceAccount specified (not using `default`)
- [ ] RBAC roles follow least privilege
- [ ] No `hostNetwork`, `hostPID`, or `hostIPC` unless justified
- [ ] Secrets not stored in ConfigMaps
- [ ] Image tags are pinned versions (not `latest`)

### Availability
- [ ] Resource requests AND limits set on every container
- [ ] Liveness probe configured with appropriate thresholds
- [ ] Readiness probe configured (separate from liveness)
- [ ] PodDisruptionBudget exists for production workloads
- [ ] `terminationGracePeriodSeconds` set (matches app shutdown time)
- [ ] Rolling update strategy configured (`maxSurge` / `maxUnavailable`)
- [ ] Anti-affinity rules for multi-replica deployments (production)

### Autoscaling
- [ ] HPA configured with appropriate min/max replicas
- [ ] HPA metrics are meaningful (CPU, memory, or custom)
- [ ] Production minReplicas >= 2

### Networking
- [ ] Services use `ClusterIP` (not `NodePort` in production)
- [ ] Port names follow convention (`http`, `grpc`, `metrics`)
- [ ] NetworkPolicy defined to restrict traffic (if applicable)

### Configuration
- [ ] Environment variables sourced from ConfigMaps/Secrets (not inline)
- [ ] ConfigMaps and Secrets managed via Kustomize generators where possible
- [ ] No sensitive data in plain text

### Observability
- [ ] Metrics port exposed (e.g., `:9090` for Prometheus)
- [ ] Annotations for Prometheus scraping if using Prometheus
- [ ] Labels consistent and meaningful (`app`, `version`, `component`)

## Output Format

For each manifest reviewed, report:
1. File path
2. Checklist results (PASS/FAIL/N/A per item)
3. Security issues with severity (CRITICAL / HIGH / MEDIUM / LOW)
4. Suggested fixes with corrected YAML snippets

Summarize with:
- Security score: X/Y security checks passed
- Availability score: X/Y availability checks passed
- Overall: List CRITICAL and HIGH issues that must be fixed before deployment.
