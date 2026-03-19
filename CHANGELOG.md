# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Features

- Initial release: CLI tool to view and summarize Kubernetes events
- Resource grouping by involved object (Pod, Deployment, etc.)
- Warning highlighting with color-coded output
- Time-based filtering (`--since`)
- Filter by namespace, kind, name, type, reason
- Multiple output formats: color, plain, JSON, markdown, table
- Watch mode for real-time event streaming (`-w`)
- Summary statistics
