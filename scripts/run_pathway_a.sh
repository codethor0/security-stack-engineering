#!/usr/bin/env bash
# Run Pathway A: L0 + L4 + L5 + L8
# No Python required. Paste-and-run.

set -e
cd "$(dirname "$0")/.."
go run ./cmd/orchestrator/main.go -pathway=a
