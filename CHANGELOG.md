# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Features

- add `--group-by` / `-g` flag to group events by resource, namespace, kind, or reason

### Tests

- add GroupEvents, ValidGroupBy, and group mode tests across all packages

### Documentation

- update all documentation for `--group-by` feature (README, USAGE, CONFIGURATION, EXAMPLES, CLAUDE.md)
- add grouping examples and demo phases

<br/>

## [v0.2.0](https://github.com/somaz94/kube-events/compare/v0.1.1...v0.2.0) (2026-03-19)

### Code Refactoring

- deduplicate event conversion, time formatting, and color constants ([28ea08b](https://github.com/somaz94/kube-events/commit/28ea08b9061ddb9a6c1420cc0ff685f94ab5615a))

### Documentation

- update documentation for refactoring and coverage improvements ([4423b03](https://github.com/somaz94/kube-events/commit/4423b036ae8671527d6af4f2b7aa3d769d4eccec))
- update changelog ([dc30e36](https://github.com/somaz94/kube-events/commit/dc30e364fe2dc26ceb4cdd9a9bfac22dbe83452a))

### Tests

- improve test coverage across all packages ([55d5517](https://github.com/somaz94/kube-events/commit/55d5517d81de1569df4ebbaf43124826e1c50904))

### Contributors

- somaz

<br/>

## [v0.1.1](https://github.com/somaz94/kube-events/compare/v0.1.0...v0.1.1) (2026-03-19)

### Features

- add brew install caveats message ([b883493](https://github.com/somaz94/kube-events/commit/b883493dec750efb6bf796ea4cc67c242f308293))

### Bug Fixes

- align goreleaser config with kube-diff structure ([cdec6bf](https://github.com/somaz94/kube-events/commit/cdec6bf5ea952fa4b02936b076b0b7409f4b7c98))

### Documentation

- README.md ([b05bc61](https://github.com/somaz94/kube-events/commit/b05bc619eec15500794ef32f339804564433aa81))
- add no-push rule to CLAUDE.md ([cc72819](https://github.com/somaz94/kube-events/commit/cc72819b85913c24431c5e0d0c3e02eac7fe0b4c))
- update changelog ([340200a](https://github.com/somaz94/kube-events/commit/340200a52417d43d1dd488d692a74dde3f0a6baa))

### Continuous Integration

- remove lint workflow ([8a3da34](https://github.com/somaz94/kube-events/commit/8a3da34c1c97c65c6157844ae6857c2bd70fb7b2))
- upgrade golangci-lint to v2.11.3 for Go 1.26 compatibility ([94d9b90](https://github.com/somaz94/kube-events/commit/94d9b90dd2d9aa0add844050f5a275026a9b2d12))
- enable lint workflow on push and pull_request triggers ([2c71f64](https://github.com/somaz94/kube-events/commit/2c71f649bcf41ad8a9cdef62d73a563c4be66de1))
- add e2e test workflow with kind cluster ([be368e0](https://github.com/somaz94/kube-events/commit/be368e0051bbe24d6ca31c4ef3d388454248059b))

### Contributors

- somaz

<br/>

## [v0.1.0](https://github.com/somaz94/kube-events/releases/tag/v0.1.0) (2026-03-19)

### Features

- add demo scripts, examples, and testdata ([9a71477](https://github.com/somaz94/kube-events/commit/9a71477176aa6b377ee4b638bd63f21794420b6c))
- initial project structure with CLI, tests, and documentation ([2a2ac0b](https://github.com/somaz94/kube-events/commit/2a2ac0bdcd359075ca8c48cd615b76baa547d662))

### Bug Fixes

- add missing krews short_description and brews metadata ([de21b9a](https://github.com/somaz94/kube-events/commit/de21b9add16ecef9ab64dbe05c4c1412e43f2b71))

### Documentation

- improve README badges, structure, and CLAUDE.md build commands ([327cf36](https://github.com/somaz94/kube-events/commit/327cf36a0f5cb9c1821b7338f3550ce7d73cd51b))
- docs/*.md ([468b5ae](https://github.com/somaz94/kube-events/commit/468b5ae7ea397677bb56558b2005e3ee4309b3da))

### Tests

- improve test coverage for client and cli packages ([b8b20c9](https://github.com/somaz94/kube-events/commit/b8b20c936e969f030ae399b33426f3711ac93fc5))

### Contributors

- somaz

<br/>

