"""Tests for L1 Asset Surface stub."""

from python.l1_asset_surface.asset_surface import AssetSurfaceStub


def test_stub_snapshot():
    stub = AssetSurfaceStub()
    snapshot = stub.generate_stub_snapshot()
    assert "snapshot_id" in snapshot
    assert "timestamp" in snapshot
    assert "assets" in snapshot
    assert "identities" in snapshot
    assert len(snapshot["assets"]) >= 1
    assert snapshot["assets"]["asset-1"]["criticality"] == "high"
