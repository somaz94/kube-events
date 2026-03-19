# Usage

Complete guide for using kube-events CLI.

<br/>

## Table of Contents

- [Basic Usage](#basic-usage)
- [Global Flags](#global-flags)
- [Output Formats](#output-formats)
- [Filtering](#filtering)
- [Watch Mode](#watch-mode)
- [CI/CD Integration](#cicd-integration)

<br/>

## Basic Usage

```bash
# Show all events from the last 1 hour (default)
kube-events

# Show events from a specific namespace
kube-events -n production

# Show only Warning events
kube-events -t Warning

# Show events from all namespaces
kube-events --all-namespaces

# Watch events in real-time
kube-events -w
```

<br/>

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--kubeconfig` | | `$KUBECONFIG` or `~/.kube/config` | Path to kubeconfig file |
| `--context` | | Current context | Kubernetes context to use |
| `--namespace` | `-n` | All | Filter by namespace (repeatable) |
| `--kind` | `-k` | All | Filter by involved object kind (repeatable) |
| `--name` | `-N` | All | Filter by involved object name (repeatable) |
| `--type` | `-t` | All | Filter by event type: `Normal`, `Warning` |
| `--reason` | `-r` | All | Filter by event reason (e.g., `BackOff`, `Unhealthy`) |
| `--since` | | `1h` | Show events newer than relative duration |
| `--output` | `-o` | `color` | Output format: `color`, `plain`, `json`, `markdown`, `table` |
| `--summary-only` | `-s` | `false` | Show summary statistics only |
| `--all-namespaces` | | `false` | Show events from all namespaces |
| `--watch` | `-w` | `false` | Watch for new events in real-time |

<br/>

## Output Formats

<br/>

### Color (default)

Human-readable output with ANSI color codes. Best for terminal use.

```bash
kube-events -o color
```

```
Pod/app-1 [default] (3 events)
  ! BackOff            2m       Back-off restarting failed container
  ! Unhealthy          5m       Readiness probe failed
    Scheduled          8m       Successfully assigned...

Deployment/api [prod] (1 event)
    ScalingUp          1m       Scaled up replica set

Summary: 4 events, 2 resources | Warning: 2 | Normal: 2
```

| Symbol | Color | Meaning |
|--------|-------|---------|
| `! ` | Yellow | Warning event |
| `  ` | Green | Normal event |

<br/>

### Plain

Same as color but without ANSI escape codes. For piping or log files.

```bash
kube-events -o plain
```

<br/>

### JSON

Machine-readable JSON output. Best for CI/CD pipelines and scripting.

```bash
kube-events -o json
```

```json
{
  "summary": {
    "totalEvents": 4,
    "warningCount": 2,
    "normalCount": 2,
    "resources": 2
  },
  "groups": [
    {
      "kind": "Pod",
      "name": "app-1",
      "namespace": "default",
      "events": [
        {
          "type": "Warning",
          "reason": "BackOff",
          "message": "Back-off restarting failed container",
          "age": "2m",
          "count": 5
        }
      ]
    }
  ]
}
```

<br/>

### Markdown

Markdown-formatted output. Best for GitHub PR comments.

```bash
kube-events -o markdown
```

<br/>

### Table

Compact tabular output.

```bash
kube-events -o table
```

```
TYPE      RESOURCE                                 REASON               AGE      MESSAGE
------------------------------------------------------------------------------------------------------------------------
Warning   [default] Pod/app-1                      BackOff              2m       Back-off restarting failed container
Warning   [default] Pod/app-1                      Unhealthy            5m       Readiness probe failed
Normal    [prod] Deployment/api                    ScalingUp            1m       Scaled up replica set
------------------------------------------------------------------------------------------------------------------------
Total: 3 events, 2 resources (Warning: 2, Normal: 1)
```

<br/>

## Filtering

<br/>

### By namespace

```bash
kube-events -n production
kube-events -n staging -n production  # multiple namespaces
```

<br/>

### By resource kind

```bash
kube-events -k Pod
kube-events -k Pod,Deployment  # multiple kinds
```

<br/>

### By resource name

```bash
kube-events -N api-server
kube-events -N api-server,worker  # multiple names
```

<br/>

### By event type

```bash
kube-events -t Warning
kube-events -t Normal
```

<br/>

### By event reason

```bash
kube-events -r BackOff
kube-events -r BackOff,Unhealthy,FailedMount
```

<br/>

### Time window

```bash
kube-events --since 5m     # last 5 minutes
kube-events --since 1h     # last 1 hour (default)
kube-events --since 24h    # last 24 hours
```

<br/>

### Combined filters

```bash
kube-events -n production -k Pod -t Warning --since 30m
```

<br/>

## Watch Mode

Monitor events in real-time:

```bash
# Watch all events
kube-events -w

# Watch warnings only in production
kube-events -w -n production -t Warning

# Watch specific pod events
kube-events -w -N my-app -k Pod
```

Press `Ctrl+C` to stop.

<br/>

## CI/CD Integration

<br/>

### GitHub Actions

```yaml
- name: Check cluster events
  run: |
    EVENTS=$(kube-events -t Warning --since 10m -o json)
    WARNING_COUNT=$(echo "$EVENTS" | jq '.summary.warningCount')
    if [ "$WARNING_COUNT" -gt 0 ]; then
      echo "::warning::$WARNING_COUNT warning events in cluster"
    fi

- name: Post events to PR
  if: github.event_name == 'pull_request'
  run: |
    REPORT=$(kube-events --since 10m -o markdown)
    gh pr comment ${{ github.event.number }} --body "$REPORT"
```

<br/>

### GitLab CI

```yaml
event-check:
  stage: validate
  script:
    - kube-events -t Warning --since 10m -o json > events.json
  artifacts:
    when: always
    paths:
      - events.json
```

<br/>

## Kubeconfig

kube-events uses the standard Kubernetes client configuration:

1. `--kubeconfig` flag (highest priority)
2. `$KUBECONFIG` environment variable
3. `~/.kube/config` (default)

```bash
# Use specific context
kube-events --context prod-cluster

# Use specific kubeconfig
kube-events --kubeconfig /path/to/config
```
