# CLAUDE.md - kube-events

CLI tool to view and summarize Kubernetes events, grouped by resource, namespace, kind or reason, with warnings highlighted.

<br/>

## Build & Test

```bash
make build           # Build binary
make test            # Unit tests with -race and coverage
make lint            # golangci-lint
make test-e2e-watch  # Watch mode across namespaces (needs a kind cluster)
make demo            # Deploy demo resources; local contexts only (`make demo-clean` removes them)
```

`make help` lists every target.

<br/>

## Key Concepts

- **Client**: Uses client-go to fetch events from Kubernetes API
- **Event**: Normalized event struct with InvolvedObject, Source, Age
- **ConvertK8sEvent**: Converts corev1.Event to internal Event (shared by client and watch)
- **FormatAge**: Formats duration to human-readable short form (5s, 3m, 2h, 1d)
- **Filter**: Filters events by time, kind, name, type, reason, and sorts them newest first
- **GroupEvents**: Groups events by resource, namespace, kind, or reason; an empty mode means resource
- **GroupByResource**: Groups events by involved object (Kind/Name/Namespace)
- **Report**: Outputs color/plain/json/markdown/table summary
- **Watch**: Opens one watch per `--namespace` (one cluster-wide watch when none is given or with `--all-namespaces`) and fans them into one stream; each event prints on its own as it arrives, in the shape the README Output Formats section describes

<br/>

## CLI Flags

Flags are registered in `cmd/cli/root.go`, and the README Flags table is the user-facing reference. Rules that are easy to break:

- `--namespace` is repeatable, in both list and watch mode.
- `--group-by` and `--summary-only` do nothing with `--watch`, but an invalid `--group-by` is still rejected there with the same message as when listing.

<br/>

## Project Structure

One role per directory:

- `cmd`: `main.go` entry point
- `cmd/cli`: Cobra command, flag wiring, the list path (`run.go`) and the watch path (`watch.go`)
- `internal/client`: client-go wrapper and kubeconfig loading
- `internal/event`: event model, conversion, filtering and grouping
- `internal/report`: output formatters
- `scripts/`: shell helpers, each run by a Makefile target

<br/>

## Important Rules

- After changing code or tests, check the docs and update them:
  - `README.md`: Quick Start, installation, flags and output formats
  - `CLAUDE.md`: Key Concepts and the CLI Flags rules
- Do not edit `CHANGELOG.md` by hand. `.github/workflows/changelog-generator.yml` regenerates it after a merged PR or a release.
