# Changelog

All notable changes to this project will be documented in this file.

## Unreleased (2026-04-17)

### Chores

- **deps:** bump the go-minor group with 3 updates (#3) ([#3](https://github.com/somaz94/kube-events/pull/3)) ([4401033](https://github.com/somaz94/kube-events/commit/440103378a19c22fb247cf45576cf14147f51b76))
- **deps:** bump actions/github-script from 8 to 9 ([41e8401](https://github.com/somaz94/kube-events/commit/41e8401707fe3b74aaf36384f84b085688bbe093))
- **deps:** bump dependabot/fetch-metadata from 2 to 3 ([45589c7](https://github.com/somaz94/kube-events/commit/45589c74f7882da153de51a66db3a6874db1f98f))

<br/>

## [v0.3.1](https://github.com/somaz94/kube-events/compare/v0.3.0...v0.3.1) (2026-04-03)

### Features

- add branch and pr workflow targets to Makefile ([aa6d73c](https://github.com/somaz94/kube-events/commit/aa6d73cd7fb8e6f18879181c0fca0dd30ecb87ae))
- add Scoop bucket support for Windows distribution ([d438872](https://github.com/somaz94/kube-events/commit/d4388725476d647f27f50fe1b91240ca74c7dee8))

### Code Refactoring

- extract shared helpers, fix truncate boundary and parseSince error handling ([edf1c70](https://github.com/somaz94/kube-events/commit/edf1c70cad5ff390964791e89b6bf24e09ad1199))

### Documentation

- remove duplicate rules covered by global CLAUDE.md ([a634092](https://github.com/somaz94/kube-events/commit/a634092457e645f600192b46d44ebcba689e4f1a))
- add specific version install instructions to README ([a5fa59f](https://github.com/somaz94/kube-events/commit/a5fa59ffbcdc3928d409949e62001e07d0c6e54d))
- add Uninstall section to README ([5a7258a](https://github.com/somaz94/kube-events/commit/5a7258a6969df5e08dced6df0183899470835d2d))

### Continuous Integration

- add changelog category groups in goreleaser config ([32cd04a](https://github.com/somaz94/kube-events/commit/32cd04abf3d33505e5d07fcf2a00112a1726f657))
- add auto-generated PR body script for make pr ([5d83c2f](https://github.com/somaz94/kube-events/commit/5d83c2fe87b056fe6b0c3f3cec0c712d3aa8ae59))

### Chores

- remove duplicate rules from CLAUDE.md (moved to global) ([03a5efe](https://github.com/somaz94/kube-events/commit/03a5efee42a4c7101967973aada59cbd1f0da235))
- add git config protection to CLAUDE.md ([4957413](https://github.com/somaz94/kube-events/commit/49574137de4d9f5d304b2b7e96b1e88b31ecea06))
- **deps:** bump the go-minor group with 3 updates (#1) ([#1](https://github.com/somaz94/kube-events/pull/1)) ([c1c6221](https://github.com/somaz94/kube-events/commit/c1c6221ee7a1ebbc0f2fd824211dcdc9a3431c79))

### Contributors

- somaz

<br/>

## [v0.3.0](https://github.com/somaz94/kube-events/compare/v0.2.0...v0.3.0) (2026-03-19)

### Features

- add --group-by flag to group events by resource, namespace, kind, or reason ([35a76c7](https://github.com/somaz94/kube-events/commit/35a76c741d2da2c62d2db3846ccf6c4fe42d0a70))

### Documentation

- update documentation and demo for --group-by feature ([b2b37c2](https://github.com/somaz94/kube-events/commit/b2b37c2b8606f8fd95c45d92d0a3057cad920260))

### Tests

- add tests for --group-by feature across all packages ([361a518](https://github.com/somaz94/kube-events/commit/361a5184fbebeeb018b69c4c0a18834b1bf22050))

### Chores

- add --group-by example to brew install caveats ([465c6f3](https://github.com/somaz94/kube-events/commit/465c6f302f0f9583854e7f2f54e571fe64ae2ac7))

### Contributors

- somaz

<br/>

## [v0.2.0](https://github.com/somaz94/kube-events/compare/v0.1.1...v0.2.0) (2026-03-19)

### Code Refactoring

- deduplicate event conversion, time formatting, and color constants ([28ea08b](https://github.com/somaz94/kube-events/commit/28ea08b9061ddb9a6c1420cc0ff685f94ab5615a))

### Documentation

- update documentation for refactoring and coverage improvements ([4423b03](https://github.com/somaz94/kube-events/commit/4423b036ae8671527d6af4f2b7aa3d769d4eccec))

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

