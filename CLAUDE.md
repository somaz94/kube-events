# CLAUDE.md - kube-events

CLI tool to view and summarize Kubernetes events with resource grouping and warning highlighting.

## Build & Test

```bash
make build           # Build binary
make test            # Run unit tests (alias for test-unit)
make test-unit       # go test ./... -v -race -cover
make cover           # Generate coverage report
make cover-html      # Open coverage in browser
make lint            # golangci-lint
make fmt             # go fmt
make vet             # go vet
```

## Commit Guidelines

- Do not include `Co-Authored-By` lines in commit messages.
- Use Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `ci:`, `chore:`)
- Do not push to remote. Only commit. The user will push manually.

## Key Concepts

- **Client**: Uses client-go to fetch events from Kubernetes API
- **Event**: Normalized event struct with InvolvedObject, Source, Age
- **Filter**: Filters events by time, kind, name, type, reason
- **GroupByResource**: Groups events by involved object (Kind/Name/Namespace)
- **Report**: Outputs color/plain/json/markdown/table summary

## CLI Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--kubeconfig` | | `~/.kube/config` | Path to kubeconfig |
| `--context` | | current | Kubernetes context |
| `--namespace` | `-n` | all | Filter by namespace |
| `--kind` | `-k` | all | Filter by involved object kind |
| `--name` | `-N` | all | Filter by involved object name |
| `--type` | `-t` | all | Filter by event type (Normal, Warning) |
| `--reason` | `-r` | all | Filter by reason (BackOff, Unhealthy, etc.) |
| `--since` | | `1h` | Show events newer than duration |
| `--output` | `-o` | `color` | Output format |
| `--summary-only` | `-s` | `false` | Summary statistics only |
| `--all-namespaces` | | `false` | All namespaces |
| `--watch` | `-w` | `false` | Watch for new events |

## Project Structure

```
cmd/
  main.go              # Entry point
  cli/
    root.go            # Cobra root command + global flags
    run.go             # Core execution logic
    version.go         # Version subcommand
internal/
  client/
    client.go          # Kubernetes client wrapper
  event/
    types.go           # Event data model
    filter.go          # Filtering and grouping logic
  report/
    summary.go         # Output formatters (color/plain/json/markdown/table)
```

## Important Rules

- **코드/테스트 수정 후 반드시 관련 문서를 확인하고 업데이트할 것.**
  - `README.md` — Quick Start, 설치 방법
  - `CHANGELOG.md` — Unreleased 섹션에 변경사항 추가
  - `CLAUDE.md` — Key Concepts, CLI Flags 테이블

## Language

- Communicate with the user in Korean.
- All documentation and code comments must be written in English.
