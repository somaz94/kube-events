# kube-events

A CLI tool to view and summarize Kubernetes events with resource grouping and warning highlighting.

[![CI](https://github.com/somaz94/kube-events/actions/workflows/ci.yml/badge.svg)](https://github.com/somaz94/kube-events/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/somaz94/kube-events)](https://goreportcard.com/report/github.com/somaz94/kube-events)
[![License](https://img.shields.io/github/license/somaz94/kube-events)](LICENSE)

<br/>

## Why kube-events?

`kubectl get events` output is flat, hard to read, and lacks grouping. **kube-events** provides:

- **Resource grouping** — events organized by Pod, Deployment, Service, etc.
- **Warning highlighting** — warnings sorted first with color-coded output
- **Time filtering** — show events from the last 5m, 1h, 24h, etc.
- **Multiple output formats** — color, plain, JSON, Markdown, table
- **Flexible filtering** — by namespace, kind, name, type, reason

<br/>

### Before (kubectl)

```
LAST SEEN   TYPE      REASON      OBJECT           MESSAGE
2m          Warning   BackOff     pod/app-1        Back-off restarting failed container
5m          Normal    Scheduled   pod/app-1        Successfully assigned...
3m          Warning   Unhealthy   pod/app-1        Readiness probe failed
1m          Normal    ScalingUp   deployment/api   Scaled up replica set
```

<br/>

### After (kube-events)

```
Pod/app-1 [default] (3 events)
  ! BackOff            2m       Back-off restarting failed container
  ! Unhealthy          3m       Readiness probe failed
    Scheduled          5m       Successfully assigned...

Deployment/api [default] (1 event)
    ScalingUp          1m       Scaled up replica set

Summary: 4 events, 2 resources | Warning: 2 | Normal: 2
```

## Quick Start

<br/>

### Install

```bash
# Homebrew
brew install somaz94/tap/kube-events

# Go install
go install github.com/somaz94/kube-events/cmd@latest

# Download binary
# See https://github.com/somaz94/kube-events/releases
```

<br/>

### Usage

```bash
# Show all events from last 1 hour (default)
kube-events

# Show only Warning events
kube-events -t Warning

# Show events from specific namespace
kube-events -n production

# Show events for Pods only
kube-events -k Pod

# Show events from last 5 minutes
kube-events --since 5m

# Filter by resource name
kube-events -N api-server

# Filter by reason
kube-events -r BackOff,Unhealthy

# All namespaces, JSON output
kube-events --all-namespaces -o json

# Summary only
kube-events -s

# Combine filters
kube-events -n prod -k Pod -t Warning --since 30m
```

<br/>

## CLI Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--kubeconfig` | | `~/.kube/config` | Path to kubeconfig |
| `--context` | | current | Kubernetes context |
| `--namespace` | `-n` | all | Filter by namespace (repeatable) |
| `--kind` | `-k` | all | Filter by involved object kind |
| `--name` | `-N` | all | Filter by involved object name |
| `--type` | `-t` | all | Event type: `Normal`, `Warning` |
| `--reason` | `-r` | all | Event reason (e.g., `BackOff`) |
| `--since` | | `1h` | Show events newer than duration |
| `--output` | `-o` | `color` | Format: `color`, `plain`, `json`, `markdown`, `table` |
| `--summary-only` | `-s` | `false` | Show summary statistics only |
| `--all-namespaces` | | `false` | Show events from all namespaces |
| `--watch` | `-w` | `false` | Watch for new events in real-time |

<br/>

## Output Formats

| Format | Flag | Use Case |
|--------|------|----------|
| Color | `-o color` | Terminal (default) |
| Plain | `-o plain` | Terminal without color support |
| JSON | `-o json` | Scripting, piping to `jq` |
| Markdown | `-o markdown` | GitHub PR comments, docs |
| Table | `-o table` | Structured terminal view |

<br/>

## Comparison with kubectl

| Feature | `kubectl get events` | `kube-events` |
|---------|---------------------|---------------|
| Resource grouping | No | Yes |
| Warning highlighting | No | Yes |
| Color output | No | Yes |
| Multiple output formats | Limited | 5 formats |
| Time-based filtering | No | `--since` |
| Reason filtering | No | `--reason` |
| Summary statistics | No | Yes |

<br/>

## Development

```bash
# Build
make build

# Run tests
make test

# Lint
make lint

# Coverage report
make cover
```

<br/>

## License

[Apache License 2.0](LICENSE)
