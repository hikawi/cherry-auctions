# CherryAuctions Backend

## Summary

This document aims to provide technical insights like a repository README file
for the backend section of this project. This does not go into reasons by certain
technical decisions are made, merely as a `CONTRIBUTING.md` or an `OVERVIEW.md`
file.

## Requirements

Technical Requirements:

- Golang Installed (at least v1.25.4)
- Gin Swagger Installed (use `go install`)
- Atlas, optional, mostly for inspecting, migrating and viewing database relationships.

For Atlas:

```bash
# MacOS/Linux
curl -sSf https://atlasgo.sh | sh

# Homebrew
brew install ariga/tap/atlas
```

Or get the pre-built binaries on each project's GitHub releases page. For Windows,
this is usually the best choice, get the `.exe` file and put a `PATH` variable
pointing to the `.exe`'s directory.
