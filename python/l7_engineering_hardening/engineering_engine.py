"""L7 Engineering & Hardening algorithm."""

import asyncio
import hashlib
from dataclasses import dataclass, asdict
from typing import Dict, List, Any, Optional
from datetime import datetime, timezone, timedelta
from enum import Enum


class ChangeType(Enum):
    PATCH = "patch"
    CONFIG_CHANGE = "config_change"
    IAM_POLICY = "iam_policy"
    IaC_UPDATE = "iac_update"
    CODE_FIX = "code_fix"


class ChangePriority(Enum):
    CRITICAL = "critical"
    HIGH = "high"
    MEDIUM = "medium"
    LOW = "low"


@dataclass
class EngineeringChange:
    change_id: str
    finding_source: str
    finding_ids: List[str]
    change_type: str
    priority: str
    description: str
    affected_assets: List[str]
    implementation_plan: Dict[str, Any]
    owner: str
    due_date: str
    status: str
    evidence_refs: List[str]
    rollback_plan: Optional[Dict[str, Any]] = None


class EngineeringHardeningEngine:
    """L7 algorithm: aggregate findings, create changes, implement, validate."""

    def __init__(
        self,
        message_bus: Any = None,
        evidence_store: Any = None,
        config: Optional[Dict] = None,
        ci_cd: Any = None,
    ):
        self.bus = message_bus
        self.store = evidence_store
        self.config = config or {}
        self.ci_cd = ci_cd
        self.pending_changes: Dict[str, EngineeringChange] = {}
        self.governance_context: Optional[Dict] = None

    async def run(self) -> None:
        """Run finding aggregation and implementation loops."""
        await asyncio.gather(
            self.finding_aggregation_loop(60),
            self.implementation_loop(30),
        )

    async def finding_aggregation_loop(self, interval_seconds: int) -> None:
        """Aggregate findings from all layer topics."""
        while True:
            await asyncio.sleep(interval_seconds)
            if self.bus and hasattr(self.bus, "get_latest"):
                self.governance_context = self.bus.get_latest("topic/governance_tokens")
            findings = self._collect_findings()
            prioritized = self._prioritize_findings(findings)
            for item in prioritized[:5]:
                change = self._create_change_for_finding(item)
                if change:
                    self.pending_changes[change.change_id] = change
                    if self.bus:
                        self.bus.publish("topic/engineering_changes", asdict(change))

    def _collect_findings(self) -> Dict[str, List[Dict]]:
        return {
            "asset_gaps": [{"id": "ag1", "severity": "high", "description": "Stub gap"}],
            "detection_gaps": [],
            "identity_risks": [],
        }

    def _prioritize_findings(self, findings: Dict) -> List[Dict]:
        scored = []
        for source, items in findings.items():
            for item in items:
                score = 7.0 if item.get("severity") == "high" else 4.0
                scored.append({"finding": item, "source": source, "score": score})
        scored.sort(key=lambda x: x["score"], reverse=True)
        return scored

    def _create_change_for_finding(self, item: Dict) -> Optional[EngineeringChange]:
        finding = item["finding"]
        source = item["source"]
        score = item["score"]
        priority = "high" if score >= 7 else "medium"
        change_id = hashlib.sha256(
            f"{datetime.now(timezone.utc).isoformat()}{source}".encode()
        ).hexdigest()[:16]
        days = {"critical": 1, "high": 7, "medium": 30, "low": 90}
        due_date = datetime.now(timezone.utc) + timedelta(days=days.get(priority, 30))
        return EngineeringChange(
            change_id=change_id,
            finding_source=source,
            finding_ids=[finding.get("id", "unknown")],
            change_type=ChangeType.CONFIG_CHANGE.value,
            priority=priority,
            description=finding.get("description", "Remediation required"),
            affected_assets=finding.get("affected_assets", []),
            implementation_plan={"steps": ["validate", "apply", "verify"]},
            owner=finding.get("owner", "security-team"),
            due_date=due_date.isoformat(),
            status="pending",
            evidence_refs=finding.get("evidence_refs", []),
        )

    async def implementation_loop(self, interval_seconds: int) -> None:
        """Process approved changes (stub: mark completed)."""
        while True:
            await asyncio.sleep(interval_seconds)
            for cid, change in list(self.pending_changes.items()):
                if change.status == "pending":
                    change.status = "completed"
                    if self.bus:
                        self.bus.publish("topic/evidence_records", {
                            "layer": "L7",
                            "action": "change_completed",
                            "change_id": cid,
                            "timestamp": datetime.now(timezone.utc).isoformat(),
                        })
                    del self.pending_changes[cid]

    def create_change_from_finding(
        self,
        finding: Dict,
        source: str,
        priority: str = "medium",
    ) -> Dict:
        """Synchronous helper for tests."""
        item = {"finding": finding, "source": source, "score": 7.0 if priority == "high" else 4.0}
        change = self._create_change_for_finding(item)
        if change:
            self.pending_changes[change.change_id] = change
            return asdict(change)
        return {}


class EngineeringEngineStub(EngineeringHardeningEngine):
    """Backward-compatible stub alias."""

    pass
