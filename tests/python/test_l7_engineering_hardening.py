"""Tests for L7 Engineering stub."""

from python.l7_engineering_hardening.engineering_engine import (
    EngineeringEngineStub,
    EngineeringChange,
)


def test_create_change_from_finding():
    engine = EngineeringEngineStub()
    finding = {
        "id": "f1",
        "description": "Missing MFA",
        "owner": "iam-team",
        "affected_assets": ["asset-1"],
    }
    change = engine.create_change_from_finding(finding, "identity_risks", "high")
    assert "change_id" in change
    assert change["finding_source"] == "identity_risks"
    assert change["priority"] == "high"
    assert change["status"] == "pending"
    assert "due_date" in change
