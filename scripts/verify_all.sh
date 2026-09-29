#!/usr/bin/env bash
# Verification script: run all checks and report status.
# Every check runs; failures are counted and the summary always prints.
set -u
cd "$(dirname "$0")/.."

echo "=== Security Stack Engineering Verification ==="
echo ""

PASS=0
FAIL=0

check() {
    if "$@" > /dev/null 2>&1; then
        echo "[OK] $1"
        PASS=$((PASS + 1))
        return 0
    else
        echo "[FAIL] $1"
        FAIL=$((FAIL + 1))
        return 1
    fi
}

echo "1. Go build"
check go build ./...

echo "2. Go tests"
check go test ./...

echo "3. Python imports"
check bash -c 'PYTHONPATH=. python3 -c "
from python.l1_asset_surface.asset_surface import AssetSurfaceAlgorithm
from python.l3_telemetry_fabric.telemetry_fabric import TelemetryFabricStub
from python.l5_detection_analytics.detection_engine import DetectionEngineStub
from python.l7_engineering_hardening.engineering_engine import EngineeringEngineStub
"'

echo "4. Attribution (README)"
if grep -q "Thor Thor" README.md && grep -q "codethor0" README.md; then
    echo "[OK] README attribution"
    PASS=$((PASS + 1))
else
    echo "[FAIL] README attribution"
    FAIL=$((FAIL + 1))
fi

echo "5. Patent notice (README)"
if grep -q "patent filings" README.md; then
    echo "[OK] Patent notice present"
    PASS=$((PASS + 1))
else
    echo "[FAIL] Patent notice missing"
    FAIL=$((FAIL + 1))
fi

echo "6. Substack link (README)"
if grep -q "uniqueviolation.substack.com" README.md; then
    echo "[OK] Substack link correct"
    PASS=$((PASS + 1))
else
    echo "[FAIL] Substack link missing or wrong"
    FAIL=$((FAIL + 1))
fi

echo "7. Required files"
for f in README.md LICENSE CONTRIBUTING.md CODE_OF_CONDUCT.md go.mod .github/workflows/go-test.yml; do
    if [ -f "$f" ]; then
        echo "[OK] $f exists"
        PASS=$((PASS + 1))
    else
        echo "[FAIL] $f missing"
        FAIL=$((FAIL + 1))
    fi
done

echo ""
echo "=== Summary ==="
echo "Passed: $PASS"
echo "Failed: $FAIL"
echo ""

if [ "$FAIL" -gt 0 ]; then
    echo "Verification FAILED"
    exit 1
fi

echo "Verification PASSED"
exit 0
