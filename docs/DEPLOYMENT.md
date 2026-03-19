# Deployment

Guide for releasing and distributing kube-events.

<br/>

## Release Flow

A single tag push triggers the entire release pipeline automatically:

```
git tag v1.0.0 && git push origin v1.0.0
    └→ GitHub Actions (release.yml)
        └→ GoReleaser
            ├→ GitHub Releases (linux/darwin/windows x amd64/arm64)
            ├→ Homebrew tap update (somaz94/homebrew-tap)
            └→ Krew manifest update (somaz94/krew-index)
```

<br/>

## Distribution Channels

### 1. GitHub Releases (Default)

Automatically built by GoReleaser.

```bash
curl -sL https://github.com/somaz94/kube-events/releases/latest/download/kube-events_linux_amd64.tar.gz | tar xz
sudo mv kube-events /usr/local/bin/
```

**Supported platforms:**

| OS | Architecture |
|----|-------------|
| Linux | amd64, arm64 |
| macOS (Darwin) | amd64, arm64 |
| Windows | amd64, arm64 |

<br/>

### 2. Homebrew (macOS / Linux)

```bash
brew install somaz94/tap/kube-events
```

<br/>

### 3. Krew (kubectl plugin)

```bash
kubectl krew install events2
kubectl events2
```

<br/>

### 4. Go Install

```bash
go install github.com/somaz94/kube-events/cmd@latest
```

<br/>

## Secrets Configuration

| Secret | Purpose | Scope |
|--------|---------|-------|
| `PAT_TOKEN` | Cross-repo write access (Homebrew tap, Krew index) | GoReleaser |
| `GITHUB_TOKEN` | Release creation, dependabot auto-merge | Auto-provided |
| `GITLAB_TOKEN` | GitLab mirror backup | GitLab personal access token |

<br/>

## Step-by-Step: First Release

1. Ensure all tests pass:
   ```bash
   make test
   make build
   ```

2. Create and push a tag:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

3. GitHub Actions triggers `release.yml` → GoReleaser runs automatically

4. Verify:
   - Check [GitHub Releases](https://github.com/somaz94/kube-events/releases) for binaries
   - Check `somaz94/homebrew-tap` for formula commit
   - Check `somaz94/krew-index` for manifest update
