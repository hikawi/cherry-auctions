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
- SQLC installed (you may use your distro's package manager or `go install`).
- Atlas, optional, mostly for inspecting, migrating and viewing database relationships.

For sqlc:

```bash
# MacOS:
brew install sqlc

# Ubuntu:
sudo snap install sqlc

# Using Go Install (should work on every platform)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Arch Linux:
sudo pacman -S sqlc
```

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
