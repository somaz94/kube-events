# Examples

Hands-on examples for using kube-events against a live cluster.

<br/>

## Table of Contents

- [Basic Examples](#basic-examples)
- [Filtering Examples](#filtering-examples)
- [Output Format Examples](#output-format-examples)
- [Watch Mode Examples](#watch-mode-examples)
- [Troubleshooting Scenarios](#troubleshooting-scenarios)

<br/>

## Basic Examples

<br/>

### Show recent events

```bash
# Default: last 1 hour, all namespaces
kube-events

# Last 5 minutes
kube-events --since 5m

# Last 24 hours
kube-events --since 24h

# Summary only
kube-events -s
```

<br/>

### Specific namespace

```bash
# Single namespace
kube-events -n production

# All namespaces explicitly
kube-events --all-namespaces
```

<br/>

## Filtering Examples

<br/>

### Warning events only

```bash
kube-events -t Warning
kube-events -t Warning --since 30m
kube-events -t Warning -n production
```

<br/>

### By resource kind

```bash
# Pod events only
kube-events -k Pod

# Deployment and StatefulSet events
kube-events -k Deployment,StatefulSet
```

<br/>

### By resource name

```bash
# Events for a specific pod
kube-events -N my-app-7f8b9-x2k4p

# Events for multiple resources
kube-events -N api-server,worker
```

<br/>

### By reason

```bash
# CrashLoopBackOff investigation
kube-events -r BackOff

# Probe failures
kube-events -r Unhealthy

# Scheduling issues
kube-events -r FailedScheduling

# Multiple reasons
kube-events -r BackOff,Unhealthy,FailedMount
```

<br/>

### Combined filters

```bash
# Warning Pod events in production from last 30 minutes
kube-events -n production -k Pod -t Warning --since 30m

# BackOff events for a specific app
kube-events -N my-app -r BackOff -t Warning
```

<br/>

## Output Format Examples

<br/>

### JSON for scripting

```bash
# Pipe to jq
kube-events -o json | jq '.summary'

# Extract warning count
kube-events -o json | jq '.summary.warningCount'

# List resources with warnings
kube-events -t Warning -o json | jq '.groups[].name'

# Save report
kube-events -o json > events-report.json
```

<br/>

### Markdown for documentation

```bash
# Generate markdown report
kube-events -o markdown > events.md

# Post to GitHub PR
REPORT=$(kube-events -t Warning -o markdown)
gh pr comment 123 --body "$REPORT"
```

<br/>

### Table for structured view

```bash
kube-events -o table
kube-events -t Warning -o table
```

<br/>

### Plain for logs

```bash
# No color codes for log files
kube-events -o plain > events.log

# Pipe to other tools
kube-events -o plain | grep "BackOff"
```

<br/>

## Watch Mode Examples

```bash
# Watch all events
kube-events -w

# Watch warnings in production
kube-events -w -n production -t Warning

# Watch pod events for specific app
kube-events -w -k Pod -N my-app

# Watch with JSON output
kube-events -w -o json
```

<br/>

## Troubleshooting Scenarios

<br/>

### CrashLoopBackOff investigation

```bash
# Find all crashing pods
kube-events -r BackOff -t Warning

# Focus on specific namespace
kube-events -n production -r BackOff -k Pod --since 1h
```

<br/>

### Image pull issues

```bash
kube-events -r ImagePullBackOff,ErrImagePull -t Warning
```

<br/>

### Node scheduling problems

```bash
kube-events -r FailedScheduling -k Pod
```

<br/>

### Volume mount failures

```bash
kube-events -r FailedMount,FailedAttachVolume -t Warning
```

<br/>

### Post-deployment health check

```bash
# After deploying, check for issues in last 5 minutes
kube-events -n my-app-ns --since 5m -t Warning

# Quick summary
kube-events -n my-app-ns --since 5m -s
```

<br/>

### Multi-cluster check

```bash
# Check staging
kube-events --context staging -t Warning -s

# Check production
kube-events --context production -t Warning -s
```
