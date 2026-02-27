"""L3 Telemetry & Data Fabric algorithm."""

import asyncio
import hashlib
import json
from dataclasses import dataclass, asdict
from typing import Dict, Any, Optional, List
from datetime import datetime, timezone
from collections import deque


@dataclass
class NormalizedEvent:
    event_id: str
    timestamp: str
    event_type: str
    source_type: str
    source_id: str
    actor: Dict[str, Any]
    target: Dict[str, Any]
    action: str
    outcome: str
    raw_hash: str
    enrichment: Dict[str, Any]
    engagement_id: Optional[str] = None


class TelemetryDataFabric:
    """L3 algorithm: normalize events, detect gaps, correlate evidence."""

    def __init__(
        self,
        message_bus: Any = None,
        evidence_store: Any = None,
        config: Optional[Dict] = None,
    ):
        self.bus = message_bus
        self.store = evidence_store
        self.config = config or {}
        self.event_buffer: deque = deque(maxlen=10000)
        self.normalization_rules = self._load_normalization_rules()

    def _load_normalization_rules(self) -> Dict:
        return {
            "o365": {"event_type": "auth", "actor_mapping": {"user": "userPrincipalName"}},
            "aws_cloudtrail": {"event_type": "cloud_api", "actor_mapping": {"user": "userIdentity"}},
            "stub": {"event_type": "unknown"},
        }

    async def run(self) -> None:
        """Run concurrent pipelines."""
        await asyncio.gather(
            self.ingestion_pipeline(),
            self.gap_detection_pipeline(30),
        )

    async def ingestion_pipeline(self) -> None:
        """Ingest and normalize events from stub sources."""
        sources = [{"name": "stub", "source_type": "stub"}]
        while True:
            for source in sources:
                raw_batch = self._poll_stub_source(source)
                for raw in raw_batch:
                    try:
                        normalized = self._normalize_event(raw, source["source_type"])
                        self.event_buffer.append(normalized)
                        if self.bus and hasattr(self.bus, "publish"):
                            self.bus.publish("topic/telemetry_view", {"events": [asdict(normalized)]})
                    except Exception as e:
                        if self.bus:
                            self.bus.publish("topic/telemetry_gaps", {
                                "gap_type": "parsing_failure",
                                "source_type": source["source_type"],
                                "severity": "medium",
                                "description": str(e),
                            })
            await asyncio.sleep(5)

    def _poll_stub_source(self, source: Dict) -> List[Dict]:
        return [{
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "event_type": "auth",
            "source_id": source["name"],
            "actor": {"identity_id": "u1"},
            "target": {"asset_id": "a1"},
            "action": "login",
            "outcome": "success",
            "engagement_id": None,
        }]

    def _normalize_event(self, raw: Dict, source_type: str) -> NormalizedEvent:
        rules = self.normalization_rules.get(source_type, {})
        raw_str = json.dumps(raw, sort_keys=True)
        raw_hash = hashlib.sha256(raw_str.encode()).hexdigest()[:32]
        event_id = hashlib.sha256(f"{raw_hash}{datetime.now(timezone.utc).isoformat()}".encode()).hexdigest()[:16]
        return NormalizedEvent(
            event_id=event_id,
            timestamp=raw.get("timestamp", datetime.now(timezone.utc).isoformat()),
            event_type=rules.get("event_type", raw.get("event_type", "unknown")),
            source_type=source_type,
            source_id=raw.get("source_id", "stub"),
            actor=raw.get("actor", {}),
            target=raw.get("target", {}),
            action=raw.get("action", "unknown"),
            outcome=raw.get("outcome", "unknown"),
            raw_hash=raw_hash,
            enrichment={},
            engagement_id=raw.get("engagement_id"),
        )

    async def gap_detection_pipeline(self, interval_seconds: int) -> None:
        """Detect telemetry gaps for critical assets."""
        while True:
            await asyncio.sleep(interval_seconds)
            critical_assets = [{"asset_id": "asset-1", "criticality": "high"}]
            for asset in critical_assets:
                recent = any(
                    e.target.get("asset_id") == asset["asset_id"]
                    for e in list(self.event_buffer)[-100:]
                )
                if not recent and self.bus:
                    self.bus.publish("topic/telemetry_gaps", {
                        "gap_type": "source_offline",
                        "source_type": "unknown",
                        "severity": "high",
                        "affected_assets": [asset["asset_id"]],
                        "description": f"No telemetry from {asset['asset_id']} in window",
                    })


class TelemetryFabricStub(TelemetryDataFabric):
    """Backward-compatible stub with sync normalize_event."""

    def normalize_event(self, raw: Dict, source_type: str) -> Dict:
        e = self._normalize_event(raw, source_type)
        return asdict(e)
