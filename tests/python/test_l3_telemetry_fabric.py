"""Tests for L3 Telemetry Fabric stub."""

from python.l3_telemetry_fabric.telemetry_fabric import TelemetryFabricStub


def test_normalize_event():
    stub = TelemetryFabricStub()
    raw = {
        "timestamp": "2025-01-01T12:00:00Z",
        "event_type": "auth",
        "actor": {"identity_id": "u1"},
        "target": {"asset_id": "a1"},
    }
    normalized = stub.normalize_event(raw, "o365")
    assert "event_id" in normalized
    assert normalized["event_type"] == "auth"
    assert normalized["source_type"] == "o365"
    assert "raw_hash" in normalized
