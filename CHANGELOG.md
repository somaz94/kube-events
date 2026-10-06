# Changelog

All notable changes to this project will be documented in this file.

## [v0.3.5](https://github.com/somaz94/kube-events/compare/v0.3.4...v0.3.5) (2026-10-06)

### Bug Fixes

- run demo commands without eval so quoted paths survive ([2ffbaee](https://github.com/somaz94/kube-events/commit/2ffbaee1d5fc919afcba11d198b3417f35fa3e48))
- report an invalid --since identically in list and watch mode ([1c2f411](https://github.com/somaz94/kube-events/commit/1c2f411af8ff2887db420664261def76e559dd5b))

### Code Refactoring

- drop the no-op toUpper on --type values ([d5d7549](https://github.com/somaz94/kube-events/commit/d5d7549d5709fa8dc1cfb24dd4677c3c4b32c6a6))

### Continuous Integration

- pin kind and scope the e2e workflow token and triggers ([707f974](https://github.com/somaz94/kube-events/commit/707f9747ce567a7e03eee7736f2661e8e3410368))
- run golangci-lint in CI ([67d03a2](https://github.com/somaz94/kube-events/commit/67d03a28c485b6bbbb180ff91a9ba2102092a860))

### Chores

- classify scoped and breaking commits in create-pr.sh ([b52d980](https://github.com/somaz94/kube-events/commit/b52d980b27bebcf103cbfc3080e08280d519f26a))
- mark cover-html as phony ([1e313d1](https://github.com/somaz94/kube-events/commit/1e313d1820d999b33bcc128eeb5176ee033128f3))
- trim boilerplate comments in repo config and scripts ([80075cc](https://github.com/somaz94/kube-events/commit/80075cc511843f2df278fd7965d346fa6d625c8e))

### Contributors

- somaz

<br/>

## [v0.3.4](https://github.com/somaz94/kube-events/compare/v0.3.3...v0.3.4) (2026-10-06)

### Bug Fixes

- print uncolored watch lines for plain, markdown and table ([876d8dd](https://github.com/somaz94/kube-events/commit/876d8dd0a29b87b61d4c977dfd9e254a38b8bf71))
- label summary counts by group mode and truncate by rune ([c5fecf6](https://github.com/somaz94/kube-events/commit/c5fecf6e2b2048075caa701ed669691bc8c0c4c6))
- group by resource when --group-by is empty ([729dd72](https://github.com/somaz94/kube-events/commit/729dd7249d2cc71e68ca81925b481ed9aa2bad18))

### Documentation

- drop stale structure listings and fix the kubeconfig default ([4c26e40](https://github.com/somaz94/kube-events/commit/4c26e40dd80543fff7d5530a73448cdaf5fc3343))
- correct stale doc comments that predate --group-by ([0cf9f89](https://github.com/somaz94/kube-events/commit/0cf9f895ab85ee738c5606cf193f26794e2dec93))

### Continuous Integration

- trim redundant comments in gitlab-mirror workflow ([ab7c711](https://github.com/somaz94/kube-events/commit/ab7c71161b22020188b1be22e2d73a91e91c4dfa))
- retry mirror pushes on transient remote failures ([73e6c37](https://github.com/somaz94/kube-events/commit/73e6c373454ae778741eddc916a6520445fc16fc))
- drop the dead issue-close trigger from changelog generation ([2463fdc](https://github.com/somaz94/kube-events/commit/2463fdcc94fb49b304bd393d2aa1b640c23bafc5))

### Chores

- drop unprinted Makefile labels and redundant script comments ([fc5738b](https://github.com/somaz94/kube-events/commit/fc5738b4f56b95ed81e9b9a69c3fde6fbf07f147))
- tighten remaining comments in Go sources and tests ([ecd2002](https://github.com/somaz94/kube-events/commit/ecd20020e12b82c2458e4d7d218d28881e81434b))
- trim redundant comments in demo and e2e scripts ([a593ac8](https://github.com/somaz94/kube-events/commit/a593ac8fb75185e63617486b13936e554155602c))
- trim and tighten comments in Go sources and tests ([7508992](https://github.com/somaz94/kube-events/commit/7508992a1e25dedb11e9051c671c39accd2ff02f))
- **deps:** bump the go-minor group with 3 updates (#14) ([#14](https://github.com/somaz94/kube-events/pull/14)) ([800c71c](https://github.com/somaz94/kube-events/commit/800c71c7f2e6428c1382efc1db96e80aa6fda3b7))
- **deps:** bump the go-minor group with 2 updates (#13) ([#13](https://github.com/somaz94/kube-events/pull/13)) ([b007215](https://github.com/somaz94/kube-events/commit/b007215c2cdc26d71a6041be1d3e8ec688b56cc6))
- **deps:** bump the go-minor group with 3 updates (#12) ([#12](https://github.com/somaz94/kube-events/pull/12)) ([ae26d25](https://github.com/somaz94/kube-events/commit/ae26d25827c5343251ccb308e223fe8ad6723eff))

### Contributors

- somaz

<br/>

## [v0.3.3](https://github.com/somaz94/kube-events/compare/v0.3.2...v0.3.3) (2026-08-14)

### Continuous Integration

- add a golangci-lint config scoped to defect-finding linters ([a5b381f](https://github.com/somaz94/kube-events/commit/a5b381ff6c250135a51d847e9314fae83da2af7b))

### Chores

- publish the Homebrew package as a cask instead of a deprecated formula ([d358de7](https://github.com/somaz94/kube-events/commit/d358de765997f6e6c582d41dabf815cadcb5694e))

### Contributors

- somaz

<br/>

## [v0.3.2](https://github.com/somaz94/kube-events/compare/v0.3.1...v0.3.2) (2026-08-04)

### Bug Fixes

- validate --group-by in watch mode and document the flags it ignores ([cf819ea](https://github.com/somaz94/kube-events/commit/cf819ea1c5bc2fb60a2f4595b4fae12c9ae11f94))
- watch every namespace given to a repeatable --namespace, not just the first ([5865df4](https://github.com/somaz94/kube-events/commit/5865df493ffc7fd6026ab467289ceaad9369584f))

### Code Refactoring

- extract the watch event loop so it can be driven by a fake stream ([c1d3ccb](https://github.com/somaz94/kube-events/commit/c1d3ccb18f971229872b7571a0e0d84070d42e54))

### Documentation

- remove DCO sign-off instructions ([247aea2](https://github.com/somaz94/kube-events/commit/247aea2257f84337eb8ee09a6c4a1e28ecdf9b48))
- document DCO sign-off requirement in CONTRIBUTING ([40c53df](https://github.com/somaz94/kube-events/commit/40c53dfefb2b5c9214cb8226f82b4f0ededdf372))

### Tests

- make the watch e2e check re-runnable by waiting out terminating namespaces ([fc52ffb](https://github.com/somaz94/kube-events/commit/fc52ffb6e3b72a8cbaa19e6f4e82fda7b99d55e5))
- verify multi-namespace watch mode against a live cluster in e2e ([4efcfd6](https://github.com/somaz94/kube-events/commit/4efcfd66f89a691819f624ea72b2d17ba4f4661b))
- check the Close error in the stream test helper ([ec48be2](https://github.com/somaz94/kube-events/commit/ec48be2c80d39d6eb42fc0768e0ae412091b55c4))
- cover the CLI wiring paths and the pre-connection flag validation ([38b8a47](https://github.com/somaz94/kube-events/commit/38b8a478a5d6274825480815415bbf50a7f5f30d))

### Continuous Integration

- remove DCO workflow ([eb1df7b](https://github.com/somaz94/kube-events/commit/eb1df7b0ac57d393b8655055de74555b9caa38f8))
- adopt semantic-pr, labels, lock-threads, PR size, and auto-assign reusables ([0f6f12c](https://github.com/somaz94/kube-events/commit/0f6f12cfc186e46d3b1ec8a13e2dc9f839a186c7))
- use reusable stale-issues workflow ([8128077](https://github.com/somaz94/kube-events/commit/8128077fb68fed3624cc01d50e747732875e0bf1))
- use reusable issue-greeting workflow ([7bf0c4b](https://github.com/somaz94/kube-events/commit/7bf0c4bbac77d0d49b0dbfb407c1676b2445ca8f))
- use reusable dependabot-auto-merge workflow ([66b02f1](https://github.com/somaz94/kube-events/commit/66b02f15a1488f838f474f0fcd662681dabb2aed))
- use reusable contributors workflow ([216c5d4](https://github.com/somaz94/kube-events/commit/216c5d4f4fe528887399e9bfb658c279b18fb3d5))
- grant ok-to-test stub the permissions its reusable requires ([9a73cd7](https://github.com/somaz94/kube-events/commit/9a73cd77d128b9205b840efa51820a3b833531e7))
- add PR welcome and ok-to-test workflow stubs ([e6ca686](https://github.com/somaz94/kube-events/commit/e6ca686f0bdefcc6634ec62979feffd5123e1281))
- add DCO check via shared reusable workflow ([21f128b](https://github.com/somaz94/kube-events/commit/21f128bee7570573aac343651b1db5576b71fc18))
- add concurrency guards to recurring workflows ([e61314b](https://github.com/somaz94/kube-events/commit/e61314bd2dd8ed47e93e58c2ce5a808543d9de65))

### Chores

- **deps:** bump the go-minor group with 3 updates (#11) ([#11](https://github.com/somaz94/kube-events/pull/11)) ([c66896f](https://github.com/somaz94/kube-events/commit/c66896f551a545e23f3d14d718174a3b634760f3))
- **deps:** bump actions/setup-go from 6 to 7 (#10) ([#10](https://github.com/somaz94/kube-events/pull/10)) ([bb65650](https://github.com/somaz94/kube-events/commit/bb6565068ccda95bf3662b3bcf1a5a514537dd7a))
- **deps:** bump actions/checkout from 6 to 7 (#7) ([#7](https://github.com/somaz94/kube-events/pull/7)) ([6c7178c](https://github.com/somaz94/kube-events/commit/6c7178cdc98e73ec58594ec49bf9245a5a24baf4))
- **deps:** bump the go-minor group with 3 updates (#8) ([#8](https://github.com/somaz94/kube-events/pull/8)) ([99e903a](https://github.com/somaz94/kube-events/commit/99e903a426e07476c2bc167e40ed255e3c69fbe5))
- **deps:** bump the go-minor group with 3 updates (#6) ([#6](https://github.com/somaz94/kube-events/pull/6)) ([d4e3908](https://github.com/somaz94/kube-events/commit/d4e3908ce2952cf3d11910201dccc6801cb7d25c))
- **deps:** bump the go-minor group with 3 updates (#5) ([#5](https://github.com/somaz94/kube-events/pull/5)) ([e417e39](https://github.com/somaz94/kube-events/commit/e417e3975ecfbbc233de53c799e25d05a2b0a186))
- **deps:** bump the go-minor group with 3 updates (#3) ([#3](https://github.com/somaz94/kube-events/pull/3)) ([4401033](https://github.com/somaz94/kube-events/commit/440103378a19c22fb247cf45576cf14147f51b76))
- **deps:** bump actions/github-script from 8 to 9 ([41e8401](https://github.com/somaz94/kube-events/commit/41e8401707fe3b74aaf36384f84b085688bbe093))
- **deps:** bump dependabot/fetch-metadata from 2 to 3 ([45589c7](https://github.com/somaz94/kube-events/commit/45589c74f7882da153de51a66db3a6874db1f98f))

### Contributors

- somaz

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

