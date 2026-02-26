# Onboarding Guide: Claude Code for Go/DevOps Engineers

> **Time to complete:** 1–2 hours
> **Goal:** Get productive with Claude Code on your first day

---

## 1. Installation and Authentication

```bash
# Install Claude Code
npm install -g @anthropic-ai/claude-code

# Verify
claude --version

# Authenticate with your Anthropic account
claude auth login
```

## 2. First Session

```bash
# Navigate to your project
cd ~/your-project

# Start Claude Code
claude
```

You are now in an interactive session. Claude can see your project files. Try these:

```
> Explain the structure of this project
> What does main.go do?
> Show me all TODO comments in the codebase
```

## 3. Understand CLAUDE.md

Every project should have a `CLAUDE.md` at the root. Claude reads it automatically at the start of each session. It tells Claude about your team's conventions.

Our `CLAUDE.md` covers:
- **Architecture:** Structured logging (logrus), YAML configs, error wrapping with `%w`
- **Coding standards:** Effective Go, errgroup, max 50-line functions, godoc on public APIs
- **Infrastructure:** Terraform remote state in S3, Kustomize (not Helm), GitHub Actions CI/CD, multi-stage Docker builds
- **Gotchas:** Validate YAML before applying, always use `context.Context` for I/O, exponential backoff (never `time.Sleep`), indexes for foreign keys

If you are starting a new project, run `/init` to generate one interactively.

## 4. Essential Commands

| Command | What It Does |
|---------|-------------|
| `/help` | Show all available commands |
| `/clear` | Reset conversation context |
| `/compact` | Compress conversation, keep summary |
| `/commit` | Stage changes and commit with generated message |
| `/review` | Review recent changes |
| `/cost` | Show token usage |
| `/init` | Generate a CLAUDE.md |
| `/memory` | Edit CLAUDE.md |
| `Shift+Tab` | Toggle Plan Mode (research only) / Act Mode (make changes) |
| `Escape` | Cancel current generation |

## 5. Your First Task

Try a small, real task on the codebase:

```
> Add a /ready endpoint to the HTTP server that checks the database
> connection. Follow our CLAUDE.md standards.
```

Watch how Claude:
1. Reads existing code to understand the patterns
2. Uses logrus for logging (per CLAUDE.md)
3. Wraps errors with `fmt.Errorf` and `%w` (per CLAUDE.md)
4. Adds a godoc comment on the public handler (per CLAUDE.md)

## 6. Plan Before You Build

For anything non-trivial, start in Plan Mode:

1. Press `Shift+Tab` — the mode indicator changes to "Plan"
2. Describe what you want to build
3. Claude researches and proposes a plan without changing files
4. Review and refine the plan
5. Press `Shift+Tab` again to switch to Act Mode
6. Say "implement the plan" and Claude executes it

## 7. Git Workflow

Claude handles git operations directly:

```
> Create a branch called feature/add-caching
> [make changes]
> Commit with a descriptive message
> Create a PR with a summary and test plan
```

Or use the shortcut:
```
/commit
```

## 8. Custom Commands

Check `.claude/commands/` for project-specific slash commands. These automate common team workflows. Use them by typing `/project:command-name`.

## 9. What Not to Do

- **Don't paste secrets** into the conversation (API keys, passwords)
- **Don't blindly accept changes** — always review diffs before committing
- **Don't skip Plan Mode** for architectural decisions — think first, code second
- **Don't fight Claude** on CLAUDE.md rules — if you need an exception, explain why

## 10. Next Steps

1. Complete **Exercise 1** in the [Training Program](Claude-Code-Training-Go-DevOps.md)
2. Read [Best Practices](../reference/Best-Practices.md) for team conventions
3. Bookmark the [Troubleshooting Guide](../reference/Troubleshooting-Guide.md)
4. Review [Boris's Patterns](../reference/Boris-Patterns.md) once you are comfortable with basics
