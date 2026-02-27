"""L1 Asset & Attack Surface algorithm."""

import asyncio
import hashlib
from datetime import datetime, timezone
from typing import Dict, List, Any, Optional

from .models import Asset, Identity, AssetGap


class AssetSurfaceAlgorithm:
    """L1 algorithm: maintains EnvironmentSnapshot, asset gaps, attack surface."""

    def __init__(
        self,
        message_bus: Any = None,
        evidence_store: Any = None,
        config: Optional[Dict] = None,
    ):
        self.bus = message_bus
        self.store = evidence_store
        self.config = config or {}
        self.governance_context: Optional[Dict] = None

    async def run(self, interval_seconds: int = 60) -> None:
        """Main loop: produce snapshots at interval."""
        while True:
            await self.execute_cycle()
            await asyncio.sleep(interval_seconds)

    async def execute_cycle(self) -> None:
        """Single cycle: ingest, normalize, derive surface, identify gaps, publish."""
        if self.bus and hasattr(self.bus, "consume_latest"):
            self.governance_context = await self.bus.consume_latest(
                "topic/governance_tokens", timeout=5
            )
        elif self.bus and hasattr(self.bus, "get_latest"):
            self.governance_context = self.bus.get_latest("topic/governance_tokens")

        if not self.governance_context:
            self._emit_evidence("governance_context_missing", "warning")
            return

        raw_assets = self._ingest_asset_data()
        raw_identities = self._ingest_identity_data()
        assets = self._normalize_assets(raw_assets)
        identities = self._normalize_identities(raw_identities)
        attack_surface = self._derive_attack_surface(assets, identities)
        snapshot = self._build_snapshot(assets, identities, attack_surface)
        gaps = self._identify_asset_gaps(snapshot)

        if not self._validate_invariants(snapshot):
            self._emit_evidence("invariant_violation", "critical")
            return

        if self.bus:
            self.bus.publish("topic/environment_snapshots", snapshot)
            for gap in gaps:
                self.bus.publish("topic/asset_gaps", gap)

        self._emit_evidence("cycle_completed", "info", {"asset_count": len(assets), "gap_count": len(gaps)})

    def _ingest_asset_data(self) -> List[Dict]:
        return [
            {
                "asset_id": "asset-1",
                "asset_type": "host",
                "owner": "platform-team",
                "criticality": "high",
                "business_processes": ["auth"],
                "data_classification": "internal",
                "exposure": "internal",
                "last_seen": datetime.now(timezone.utc).isoformat(),
            }
        ]

    def _ingest_identity_data(self) -> List[Dict]:
        return [
            {
                "identity_id": "id-1",
                "identity_type": "user",
                "owner": "admin",
                "privileged": True,
                "last_auth": datetime.now(timezone.utc).isoformat(),
                "orphaned": False,
            }
        ]

    def _normalize_assets(self, raw: List[Dict]) -> List[Dict]:
        return [{"asset_id": r.get("asset_id", ""), **r} for r in raw]

    def _normalize_identities(self, raw: List[Dict]) -> List[Dict]:
        return [{"identity_id": r.get("identity_id", ""), **r} for r in raw]

    def _derive_attack_surface(
        self, assets: List[Dict], identities: List[Dict]
    ) -> List[Dict]:
        surface = []
        for a in assets:
            if a.get("criticality") == "critical" and a.get("exposure") == "internet":
                surface.append({
                    "entry_point": a.get("asset_id"),
                    "paths": [],
                    "blast_radius": 1,
                })
        return surface

    def _build_snapshot(
        self,
        assets: List[Dict],
        identities: List[Dict],
        attack_surface: List[Dict],
    ) -> Dict:
        now = datetime.now(timezone.utc).isoformat()
        snapshot_id = hashlib.sha256(now.encode()).hexdigest()[:16]
        return {
            "snapshot_id": snapshot_id,
            "timestamp": now,
            "assets": {a["asset_id"]: a for a in assets},
            "identities": {i["identity_id"]: i for i in identities},
            "data_flows": [],
            "attack_surface": attack_surface,
        }

    def _identify_asset_gaps(self, snapshot: Dict) -> List[Dict]:
        gaps = []
        critical_processes = self.governance_context.get("critical_processes", [])
        covered = set()
        for a in snapshot.get("assets", {}).values():
            for p in a.get("business_processes", []):
                covered.add(p)
        for proc in critical_processes:
            if proc not in covered:
                gaps.append({
                    "gap_type": "missing_critical_coverage",
                    "severity": "critical",
                    "description": f"No assets mapped to critical process: {proc}",
                    "remediation_owner": "asset_management_team",
                })
        for i in snapshot.get("identities", {}).values():
            if i.get("orphaned") and i.get("privileged"):
                gaps.append({
                    "gap_type": "orphaned_privileged_identity",
                    "severity": "high",
                    "description": f"Privileged identity {i.get('identity_id')} has no owner",
                    "remediation_owner": "iam_team",
                })
        return gaps

    def _validate_invariants(self, snapshot: Dict) -> bool:
        for a in snapshot.get("assets", {}).values():
            last_seen = a.get("last_seen", "")
            if last_seen:
                try:
                    dt = datetime.fromisoformat(last_seen.replace("Z", "+00:00"))
                    age = datetime.now(timezone.utc) - dt
                    if age.days > 7 and a.get("criticality") == "critical":
                        return False
                except (ValueError, TypeError):
                    pass
        return True

    def _emit_evidence(self, action: str, level: str, details: Optional[Dict] = None) -> None:
        record = {
            "layer": "L1",
            "action": action,
            "level": level,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "details": details or {},
        }
        if self.bus and hasattr(self.bus, "publish"):
            self.bus.publish("topic/evidence_records", record)


class AssetSurfaceStub(AssetSurfaceAlgorithm):
    """Backward-compatible stub alias."""

    def generate_stub_snapshot(self) -> Dict:
        self.governance_context = {"critical_processes": []}
        raw_a = self._ingest_asset_data()
        raw_i = self._ingest_identity_data()
        assets = self._normalize_assets(raw_a)
        identities = self._normalize_identities(raw_i)
        surface = self._derive_attack_surface(assets, identities)
        return self._build_snapshot(assets, identities, surface)
