# kube-events

[![CI](https://github.com/somaz94/kube-events/actions/workflows/ci.yml/badge.svg)](https://github.com/somaz94/kube-events/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/somaz94/kube-events)](https://goreportcard.com/report/github.com/somaz94/kube-events)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Latest Tag](https://img.shields.io/github/v/tag/somaz94/kube-events)](https://github.com/somaz94/kube-events/tags)
[![Top Language](https://img.shields.io/github/languages/top/somaz94/kube-events)](https://github.com/somaz94/kube-events)

A CLI tool to view and summarize Kubernetes events with resource grouping and warning highlighting.

> For detailed documentation, see the [docs/](docs/) folder:
>
> [Usage](docs/USAGE.md) |
> [Configuration](docs/CONFIGURATION.md) |
> [Examples](docs/EXAMPLES.md) |
> [Deployment](docs/DEPLOYMENT.md) |
> [Development](docs/DEVELOPMENT.md) |
> [Use Cases](docs/USE-CASES.md)

<br/>

## Why kube-events?

| | `kubectl get events` | `kube-events` |
|---|---|---|
| **Output** | Flat, unsorted list | Resource-grouped, warning-first display |
| **Highlighting** | None | Color-coded warnings with icons |
| **Filtering** | Limited (`--field-selector`) | Namespace, kind, name, type, reason |
| **Time window** | Not supported | `--since 5m`, `1h`, `24h` |
| **Output formats** | Text only | Color, plain, JSON, Markdown, table |
| **Summary** | Not supported | Event counts, resource counts, warning ratio |
| **Watch mode** | `--watch` (raw) | Filtered, formatted real-time stream |
| **CI integration** | Not designed for CI | JSON/Markdown output for pipelines |

<br/>

## Quick Start

### Install

```bash
# Homebrew
brew install somaz94/tap/kube-events

# Krew (kubectl plugin)
kubectl krew install events2

# Binary
curl -sL https://github.com/somaz94/kube-events/releases/latest/download/kube-events_linux_amd64.tar.gz | tar xz
sudo mv kube-events /usr/local/bin/

# From source
go install github.com/somaz94/kube-events/cmd@latest
```

### Upgrade

```bash
# Homebrew
brew update && brew upgrade kube-events

# Krew
kubectl krew upgrade events2

# From source
go install github.com/somaz94/kube-events/cmd@latest
```

### Basic Usage

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

# Watch events in real-time
kube-events -w -n production -t Warning
```

### Example Output

```
Pod/app-1 [default] (3 events)
  ! BackOff            2m       Back-off restarting failed container
  ! Unhealthy          3m       Readiness probe failed
    Scheduled          5m       Successfully assigned...

Deployment/api [default] (1 event)
    ScalingUp          1m       Scaled up replica set

Summary: 4 events, 2 resources | Warning: 2 | Normal: 2
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

## Project Structure

```
cmd/                    # CLI entry point & Cobra commands
internal/
  client/               # Kubernetes client wrapper (EventLister interface)
  event/                # Event model, filtering, grouping
  report/               # Color/JSON/Markdown/Table output
```

<br/>

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

<br/>

## License

This project is licensed under the Apache License 2.0 — see the [LICENSE](LICENSE) file for details.
