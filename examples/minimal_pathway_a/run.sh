#!/usr/bin/env bash
# Run Pathway A from project root

set -e
cd "$(dirname "$0")/../.."
go run ./cmd/orchestrator/main.go -pathway=a
