Review Go code for quality, correctness, and adherence to project standards.

Target: $ARGUMENTS (defaults to all .go files if empty)

## Review Checklist

Run through every item and report PASS, FAIL, or N/A for each.

### Error Handling (CLAUDE.md: "explicit error handling over panics")
- [ ] All errors are wrapped with context using `fmt.Errorf("doing X: %w", err)`
- [ ] No ignored errors (`_ = fn()` where fn returns error)
- [ ] No bare `panic()` or `log.Fatal()` outside of main
- [ ] Custom error types implement the `error` interface correctly
- [ ] `errors.Is` / `errors.As` used for error inspection (not type assertions)
- [ ] Sentinel errors are `var`, not `const`, and use `errors.New` or `fmt.Errorf`

### Context Usage (CLAUDE.md: "Use contexts for timeout handling and cancellation")
- [ ] All functions that do I/O accept `context.Context` as first parameter
- [ ] Contexts are propagated through the call chain (not dropped)
- [ ] `context.Background()` only used at top-level entry points
- [ ] Timeout/cancel contexts used for external calls (DB, HTTP, Redis)

### Concurrency (CLAUDE.md: "Use errgroup for concurrent operations")
- [ ] `errgroup` used for managing goroutine lifecycles
- [ ] No goroutine leaks (all goroutines have a termination path)
- [ ] No `time.Sleep()` for retries — exponential backoff used instead
- [ ] Shared state protected by mutex or channels
- [ ] Race conditions: would `go test -race` pass?

### Code Quality (CLAUDE.md: "Follow Effective Go guidelines")
- [ ] All public functions and types have godoc comments
- [ ] Functions are under 50 lines
- [ ] No global mutable state
- [ ] Interfaces defined where consumed, not where implemented
- [ ] Receiver names are consistent and short (not `this` or `self`)

### Logging (CLAUDE.md: "structured logging with logrus")
- [ ] `logrus` used (not `fmt.Println` or `log` stdlib for app logging)
- [ ] Structured fields used (`WithField`/`WithFields`), not string concatenation
- [ ] Appropriate log levels: Debug, Info, Warn, Error
- [ ] Sensitive data not logged (passwords, tokens, PII)

### Testing
- [ ] Table-driven tests for functions with multiple cases
- [ ] Test names follow `TestFunctionName_Scenario` convention
- [ ] Edge cases covered (nil input, empty strings, context cancellation)
- [ ] No `time.Sleep` in tests — use channels or `testing.T` helpers

## Output Format

For each file reviewed, report:
1. File path and package
2. Checklist results (PASS/FAIL/N/A per item)
3. Specific issues with line numbers and severity (CRITICAL / WARNING / INFO)
4. Suggested fixes with corrected code snippets

Summarize with: X/Y checks passed. List CRITICAL issues first, then WARNING, then INFO.
