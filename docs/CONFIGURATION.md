# Configuration

Reference for all kube-events configuration options.

<br/>

## Table of Contents

- [CLI Flags](#cli-flags)
- [Environment Variables](#environment-variables)
- [Event Types](#event-types)
- [Common Event Reasons](#common-event-reasons)

<br/>

## CLI Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--kubeconfig` | | string | `$KUBECONFIG` or `~/.kube/config` | Path to kubeconfig file |
| `--context` | | string | Current context | Kubernetes context to use |
| `--namespace` | `-n` | []string | (all) | Filter by namespace (repeatable) |
| `--kind` | `-k` | []string | (all) | Filter by involved object kind (repeatable) |
| `--name` | `-N` | []string | (all) | Filter by involved object name (repeatable) |
| `--type` | `-t` | []string | (all) | Filter by event type: `Normal`, `Warning` |
| `--reason` | `-r` | []string | (all) | Filter by event reason (repeatable) |
| `--since` | | string | `1h` | Show events newer than relative duration |
| `--output` | `-o` | string | `color` | Output format: `color`, `plain`, `json`, `markdown`, `table` |
| `--group-by` | `-g` | string | `resource` | Group events by: `resource`, `namespace`, `kind`, `reason` |
| `--summary-only` | `-s` | bool | `false` | Show summary statistics only |
| `--all-namespaces` | | bool | `false` | Show events from all namespaces |
| `--watch` | `-w` | bool | `false` | Watch for new events in real-time |

<br/>

## Environment Variables

| Variable | Description |
|----------|-------------|
| `KUBECONFIG` | Path to kubeconfig file (overridden by `--kubeconfig` flag) |

<br/>

## Event Types

Kubernetes events have two types:

| Type | Description |
|------|-------------|
| `Normal` | Routine operations (scheduling, pulling, scaling) |
| `Warning` | Issues that may need attention (probe failures, back-offs, mount errors) |

<br/>

## Common Event Reasons

<br/>

### Warning Reasons

| Reason | Kind | Description |
|--------|------|-------------|
| `BackOff` | Pod | Container crash loop back-off |
| `Unhealthy` | Pod | Liveness/readiness probe failed |
| `FailedMount` | Pod | Volume mount failure |
| `FailedScheduling` | Pod | Cannot schedule pod (insufficient resources, node affinity, etc.) |
| `ImagePullBackOff` | Pod | Cannot pull container image |
| `OOMKilled` | Pod | Container killed due to out-of-memory |
| `FailedCreate` | ReplicaSet | Cannot create pod |
| `Evicted` | Pod | Pod evicted from node |

<br/>

### Normal Reasons

| Reason | Kind | Description |
|--------|------|-------------|
| `Scheduled` | Pod | Pod assigned to node |
| `Pulling` | Pod | Pulling container image |
| `Pulled` | Pod | Image successfully pulled |
| `Created` | Pod | Container created |
| `Started` | Pod | Container started |
| `Killing` | Pod | Container stopping |
| `ScalingReplicaSet` | Deployment | Replica set scaled up/down |
| `SuccessfulCreate` | ReplicaSet/Job | Pod created successfully |
