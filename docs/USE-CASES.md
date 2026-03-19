# Use Cases

Real-world scenarios where kube-events helps.

<br/>

## Table of Contents

- [Incident Response](#incident-response)
- [Post-deployment Monitoring](#post-deployment-monitoring)
- [Cluster Health Dashboard](#cluster-health-dashboard)
- [CI/CD Pipeline Integration](#cicd-pipeline-integration)
- [Multi-cluster Operations](#multi-cluster-operations)

<br/>

## Incident Response

<br/>

### Quick triage during an outage

When pods start crashing, get a fast overview of what's happening:

```bash
# What's going wrong in production right now?
kube-events -n production -t Warning --since 5m

# Focus on specific app
kube-events -n production -N payment-service -t Warning
```

**When to use**: First responder during incidents. Get structured event overview instead of scrolling through raw `kubectl get events`.

<br/>

### Post-incident audit

After resolving an incident, review what happened:

```bash
# Full 24-hour event history
kube-events -n production --since 24h -o json > incident-events.json

# Warning timeline
kube-events -n production -t Warning --since 24h -o table
```

**When to use**: During incident post-mortems to build a timeline of events.

<br/>

## Post-deployment Monitoring

<br/>

### Verify deployment health

After deploying a new version, check for issues:

```bash
# Check for warnings in last 5 minutes after deploy
kube-events -n my-app --since 5m -t Warning

# Quick pass/fail check
WARNING_COUNT=$(kube-events -n my-app --since 5m -t Warning -o json | jq '.summary.warningCount')
if [ "$WARNING_COUNT" -gt 0 ]; then
  echo "WARN: $WARNING_COUNT warning events after deployment"
fi
```

**When to use**: As a post-deploy smoke test in CI/CD pipelines.

<br/>

### Watch rollout in real-time

```bash
# Watch events during deployment rollout
kube-events -w -n my-app -k Pod,Deployment

# Watch only warnings
kube-events -w -n my-app -t Warning
```

**When to use**: During manual deployments or when debugging rollout issues.

<br/>

## Cluster Health Dashboard

<br/>

### Morning health check

Quick daily cluster review:

```bash
# Summary of all warnings in last 12 hours
kube-events -t Warning --since 12h -s

# Detailed warnings grouped by resource
kube-events -t Warning --since 12h

# Export for team review
kube-events -t Warning --since 12h -o markdown > daily-events.md
```

**When to use**: As part of daily operational routine.

<br/>

### Resource-specific monitoring

```bash
# Monitor CronJob execution
kube-events -k Job,CronJob --since 24h

# Check StatefulSet stability
kube-events -k StatefulSet,Pod -N my-database --since 1h

# Node-level events
kube-events -k Node --since 1h
```

<br/>

## CI/CD Pipeline Integration

<br/>

### GitHub Actions post-deploy check

```yaml
- name: Post-deploy event check
  run: |
    sleep 60  # Wait for events to propagate
    EVENTS=$(kube-events -n ${{ env.NAMESPACE }} --since 5m -t Warning -o json)
    WARNING_COUNT=$(echo "$EVENTS" | jq '.summary.warningCount')
    if [ "$WARNING_COUNT" -gt 0 ]; then
      echo "::warning::$WARNING_COUNT warning events detected after deployment"
      echo "$EVENTS" | jq '.groups[] | "\(.kind)/\(.name): \(.events[0].reason) - \(.events[0].message)"'
    fi
```

<br/>

### Scheduled cluster audit

```yaml
name: Cluster Event Audit
on:
  schedule:
    - cron: '0 8 * * *'  # Daily at 8 AM

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - name: Check warnings
        run: |
          REPORT=$(kube-events -t Warning --since 24h -o markdown)
          if [ -n "$REPORT" ]; then
            curl -X POST "${{ secrets.SLACK_WEBHOOK }}" \
              -d "{\"text\": \"Daily cluster event report:\n$REPORT\"}"
          fi
```

<br/>

## Multi-cluster Operations

<br/>

### Compare cluster health across environments

```bash
# Staging
echo "=== Staging ===" && kube-events --context staging -t Warning --since 1h -s

# Production
echo "=== Production ===" && kube-events --context production -t Warning --since 1h -s
```

<br/>

### Environment-specific monitoring

```bash
# Check all environments in a loop
for ctx in staging production; do
  echo "--- $ctx ---"
  kube-events --context $ctx -t Warning --since 1h -o table
done
```

**When to use**: When managing multiple clusters and need a quick health comparison.
