#!/usr/bin/env bash
# Run Pathway B: L0 + L1(stub) + L2 + L3(stub) + L4 + L5 + L6 + L8

set -e
cd "$(dirname "$0")/.."
go run ./cmd/orchestrator/main.go -pathway=b
