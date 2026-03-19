# Development

Guide for building, testing, and contributing to kube-events.

<br/>

## Table of Contents

- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Build](#build)
- [Testing](#testing)
- [CI/CD Workflows](#cicd-workflows)
- [Conventions](#conventions)

<br/>

## Prerequisites

- Go 1.26+
- Make
- kubectl configured (for live testing)
- golangci-lint (for linting)

<br/>

## Project Structure

```
.
├── cmd/
│   ├── main.go                    # Entry point
│   └── cli/
│       ├── root.go                # Root command with global flags
│       ├── run.go                 # Core execution logic
│       ├── watch.go               # Watch mode implementation
│       ├── version.go             # Version subcommand
│       └── cli_test.go            # CLI tests
├── internal/
│   ├── client/
│   │   └── client.go             # Kubernetes client wrapper
│   ├── event/
│   │   ├── types.go              # Event, ResourceKey, ResourceGroup models
│   │   ├── filter.go             # Filter() + GroupByResource()
│   │   └── filter_test.go        # Filter and grouping tests
│   └── report/
│       ├── summary.go            # Output formatters (5 formats)
│       └── summary_test.go       # Report tests
├── docs/                         # Documentation
├── .github/
│   ├── workflows/                # CI/CD workflows
│   ├── dependabot.yml            # Dependency updates
│   └── release.yml               # Release note categories
├── .goreleaser.yml               # Multi-platform build + Krew + Homebrew
├── Makefile                      # Build, test, lint
├── CODEOWNERS                    # Repository ownership
└── go.mod
```

<br/>

### Key Directories

| Directory | Description |
|-----------|-------------|
| `cmd/cli/` | Cobra CLI commands and flag definitions |
| `internal/client/` | Kubernetes client-go wrapper for event fetching |
| `internal/event/` | Event data model, filtering, and resource grouping |
| `internal/report/` | Output formatting (color, plain, JSON, markdown, table) |

<br/>

## Build

```bash
make build           # Build binary → ./kube-events
make clean           # Remove build artifacts
```

<br/>

## Testing

```bash
make test            # Run unit tests (alias)
make test-unit       # go test ./... -v -race -cover
make cover           # Generate coverage report
make cover-html      # Open coverage report in browser
```

<br/>

### Test Coverage

| Package | Coverage |
|---------|----------|
| `internal/event` | 100% |
| `internal/report` | 97.2% |
| `cmd/cli` | 24.7% |

<br/>

## CI/CD Workflows

| Workflow | Trigger | Description |
|----------|---------|-------------|
| `ci.yml` | push, PR, dispatch | Unit tests → Build → Version verify |
| `lint.yml` | dispatch | golangci-lint |
| `release.yml` | tag push `v*` | GoReleaser (binaries + Homebrew + Krew) |
| `changelog-generator.yml` | after release, PR merge | Auto-generate CHANGELOG.md |
| `contributors.yml` | after changelog | Auto-generate CONTRIBUTORS.md |
| `gitlab-mirror.yml` | push(main) | Backup to GitLab |
| `stale-issues.yml` | daily cron | Auto-close stale issues |
| `dependabot-auto-merge.yml` | PR (dependabot) | Auto-merge minor/patch updates |
| `issue-greeting.yml` | issue opened | Welcome message |

<br/>

### Workflow Chain

```
tag push v* → Create release (GoReleaser)
                └→ Generate changelog
                      └→ Generate Contributors
```

<br/>

## Conventions

- **Commits**: Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `ci:`, `chore:`)
- **Secrets**: `PAT_TOKEN` (cross-repo ops), `GITHUB_TOKEN` (releases), `GITLAB_TOKEN` (mirror)
- **paths-ignore**: `.github/workflows/**`, `**/*.md`
