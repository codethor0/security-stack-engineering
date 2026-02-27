#!/usr/bin/env bash
# Run Pathway C: full L0-L8 (same as B; Python layers run via run_python_demo.sh)

set -e
cd "$(dirname "$0")/.."
go run ./cmd/orchestrator/main.go -pathway=c
