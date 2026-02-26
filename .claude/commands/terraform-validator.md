Review Terraform code for infrastructure best practices and security.

Target: $ARGUMENTS (defaults to all .tf files if empty)

## Validation Checklist

Run through every item and report PASS, FAIL, or N/A for each.

### Structure and State
- [ ] Remote backend configured (S3 + DynamoDB locking per CLAUDE.md)
- [ ] Provider versions pinned in `required_providers`
- [ ] Terraform version constraint set in `required_version`
- [ ] Modules used for reusable components (not inline copy-paste)

### Variables and Outputs
- [ ] All variables have `description` and `type`
- [ ] Sensitive variables marked with `sensitive = true`
- [ ] Input validation blocks on variables where applicable
- [ ] All outputs have `description`
- [ ] Sensitive outputs marked with `sensitive = true`

### Security
- [ ] No hardcoded secrets, passwords, or API keys
- [ ] IAM policies follow least privilege (no `"Action": "*"` or `"Resource": "*"`)
- [ ] S3 buckets have encryption enabled and public access blocked
- [ ] Security groups have no unrestricted ingress (`0.0.0.0/0` on sensitive ports)
- [ ] RDS/database instances have `storage_encrypted = true`
- [ ] RDS/database instances have `deletion_protection` in production

### Best Practices
- [ ] Resources tagged with project, environment, and managed-by
- [ ] Lifecycle rules set on stateful resources (databases, storage)
- [ ] `create_before_destroy` used where appropriate
- [ ] No deprecated resource types or arguments
- [ ] Data sources used instead of hardcoded IDs/ARNs

### Naming
- [ ] Consistent resource naming convention: `${project}-${environment}-${purpose}`
- [ ] Local names are descriptive (`main`, `this` only when unambiguous)

## Output Format

For each file reviewed, report:
1. File path
2. Checklist results (PASS/FAIL/N/A per item)
3. Specific issues found with line numbers
4. Suggested fixes (provide the corrected code)

Summarize with a score: X/Y checks passed. List critical issues first.
