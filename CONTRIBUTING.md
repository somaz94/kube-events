# Contributing

Thank you for your interest in contributing to kube-events!

<br/>

## Getting Started

### Prerequisites

- Go 1.26+
- Access to a Kubernetes cluster (for integration testing)
- kubectl configured

### Setup

```bash
git clone https://github.com/somaz94/kube-events.git
cd kube-events
make build
```

<br/>

## Development Workflow

### 1. Create a branch

```bash
git checkout -b feat/your-feature
```

### 2. Make changes and verify

```bash
# Format and lint
make fmt
make vet
make lint

# Run tests
make test

# Build binary
make build

# Verify
./kube-events version
```

### 3. Commit with conventional commits

We use [Conventional Commits](https://www.conventionalcommits.org/):

| Prefix | Usage |
|--------|-------|
| `feat:` | New feature |
| `fix:` | Bug fix |
| `docs:` | Documentation only |
| `ci:` | CI/CD changes |
| `chore:` | Maintenance (deps, version bumps) |
| `refactor:` | Code restructuring |
| `test:` | Test additions/changes |

```bash
git commit -m "feat: add event count threshold filter"
```

### 4. Push and create a PR

```bash
git push origin feat/your-feature
```

Then create a Pull Request on GitHub.

<br/>

## Code Structure

```
cmd/
  main.go              # Entry point
  cli/                 # Cobra CLI commands (root, version)
internal/
  client/              # Kubernetes client wrapper
  event/               # Event data model, filtering, grouping
  report/              # Output formatting (color, plain, json, markdown, table)
```

<br/>

## Running Tests

```bash
make test              # Unit tests with race detection
make cover             # Coverage report
make cover-html        # Open coverage in browser
```

<br/>

## Linting

```bash
make lint              # golangci-lint
make vet               # go vet
make fmt               # go fmt
```

<br/>

## Questions?

Open an [issue](https://github.com/somaz94/kube-events/issues) for questions or discussion.
