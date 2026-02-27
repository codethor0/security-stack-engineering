#!/usr/bin/env bash
# One-shot init script for security-stack-engineering repo
# Run from repo root after clone

set -e

echo "Bootstrapping security-stack-engineering..."

# Ensure Go deps
if command -v go &>/dev/null; then
    go mod tidy
    go build ./...
    echo "Go build OK"
fi

# Ensure Python structure
if command -v python3 &>/dev/null; then
    python3 -c "import sys; assert sys.version_info >= (3, 11)" 2>/dev/null || true
fi

# Make scripts executable
chmod +x scripts/*.sh 2>/dev/null || true

echo "Bootstrap complete"
